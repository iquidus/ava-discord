package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/iquidus/ava-discord/util"
	. "github.com/iquidus/ava-discord/watchlist"
	. "github.com/iquidus/ava-discord/crawler"
)

var (
	height uint64
	currentBlock uint64
	client = &http.Client{Timeout: 60 * time.Second}
	// channel to broadcast alerts to
	broadcastChannelId = "535601803459428362"
)

type statusResponse struct {
	LatestBlock *Block       `json:"latestBlock"`
}

func handleMessage(s *discordgo.Session, m *discordgo.MessageCreate) *string {

	vals := &m.Content
	valSplit := strings.Split(*vals, " ")
	message := ""

	if len(*vals) == 0 {
		return nil
	}

	command := valSplit[0]
	// arguments := valSplit[1:]

	switch command {
	case "?commands":
		message = "Ava commands\n"
		message += "\n"
		message += "__**General**__\n\n"
		message += "`?watchlist` - Returns current watchlist\n"
		message += "`?info` - Returns Ava info"
	case "?info":
		message = "Ava monitors the Ubiq blockchain and alerts on certian activity\n"
		message += "github: https://github.com/iquidus/ava-discord"
	case "?watchlist":
		var watchlist []Address
		watchlist = GetWatchlist()
		message +="```"
		for i := 0; i < len(watchlist); i++ {
			message += watchlist[i].Hash
			message += " ("
			message += watchlist[i].Label
			message += ")\n"
		}
		message +="```"
	default:
	}

	return &message
}

func syncLoop(ds *discordgo.Session) {
	var status *statusResponse
  err := util.GetJson(client, "https://v3.ubiqscan.io/status", &status)
  if err != nil {
    fmt.Println("Unable to get status: ", err)
  }
  height = status.LatestBlock.Number
	fmt.Println("height:", height)
	if currentBlock == 0 {
		currentBlock = height - 10 // start from 10 blocks ago.
	}

  flagged, localHead := Sync(height, currentBlock)
	currentBlock = localHead
	if len(flagged) > 0 {
		for i := 0; i < len(flagged); i++ {
			message := *generateAlertMessage(flagged[i].Address, flagged[i].Label, flagged[i].Hash)
			fmt.Println(message)
			ds.ChannelMessageSend(broadcastChannelId, message)
		}
	}
}

func main() {
	//godotenv.Load(defaultConfigFile)
	token := os.Getenv("DISCORD_API_TOKEN")
  currentBlock = 0

	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(messageCreate)

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	ticker := time.NewTicker(60 * time.Second)
	quit := make(chan struct{})

	go func() {
		for {
			select {
				case <- ticker.C:
					go syncLoop(dg)
				case <- quit:
					ticker.Stop()
					return
			}
		}
	}()
	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	close(quit)
	dg.Close()
}

func generateAlertMessage(address string, label string, txid string) *string {
	message := ""
	prefix := "@everyone\n**ALERT**: Flagged address in use!```"
	suffix := "```**Ubiqscan**: https://ubiqscan.io/tx/"
	suffix += txid
	addressLine := "Address: "
	addressLine += address
	addressLine += " - "
	addressLine += label
	txnLine := "Txid   : "
	txnLine += txid
	message = fmt.Sprintf("%s\n%s\n%s\n%s", prefix, addressLine, txnLine, suffix)

	return &message
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the authenticated bot has access to.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}

	if len(m.Content) > 0 {
		message := handleMessage(s, m)
		s.ChannelMessageSend(m.ChannelID, *message)
	}
}
