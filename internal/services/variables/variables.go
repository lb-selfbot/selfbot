package variables

import (
	"fmt"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/LightningDev1/go-foreground"
	"github.com/lb-selfbot/discordgo"
	"github.com/lb-selfbot/selfbot/internal/api"
	"github.com/lb-selfbot/selfbot/internal/constants"
	"github.com/lb-selfbot/selfbot/internal/fn"
	"github.com/lb-selfbot/selfbot/internal/formatting"
	"github.com/lb-selfbot/selfbot/internal/services/eventbroker"
	"github.com/lb-selfbot/selfbot/internal/services/extensions"
	"github.com/lb-selfbot/selfbot/internal/services/router"
	"github.com/lb-selfbot/selfbot/internal/services/spotify"
	"github.com/lb-selfbot/selfbot/internal/utils"
	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/disk"
	"github.com/shirou/gopsutil/v4/mem"
)

// Service handles variables in user-defined strings.“
type Service interface {
	// Parse parses the given string and replaces variables with their values.
	Parse(string) string
}

type variablesService struct {
	eventBroker   eventbroker.Service
	routerSvc     router.Service
	extensionsSvc extensions.Service
	spotifySvc    spotify.Service

	session *discordgo.Session
}

// NewService creates a new instance of VariablesService.
func NewService(eventBroker eventbroker.Service, routerSvc router.Service, extensionsSvc extensions.Service, spotifySvc spotify.Service) Service {
	svc := &variablesService{eventBroker: eventBroker, routerSvc: routerSvc, extensionsSvc: extensionsSvc, spotifySvc: spotifySvc}
	svc.subscribe()
	return svc
}

func (v *variablesService) subscribe() {
	v.eventBroker.Subscribe(eventbroker.EventSessionReady, func(payload any) {
		session, ok := payload.(*discordgo.Session)
		if !ok {
			return
		}

		v.session = session
	})
}

func (v *variablesService) Parse(str string) string {
	for _, variable := range variables {
		matches := variable.re.FindAllStringSubmatch(str, -1)
		if len(matches) == 0 {
			continue
		}

		match := matches[0]

		if len(match) == 0 {
			continue
		}

		str = strings.ReplaceAll(str, match[0], v.getValue(variable.name, match[1:]))
	}

	return str
}

