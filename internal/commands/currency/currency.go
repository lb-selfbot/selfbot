package currency

import (
	"fmt"
	"math"
	"net/url"
	"strconv"
	"strings"

	"github.com/lb-selfbot/selfbot/internal/http"
	"github.com/lb-selfbot/selfbot/internal/services/router"

	"github.com/lb-selfbot/selfbot/internal/services"
	"github.com/lb-selfbot/selfbot/internal/utils"
)

type CurrencyCommands struct{ services *services.Services }

func NewCurrencyCommands(s *services.Services) *CurrencyCommands {
	return &CurrencyCommands{services: s}
}

func (c *CurrencyCommands) Register() {
	c.services.Router.RegisterCategory(&router.Category{
		Name:        "Currency",
		Description: "Info about (crypto)currencies",
	}, []*router.Command{
		{
			Name:        "btcaddress",
			Description: "Get information about a Bitcoin address",
			Arguments:   "<address>",
			Handler:     c.generateAddressFunc("Bitcoin", "BTC", 8),
		},
		{
			Name:        "btctransaction",
			Description: "Get information about a Bitcoin transaction",
			Arguments:   "<hash>",
			Handler:     c.generateTransactionFunc("Bitcoin", "BTC", 8),
		},

		{
			Name:        "ethaddress",
			Description: "Get information about an Ethereum address",
			Arguments:   "<address>",
			Handler:     c.generateAddressFunc("Ethereum", "ETH", 18),
		},
		{
			Name:        "ethtransaction",
			Description: "Get information about an Ethereum transaction",
			Arguments:   "<hash>",
			Handler:     c.generateTransactionFunc("Ethereum", "ETH", 18),
		},

		{
			Name:        "ltcaddress",
			Description: "Get information about a Litecoin address",
			Arguments:   "<address>",
			Handler:     c.generateAddressFunc("Litecoin", "LTC", 8),
		},
		{
			Name:        "ltctransaction",
			Description: "Get information about a Litecoin transaction",
			Arguments:   "<address>",
			Handler:     c.generateTransactionFunc("Litecoin", "LTC", 8),
		},

		{
			Name:        "dashaddress",
			Description: "Get information about a Dash address",
			Arguments:   "<address>",
			Handler:     c.generateAddressFunc("Dash", "DASH", 8),
		},
		{
			Name:        "dashtransaction",
			Description: "Get information about a Dash transaction",
			Arguments:   "<address>",
			Handler:     c.generateTransactionFunc("Dash", "DASH", 8),
		},

		{
			Name:        "dogeaddress",
			Description: "Get information about a Dogecoin address",
			Arguments:   "<address>",
			Handler:     c.generateAddressFunc("Dogecoin", "DOGE", 8),
		},
		{
			Name:        "dogetransaction",
			Description: "Get information about a Dogecoin transaction",
			Arguments:   "<address>",
			Handler:     c.generateTransactionFunc("Dogecoin", "DOGE", 8),
		},

		{
			Name:        "priceinfo",
			Description: "Get price information from a (crypto)currency",
			Arguments:   "<from symbol> [to symbol]",
			Handler:     c.priceinfo,
		},

		{
			Name:        "searchsymbol",
			Description: "Search for a stock market symbol",
			Arguments:   "<query>",
			Handler:     c.searchsymbol,
		},
		{
			Name:        "symbolinfo",
			Description: "Get information about a stock market symbol",
			Arguments:   "<symbol>",
			Handler:     c.symbolinfo,
		},
	})
}

func (c *CurrencyCommands) generateAddressFunc(name, shortCode string, power int) func(ctx *router.Ctx) {
	return func(ctx *router.Ctx) {
		address := ctx.Arguments.Get(0).Raw()
		if address == "" {
			ctx.RespondText("Please provide an address.")
			return
		}

		url := fmt.Sprintf("https://api.blockcypher.com/v1/%s/main/addrs/%s", strings.ToLower(shortCode), url.PathEscape(address))

		httpResponse := http.Get(url)
		if httpResponse.Error != nil {
			ctx.RespondText("An error has occurred.")
			return
		}

		json := httpResponse.Json()
		balance := float64(json.Get("final_balance").MustInt64()) / math.Pow10(power)
		transactions := fmt.Sprint(json.Get("final_n_tx").MustInt64())
		totalReceived := float64(json.Get("total_received").MustInt64()) / math.Pow10(power)
		totalSent := float64(json.Get("total_sent").MustInt64()) / math.Pow10(power)

		ctx.NewEmbed().
			SetDescription("%s Address Info", name).
			AddField("Address", address, true).
			AddField("Balance", fmt.Sprintf("%.8f %s", balance, shortCode), true).
			AddField("Transactions", transactions, true).
			AddField("Total Received", fmt.Sprintf("%.8f %s", totalReceived, shortCode), true).
			AddField("Total Sent", fmt.Sprintf("%.8f %s", totalSent, shortCode), true).
			Send()
	}
}

