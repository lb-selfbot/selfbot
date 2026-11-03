module github.com/lb-selfbot/selfbot

go 1.27

require (
	github.com/LightningDev1/beeep v0.1.0
	github.com/LightningDev1/go-foreground v0.0.0-20250402184805-913886c46a00
	github.com/LightningDev1/go-privnote v0.1.1
	github.com/agnivade/levenshtein v1.2.2
	github.com/bogdanfinn/fhttp v0.7.3
	github.com/bogdanfinn/tls-client v1.18.0
	github.com/charmbracelet/bubbles v1.0.0
	github.com/charmbracelet/bubbletea v1.3.10
	github.com/charmbracelet/lipgloss v1.1.0
	github.com/chromedp/chromedp v0.21.0
	github.com/common-nighthawk/go-figure v0.0.0-20210622060536-734e95fb86be
	github.com/darkwyrm/b85 v0.0.0-20210908110113-bdd55eedeb3c
	github.com/deadshot465/owoify-go v1.0.1
	github.com/getsentry/sentry-go v0.50.0
	github.com/goccy/go-json v0.12.0
	github.com/gorilla/websocket v1.5.3
	github.com/hugolgst/rich-go v0.0.0-20240715122152-74618cc1ace2
	github.com/kbinani/screenshot v0.0.0-20250118074034-a3924b7bbc8c
	github.com/lb-selfbot/discordgo v0.26.6-0.20250421075913-0875e9737a1e
	github.com/lb-selfbot/go-jester v0.3.0
	github.com/mr-tron/base58 v1.3.0
	github.com/multiformats/go-base36 v0.2.0
	github.com/robert-nix/ansihtml v1.0.1
	github.com/shirou/gopsutil/v4 v4.26.10
	golang.org/x/crypto v0.57.0
	golang.org/x/exp v0.0.0-20250218142911-aa4b98e5adaa
	golang.org/x/sys v0.48.0
	golang.org/x/term v0.46.0
	golang.org/x/text v0.42.0
	google.golang.org/protobuf v1.36.12
)

require (
	github.com/andybalholm/brotli v1.2.7 // indirect
	github.com/atotto/clipboard v0.1.4 // indirect
	github.com/aymanbagabas/go-osc52/v2 v2.0.1 // indirect
	github.com/bdandy/go-errors v1.2.2 // indirect
	github.com/bdandy/go-socks4 v1.2.3 // indirect
	github.com/bogdanfinn/quic-go-utls v1.0.13-utls // indirect
	github.com/bogdanfinn/utls v1.8.2-barnius // indirect
	github.com/bogdanfinn/websocket v1.5.9-barnius // indirect
	github.com/charmbracelet/colorprofile v0.4.1 // indirect
	github.com/charmbracelet/x/ansi v0.11.6 // indirect
	github.com/charmbracelet/x/cellbuf v0.0.15 // indirect
	github.com/charmbracelet/x/term v0.2.2 // indirect
	github.com/chromedp/cdproto v0.157.8 // indirect
	github.com/clipperhouse/displaywidth v0.9.0 // indirect
	github.com/clipperhouse/stringish v0.1.1 // indirect
	github.com/clipperhouse/uax29/v2 v2.5.0 // indirect
	github.com/ebitengine/purego v0.11.2 // indirect
	github.com/erikgeiser/coninput v0.0.0-20211004153227-1c3628e74d0f // indirect
	github.com/gen2brain/shm v0.1.1 // indirect
	github.com/go-json-experiment/json v0.0.0-20260213210345-44df1a37e875 // indirect
	github.com/go-ole/go-ole v1.3.0 // indirect
	github.com/go-toast/toast v0.0.0-20190211030409-01e6764cf0a4 // indirect
	github.com/godbus/dbus/v5 v5.1.0 // indirect
	github.com/jezek/xgb v1.1.1 // indirect
	github.com/klauspost/compress v1.20.1 // indirect
	github.com/lucasb-eyer/go-colorful v1.3.0 // indirect
	github.com/lufia/plan9stats v0.0.0-20240909124753-873cd0166683 // indirect
	github.com/lxn/win v0.0.0-20210218163916-a377121e959e // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mattn/go-localereader v0.0.1 // indirect
	github.com/mattn/go-runewidth v0.0.19 // indirect
	github.com/molecule-man/go-brrr v1.2.0 // indirect
	github.com/muesli/ansi v0.0.0-20230316100256-276c6243b2f6 // indirect
	github.com/muesli/cancelreader v0.2.2 // indirect
	github.com/muesli/termenv v0.16.0 // indirect
	github.com/nu7hatch/gouuid v0.0.0-20131221200532-179d4d0c4d8d // indirect
	github.com/power-devops/perfstat v0.0.0-20260805114148-88456608a4f6 // indirect
	github.com/quic-go/qpack v0.6.0 // indirect
	github.com/rivo/uniseg v0.4.7 // indirect
	github.com/sahilm/fuzzy v0.1.1 // indirect
	github.com/tadvi/systray v0.0.0-20190226123456-11a2b8fa57af // indirect
	github.com/tam7t/hpkp v0.0.0-20160821193359-2b70b4024ed5 // indirect
	github.com/tklauser/go-sysconf v0.4.0 // indirect
	github.com/tklauser/numcpus v0.12.0 // indirect
	github.com/xo/terminfo v0.0.0-20220910002029-abceb7e1c41e // indirect
	github.com/yusufpapurcu/wmi v1.2.4 // indirect
	golang.org/x/net v0.59.0 // indirect
	gopkg.in/natefinch/npipe.v2 v2.0.0-20160621034901-c1b8fa8bdcce // indirect
)

// Add npipe support for Windows ARM64
replace gopkg.in/natefinch/npipe.v2 => github.com/ph4r05/npipe v0.0.0-20240527105247-e977e49ae283