func (v *variablesService) getValue(varName string, args []string) string {
	switch varName {
	case "ver":
		return constants.VERSION

	case "date":
		return time.Now().Local().Format("2006-01-02")

	case "time":
		return time.Now().Local().Format("15:04:05")

	case "tz":
		if len(args) == 0 {
			return ""
		}
		loc, err := time.LoadLocation(args[0])
		if err != nil {
			return ""
		}
		return time.Now().In(loc).Format("15:04:05")

	case "uptime":
		return formatting.DurationShort(time.Since(constants.START_TIME))

	case "activewindow":
		title, err := foreground.GetForegroundTitle()
		if err != nil {
			return ""
		}
		return title

	case "usedcmds":
		return strconv.Itoa(len(v.routerSvc.GetCommands()))

	case "totalcmds":
		return strconv.Itoa(v.routerSvc.GetNumUsedCommands())

	case "exts":
		return strconv.Itoa(v.extensionsSvc.GetNumExtensions())

	case "servers":
		if v.session == nil {
			return ""
		}
		return strconv.Itoa(len(v.session.State.Guilds))

	case "friends":
		if v.session == nil {
			return ""
		}
		friends := fn.Filter(v.session.State.Relationships, func(r *discordgo.Relationship) bool {
			return r.Type == 1
		})
		return strconv.Itoa(len(friends))

	case "av":
		if v.session == nil {
			return ""
		}
		return strings.ReplaceAll(v.session.State.User.AvatarURL(""), "https://cdn.discordapp.com/", "mp:")

	case "os":
		switch runtime.GOOS {
		case "windows":
			return "Windows"
		case "linux":
			return "Linux"
		case "darwin":
			return "macOS"
		default:
			return runtime.GOOS
		}

	case "arch":
		switch runtime.GOARCH {
		case "amd64":
			return "x64"
		case "386":
			return "x86"
		case "arm64":
			return "ARM64"
		case "arm":
			return "ARM"
		default:
			return runtime.GOARCH
		}

	case "cpu%":
		percentages, err := cpu.Percent(0, false)
		if err != nil || len(percentages) == 0 {
			return "Unknown"
		}
		return strconv.FormatFloat(percentages[0], 'f', 2, 64) + "%"

	case "cpu":
		info, err := cpu.Info()
		if err != nil || len(info) == 0 {
			return "Unknown"
		}
		return fmt.Sprintf("%s %.f", info[0].ModelName, info[0].Mhz)

	case "ram":
		mem, err := mem.VirtualMemory()
		if err != nil {
			return "Unknown"
		}
		return fmt.Sprintf("%.2f of %.2f GB", float64(mem.Used)/1e9, float64(mem.Total)/1e9)

	case "ram%":
		mem, err := mem.VirtualMemory()
		if err != nil {
			return "Unknown"
		}
		return strconv.FormatFloat(mem.UsedPercent, 'f', 2, 64) + "%"

	case "disk":
		path := "/"
		if runtime.GOOS == "windows" {
			path = "C:"
		}
		usage, err := disk.Usage(path)
		if err != nil {
			return "Unknown"
		}
		return fmt.Sprintf("%.2f of %.2f GB", float64(usage.Used)/1e9, float64(usage.Total)/1e9)

	case "disk%":
		path := "/"
		if runtime.GOOS == "windows" {
			path = "C:"
		}
		usage, err := disk.Usage(path)
		if err != nil {
			return "Unknown"
		}
		return strconv.FormatFloat(usage.UsedPercent, 'f', 2, 64) + "%"

	case "sname":
		if v.session == nil {
			return ""
		}
		song, err := v.spotifySvc.GetCurrentlyPlaying(v.session)
		if err != nil {
			return ""
		}
		return song.Get("name").MustString()

	case "sartist":
		if v.session == nil {
			return ""
		}
		song, err := v.spotifySvc.GetCurrentlyPlaying(v.session)
		if err != nil {
			return ""
		}
		return strings.Join(v.spotifySvc.GetSongArtists(song), ", ")

	case "salbum":
		if v.session == nil {
			return ""
		}
		song, err := v.spotifySvc.GetCurrentlyPlaying(v.session)
		if err != nil {
			return ""
		}
		return song.Get("album", "name").MustString()

	case "sduration":
		if v.session == nil {
			return ""
		}
		song, err := v.spotifySvc.GetCurrentlyPlaying(v.session)
		if err != nil {
			return ""
		}
		duration := song.Get("duration_ms").MustInt(0) / 1000
		return fmt.Sprintf("%02d:%02d", duration/60, duration%60)

	case "scover":
		if v.session == nil {
			return ""
		}
		song, err := v.spotifySvc.GetCurrentlyPlaying(v.session)
		if err != nil {
			return ""
		}
		imgs := song.Get("album", "images")
		if imgs.Len() == 0 {
			return ""
		}
		url := imgs.Get(0, "url").MustString()
		urlParts := strings.Split(url, "/")
		return "'spotify:" + urlParts[len(urlParts)-1]

	case "slink":
		if v.session == nil {
			return ""
		}
		song, err := v.spotifySvc.GetCurrentlyPlaying(v.session)
		if err != nil {
			return ""
		}
		return song.Get("external_urls", "spotify").MustString()

	case "weather":
		if len(args) == 0 {
			return ""
		}
		weather, err := api.GetWeather(args[0])
		if err != nil {
			return ""
		}
		return utils.String.Capitalize(weather.Get("weather", 0, "description").MustString())

	case "tempc":
		if len(args) == 0 {
			return ""
		}
		weather, err := api.GetWeather(args[0])
		if err != nil {
			return ""
		}
		tempC := weather.Get("main", "temp").MustFloat64() - 273.15
		return fmt.Sprintf("%.1f °C", tempC)

	case "tempf":
		if len(args) == 0 {
			return ""
		}
		weather, err := api.GetWeather(args[0])
		if err != nil {
			return ""
		}
		tempC := weather.Get("main", "temp").MustFloat64() - 273.15
		tempF := tempC*9/5 + 32
		return fmt.Sprintf("%.1f °F", tempF)

	default:
		return ""
	}
}