func (c *CurrencyCommands) generateTransactionFunc(name, shortCode string, power int) func(ctx *router.Ctx) {
	return func(ctx *router.Ctx) {
		hash := ctx.Arguments.Get(0).Raw()
		if hash == "" {
			ctx.RespondText("Please provide an transaction hash.")
			return
		}

		url := fmt.Sprintf("https://api.blockcypher.com/v1/%s/main/txs/%s", strings.ToLower(shortCode), url.PathEscape(hash))

		httpResponse := http.Get(url)
		if httpResponse.Error != nil {
			ctx.RespondText("An error has occurred.")
			return
		}

		json := httpResponse.Json()
		total := float64(json.Get("total").MustInt64()) / math.Pow10(power)
		fees := float64(json.Get("fees").MustInt64()) / math.Pow10(power)
		confirmations := fmt.Sprint(json.Get("confirmations").MustInt64())
		received := json.Get("received").MustString()
		addresses := strings.Join(json.Get("addresses").MustStringArray(), ",\n")

		ctx.NewEmbed().
			SetDescription("%s Transaction Info", name).
			AddField("Transaction Hash", hash, true).
			AddField("Total Amount", fmt.Sprintf("%.8f %s", total, shortCode), true).
			AddField("Fees", fmt.Sprintf("%.8f %s", fees, shortCode), true).
			AddField("Confirmations", confirmations, true).
			AddField("Received At", received, true).
			AddField("Addresses", addresses, true).
			Send()
	}
}

func (c *CurrencyCommands) priceinfo(ctx *router.Ctx) {
	from, to := ctx.Arguments.Get(0).Raw(), ctx.Arguments.Get(1).Raw()
	if to == "" {
		to = "USD,EUR"
	}

	httpResponse := http.Get(fmt.Sprintf("https://min-api.cryptocompare.com/data/pricemultifull?fsyms=%s&tsyms=%s", from, to))
	if httpResponse.Error != nil {
		ctx.RespondText("An error has occurred: %v", httpResponse.Error)
		return
	}

	json := httpResponse.Json()

	if msg := json.Get("Message").MustString(); msg != "" {
		ctx.RespondText("An error has occurred: " + msg)
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

	emb := ctx.NewEmbed()

	for k, v := range parsed {
		emb.AddField(k, v)
	}

	emb.Send()
}

func (c *CurrencyCommands) searchsymbol(ctx *router.Ctx) {
	uri := fmt.Sprintf(
		"https://www.alphavantage.co/query?function=SYMBOL_SEARCH&keywords=%s&apikey=KBJ3RJS2JAYNH3WD",
		url.QueryEscape(ctx.Arguments.Raw()),
	)

	httpResponse := http.Get(uri)
	if httpResponse.Error != nil {
		ctx.RespondText("An error occurred: %v", httpResponse.Error)
		return
	}

	bestMatches := httpResponse.Json().Get("bestMatches")

	if len(bestMatches.MustArray()) == 0 {
		ctx.RespondText("No results found")
		return
	}

	emb := ctx.NewEmbed()

	for i := range bestMatches.MustArray() {
		stock := bestMatches.GetIndex(i)

		info := fmt.Sprintf(
			"Symbol: %s\nName: %s",
			stock.Get("1. symbol").MustString(),
			stock.Get("2. name").MustString(),
		)

		emb.AddField("Result "+strconv.Itoa(i+1), info)
	}

	emb.Send()
}

func (c *CurrencyCommands) symbolinfo(ctx *router.Ctx) {
	symbol := strings.ToUpper(ctx.Arguments.Raw())

	uri := fmt.Sprintf(
		"https://www.alphavantage.co/query?function=GLOBAL_QUOTE&symbol=%s&apikey=KBJ3RJS2JAYNH3WD",
		symbol,
	)

	httpResponse := http.Get(uri)
	if httpResponse.Error != nil {
		ctx.RespondText("An error occurred: %v", httpResponse.Error)
		return
	}

	globalQuote := httpResponse.Json().Get("Global Quote")

	if len(globalQuote.MustMap()) == 0 {
		ctx.RespondText("No results found")
		return
	}

	emb := ctx.NewEmbed().
		SetDescription("Stock info for %s", symbol)

	for key := range globalQuote.MustMap() {
		data := globalQuote.Get(key).MustString()

		key = strings.Split(key, ". ")[1]

		emb.AddField(utils.String.Capitalize(key), data)
	}

	emb.Send()
}
