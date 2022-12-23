package commands

import (
	"strings"
	"time"

	"github.com/hugolgst/rich-go/client"

	"github.com/LightningDev1/LB-Selfbot/config"
	"github.com/LightningDev1/LB-Selfbot/constants"
	"github.com/LightningDev1/LB-Selfbot/utils"
	"github.com/LightningDev1/dgc"
)

type RPCCommands struct{}

func (r RPCCommands) RegisterCommands(router *dgc.Router) {
	router.StartCategory("RPC", "RPC commands")
	router.RegisterCmd(&dgc.Command{
		Name:        "rpc",
		Description: "Turn Rich Presence on or off",
		Usage:       "[p]rpc <on/off>",
		Handler:     r.rpc,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "rpcid",
		Description: "Set the RPC client id",
		Usage:       "[p]rpcid <client id>",
		Handler:     r.rpcid,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "rpcstate",
		Description: "Set the RPC state",
		Usage:       "[p]rpcstate <text>",
		Handler:     r.rpcstate,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "rpcdetails",
		Description: "Set the RPC details",
		Usage:       "[p]rpcdetails <text>",
		Handler:     r.rpcdetails,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "rpctime",
		Description: "Set RPC time on or off",
		Usage:       "[p]rpctime <on/off>",
		Handler:     r.rpctime,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "rpcinvitable",
		Description: "Set if RPC is invitable",
		Usage:       "[p]rpcinvitable <on/off>",
		Handler:     r.rpcinvitable,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "rpclargeimage",
		Description: "Set the RPC large image",
		Usage:       "[p]rpclargeimage <image name>",
		Handler:     r.rpclargeimage,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "rpclargetext",
		Description: "Set the RPC large image text",
		Usage:       "[p]rpclargetext <text>",
		Handler:     r.rpclargetext,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "rpcsmallimage",
		Description: "Set the RPC small image",
		Usage:       "[p]rpcsmallimage <image name>",
		Handler:     r.rpcsmallimage,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "rpcsmalltext",
		Description: "Set the RPC small image text",
		Usage:       "[p]rpcsmalltext <text>",
		Handler:     r.rpcsmalltext,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "rpcbutton1text",
		Description: "Set the RPC first button text",
		Usage:       "[p]rpcbutton1text <text>",
		Handler:     r.rpcbutton1text,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "rpcbutton1url",
		Description: "Set the RPC first button url",
		Usage:       "[p]rpcbutton1url <url>",
		Handler:     r.rpcbutton1url,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "rpcbutton2text",
		Description: "Set the RPC second button text",
		Usage:       "[p]rpcbutton2text <text>",
		Handler:     r.rpcbutton2text,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "rpcbutton2url",
		Description: "Set the RPC second button url",
		Usage:       "[p]rpcbutton2url <url>",
		Handler:     r.rpcbutton2url,
	})

	go func() {
		startTime := time.Now()
		for {
			time.Sleep(time.Second * 5)
			cfg, err := config.GetConfig()
			if err != nil {
				utils.Logging.Error(err)
				continue
			}
			if cfg.RPC.Enabled {
				err = client.Login(cfg.RPC.ClientID)
				if err != nil {
					utils.Logging.Warn("RPC login failed: ", err)
					continue
				}
				err = client.SetActivity(client.Activity{
					State:      strings.ReplaceAll(cfg.RPC.State, "$VERSION", constants.VERSION),
					Details:    strings.ReplaceAll(cfg.RPC.Details, "$VERSION", constants.VERSION),
					LargeImage: cfg.RPC.LargeImage,
					LargeText:  strings.ReplaceAll(cfg.RPC.LargeText, "$VERSION", constants.VERSION),
					SmallImage: cfg.RPC.SmallImage,
					SmallText:  strings.ReplaceAll(cfg.RPC.SmallText, "$VERSION", constants.VERSION),
					Buttons:    r.getButtons(cfg.RPC),
					Timestamps: r.getTimestamps(cfg.RPC.Time, startTime),
					Party:      r.getParty(cfg.RPC),
					Secrets:    r.getSecrets(cfg.RPC),
				})
				if err != nil {
					utils.Logging.Warn("RPC set activity failed: ", err)
					continue
				}
			} else {
				client.Logout()
			}
		}
	}()
}

func (RPCCommands) getButtons(rpcData config.RPCData) []*client.Button {
	if rpcData.Invitable {
		return []*client.Button{}
	}
	var buttonsToReturn []*client.Button
	for _, button := range rpcData.Buttons {
		buttonsToReturn = append(buttonsToReturn, &client.Button{
			Label: strings.ReplaceAll(button.Label, "$VERSION", constants.VERSION),
			Url:   button.URL,
		})
	}
	return buttonsToReturn
}

func (RPCCommands) getTimestamps(timeEnabled bool, startTime time.Time) *client.Timestamps {
	if timeEnabled {
		return &client.Timestamps{
			Start: &startTime,
		}
	}
	return nil
}

