package browser

import (
	"context"

	"github.com/LightningDev1/discordgo"
	"github.com/chromedp/chromedp"
)

var getTokenScript = `(webpackChunkdiscord_app.push([[''],{},e=>{m=[];for(let c in e.c)m.push(e.c[c])}]),m).find(m=>m?.exports?.default?.getToken!==void 0).exports.default.getToken()`

func RunWebLogin(ctx context.Context) (string, error) {
	browser := NewBrowser(ctx, BrowserOptions{
		UserAgent: discordgo.UserAgentDesktop,
		Headless:  false,
		Size:      [2]int{900, 700},
	})

	tokenChan := make(chan string)
	errorChan := make(chan error)

	go func() {
		token := ""
		err := browser.RunAutomation(
			// Navigate to discord login page
			chromedp.Navigate(`https://discord.com/login`),

			// Wait until logged in
			chromedp.Sleep(2),
			chromedp.WaitVisible("[class^=panels]"),

			// Get the token
			chromedp.Evaluate(getTokenScript, &token),
		)
		if err != nil {
			errorChan <- err
			return
		}

		tokenChan <- token
	}()

	select {
	case token := <-tokenChan:
		return token, nil
	case err := <-errorChan:
		return "", err
	case <-ctx.Done():
		return "", ctx.Err()
	}
}
