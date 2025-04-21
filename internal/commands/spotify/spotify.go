package spotify

import (
	"fmt"
	"strings"

	"github.com/lb-selfbot/discordgo"
	"github.com/lb-selfbot/selfbot/internal/fn"
	"github.com/lb-selfbot/selfbot/internal/formatting"
	"github.com/lb-selfbot/selfbot/internal/services"
	"github.com/lb-selfbot/selfbot/internal/services/router"
)

type SpotifyCommands struct{ services *services.Services }

func NewSpotifyCommands(s *services.Services) *SpotifyCommands {
	return &SpotifyCommands{services: s}
}

func (c *SpotifyCommands) Register() {
	c.services.Router.RegisterCategory(&router.Category{
		Name:        "Spotify",
		Description: "Control your Spotify music",
	}, []*router.Command{
		{
			Name:        "spause",
			Description: "Pause Spotify Music",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous},
			Handler:     c.spause,
		},
		{
			Name:        "sresume",
			Description: "Resume Spotify Music",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous},
			Handler:     c.sresume,
		},
		{
			Name:        "sprevious",
			Description: "Play the previous song",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous},
			Handler:     c.sprevious,
		},
		{
			Name:        "snext",
			Description: "Play the next song",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous},
			Handler:     c.snext,
		},
		{
			Name:        "splay",
			Description: "Play a song on Spotify",
			Arguments:   "<song name>",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous},
			Handler:     c.splay,
		},
		{
			Name:        "squeue",
			Description: "Queue a song to play after the current one",
			Arguments:   "<song name>",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous},
			Handler:     c.squeue,
		},
		{
			Name:        "splaying",
			Description: "Current playing song info",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous},
			Handler:     c.splaying,
		},
		{
			Name:        "svolume",
			Description: "Change Spotify volume",
			Arguments:   "<volume>",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous},
			Handler:     c.svolume,
		},
		{
			Name:        "sshuffle",
			Description: "Turn shuffle on or off",
			Arguments:   "<on/off>",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous},
			Handler:     c.sshuffle,
		},
		{
			Name:        "srepeat",
			Description: "Change the Spotify repeat mode",
			Arguments:   "<track/context/off>",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous},
			Handler:     c.srepeat,
		},
		{
			Name:        "sseek",
			Description: "Seek to a specific time in the song",
			Arguments:   "<timestamp>",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous},
			Handler:     c.sseek,
		},
		{
			Name:        "sdevices",
			Description: "Get your available Spotify devices",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous},
			Handler:     c.sdevices,
		},
		{
			Name:        "sprofile",
			Description: "Get your Spotify profile info",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous},
			Handler:     c.sprofile,
		},
		{
			Name:        "sbearertoken",
			Description: "Get your Spotify bearer token used for the API",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous},
			Handler:     c.sbearertoken,
		},
		{
			Name:        "slistenalong",
			Description: "Send a Spotify listen along invite",
			Arguments:   "[content]",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous},
			Handler:     c.slistenalong,
		},
	})
}

func (c *SpotifyCommands) spause(ctx *router.Ctx) error {
	if err := c.services.Spotify.Pause(ctx.Session); err != nil {
		return fmt.Errorf("A Spotify error has occurred: %v", err)
	}

	c.services.Embed.Respond(ctx, "Successfully paused Spotify music.")
	return nil
}

func (c *SpotifyCommands) sresume(ctx *router.Ctx) error {
	if err := c.services.Spotify.Resume(ctx.Session); err != nil {
		return fmt.Errorf("A Spotify error has occurred: %v", err)
	}

	c.services.Embed.Respond(ctx, "Successfully resumed Spotify music.")
	return nil
}

func (c *SpotifyCommands) sprevious(ctx *router.Ctx) error {
	if err := c.services.Spotify.Previous(ctx.Session); err != nil {
		return fmt.Errorf("A Spotify error has occurred: %v", err)
	}

	c.services.Embed.Respond(ctx, "Successfully started playing the previous Spotify track.")
	return nil
}

