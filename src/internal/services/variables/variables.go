package variables

import (
	"fmt"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/LightningDev1/LB-Selfbot/internal/constants"
	"github.com/LightningDev1/discordgo"

	"github.com/LightningDev1/LB-Selfbot/internal/api"
	"github.com/LightningDev1/LB-Selfbot/internal/services/interfaces"
	"github.com/LightningDev1/LB-Selfbot/internal/utils"
	"github.com/LightningDev1/go-foreground"
	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/disk"
	"github.com/shirou/gopsutil/v4/mem"
)

type variable struct {
	regex     *regexp.Regexp
	valueFunc func(arg string) string
}

type VariablesServiceImpl struct {
	discordSvc        interfaces.DiscordService
	eventReactionsSvc interfaces.EventReactionsService
	extensionsSvc     interfaces.ExtensionsService
	spotifySvc        interfaces.SpotifyService

	variables []variable
}

func (v *VariablesServiceImpl) Inject(
	discordSvc interfaces.DiscordService,
	eventReactionsSvc interfaces.EventReactionsService,
	extensionsSvc interfaces.ExtensionsService,
	spotifySvc interfaces.SpotifyService,
) {
	v.discordSvc = discordSvc
	v.eventReactionsSvc = eventReactionsSvc
	v.extensionsSvc = extensionsSvc
	v.spotifySvc = spotifySvc
}

