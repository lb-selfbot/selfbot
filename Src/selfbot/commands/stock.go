package commands

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/LightningDev1/LB-Selfbot/selfbot/embed"
	"github.com/LightningDev1/LB-Selfbot/selfbot/http"
	"github.com/LightningDev1/LB-Selfbot/selfbot/utils"
	"github.com/LightningDev1/dgc"
)

type StockCommands struct{}

func (s StockCommands) RegisterCommands(router *dgc.Router) {
	router.StartCategory("Stock", "Stock market commands")
	router.RegisterCmd(&dgc.Command{
		Name:        "searchsymbol",
		Description: "Search for a stock symbol",
		Usage:       "[p]searchsymbol <query>",
		Handler:     s.searchsymbol,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "symbolinfo",
		Description: "Get information about a stock symbol",
		Usage:       "[p]symbolinfo <symbol>",
		Handler:     s.symbolinfo,
	})

	stockSymbols := [][]string{
		{"AAPL", "Apple Inc.", "apple"},
		{"MSFT", "Microsoft Corporation", "microsoft"},
		{"AMZN", "Amazon.com, Inc.", "amazon"},
		{"META", "Meta Platforms, Inc.", "meta"},
		{"GOOGL", "Alphabet Inc.", "google"},
		{"TSLA", "Tesla, Inc.", "tesla"},
		{"BRK.A", "Berkshire Hathaway Inc.", "berkshirea"},
		{"BRK.B", "Berkshire Hathaway Inc.", "berkshireb"},
		{"JPM", "JPMorgan Chase & Co.", "jpmorgan"},
		{"JNJ", "Johnson & Johnson", "johnson"},
		{"V", "Visa Inc.", "visa"},
		{"PG", "Procter & Gamble Company", "procter"},
		{"UNH", "UnitedHealth Group Incorporated", "unitedhealth"},
		{"MA", "Mastercard Incorporated", "mastercard"},
		{"HD", "The Home Depot, Inc.", "homedepot"},
		{"DIS", "The Walt Disney Company", "disney"},
		{"VZ", "Verizon Communications Inc.", "verizon"},
		{"NVDA", "NVIDIA Corporation", "nvidia"},
		{"PYPL", "PayPal Holdings, Inc.", "paypal"},
		{"ADBE", "Adobe Inc.", "adobe"},
		{"INTC", "Intel Corporation", "intel"},
		{"CMCSA", "Comcast Corporation", "comcast"},
		{"NFLX", "Netflix, Inc.", "netflix"},
		{"PEP", "PepsiCo, Inc.", "pepsi"},
		{"KO", "The Coca-Cola Company", "cocacola"},
		{"T", "AT&T Inc.", "att"},
		{"CSCO", "Cisco Systems, Inc.", "cisco"},
		{"PFE", "Pfizer Inc.", "pfizer"},
		{"CRM", "Salesforce.com, Inc.", "salesforce"},
		{"ABT", "Abbott Laboratories", "abbott"},
		{"NKE", "NIKE, Inc.", "nike"},
		{"AVGO", "Broadcom Inc.", "broadcom"},
		{"COST", "Costco Wholesale Corporation", "costco"},
		{"TMO", "Thermo Fisher Scientific Inc.", "thermo"},
		{"MRK", "Merck & Co., Inc.", "merck"},
		{"ACN", "Accenture plc", "accenture"},
		{"ORCL", "Oracle Corporation", "oracle"},
		{"QCOM", "QUALCOMM Incorporated", "qualcomm"},
		{"MCD", "McDonald's Corporation", "mcdonalds"},
		{"WMT", "Walmart Inc.", "walmart"},
		{"BAC", "Bank of America Corporation", "bankofamerica"},
		{"MDT", "Medtronic plc", "medtronic"},
		{"UNP", "Union Pacific Corporation", "unionpacific"},
		{"LLY", "Eli Lilly and Company", "lilly"},
		{"NEE", "NextEra Energy, Inc.", "nextera"},
		{"AMGN", "Amgen Inc.", "amgen"},
		{"TXN", "Texas Instruments Incorporated", "texasinstruments"},
	}

	for _, stock := range stockSymbols {
		symbol := stock[0]
		fullName := stock[1]
		shortName := stock[2] + "stock"

		router.RegisterCmd(&dgc.Command{
			Name:        shortName,
			Description: fmt.Sprintf("Get %s stock info", fullName),
			Usage:       fmt.Sprintf("[p]%s", shortName),
			Handler:     s.generateStockInfoFunc(fullName, symbol),
		})
	}
}

func (StockCommands) generateStockInfoFunc(fullName, symbol string) func(ctx *dgc.Ctx) {
	return func(ctx *dgc.Ctx) {
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
			SetDescription(fmt.Sprintf("Stock info for %s", fullName))

		for key := range globalQuote.MustMap() {
			data := globalQuote.Get(key).MustString()

			key = strings.Split(key, ". ")[1]

			emb.AddField(utils.String.Capitalize(key), data)
		}

		emb.Send(ctx)
	}
}

func (StockCommands) searchsymbol(ctx *dgc.Ctx) {
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

func (s StockCommands) symbolinfo(ctx *dgc.Ctx) {
	s.generateStockInfoFunc(ctx.Arguments.Raw(), ctx.Arguments.Raw())(ctx)
}