func (c *SpotifyCommands) snext(ctx *router.Ctx) error {
	if err := c.services.Spotify.Next(ctx.Session); err != nil {
		return fmt.Errorf("A Spotify error has occurred: %v", err)
	}

	c.services.Embed.Respond(ctx, "Successfully started playing the next Spotify track.")
	return nil
}

func (c *SpotifyCommands) splay(ctx *router.Ctx) error {
	track, err := c.services.Spotify.Play(ctx.Session, ctx.Arguments.Raw())
	if err != nil {
		return fmt.Errorf("A Spotify error has occurred: %v", err)
	}

	artists := c.services.Spotify.GetSongArtists(track)

	c.services.Embed.Respond(ctx, "Successfully started playing %s - %s", strings.Join(artists, ", "), track.Get("name").MustString())
	return nil
}

func (c *SpotifyCommands) squeue(ctx *router.Ctx) error {
	track, err := c.services.Spotify.Queue(ctx.Session, ctx.Arguments.Raw())
	if err != nil {
		return fmt.Errorf("A Spotify error has occurred: %v", err)
	}

	artists := c.services.Spotify.GetSongArtists(track)

	c.services.Embed.Respond(ctx, "Successfully queued %s - %s", strings.Join(artists, ", "), track.Get("name").MustString())
	return nil
}

func (c *SpotifyCommands) splaying(ctx *router.Ctx) error {
	song, err := c.services.Spotify.GetCurrentlyPlaying(ctx.Session)
	if err != nil {
		return fmt.Errorf("A Spotify error has occurred: %v", err)
	}

	songArtists := c.services.Spotify.GetSongArtists(song)
	songArtistsString := strings.Join(songArtists, ", ")
	songReleaseDate := song.Get("album", "release_date").MustString("Unknown")
	songExplicit := formatting.BoolYesNo(song.Get("explicit").MustBool(false))
	duration := int(song.Get("duration_ms").MustInt(0)) / 1000
	songDuration := fmt.Sprintf("%d:%02d", duration/60, duration&60)
	songUrl := song.Get("external_urls", "spotify").MustString("Unknown")
	songImageUrl := song.Get("album", "images").Get(0).Get("url").MustString("")

	c.services.Embed.New(ctx).
		AddField("Title", song.Get("name").MustString()).
		AddField("Artist", songArtistsString).
		AddField("Release Date", songReleaseDate).
		AddField("Explicit", songExplicit).
		AddField("Duration", songDuration).
		AddField("URL", songUrl).
		SetImage(songImageUrl).
		Send()
	return nil
}

func (c *SpotifyCommands) svolume(ctx *router.Ctx) error {
	volume, err := ctx.Arguments.Get(0).Int()
	volume = fn.Clamp(volume, 0, 100)
	if err != nil {
		return fmt.Errorf("Argument <volume> must be a number.")
	}

	if err := c.services.Spotify.SetVolume(ctx.Session, volume); err != nil {
		return fmt.Errorf("A Spotify error has occurred: %v", err)
	}

	c.services.Embed.Respond(ctx, "Successfully set Spotify volume to %d%%", volume)
	return nil
}

func (c *SpotifyCommands) sshuffle(ctx *router.Ctx) error {
	toggle, err := ctx.Arguments.Get(0).Bool()
	if err != nil {
		return fmt.Errorf("Invalid value. Use `on` or `off`.")
	}

	if err := c.services.Spotify.SetShuffle(ctx.Session, toggle); err != nil {
		return fmt.Errorf("A Spotify error has occurred: %v", err)
	}

	c.services.Embed.Respond(ctx, "Successfully turned shuffle %s.", formatting.BoolOnOff(toggle))
	return nil
}

func (c *SpotifyCommands) srepeat(ctx *router.Ctx) error {
	state := ctx.Arguments.Get(0).Raw()

	if err := c.services.Spotify.SetRepeat(ctx.Session, state); err != nil {
		return fmt.Errorf("A Spotify error has occurred: %v", err)
	}

	c.services.Embed.Respond(ctx, "Successfully set Spotify repeat mode to `%s`.", state)
	return nil
}

