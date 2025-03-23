package impl

import (
	"fmt"
	"strings"
	"time"

	"github.com/LightningDev1/LB-Selfbot/internal/api"
	"github.com/LightningDev1/LB-Selfbot/internal/http"
	"github.com/LightningDev1/LB-Selfbot/internal/services/eventbroker"
	"github.com/LightningDev1/LB-Selfbot/internal/services/eventreactions/types"

	"github.com/LightningDev1/LB-Selfbot/internal/utils"
	"github.com/LightningDev1/discordgo"
)

type Ping struct {
	eventBroker eventbroker.Service

	ReactionType types.EventReactionType
	Guild        *discordgo.Guild
	Channel      *discordgo.Channel
	System       string
	Response     string
	Reaction     string
}

func NewPingReaction(eventBroker eventbroker.Service, reactionType types.EventReactionType, guild *discordgo.Guild, channel *discordgo.Channel) *Ping {
	return &Ping{
		eventBroker:  eventBroker,
		ReactionType: reactionType,
		Guild:        guild,
		Channel:      channel,
	}
}

func NewPingFromData(eventBroker eventbroker.Service, session *discordgo.Session, reactionType types.EventReactionType, data map[string]any) (*Ping, error) {
	var guild *discordgo.Guild
	var channel *discordgo.Channel

	guildID, ok := data["guild"].(string)
	if ok {
		guild, _ = utils.Discord.GetGuild(session, guildID)
	}

	channelID, ok := data["channel"].(string)
	if ok {
		channel, _ = utils.Discord.GetChannel(session, channelID)
	}

	if guild == nil && channel == nil {
		utils.Logging.Error("Ping.NewPingFromData: Invalid ping data")

		return nil, fmt.Errorf("invalid ping data")
	}

	reaction := NewPingReaction(eventBroker, reactionType, guild, channel)

	reaction.System, _ = data["system"].(string)
	reaction.Response, _ = data["response"].(string)
	reaction.Reaction, _ = data["reaction"].(string)

	return reaction, nil
}

func (p *Ping) Type() types.EventReactionType {
	return p.ReactionType
}

func (p *Ping) Description() string {
	var locations []string

	if p.Channel != nil {
		locations = append(locations, p.Channel.Name)
	}
	if p.Guild != nil {
		locations = append(locations, p.Guild.Name)
	}
	if p.System != "" {
		locations = append(locations, p.System)
	}
	if p.Response != "" {
		locations = append(locations, p.Response)
	}
	if p.Reaction != "" {
		locations = append(locations, "<"+p.Reaction+">")
	}

	location := " (" + strings.Join(locations, " - ") + ")"

	switch p.ReactionType {
	case types.PingPing:
		return "Ping → Ping" + location
	case types.PingKick:
		return "Ping → Kick" + location
	case types.PingBan:
		return "Ping → Ban" + location
	case types.PingMute:
		return "Ping → Mute" + location
	case types.PingTimeout:
		return "Ping → Timeout" + location
	case types.PingInsult:
		return "Ping → Insult" + location
	case types.PingReact:
		return "Ping → React" + location
	case types.PingResponse:
		return "Ping → Response" + location
	case types.PingAIResponse:
		return "Ping → AI Response" + location
	default:
		utils.Logging.Debug("Ping.Description: Unknown reaction type")
		return "Ping reaction" + location
	}
}

func (p *Ping) Equals(other types.EventReaction) bool {
	otherPing, ok := other.(*Ping)
	if !ok {
		utils.Logging.Debug("Ping.Equals: Other is not a Ping")
		return false
	}

	return p.ReactionType == otherPing.ReactionType &&
		p.Guild == otherPing.Guild &&
		p.Channel == otherPing.Channel &&
		p.Reaction == otherPing.Reaction
}

func (p *Ping) Data() map[string]any {
	data := make(map[string]any)

	data["type"] = p.ReactionType
	if p.Guild != nil {
		data["guild"] = p.Guild.ID
	}
	if p.Channel != nil {
		data["channel"] = p.Channel.ID
	}
	data["system"] = p.System
	data["response"] = p.Response
	data["reaction"] = p.Reaction

	return data
}

func (p *Ping) Handler(session *discordgo.Session, event any) {
	message, ok := event.(*discordgo.MessageCreate)
	if !ok {
		utils.Logging.Debug("Ping.Handler: Event is not a message create event")
		return
	}

	if message.Author.ID == session.State.User.ID {
		return
	}

	if p.Guild != nil && p.Guild.ID != message.GuildID {
		return
	} else if p.Channel != nil && p.Channel.ID != message.ChannelID {
		return
	}

	if !utils.Discord.IsUserMentioned(message.Message, session.State.User.ID) {
		return
	}

	switch p.ReactionType {
	case types.PingPing:
		m, err := session.ChannelMessageSend(message.ChannelID, message.Author.Mention())
		if err != nil {
			break
		}

		p.eventBroker.Publish(eventbroker.EventHandleReactionMessage, m)
	case types.PingKick:
		_ = session.GuildMemberDeleteWithReason(message.GuildID, message.Author.ID, "Ping Kick Event Reaction")
	case types.PingBan:
		_ = session.GuildBanCreateWithReason(message.GuildID, message.Author.ID, "Ping Ban Event Reaction", 0)
	case types.PingMute:
		member, ok := utils.Discord.GetMember(session, message.GuildID, message.Author.ID)
		if !ok {
			return
		}

		_ = utils.Discord.MuteUser(session, member)
	case types.PingTimeout:
		member, ok := utils.Discord.GetMember(session, message.GuildID, message.Author.ID)
		if !ok {
			return
		}

		until := time.Now().Add(5 * time.Minute)
		_ = session.GuildMemberTimeout(member.GuildID, member.User.ID, &until)
	case types.PingInsult:
		httpResponse := http.Get("https://evilinsult.com/generate_insult.php?lang=en")
		insult := message.Author.Mention() + " " + string(httpResponse.Body)

		m, err := session.ChannelMessageSend(message.ChannelID, insult)
		if err != nil {
			break
		}

		p.eventBroker.Publish(eventbroker.EventHandleReactionMessage, m)
	case types.PingReact:
		_ = session.MessageReactionAdd(message.ChannelID, message.ID, p.Reaction)
	case types.PingResponse:
		m, err := session.ChannelMessageSend(message.ChannelID, p.Response)
		if err != nil {
			break
		}

		p.eventBroker.Publish(eventbroker.EventHandleReactionMessage, m)
	case types.PingAIResponse:
		system := utils.Misc.CreateAISystemPing(p.System, session, message.Message)
		response, err := api.GPT(system, message.Content)
		if err != nil {
			utils.Logging.Error("Ping.Handler: Failed to get AI response:", err)
			return
		}

		m, err := session.ChannelMessageSend(message.ChannelID, response)
		if err != nil {
			break
		}

		p.eventBroker.Publish(eventbroker.EventHandleReactionMessage, m)
	}
}
