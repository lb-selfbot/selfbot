package variables

import (
	"fmt"
	"regexp"
	"runtime"
	"strings"
	"time"

	"github.com/LightningDev1/LB-Selfbot/selfbot/api"
	"github.com/LightningDev1/LB-Selfbot/selfbot/constants"
	"github.com/LightningDev1/LB-Selfbot/selfbot/data"
	"github.com/LightningDev1/LB-Selfbot/selfbot/spotify"
	"github.com/LightningDev1/LB-Selfbot/selfbot/utils"
	"github.com/LightningDev1/go-foreground"
	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/disk"
	"github.com/shirou/gopsutil/v4/mem"
)

var regexes = map[string]*regexp.Regexp{}

func parseVariable(text, variable string, value func(arg string) string) string {
	if _, ok := regexes[variable]; !ok {
		regexes[variable] = regexp.MustCompile(variable)
	}

	matches := regexes[variable].FindAllStringSubmatch(text, -1)
	if len(matches) == 0 {
		return text
	}

	match := matches[0]

	if len(match) == 1 {
		text = strings.ReplaceAll(text, match[0], value(""))
	} else if len(match) == 2 {
		text = strings.ReplaceAll(text, match[0], value(match[1]))
	}

	return text
}

func Parse(text string) string {
	// Constant variables
	text = strings.ReplaceAll(text, "$VERSION", constants.VERSION)

	// Dynamic variables

	text = parseVariable(text, "\\$DATE", func(arg string) string {
		return time.Now().Local().Format("2006-01-02")
	})
	text = parseVariable(text, "\\$TIME\\((.*?)\\)", func(arg string) string {
		loc, err := time.LoadLocation(arg)
		if err != nil {
			return time.Now().Local().Format("15:04:05") + " (invalid timezone)"
		}

		return time.Now().In(loc).Format("15:04:05")
	})
	text = parseVariable(text, "\\$TIME", func(arg string) string {
		return time.Now().Local().Format("15:04:05")
	})
	text = parseVariable(text, "\\$UPTIME", func(arg string) string {
		return utils.Misc.FormatDurationShort(time.Since(constants.START_TIME))
	})
	text = parseVariable(text, "\\$ACTIVEWINDOW", func(arg string) string {
		title, _ := foreground.GetForegroundTitle()
		return title
	})
	text = parseVariable(text, "\\$USEDCOMMANDS", func(arg string) string {
		return fmt.Sprint(data.UsedCommands)
	})
	text = parseVariable(text, "\\$TOTALCOMMANDS", func(arg string) string {
		return fmt.Sprint(data.TotalCommands)
	})
	text = parseVariable(text, "\\$EXTENSIONS", func(arg string) string {
		return fmt.Sprint(data.LoadedExtensions)
	})
	text = parseVariable(text, "\\$SERVERS", func(arg string) string {
		return fmt.Sprint(data.Guilds)
	})
	text = parseVariable(text, "\\$FRIENDS", func(arg string) string {
		return fmt.Sprint(data.Friends)
	})
	text = parseVariable(text, "\\$AVATAR", func(arg string) string {
		return fmt.Sprint(data.AvatarImageKey)
	})

	// System Info
	text = parseVariable(text, "\\$OS", func(arg string) string {
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
	})
	text = parseVariable(text, "\\$ARCH", func(arg string) string {
		return runtime.GOARCH
	})
	text = parseVariable(text, "\\$CPUPERCENT", func(arg string) string {
		cpuUsage, err := cpu.Percent(0, false)
		if err != nil || len(cpuUsage) == 0 {
			return "Unknown"
		}

		return fmt.Sprintf("%.2f%%", cpuUsage[0])
	})
	text = parseVariable(text, "\\$CPU", func(arg string) string {
		cpuInfo, err := cpu.Info()
		if err != nil || len(cpuInfo) == 0 {
			return "Unknown"
		}

		return fmt.Sprintf("%s %.f", cpuInfo[0].ModelName, cpuInfo[0].Mhz)
	})
	text = parseVariable(text, "\\$RAMUSAGE", func(arg string) string {
		memory, err := mem.VirtualMemory()
		if err != nil {
			return "Unknown"
		}

		return fmt.Sprintf("%.2f/%.2f GB", float64(memory.Used)/1e9, float64(memory.Total)/1e9)
	})
	text = parseVariable(text, "\\$RAMPERCENT", func(arg string) string {
		memory, err := mem.VirtualMemory()
		if err != nil {
			return "Unknown"
		}

		return fmt.Sprintf("%.2f%%", float64(memory.Used)/float64(memory.Total)*100)
	})
	text = parseVariable(text, "\\$DISKUSAGE", func(arg string) string {
		path := "/"
		if runtime.GOOS == "windows" {
			path = "C:"
		}

		diskUsage, err := disk.Usage(path)
		if err != nil {
			return "Unknown"
		}

		return fmt.Sprintf("%.2f/%.2f GB", float64(diskUsage.Used)/1e9, float64(diskUsage.Total)/1e9)
	})
	text = parseVariable(text, "\\$DISKPERCENT", func(arg string) string {
		path := "/"
		if runtime.GOOS == "windows" {
			path = "C:"
		}

		diskUsage, err := disk.Usage(path)
		if err != nil {
			return "Unknown"
		}

		return fmt.Sprintf("%.2f%%", float64(diskUsage.Used)/float64(diskUsage.Total)*100)
	})

	// Spotify
	text = parseVariable(text, "\\$SONGNAME", func(arg string) string {
		song, err := spotify.API.GetCurrentlyPlaying(data.Token)
		if err != nil {
			return ""
		}

		return song.Get("name").MustString()
	})
	text = parseVariable(text, "\\$SONGARTIST", func(arg string) string {
		song, err := spotify.API.GetCurrentlyPlaying(data.Token)
		if err != nil {
			return ""
		}

		return strings.Join(utils.Spotify.GetSongArtists(song), ", ")
	})
	text = parseVariable(text, "\\$SONGALBUM", func(arg string) string {
		song, err := spotify.API.GetCurrentlyPlaying(data.Token)
		if err != nil {
			return ""
		}

		return song.GetPath("album", "name").MustString()
	})
	text = parseVariable(text, "\\$SONGDURATION", func(arg string) string {
		song, err := spotify.API.GetCurrentlyPlaying(data.Token)
		if err != nil {
			return ""
		}

		duration := int(song.Get("duration_ms").MustInt(0) / 1000)

		return fmt.Sprintf("%d:%02d", duration/60, duration%60)
	})
	text = parseVariable(text, "\\$SONGCOVER", func(arg string) string {
		song, err := spotify.API.GetCurrentlyPlaying(data.Token)
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
	})
	text = parseVariable(text, "\\$SONGLINK", func(arg string) string {
		song, err := spotify.API.GetCurrentlyPlaying(data.Token)
		if err != nil {
			return ""
		}

		return song.Get("external_urls").Get("spotify").MustString()
	})

	// Weather
	text = parseVariable(text, "\\$WEATHER\\((.*?)\\)", func(arg string) string {
		weather, err := api.GetWeather(arg)
		if err != nil {
			return ""
		}

		return utils.String.Capitalize(weather.Get("weather").GetIndex(0).Get("description").MustString())
	})
	text = parseVariable(text, "\\$TEMPC\\((.*?)\\)", func(arg string) string {
		weather, err := api.GetWeather(arg)
		if err != nil {
			return ""
		}

		tempC := weather.GetPath("main", "temp").MustFloat64() - 273.15

		return fmt.Sprintf("%.1f °C", tempC)
	})
	text = parseVariable(text, "\\$TEMPF\\((.*?)\\)", func(arg string) string {
		weather, err := api.GetWeather(arg)
		if err != nil {
			return ""
		}

		tempC := weather.GetPath("main", "temp").MustFloat64() - 273.15
		tempF := tempC*9/5 + 32

		return fmt.Sprintf("%.1f °F", tempF)
	})

	return text
}
