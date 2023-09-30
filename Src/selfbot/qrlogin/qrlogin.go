package qrlogin

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/png"
	"regexp"
	"strings"

	chromedp "github.com/LightningDev1/chromiumdp"
)

var getStylesScript = `
[...document.styleSheets[0].cssRules].map(r=>r.selectorText).join('\n')
`

var tokenScript = `
(webpackChunkdiscord_app.push([[''],{},e=>{m=[];for(let c in e.c)m.push(e.c[c])}]),m).find(m=>m?.exports?.default?.getToken!==void 0).exports.default.getToken()
`

// Login is used to login to discord using a QR code and retrieve the token
type Login struct {
	// VisibleBrowser determines if the browser should be visible or not
	VisibleBrowser     bool
	ImageReadyCallback func(image.Image, []byte) error
	TokenCallback      func(string) error

	ctx    context.Context
	cancel context.CancelFunc
	opts   []chromedp.ExecAllocatorOption

	qrImage []byte
	styles  string
	token   string
}

// Start starts the login process
func (l *Login) Start() error {
	// Options for chrome.
	l.opts = append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/116.0.0.0 Safari/537.36"),
		chromedp.Flag("headless", !l.VisibleBrowser),
		chromedp.Flag("disable-gpu", !l.VisibleBrowser),
		chromedp.Flag("no-sandbox", true),
		chromedp.DisableGPU,
	)

	// Create allocator with options
	l.ctx, l.cancel = chromedp.NewExecAllocator(context.Background(), l.opts...)

	// Create chrome context
	l.ctx, l.cancel = chromedp.NewContext(l.ctx)

	// Run browser automation
	return l.runAutomation()

	// return nil
}

// GetQRImage returns the QR code image
func (l *Login) GetQRImage() (image.Image, error) {
	return png.Decode(bytes.NewReader(l.qrImage))
}

func (l *Login) getClassName(query string) string {
	classes := regexp.MustCompile(fmt.Sprintf(`%s-(.*?)[, \n]`, query)).FindStringSubmatch(l.styles)
	if len(classes) < 2 {
		return ""
	}

	return strings.TrimSpace(classes[0])
}

func (l *Login) runAutomation() error {
	// First get the login page styles to get the classnames of the elements we need
	err := chromedp.Run(l.ctx,
		// Navigate to discord login page
		chromedp.Navigate(`https://discord.com/login`),

		// Get the styles
		chromedp.Evaluate(getStylesScript, &l.styles),
	)
	if err != nil {
		return err
	}

	err = chromedp.Run(l.ctx,
		// Set the viewport size
		chromedp.EmulateViewport(1920, 1080),

		// Navigate to discord login page
		chromedp.Navigate(`https://discord.com/login`),

		// Wait for login page to load
		chromedp.WaitVisible(l.getClassName("authBox")),
		// chromedp.Sleep(500*time.Millisecond),
		chromedp.WaitVisible(l.getClassName("authBox")),
		// chromedp.Sleep(500*time.Millisecond),

		// Wait for the QR code to appear
		chromedp.WaitNotPresent(l.getClassName("wanderingCubes")),
		// chromedp.Sleep(500*time.Millisecond),
		chromedp.WaitNotPresent(l.getClassName("wanderingCubes")),
		// chromedp.Sleep(500*time.Millisecond),

		// Take a screenshot of the QR code
		chromedp.Screenshot(l.getClassName("qrCode"), &l.qrImage, chromedp.NodeVisible),
		chromedp.ActionFunc(func(ctx context.Context) error {
			image, err := png.Decode(bytes.NewReader(l.qrImage))
			if err != nil {
				return err
			}

			return l.ImageReadyCallback(image, l.qrImage)
		}),

		// Wait for Discord to load
		chromedp.WaitVisible(l.getClassName("content")),

		// Get the token
		chromedp.Evaluate(tokenScript, &l.token),
		chromedp.ActionFunc(func(ctx context.Context) error {
			return l.TokenCallback(l.token)
		}),
	)
	if err != nil {
		return err
	}

	return chromedp.Cancel(l.ctx)
}