func (RPCCommands) getParty(rpcData config.RPCData) *client.Party {
	if rpcData.Invitable {
		return &client.Party{
			ID:         "party_id",
			Players:    1,
			MaxPlayers: 69,
		}
	}
	return nil
}

func (RPCCommands) getSecrets(rpcData config.RPCData) *client.Secrets {
	if rpcData.Invitable {
		return &client.Secrets{
			Join:     "join_secret",
			Spectate: "spectate_secret",
			Match:    "match_secret",
		}
	}
	return nil
}

func (RPCCommands) rpc(ctx *dgc.Ctx) {
	toggle := strings.ToLower(utils.Arguments["<on/off>"])
	if toggle != "on" && toggle != "off" {
		_ = ctx.RespondText("Invalid value. Use `on` or `off`.")
		return
	}
	cfg, _ := config.GetConfig()
	cfg.RPC.Enabled = toggle == "on"
	_ = cfg.Save()
}

func (RPCCommands) rpcid(*dgc.Ctx) {
	cfg, _ := config.GetConfig()
	cfg.RPC.ClientID = utils.Arguments["<client id>"]
	_ = cfg.Save()
}

func (RPCCommands) rpcstate(*dgc.Ctx) {
	cfg, _ := config.GetConfig()
	cfg.RPC.State = utils.Arguments["<text>"]
	_ = cfg.Save()
}

func (RPCCommands) rpcdetails(*dgc.Ctx) {
	cfg, _ := config.GetConfig()
	cfg.RPC.Details = utils.Arguments["<text>"]
	_ = cfg.Save()
}

func (RPCCommands) rpctime(ctx *dgc.Ctx) {
	toggle := strings.ToLower(utils.Arguments["<on/off>"])
	if toggle != "on" && toggle != "off" {
		_ = ctx.RespondText("Invalid value. Use `on` or `off`.")
		return
	}
	cfg, _ := config.GetConfig()
	cfg.RPC.Time = toggle == "on"
	_ = cfg.Save()
}

func (RPCCommands) rpcinvitable(ctx *dgc.Ctx) {
	toggle := strings.ToLower(utils.Arguments["<on/off>"])
	if toggle != "on" && toggle != "off" {
		_ = ctx.RespondText("Invalid value. Use `on` or `off`.")
		return
	}
	cfg, _ := config.GetConfig()
	cfg.RPC.Invitable = toggle == "on"
	_ = cfg.Save()
}

func (RPCCommands) rpclargeimage(*dgc.Ctx) {
	cfg, _ := config.GetConfig()
	cfg.RPC.LargeImage = utils.Arguments["<image name>"]
	_ = cfg.Save()
}

func (RPCCommands) rpclargetext(*dgc.Ctx) {
	cfg, _ := config.GetConfig()
	cfg.RPC.LargeText = utils.Arguments["<text>"]
	_ = cfg.Save()
}

func (RPCCommands) rpcsmallimage(*dgc.Ctx) {
	cfg, _ := config.GetConfig()
	cfg.RPC.LargeImage = utils.Arguments["<image name>"]
	_ = cfg.Save()
}

func (RPCCommands) rpcsmalltext(*dgc.Ctx) {
	cfg, _ := config.GetConfig()
	cfg.RPC.LargeText = utils.Arguments["<text>"]
	_ = cfg.Save()
}

func (RPCCommands) rpcbutton1text(*dgc.Ctx) {
	cfg, _ := config.GetConfig()
	if len(cfg.RPC.Buttons) < 1 {
		cfg.RPC.Buttons[0] = config.RPCButton{}
	}
	cfg.RPC.Buttons[0].Label = utils.Arguments["<text>"]
	_ = cfg.Save()
}

func (RPCCommands) rpcbutton1url(*dgc.Ctx) {
	cfg, _ := config.GetConfig()
	if len(cfg.RPC.Buttons) < 1 {
		cfg.RPC.Buttons[0] = config.RPCButton{}
	}
	cfg.RPC.Buttons[0].URL = utils.Arguments["<url>"]
	_ = cfg.Save()
}

func (RPCCommands) rpcbutton2text(*dgc.Ctx) {
	cfg, _ := config.GetConfig()
	if len(cfg.RPC.Buttons) < 1 {
		cfg.RPC.Buttons[0] = config.RPCButton{}
	}
	if len(cfg.RPC.Buttons) < 2 {
		cfg.RPC.Buttons[1] = config.RPCButton{}
	}
	cfg.RPC.Buttons[1].Label = utils.Arguments["<text>"]
	_ = cfg.Save()
}

func (RPCCommands) rpcbutton2url(*dgc.Ctx) {
	cfg, _ := config.GetConfig()
	if len(cfg.RPC.Buttons) < 1 {
		cfg.RPC.Buttons[0] = config.RPCButton{}
	}
	if len(cfg.RPC.Buttons) < 2 {
		cfg.RPC.Buttons[1] = config.RPCButton{}
	}
	cfg.RPC.Buttons[1].URL = utils.Arguments["<url>"]
	_ = cfg.Save()
}
