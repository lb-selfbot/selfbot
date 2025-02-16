package browser

import (
	"context"

	"github.com/chromedp/chromedp"
)

// Browser is used to interact with a browser
type Browser struct {
	ctx     context.Context
	cancel  context.CancelFunc
	options []chromedp.ExecAllocatorOption
}

// BrowserOptions are the options for creating a new browser
type BrowserOptions struct {
	// UserAgent is the user agent to use
	UserAgent string

	// Headless is whether to run the browser in headless mode
	Headless bool

	// Size is the window size
	Size [2]int
}

// NewBrowser creates a new browser instance
func NewBrowser(ctx context.Context, opts BrowserOptions) *Browser {
	options := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("user-agent", opts.UserAgent),
		chromedp.Flag("headless", opts.Headless),
		chromedp.WindowSize(opts.Size[0], opts.Size[1]),
		chromedp.DisableGPU,
	)

	// Create browser context
	ctx, _ = chromedp.NewExecAllocator(ctx, options...)
	ctx, cancel := chromedp.NewContext(ctx)

	return &Browser{
		ctx:     ctx,
		cancel:  cancel,
		options: options,
	}
}

func (browser *Browser) RunAutomation(actions ...chromedp.Action) error {
	defer chromedp.Cancel(browser.ctx)

	return chromedp.Run(browser.ctx, actions...)
}