func (v *VariablesServiceImpl) Init() {
	v.variables = []variable{
		{
			regex:     regexp.MustCompile(`\$VERSION`),
			valueFunc: func(arg string) string { return constants.VERSION },
		},

		{
			regex:     regexp.MustCompile(`\$DATE`),
			valueFunc: func(arg string) string { return time.Now().Local().Format("2006-01-02") },
		},
		{
			regex: regexp.MustCompile(`\$TIME\((.*?)\)`),
			valueFunc: func(arg string) string {
				loc, err := time.LoadLocation(arg)
				if err != nil {
					return time.Now().Local().Format("15:04:05") + " (invalid timezone)"
				}

				return time.Now().In(loc).Format("15:04:05")
			},
		},

		{
			regex:     regexp.MustCompile(`\$UPTIME`),
			valueFunc: func(arg string) string { return utils.Misc.FormatDurationShort(time.Since(constants.START_TIME)) },
		},

		{
			regex: regexp.MustCompile(`\$ACTIVEWINDOW`),
			valueFunc: func(arg string) string {
				title, _ := foreground.GetForegroundTitle()
				return title
			},
		},

		{
			regex: regexp.MustCompile(`\$USEDCOMMANDS`),
			valueFunc: func(arg string) string {
				return strconv.Itoa(len(v.eventReactionsSvc.GetUsedCommands()))
			},
		},
		{
			regex: regexp.MustCompile(`\$TOTALCOMMANDS`),
			valueFunc: func(arg string) string {
				return strconv.Itoa(len(v.discordSvc.GetRouter().Commands))
			},
		},
		{
			regex: regexp.MustCompile(`\$EXTENSIONS`),
			valueFunc: func(arg string) string {
				return strconv.Itoa(v.extensionsSvc.GetLoadedExtensions())
			},
		},
		{
			regex: regexp.MustCompile(`\$SERVERS`),
			valueFunc: func(arg string) string {
				return strconv.Itoa(len(v.discordSvc.GetSession().State.Guilds))
			},
		},
		{
			regex: regexp.MustCompile(`\$FRIENDS`),
			valueFunc: func(arg string) string {
				friends := utils.Filter(v.discordSvc.GetSession().State.Relationships, func(f *discordgo.Relationship) bool {
					return f.Type == 1
				})

				return strconv.Itoa(len(friends))
			},
		},
		{
			regex: regexp.MustCompile(`\$AVATAR`),
			valueFunc: func(arg string) string {
				return strings.ReplaceAll(v.discordSvc.GetSession().State.User.AvatarURL(""), "https://cdn.discordapp.com/", "mp:")
			},
		},
		{
			regex: regexp.MustCompile(`\$OS`),
			valueFunc: func(arg string) string {
				switch runtime.GOOS {
				case "windows":
					return "Windows"
				case "darwin":
					return "macOS"
				case "linux":
					return "Linux"
				default:
					return "Unknown"
				}
			},
		},
		{
			regex: regexp.MustCompile(`\$ARCH`),
			valueFunc: func(arg string) string {
				return runtime.GOARCH
			},
		},
		{
			regex: regexp.MustCompile(`\$CPUPERCENT`),
			valueFunc: func(arg string) string {
				cpuUsage, err := cpu.Percent(0, false)
				if err != nil || len(cpuUsage) == 0 {
					return "Unknown"
				}
				return fmt.Sprintf("%.2f%%", cpuUsage[0])
			},
		},
		{
			regex: regexp.MustCompile(`\$CPU`),
			valueFunc: func(arg string) string {
				cpuInfo, err := cpu.Info()
				if err != nil || len(cpuInfo) == 0 {
					return "Unknown"
				}
				return fmt.Sprintf("%s %.f", cpuInfo[0].ModelName, cpuInfo[0].Mhz)
			},
		},
		{
			regex: regexp.MustCompile(`\$RAMUSAGE`),
			valueFunc: func(arg string) string {
				memory, err := mem.VirtualMemory()
				if err != nil {
					return "Unknown"
				}
				return fmt.Sprintf("%.2f/%.2f GB", float64(memory.Used)/1e9, float64(memory.Total)/1e9)
			},
		},
		{
			regex: regexp.MustCompile(`\$RAMPERCENT`),
			valueFunc: func(arg string) string {
				memory, err := mem.VirtualMemory()
				if err != nil {
					return "Unknown"
				}
				return fmt.Sprintf("%.2f%%", float64(memory.Used)/float64(memory.Total)*100)
			},
		},
		{
			regex: regexp.MustCompile(`\$DISKUSAGE`),
			valueFunc: func(arg string) string {
				path := "/"
				if runtime.GOOS == "windows" {
					path = "C:"
				}
				diskUsage, err := disk.Usage(path)
				if err != nil {
					return "Unknown"
				}
				return fmt.Sprintf("%.2f/%.2f GB", float64(diskUsage.Used)/1e9, float64(diskUsage.Total)/1e9)
			},
		},
		{
			regex: regexp.MustCompile(`\$DISKPERCENT`),
			valueFunc: func(arg string) string {
				path := "/"
				if runtime.GOOS == "windows" {
					path = "C:"
				}
				diskUsage, err := disk.Usage(path)
				if err != nil {
					return "Unknown"
				}
				return fmt.Sprintf("%.2f%%", float64(diskUsage.Used)/float64(diskUsage.Total)*100)
			},
		},
		{
			regex: regexp.MustCompile(`\$SONGNAME`),
			valueFunc: func(arg string) string {
				song, err := v.spotifySvc.GetCurrentlyPlaying()
				if err != nil {
					return ""
				}
				return song.Get("name").MustString()
			},
		},
		{
			regex: regexp.MustCompile(`\$SONGARTIST`),
			valueFunc: func(arg string) string {
				song, err := v.spotifySvc.GetCurrentlyPlaying()
				if err != nil {
					return ""
				}
				return strings.Join(utils.Spotify.GetSongArtists(song), ", ")
			},
		},
		{
			regex: regexp.MustCompile(`\$SONGALBUM`),
			valueFunc: func(arg string) string {
				song, err := v.spotifySvc.GetCurrentlyPlaying()
				if err != nil {
					return ""
				}
				return song.GetPath("album", "name").MustString()
			},
		},
		{
			regex: regexp.MustCompile(`\$SONGDURATION`),
			valueFunc: func(arg string) string {
				song, err := v.spotifySvc.GetCurrentlyPlaying()
				if err != nil {
					return ""
				}
				duration := int(song.Get("duration_ms").MustInt(0) / 1000)
				return fmt.Sprintf("%d:%02d", duration/60, duration%60)
			},
		},
		{
			regex: regexp.MustCompile(`\$SONGCOVER`),
			valueFunc: func(arg string) string {
				song, err := v.spotifySvc.GetCurrentlyPlaying()
				if err != nil {
					return ""
				}
				images := song.GetPath("album", "images")
				if len(images.MustArray()) == 0 {
					return ""
				}
				url := images.GetIndex(0).Get("url").MustString()
				urlSplit := strings.Split(url, "/")
				return "spotify:" + urlSplit[len(urlSplit)-1]
			},
		},
		{
			regex: regexp.MustCompile(`\$SONGLINK`),
			valueFunc: func(arg string) string {
				song, err := v.spotifySvc.GetCurrentlyPlaying()
				if err != nil {
					return ""
				}
				return song.Get("external_urls").Get("spotify").MustString()
			},
		},
		{
			regex: regexp.MustCompile(`\$WEATHER\((.*?)\)`),
			valueFunc: func(arg string) string {
				weather, err := api.GetWeather(arg)
				if err != nil {
					return ""
				}
				return utils.String.Capitalize(weather.Get("weather").GetIndex(0).Get("description").MustString())
			},
		},
		{
			regex: regexp.MustCompile(`\$TEMPC\((.*?)\)`),
			valueFunc: func(arg string) string {
				weather, err := api.GetWeather(arg)
				if err != nil {
					return ""
				}
				tempC := weather.GetPath("main", "temp").MustFloat64() - 273.15
				return fmt.Sprintf("%.1f °C", tempC)
			},
		},
		{
			regex: regexp.MustCompile(`\$TEMPF\((.*?)\)`),
			valueFunc: func(arg string) string {
				weather, err := api.GetWeather(arg)
				if err != nil {
					return ""
				}
				tempC := weather.GetPath("main", "temp").MustFloat64() - 273.15
				tempF := tempC*9/5 + 32
				return fmt.Sprintf("%.1f °F", tempF)
			},
		},
	}
}

func (v *VariablesServiceImpl) Parse(text string) string {
	for _, variable := range v.variables {
		text = v.parseVariable(text, variable)
	}

	return text
}

func (v *VariablesServiceImpl) parseVariable(text string, variable variable) string {
	matches := variable.regex.FindAllStringSubmatch(text, -1)
	if len(matches) == 0 {
		return text
	}

	match := matches[0]

	if len(match) == 1 {
		text = strings.ReplaceAll(text, match[0], variable.valueFunc(""))
	} else if len(match) == 2 {
		text = strings.ReplaceAll(text, match[0], variable.valueFunc(match[1]))
	}

	return text
}
