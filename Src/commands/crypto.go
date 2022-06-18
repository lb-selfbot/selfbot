package commands

import (
	"fmt"

	"github.com/LightningDev/LB-Selfbot/embed"
	"github.com/LightningDev/LB-Selfbot/http"
	"github.com/LightningDev1/dgc"
)

type CryptoCommands struct{}

func (t CryptoCommands) RegisterCommands(router *dgc.Router) {
	cryptoCommands := [][]string{
		{"btc", "BTC", "Bitcoin", "1"},
		{"eth", "ETH", "Ethereum", "1027"},
		{"bnb", "BNB", "Binance Coin", "1839"},
		{"dogecoin", "DOGE", "Dogecoin", "74"},
		{"tether", "USDT", "Tether", "825"},
		{"cardano", "ADA", "Cardano", "2010"},
		{"ripple", "XRP", "Ripple", "52"},
		{"polkadot", "DOT", "Polkadot", "6636"},
		{"bitcoincash", "BCH", "Bitcoin Cash", "1831"},
		{"litecoin", "LTC", "Litecoin", "2"},
		{"chainlink", "LINK", "Chainlink", "1975"},
		{"uniswap", "UNI", "UniSwap", "7083"},
		{"stellar", "XLM", "Stellar", "512"},
		{"usdcoin", "USDC", "USD Coin", "3408"},
		{"vechain", "VET", "VeChain", "3077"},
		{"ethereumclassic", "ETC", "Ehereum Classic", "3077"},
		{"solana", "SOL", "Solana", "5426"},
		{"theta", "THETA", "Theta", "2416"},
		{"shibainu", "SHIB", "Shiba Inu", "5994"},
		{"tron", "TRX", "TRON", "1958"},
		{"filecoin", "FIL", "FileCoin", "2280"},
		{"wrappedbitcoin", "WBTC", "Wrapped Bitcoin", "3717"},
		{"eos", "EOS", "EOS", "1765"},
		{"monero", "XMR", "Monero", "328"},
		{"neo", "NEO", "Neo", "1376"},
		{"binanceusd", "BUSD", "Binance USD", "4687"},
		{"bitcoinsv", "BSV", "Bitcoin SV", "3602"},
		{"terra", "LUNA", "Terra", "4172"},
		{"pancakeswap", "CAKE", "PancakeSwap", "7186"},
		{"iota", "MIOTA", "IOTA", "1720"},
		{"aave", "AAVE", "Aave", "7278"},
		{"ftxtoken", "FTT", "FTX Token", "4195"},
		{"cosmos", "ATOM", "Cosmos", "3794"},
		{"bittorrent", "BTT", "BitTorrent", "3718"},
		{"cryptocomcoin", "CRO", "Crypto.com Coin", "3635"},
		{"dash", "DASH", "Dash", "131"},
		{"zcash", "ZEC", "Zcash", "1437"},
		{"chiliz", "CHZ", "Chiliz", "4066"},
		{"digibyte", "DGB", "DigiByte", "109"},
		{"sushiswap", "SUSHI", "SushiSwap", "6758"},
		{"icon", "ICX", "ICON", "2099"},
		{"zrx", "ZRX", "0x", "1896"},
		{"ravencoin", "RVN", "RavenCoin", "2577"},
		{"nano", "NANO", "Nano", "1567"},
	}
	for _, element := range cryptoCommands {
		router.RegisterCmd(&dgc.Command{
			Name:        element[0],
			Description: "Sends the current " + element[2] + " price in USD and Euro",
			Usage:       "[p]" + element[0],
			Handler:     t.generateFunc(element),
		})
	}
}

func (t CryptoCommands) generateFunc(element []string) func(ctx *dgc.Ctx) {
	return func(ctx *dgc.Ctx) {
		json, err := http.GetJson("https://min-api.cryptocompare.com/data/price?fsym=" + element[1] + "&tsyms=USD,EUR")
		if err != nil {
			ctx.RespondText("An error has occurred.")
		}
		usd := json["USD"].(float64)
		eur := json["EUR"].(float64)
		embed.NewEmbed().
			SetDescription(fmt.Sprintf("USD: %.2f$\nEUR: %.2f€", usd, eur)).
			SetAuthor(element[2], fmt.Sprintf("https://s2.coinmarketcap.com/static/img/coins/200x200/%s.png", element[3])).
			Send(ctx)
	}
}
