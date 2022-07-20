package commands

import (
	"time"

	"github.com/hugolgst/rich-go/client"

	"github.com/LightningDev1/LB-Selfbot/config"
	"github.com/LightningDev1/dgc"
)

type RPCCommands struct{}

func (r RPCCommands) RegisterCommands(router *dgc.Router) {
	/*router.RegisterCmd(&dgc.Command{
		Name:        "tweet",
		Description: "Tweet a message",
		Usage:       "[p]tweet <username> <text>",
		Handler:     i.tweet,
	})*/
	
	go func() {
		startTime := time.Now()
		for {
			time.Sleep(time.Second * 5)
			config, err := config.GetConfig()
			if err != nil {
				continue
			}
			if config.RPC.Enabled {
				client.Login(config.RPC.ClientID)
				client.SetActivity(client.Activity{
					State:      config.RPC.State,
					Details:    config.RPC.Details,
					LargeImage: config.RPC.LargeImage,
					LargeText:  config.RPC.LargeText,
					SmallImage: config.RPC.SmallImage,
					SmallText:  config.RPC.SmallText,
					Buttons:    r.getButtons(config.RPC),
					Timestamps: r.getTimestamps(config.RPC.Time, startTime),
					Party:      r.getParty(config.RPC),
					Secrets:    r.getSecrets(config.RPC),
				})
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
			Label: button.Label,
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
