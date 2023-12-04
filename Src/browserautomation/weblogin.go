package browserautomation

import (
	"context"
	"time"

	chromedp "github.com/LightningDev1/chromiumdp"
	"github.com/LightningDev1/discordgo"
)

// WebLogin is used to login to Discord using a browser
type WebLogin struct {
	// TokenCallback is called when the token is ready
	TokenCallback func(string) error

	ctx    context.Context
	cancel context.CancelFunc
	opts   []chromedp.ExecAllocatorOption

	token string
}

// Start starts the login process
func (login *WebLogin) Start() error {
	// Options for chrome.
	login.opts = append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("user-agent", discordgo.UserAgentDesktop),
		chromedp.Flag("headless", false),
		chromedp.Flag("enable-automation", false),
		chromedp.WindowSize(900, 700),
		chromedp.DisableGPU,
	)

	// Create ChromeDP context
	login.ctx, login.cancel = chromedp.NewExecAllocator(context.Background(), login.opts...)
	login.ctx, login.cancel = chromedp.NewContext(login.ctx)
	defer login.cancel()

	// Run browser automation
	return login.runAutomation()
}

// Run automation
func (login *WebLogin) runAutomation() error {
	err := chromedp.Run(login.ctx,
		// Navigate to discord login page
		chromedp.Navigate(`https://discord.com/login`),

		chromedp.Sleep(2*time.Second),

		// Wait until logged in
		chromedp.WaitVisible("[class^=panels]"),

		// Get the token
		chromedp.Evaluate(getTokenScript, &login.token),
		chromedp.ActionFunc(func(ctx context.Context) error {
			return login.TokenCallback(login.token)
		}),
	)
	if err != nil {
		return err
	}

	return chromedp.Cancel(login.ctx)
}
