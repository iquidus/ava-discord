package watchlist

type Address struct {
  Hash  string
  Label string
}

func GetWatchlist() ([]Address) {
  var addresses []Address
  //addresses = append(addresses, Address{Hash: "0xda904bc07fd95e39661941b3f6daded1b8a38c71", Label: "Test"})
  addresses = append(addresses, Address{Hash: "0xde89c4687984d7cb91cacdd084003ffdf36e493a", Label: "Cryptopia - UBQ - OLD"})
  addresses = append(addresses, Address{Hash: "0x6b7bcaebcbe0b92f879cfe5ed2cdb34247d49f0d", Label: "Cryptopia - ERC20 - OLD"})
  addresses = append(addresses, Address{Hash: "0xabee6c9855af9202f995efa7eea46c7819eafe09", Label: "Cryptopia - UBQ"})
  addresses = append(addresses, Address{Hash: "0x81e8416fabcfb122964b61e24e8b005fb1c7081b", Label: "Cryptopia - ERC20"})
  return addresses
}
