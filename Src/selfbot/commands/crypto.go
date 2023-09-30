package commands

import (
	"fmt"
	"math"
	"net/url"
	"strings"

	"github.com/LightningDev1/LB-Selfbot/selfbot/embed"
	"github.com/LightningDev1/LB-Selfbot/selfbot/http"
	"github.com/LightningDev1/dgc"
)

type CryptoCommands struct{}

func (c CryptoCommands) RegisterCommands(router *dgc.Router) {
	router.StartCategory("Crypto", "Check cryptocurrency prices")

	router.RegisterCmd(&dgc.Command{
		Name:        "btcaddress",
		Description: "Get information about a Bitcoin address",
		Usage:       "[p]btcaddress <address>",
		Handler:     c.generateAddressFunc("Bitcoin", "BTC", 8),
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "btctransaction",
		Description: "Get information about a Bitcoin transaction",
		Usage:       "[p]btctransaction <hash>",
		Handler:     c.generateTransactionFunc("Bitcoin", "BTC", 8),
	})

	router.RegisterCmd(&dgc.Command{
		Name:        "ethaddress",
		Description: "Get information about an Ethereum address",
		Usage:       "[p]ethaddress <address>",
		Handler:     c.generateAddressFunc("Ethereum", "ETH", 18),
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "ethtransaction",
		Description: "Get information about an Ethereum transaction",
		Usage:       "[p]ethtransaction <hash>",
		Handler:     c.generateTransactionFunc("Ethereum", "ETH", 18),
	})

	router.RegisterCmd(&dgc.Command{
		Name:        "ltcaddress",
		Description: "Get information about a Litecoin address",
		Usage:       "[p]ltcaddress <address>",
		Handler:     c.generateAddressFunc("Litecoin", "LTC", 8),
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "ltctransaction",
		Description: "Get information about a Litecoin transaction",
		Usage:       "[p]ltctransaction <address>",
		Handler:     c.generateTransactionFunc("Litecoin", "LTC", 8),
	})

	router.RegisterCmd(&dgc.Command{
		Name:        "dashaddress",
		Description: "Get information about a Dash address",
		Usage:       "[p]dashaddress <address>",
		Handler:     c.generateAddressFunc("Dash", "DASH", 8),
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "dashtransaction",
		Description: "Get information about a Dash transaction",
		Usage:       "[p]dashtransaction <address>",
		Handler:     c.generateTransactionFunc("Dash", "DASH", 8),
	})

	router.RegisterCmd(&dgc.Command{
		Name:        "dogeaddress",
		Description: "Get information about a Dogecoin address",
		Usage:       "[p]dogeaddress <address>",
		Handler:     c.generateAddressFunc("Dogecoin", "DOGE", 8),
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "dogetransaction",
		Description: "Get information about a Dogecoin transaction",
		Usage:       "[p]dogetransaction <address>",
		Handler:     c.generateTransactionFunc("Dogecoin", "DOGE", 8),
	})

	cryptoCurrencies := [][]string{
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
		{"polygon", "MATIC", "Polygon", "3890"},
		{"avalanche", "AVAX", "Avalanche", "5805"},
		{"algorand", "ALGO", "Algorand", "4030"},
	}
	for _, element := range cryptoCurrencies {
		router.RegisterCmd(&dgc.Command{
			Name:        element[0],
			Description: "Sends the current " + element[2] + " price in USD and Euro",
			Usage:       "[p]" + element[0],
			Handler:     c.generatePriceFunc(element),
		})
	}
}

func (CryptoCommands) generateAddressFunc(name, shortCode string, power int) func(ctx *dgc.Ctx) {
	return func(ctx *dgc.Ctx) {
		address := ctx.Arguments.Get(0).Raw()
		if address == "" {
			_ = ctx.RespondText("Please provide an address.")
			return
		}

		url := fmt.Sprintf("https://api.blockcypher.com/v1/%s/main/addrs/%s", strings.ToLower(shortCode), url.PathEscape(address))

		httpResponse := http.Get(url)
		if httpResponse.Error != nil {
			_ = ctx.RespondText("An error has occurred.")
			return
		}

		json := httpResponse.Json()
		balance := float64(json.Get("final_balance").MustInt64()) / math.Pow10(power)
		transactions := fmt.Sprint(json.Get("final_n_tx").MustInt64())
		totalReceived := float64(json.Get("total_received").MustInt64()) / math.Pow10(power)
		totalSent := float64(json.Get("total_sent").MustInt64()) / math.Pow10(power)

		embed.NewEmbed().
			SetDescription(fmt.Sprintf("%s Address Info", name)).
			AddField("Address", address, true).
			AddField("Balance", fmt.Sprintf("%.8f %s", balance, shortCode), true).
			AddField("Transactions", transactions, true).
			AddField("Total Received", fmt.Sprintf("%.8f %s", totalReceived, shortCode), true).
			AddField("Total Sent", fmt.Sprintf("%.8f %s", totalSent, shortCode), true).
			Send(ctx)
	}
}

func (CryptoCommands) generateTransactionFunc(name, shortCode string, power int) func(ctx *dgc.Ctx) {
	return func(ctx *dgc.Ctx) {
		hash := ctx.Arguments.Get(0).Raw()
		if hash == "" {
			_ = ctx.RespondText("Please provide an transaction hash.")
			return
		}

		url := fmt.Sprintf("https://api.blockcypher.com/v1/%s/main/txs/%s", strings.ToLower(shortCode), url.PathEscape(hash))

		httpResponse := http.Get(url)
		if httpResponse.Error != nil {
			_ = ctx.RespondText("An error has occurred.")
			return
		}

		json := httpResponse.Json()
		total := float64(json.Get("total").MustInt64()) / math.Pow10(power)
		fees := float64(json.Get("fees").MustInt64()) / math.Pow10(power)
		confirmations := fmt.Sprint(json.Get("confirmations").MustInt64())
		received := json.Get("received").MustString()
		addresses := strings.Join(json.Get("addresses").MustStringArray(), ",\n")

		embed.NewEmbed().
			SetDescription(fmt.Sprintf("%s Transaction Info", name)).
			AddField("Transaction Hash", hash, true).
			AddField("Total Amount", fmt.Sprintf("%.8f %s", total, shortCode), true).
			AddField("Fees", fmt.Sprintf("%.8f %s", fees, shortCode), true).
			AddField("Confirmations", confirmations, true).
			AddField("Received At", received, true).
			AddField("Addresses", addresses, true).
			Send(ctx)
	}
}

func (CryptoCommands) generatePriceFunc(element []string) func(ctx *dgc.Ctx) {
	// This makes sure that the function is not overwritten by the next command.
	// If you put the function in the RegisterCommands function, the element
	// variable will always be the last element in the command array.
	return func(ctx *dgc.Ctx) {
		url := fmt.Sprintf("https://min-api.cryptocompare.com/data/price?fsym=%s&tsyms=USD,EUR", element[1])

		httpResponse := http.Get(url)
		if httpResponse.Error != nil {
			_ = ctx.RespondText("An error has occurred.")
			return
		}

		json := httpResponse.Json()
		usd := json.Get("USD").MustFloat64()
		eur := json.Get("EUR").MustFloat64()

		embed.NewEmbed().
			SetDescription(fmt.Sprintf("USD: %.2f$\nEUR: %.2f€", usd, eur)).
			SetAuthor(element[2], fmt.Sprintf("https://s2.coinmarketcap.com/static/img/coins/200x200/%s.png", element[3])).
			Send(ctx)
	}
}
