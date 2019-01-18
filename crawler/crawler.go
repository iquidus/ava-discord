package crawler

import (
  "fmt"
  "net/http"
  "strconv"
	"time"
	"github.com/iquidus/ava-discord/util"
  . "github.com/iquidus/ava-discord/watchlist"
)

type Block struct {
	Number  uint64           `json:"number"`
}

type transaction struct {
  BlockNumber      uint64       `json:"blockNumber"`
  From             string       `json:"from"`
  Hash             string       `json:"hash"`
  Value            string       `json:"value"`
  Timestamp        uint64       `json:"timestamp"`
  Input            string       `json:"input"`
  Gas              uint64       `json:"gas"`
  GasUsed          uint64       `json:"gasUsed"`
  GasPrice         string       `json:"gasPrice"`
  Nonce            uint64       `json:"nonce"`
  TransactionIndex uint64       `json:"transactionIndex"`
  To               string       `json:"to"`
  ContractAddress  string       `json:"contractAddress"`
}

type FlaggedTxn struct {
  Hash         string
  Address      string
  Label        string
}

var client = &http.Client{Timeout: 60 * time.Second}

func Sync(height uint64, currentBlock uint64) ([]FlaggedTxn, uint64) {
  var flaggedTxns []FlaggedTxn

  var watchlist []Address
  watchlist = GetWatchlist()

  var txns []transaction
  localHead := currentBlock

  if currentBlock != height {
    for n := currentBlock + 1; n <= height; n++ {
      url := "https://v3.ubiqscan.io/blocktransactions/"
      url += strconv.FormatInt(int64(n), 10)
      err := util.GetJson(client, url, &txns)
      if err != nil {
        fmt.Println("Unable to get blocktransactions: ", err)
        return flaggedTxns, n
      }

      fmt.Println("blockNumber: ", n)

      for i := 0; i < len(txns); i++ {
        for x := 0; x < len(watchlist); x++ {
          if (txns[i].From == watchlist[x].Hash) {
            flaggedTxns = append(flaggedTxns, FlaggedTxn{Hash: txns[i].Hash, Address: watchlist[x].Hash, Label: watchlist[x].Label})
          }
        }
      }
      localHead = n
    }
  }
  return flaggedTxns, localHead
}