func (c *SpotifyCommands) sseek(ctx *router.Ctx) error {
	timestamp, err := ctx.Arguments.Get(0).Duration()
	if err != nil {
		return fmt.Errorf("Invalid value. Use a timestamp like `1m30s`.")
	}

	if err := c.services.Spotify.Seek(ctx.Session, timestamp); err != nil {
		return fmt.Errorf("A Spotify error has occurred: %v", err)
	}

	c.services.Embed.Respond(ctx, "Successfully seeked to %s.", timestamp)
	return nil
}

func (c *SpotifyCommands) sdevices(ctx *router.Ctx) error {
	devices, err := c.services.Spotify.GetDevices(ctx.Session)
	if err != nil {
		return fmt.Errorf("A Spotify error has occurred: %v", err)
	}

	var deviceList []string

	for device := range devices.Iterator() {
		deviceList = append(deviceList, fmt.Sprintf(
			"**%s**\nType: %s\nActive: %s\nVolume: %d%%\n",
			device.Get("name").MustString(),
			device.Get("type").MustString(),
			formatting.BoolYesNo(device.Get("is_active").MustBool(false)),
			device.Get("volume_percent").MustInt(0),
		))
	}

	c.services.Embed.RespondField(ctx, "Devices", strings.Join(deviceList, "\n"))
	return nil
}

func (c *SpotifyCommands) sprofile(ctx *router.Ctx) error {
	profile, err := c.services.Spotify.GetProfile(ctx.Session)
	if err != nil {
		return fmt.Errorf("A Spotify error has occurred: %v", err)
	}

	c.services.Embed.New(ctx).
		AddField("Username", profile.Get("display_name").MustString()).
		AddField("ID", profile.Get("id").MustString()).
		AddField("Country", profile.Get("country").MustString()).
		AddField("Product", profile.Get("product").MustString()).
		AddField("Followers", fmt.Sprint(profile.Get("followers", "total").MustInt())).
		AddField("URL", profile.Get("href").MustString()).
		Send()
	return nil
}

func (c *SpotifyCommands) sbearertoken(ctx *router.Ctx) error {
	token, err := c.services.Spotify.GetBearerToken(ctx.Session)
	if err != nil {
		return fmt.Errorf("A Spotify error has occurred: %v", err)
	}

	profile, err := c.services.Spotify.GetProfile(ctx.Session)
	if err != nil {
		return fmt.Errorf("A Spotify error has occurred: %v", err)
	}

	confirmed, err := c.services.Discord.ShowConfirmDialog(ctx.Message.ChannelID, "Please confirm you want to show your bearer token in this channel by removing the reaction.\n\nWarning: Your bearer token can be used to control your Spotify account.")
	if err != nil {
		return fmt.Errorf("Error while showing confirm dialog: %v", err)
	}

	if confirmed {
		c.services.Embed.New(ctx).
			AddField("Spotify Account", profile.Get("display_name").MustString()).
			AddField("Bearer Token", token).
			Send()
	}

	return nil
}

func (c *SpotifyCommands) slistenalong(ctx *router.Ctx) error {
	spotifySessions := fn.Filter(c.services.Sessions.GetGatewaySessions(), func(s *discordgo.GatewaySession) bool {
		return len(s.Activities) != 0 && s.Activities[0].Name == "Spotify"
	})

	if len(spotifySessions) == 0 {
		return fmt.Errorf("No Spotify sessions found.")
	}

	_, err := ctx.Session.ChannelMessageSendComplex(ctx.Message.ChannelID, &discordgo.MessageSend{
		Content: ctx.Arguments.Raw(),
		Activity: &discordgo.MessageActivity{
			Type:      discordgo.MessageActivityTypeListen,
			SessionID: spotifySessions[0].SessionID,
			PartyID:   spotifySessions[0].Activities[0].Party.ID,
		},
		Nonce: discordgo.GenerateNonce(),
		Flags: 0,
	})
	if err != nil {
		return fmt.Errorf("An error has occurred while sending the invite: %v", err)
	}
	return nil
}
