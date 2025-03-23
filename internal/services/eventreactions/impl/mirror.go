package impl

import (
	"fmt"

	"github.com/lb-selfbot/discordgo"
	"github.com/lb-selfbot/selfbot/internal/services/eventbroker"
	"github.com/lb-selfbot/selfbot/internal/services/eventreactions/types"
	"github.com/lb-selfbot/selfbot/internal/utils"
)

type Mirror struct {
	eventBroker eventbroker.Service

	SourceGuild   *discordgo.Guild
	DestGuild     *discordgo.Guild
	SourceChannel *discordgo.Channel
	DestChannel   *discordgo.Channel
}

func NewMirrorGuild(eventBroker eventbroker.Service, source, dest *discordgo.Guild) *Mirror {
	return &Mirror{eventBroker: eventBroker, SourceGuild: source, DestGuild: dest}
}

func NewMirrorChannel(eventBroker eventbroker.Service, source, dest *discordgo.Channel) *Mirror {
	return &Mirror{eventBroker: eventBroker, SourceChannel: source, DestChannel: dest}
}

func NewMirrorFromData(eventBroker eventbroker.Service, session *discordgo.Session, data map[string]any) (*Mirror, error) {
	sourceGuildID, ok := data["source_guild"].(string)
	if ok {
		sourceGuild, err := utils.Discord.GetGuild(session, sourceGuildID)
		if err != nil {
			sourceGuild = &discordgo.Guild{ID: sourceGuildID}
		}

		destGuildID, ok := data["dest_guild"].(string)
		if ok {
			destGuild, err := utils.Discord.GetGuild(session, destGuildID)
			if err != nil {
				destGuild = &discordgo.Guild{ID: destGuildID}
			}

			return NewMirrorGuild(eventBroker, sourceGuild, destGuild), nil
		} else {
			utils.Logging.Error("Mirror.NewMirrorFromData: Invalid mirror data")
			return nil, fmt.Errorf("invalid mirror data")
		}
	}

	sourceChannelID, ok := data["source_channel"].(string)
	if ok {
		sourceChannel, err := utils.Discord.GetChannel(session, sourceChannelID)
		if err != nil {
			sourceChannel = &discordgo.Channel{ID: sourceChannelID}
		}

		destChannelID, ok := data["dest_channel"].(string)
		if ok {
			destChannel, err := utils.Discord.GetChannel(session, destChannelID)
			if err != nil {
				destChannel = &discordgo.Channel{ID: destChannelID}
			}

			return NewMirrorChannel(eventBroker, sourceChannel, destChannel), nil
		} else {
			utils.Logging.Error("Mirror.NewMirrorFromData: Invalid mirror data")
			return nil, fmt.Errorf("invalid mirror data")
		}
	}

	utils.Logging.Error("Mirror.NewMirrorFromData: Invalid mirror data")
	return nil, fmt.Errorf("invalid mirror data")
}

func (m *Mirror) Type() types.EventReactionType {
	if m.SourceGuild != nil {
		return types.MirrorGuild
	} else if m.SourceChannel != nil {
		return types.MirrorChannel
	} else {
		utils.Logging.Debug("Mirror.Type: Invalid ignore type")
		return -1
	}
}

func (m *Mirror) Description() string {
	if m.SourceGuild != nil && m.DestGuild != nil {
		return fmt.Sprintf("Mirror Server - %s --> %s", m.SourceGuild.Name, m.DestGuild.Name)
	} else if m.SourceChannel != nil && m.DestChannel != nil {
		return fmt.Sprintf("Mirror Channel - %s --> %s", m.SourceChannel.Name, m.DestChannel.Name)
	} else {
		utils.Logging.Debug("Mirror.Description: Invalid ignore type")
		return "Mirror"
	}
}

func (m *Mirror) Equals(other types.EventReaction) bool {
	otherMirror, ok := other.(*Mirror)
	if !ok {
		utils.Logging.Debug("Mirror.Equals: Other is not a Mirror")
		return false
	}

	return m.SourceGuild == otherMirror.SourceGuild &&
		m.DestGuild == otherMirror.DestGuild &&
		m.SourceChannel == otherMirror.SourceChannel &&
		m.DestChannel == otherMirror.DestChannel
}

func (m *Mirror) Data() map[string]any {
	data := make(map[string]any)

	if m.SourceGuild != nil {
		data["source_guild"] = m.SourceGuild.ID
	}
	if m.DestGuild != nil {
		data["dest_guild"] = m.DestGuild.ID
	}
	if m.SourceChannel != nil {
		data["source_channel"] = m.SourceChannel.ID
	}
	if m.DestChannel != nil {
		data["dest_channel"] = m.DestChannel.ID
	}

	return data
}

func (m *Mirror) Handler(session *discordgo.Session, event any) {
	message, ok := event.(*discordgo.MessageCreate)
	if !ok {
		utils.Logging.Debug("Mirror.Handler: Event is not a message create event")
		return
	}

	if m.SourceGuild != nil && message.GuildID != m.SourceGuild.ID {
		return
	} else if m.SourceChannel != nil && message.ChannelID != m.SourceChannel.ID {
		return
	}

	switch {
	case m.DestGuild != nil:
		ogChannel, err := session.State.Channel(message.ChannelID)
		if err != nil {
			utils.Logging.Debug("Mirror.Handler: Error getting channel: " + err.Error())
			return
		}

		var channel *discordgo.Channel
		for _, c := range m.DestGuild.Channels {
			if c.Name == ogChannel.Name {
				channel = c
				break
			}
		}

		if channel == nil {
			utils.Logging.Error("Mirror.Handler: Could not find channel")
			return
		}

		msg, err := utils.Discord.ImpersonateMessage(session, message.Author, message.Message, channel.ID)
		if err != nil {
			utils.Logging.Debug("Mirror.Handler: Error impersonating: " + err.Error())
			return
		}

		m.eventBroker.Publish(eventbroker.EventHandleReactionMessage, msg)

	case m.DestChannel != nil:
		msg, err := utils.Discord.ImpersonateMessage(session, message.Author, message.Message, m.DestChannel.ID)
		if err != nil {
			utils.Logging.Debug("Mirror.Handler: Error impersonating: " + err.Error())
			return
		}

		m.eventBroker.Publish(eventbroker.EventHandleReactionMessage, msg)

	default:
		utils.Logging.Debug("Mirror.Handler: Invalid ignore type")
		return
	}
}
