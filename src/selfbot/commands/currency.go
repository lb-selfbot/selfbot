package commands

import (
	"fmt"
	"math"
	"net/url"
	"strconv"
	"strings"

	"github.com/LightningDev1/LB-Selfbot/internal/http"
	"github.com/LightningDev1/LB-Selfbot/internal/utils"
	"github.com/LightningDev1/LB-Selfbot/selfbot/embed"
	"github.com/LightningDev1/dgc"
)

type CurrencyCommands struct{}

func (c CurrencyCommands) RegisterCommands(router *dgc.Router) {
	router.StartCategory("Currency", "Information about (crypto)currencies")

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

	router.RegisterCmd(&dgc.Command{
		Name:        "priceinfo",
		Description: "Get price information from a (crypto)currency",
		Usage:       "[p]priceinfo <from symbol> [to symbol]",
		Handler:     c.priceinfo,
	})

	router.RegisterCmd(&dgc.Command{
		Name:        "searchsymbol",
		Description: "Search for a stock market symbol",
		Usage:       "[p]searchsymbol <query>",
		Handler:     c.searchsymbol,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "symbolinfo",
		Description: "Get information about a stock market symbol",
		Usage:       "[p]symbolinfo <symbol>",
		Handler:     c.symbolinfo,
	})
}

func (CurrencyCommands) generateAddressFunc(name, shortCode string, power int) func(ctx *dgc.Ctx) {
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

func (CurrencyCommands) generateTransactionFunc(name, shortCode string, power int) func(ctx *dgc.Ctx) {
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

func (CurrencyCommands) priceinfo(ctx *dgc.Ctx) {
	from, to := utils.Arguments["<from symbol>"], utils.Arguments["[to symbol]"]
	if to == "" {
		to = "USD,EUR"
	}

	httpResponse := http.Get(fmt.Sprintf("https://min-api.cryptocompare.com/data/pricemultifull?fsyms=%s&tsyms=%s", from, to))
	if httpResponse.Error != nil {
		_ = ctx.RespondText("An error has occurred: " + httpResponse.Error.Error())
		return
	}

	json := httpResponse.Json()

	if msg := json.Get("Message").MustString(); msg != "" {
		_ = ctx.RespondText("An error has occurred: " + msg)
		return
	}

	raw := json.Get("DISPLAY").MustMap()
	parsed := make(map[string]string)

	for from, data := range raw {
		data, ok := data.(map[string]any)
		if !ok {
			continue
		}

		for to, data := range data {
			data, ok := data.(map[string]any)
			if !ok {
				continue
			}

			parsed[from+"/"+to] = fmt.Sprintf(
				"Price: %s\nVolume (24h): %s (%s)\nChange (24h): %s (%s%%)\nMarket cap: %s",
				data["PRICE"],
				data["VOLUME24HOURTO"], data["VOLUME24HOUR"],
				data["CHANGE24HOUR"], data["CHANGEPCT24HOUR"],
				data["MKTCAP"],
			)
		}
	}

	emb := embed.NewEmbed()

	for k, v := range parsed {
		emb.AddField(k, v)
	}

	emb.Send(ctx)
}

func (CurrencyCommands) searchsymbol(ctx *dgc.Ctx) {
	uri := fmt.Sprintf(
		"https://www.alphavantage.co/query?function=SYMBOL_SEARCH&keywords=%s&apikey=KBJ3RJS2JAYNH3WD",
		url.QueryEscape(ctx.Arguments.Raw()),
	)

	httpResponse := http.Get(uri)
	if httpResponse.Error != nil {
		ctx.RespondText("An error occurred: " + httpResponse.Error.Error())
		return
	}

	bestMatches := httpResponse.Json().Get("bestMatches")

	if len(bestMatches.MustArray()) == 0 {
		ctx.RespondText("No results found")
		return
	}

	emb := embed.NewEmbed()

	for i := range bestMatches.MustArray() {
		stock := bestMatches.GetIndex(i)

		info := fmt.Sprintf(
			"Symbol: %s\nName: %s",
			stock.Get("1. symbol").MustString(),
			stock.Get("2. name").MustString(),
		)

		emb.AddField("Result "+strconv.Itoa(i+1), info)
	}

	emb.Send(ctx)
}

func (CurrencyCommands) symbolinfo(ctx *dgc.Ctx) {
	symbol := strings.ToUpper(ctx.Arguments.Raw())

	uri := fmt.Sprintf(
		"https://www.alphavantage.co/query?function=GLOBAL_QUOTE&symbol=%s&apikey=KBJ3RJS2JAYNH3WD",
		symbol,
	)

	httpResponse := http.Get(uri)
	if httpResponse.Error != nil {
		ctx.RespondText("An error occurred: " + httpResponse.Error.Error())
		return
	}

	globalQuote := httpResponse.Json().Get("Global Quote")

	if len(globalQuote.MustMap()) == 0 {
		ctx.RespondText("No results found")
		return
	}

	emb := embed.NewEmbed().
		SetDescription(fmt.Sprintf("Stock info for %s", symbol))

	for key := range globalQuote.MustMap() {
		data := globalQuote.Get(key).MustString()

		key = strings.Split(key, ". ")[1]

		emb.AddField(utils.String.Capitalize(key), data)
	}

	emb.Send(ctx)
}
