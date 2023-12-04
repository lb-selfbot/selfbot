package browserautomation

import (
	"context"
	"time"

	chromedp "github.com/LightningDev1/chromiumdp"
	"github.com/LightningDev1/discordgo"
)

// QRLogin is used to login to Discord using a QR code
type QRLogin struct {
	// VisibleBrowser determines if the browser should be visible or not
	VisibleBrowser bool

	// ImageReadyCallback is called when the QR code image is ready
	ImageReadyCallback func([]byte) error

	// TokenCallback is called when the token is ready
	TokenCallback func(string) error

	ctx    context.Context
	cancel context.CancelFunc
	opts   []chromedp.ExecAllocatorOption

	qrImage []byte
	token   string
}

// Start starts the login process
func (login *QRLogin) Start() error {
	// Options for chrome.
	login.opts = append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("user-agent", discordgo.UserAgentDesktop),
		chromedp.Flag("headless", !login.VisibleBrowser),
		chromedp.Flag("enable-automation", false),
		chromedp.DisableGPU,
	)

	// Create ChromeDP context
	login.ctx, login.cancel = chromedp.NewExecAllocator(context.Background(), login.opts...)
	login.ctx, login.cancel = chromedp.NewContext(login.ctx)

	// Run browser automation
	return login.runAutomation()
}

// Run automation
func (login *QRLogin) runAutomation() error {
	err := chromedp.Run(login.ctx,
		// Set the viewport size
		chromedp.EmulateViewport(1920, 1080),

		// Navigate to discord login page
		chromedp.Navigate(`https://discord.com/login`),

		// Wait for login page to load
		chromedp.WaitVisible("[class^=authBox]"),

		// Wait for the QR code to appear
		chromedp.WaitNotPresent("[class^=wanderingCubes]"),
		chromedp.WaitVisible("[class^=qrCode]"),

		// Wait for the QR code to be ready
		chromedp.Sleep(2*time.Second),

		// Wait for the QR code to appear
		chromedp.WaitNotPresent("[class^=wanderingCubes]"),
		chromedp.WaitVisible("[class^=qrCode]"),

		// Take a screenshot of the QR code
		chromedp.Screenshot("[class^=qrCode]", &login.qrImage, chromedp.NodeVisible),
		chromedp.ActionFunc(func(ctx context.Context) error {
			return login.ImageReadyCallback(login.qrImage)
		}),

		// Wait for Discord to load
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
