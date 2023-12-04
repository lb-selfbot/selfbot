package eventreactions

import (
	"fmt"
	"time"

	"github.com/LightningDev1/LB-Selfbot/selfbot/api"
	"github.com/LightningDev1/LB-Selfbot/selfbot/http"
	"github.com/LightningDev1/LB-Selfbot/selfbot/utils"
	"github.com/LightningDev1/discordgo"
)

type Ping struct {
	ReactionType EventReactionType
	Guild        *discordgo.Guild
	Channel      *discordgo.Channel
	System       string
	Response     string
}

func NewPingReaction(reactionType EventReactionType, guild *discordgo.Guild, channel *discordgo.Channel) *Ping {
	return &Ping{ReactionType: reactionType, Guild: guild, Channel: channel}
}

func NewPingFromData(session *discordgo.Session, reactionType EventReactionType, data map[string]any) (*Ping, error) {
	guildID, ok := data["guild"].(string)
	if ok {
		guild, err := utils.Discord.GetGuild(session, guildID)
		if err != nil {
			guild = &discordgo.Guild{ID: guildID}
		}

		return NewPingReaction(reactionType, guild, nil), nil
	}

	channelID, ok := data["channel"].(string)
	if ok {
		channel, err := utils.Discord.GetChannel(session, channelID)
		if err != nil {
			channel = &discordgo.Channel{ID: channelID}
		}

		return NewPingReaction(reactionType, nil, channel), nil
	}

	utils.Logging.Error("Ping.NewPingFromData: Invalid ping data")

	return nil, fmt.Errorf("invalid ping data")
}

func (p *Ping) Type() EventReactionType {
	return p.ReactionType
}

func (p *Ping) Description() string {
	var location string
	if p.Guild != nil {
		location = p.Guild.Name
	} else {
		location = p.Channel.Name
	}
	location = " (" + location + ")"

	switch p.ReactionType {
	case PING_PING:
		return "Ping --> Ping" + location
	case PING_KICK:
		return "Ping --> Kick" + location
	case PING_BAN:
		return "Ping --> Ban" + location
	case PING_MUTE:
		return "Ping --> Mute" + location
	case PING_TIMEOUT:
		return "Ping --> Timeout" + location
	case PING_INSULT:
		return "Ping --> Insult" + location
	case PING_RESPONSE:
		return "Ping --> Response" + location
	case PING_AI_RESPONSE:
		return "Ping --> AI Response" + location
	default:
		utils.Logging.Debug("Ping.Description: Unknown reaction type")
		return "Ping reaction" + location
	}
}

func (p *Ping) Equals(other EventReaction) bool {
	otherPing, ok := other.(*Ping)
	if !ok {
		utils.Logging.Debug("Ping.Equals: Other is not a Ping")
		return false
	}

	return p.ReactionType == otherPing.ReactionType &&
		p.Guild == otherPing.Guild &&
		p.Channel == otherPing.Channel
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

	pinged := false
	for _, mention := range message.Mentions {
		if mention.ID == session.State.User.ID {
			pinged = true
		}
	}

	if !pinged {
		return
	}

	switch p.ReactionType {
	case PING_PING:
		m, err := session.ChannelMessageSend(message.ChannelID, message.Author.Mention())
		HandleMessage(m, err)
	case PING_KICK:
		_ = session.GuildMemberDeleteWithReason(message.GuildID, message.Author.ID, "Ping Kick Event Reaction")
	case PING_BAN:
		_ = session.GuildBanCreateWithReason(message.GuildID, message.Author.ID, "Ping Ban Event Reaction", 0)
	case PING_MUTE:
		member, ok := utils.Discord.GetMember(session, message.GuildID, message.Author.ID)
		if !ok {
			return
		}

		_ = utils.Discord.MuteUser(session, member)
	case PING_TIMEOUT:
		member, ok := utils.Discord.GetMember(session, message.GuildID, message.Author.ID)
		if !ok {
			return
		}

		until := time.Now().Add(5 * time.Minute)
		_ = session.GuildMemberTimeout(member.GuildID, member.User.ID, &until)
	case PING_INSULT:
		httpResponse := http.Get("https://evilinsult.com/generate_insult.php?lang=en")
		m, err := session.ChannelMessageSend(message.ChannelID, message.Author.Mention()+" "+httpResponse.Body)
		HandleMessage(m, err)
	case PING_RESPONSE:
		m, err := session.ChannelMessageSend(message.ChannelID, p.Response)
		HandleMessage(m, err)
	case PING_AI_RESPONSE:
		system := utils.Misc.CreateAISystemPing(p.System, session, message.Message)
		response, err := api.GPT3(system, message.Content)
		if err != nil {
			utils.Logging.Error("Ping.Handler: Failed to get AI response:", err)
			return
		}

		m, err := session.ChannelMessageSend(message.ChannelID, response)
		HandleMessage(m, err)
	}
}
