package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/lb-selfbot/discordgo"
	"github.com/lb-selfbot/selfbot/internal/constants"
	"github.com/lb-selfbot/selfbot/internal/logging"
	"github.com/lb-selfbot/selfbot/internal/services/config"
	"github.com/lb-selfbot/selfbot/internal/services/eventbroker"
	"github.com/lb-selfbot/selfbot/internal/services/ui/featurelogs"
	"github.com/lb-selfbot/selfbot/internal/utils"
)

const consoleMsgColor = "#C1EFFF"

func (ui *uiService) logFeature(payload any) {
	switch data := payload.(type) {
	case *featurelogs.CmdUsed:
		ui.logCmdUsed(*data)
	case *featurelogs.Ready:
		ui.logReady(data.User, data.Guilds, data.NumCmds)
	case *featurelogs.NitroSniper:
		ui.logNitroSniper(data.Message, data.Elapsed, data.Status, data.Code, data.Account)
	case *featurelogs.NitroSniperFake:
		ui.logNitroSniperFake(data.Message, data.Code, data.Account)
	case *featurelogs.PrivnoteSniper:
		ui.logPrivnoteSniper(data.Message, data.Elapsed, data.Status, data.Path, data.Account)
	case *featurelogs.GiveawaySniperFound:
		ui.logGiveawaySniperFound(data.Message, data.Prize, data.Host, data.EndsIn, data.Account)
	case *featurelogs.GiveawaySniperWon:
		ui.logGiveawaySniperWon(data.Message, data.Prize, data.Account)
	case *featurelogs.PokeTwoSniper:
		ui.logPokeTwoSniper(data.Message, data.Elapsed, data.Confidence, data.Pokemon, data.Received, data.IV, data.Account)
	case *featurelogs.BanNotifier:
		ui.logBanNotifier(data.GuildID)
	case *featurelogs.DMNotifier:
		ui.logDMNotifier(data.Message)
	case *featurelogs.DMEditNotifier:
		ui.logDMEditNotifier(data.Message, data.Before, data.After)
	case *featurelogs.DMDeleteNotifier:
		ui.logDMDeleteNotifier(data.Message, data.Attachments)
	case *featurelogs.DMTypingNotifier:
		ui.logDMTypingNotifier(data.Username)
	case *featurelogs.MessageDeleteNotifier:
		ui.logMessageDeleteNotifier(data.Message, data.Attachments)
	case *featurelogs.DiscordStaffDetection:
		ui.logDiscordStaffDetection(data.Message)
	case *featurelogs.GhostPingDetection:
		ui.logGhostPingDetection(data.Message)
	case *featurelogs.SelfbotDetection:
		ui.logSelfbotDetection(data.Message, data.Selfbot)
	case *featurelogs.TicketNotifier:
		ui.logTicketNotifier(data.Channel)
	case *featurelogs.KeywordNotifier:
		ui.logKeywordNotifier(data.Message, data.Keyword)
	case *featurelogs.PhishingDetection:
		ui.logPhishingDetection(data.Message, data.Domain)
	case *featurelogs.MassDMBlocked:
		ui.logMassDMBlocked()
	case *featurelogs.FriendRequestSpamBlocked:
		ui.logFriendRequestSpamBlocked()
	case *featurelogs.AntiDelete:
		ui.logAntiDelete(data.Message)
	case *featurelogs.AFKMessage:
		ui.logAFKMessage(data.Message, data.Response)
	case *featurelogs.MultiTokenSniper:
		ui.logMultiTokenSniper(data.NumTokens, data.NumGuilds)
	case *featurelogs.SessionConnected:
		ui.logSessionConnected(data.ID, data.OS, data.Client)
	case *featurelogs.SessionDisconnected:
		ui.logSessionDisconnected(data.ID, data.OS, data.Client)
	case *featurelogs.SCJobs:
		ui.logSCJobs(data.NumJobs)
	case *featurelogs.EventReactions:
		ui.logEventReactions(data.NumReactions)
	case *featurelogs.DangerousCmdBlocked:
		ui.logDangerousCmdBlocked(data.User, data.Cmd)
	case *featurelogs.DangerousHostingCmdBlocked:
		ui.logDangerousHostingCmdBlocked(data.Cmd)
	case *featurelogs.Reminder:
		ui.logReminder(data.Reminder)
	default:
		logging.Debug("Unknown feature logging payload type:", data)
	}
}

func (ui *uiService) logCmdUsed(data featurelogs.CmdUsed) {
	_ = ui.ShowNotification("You used the %s command.", data.Name)
	_ = ui.NewConsoleMessage().
		SetTitle("Command Used").
		SetColor(consoleMsgColor).
		AddInfo("Command", data.Name).
		Show()
}

// Updated to take guilds slice instead of event
func (ui *uiService) logReady(user *discordgo.User, guilds []*discordgo.Guild, numCmds int) {
	cfg := ui.configSvc.GetConfig()

	_ = ui.ShowNotification("Successfully connected as %s", user.DisplayName)

	// Count friends (no relationships in the new struct, so we use session)
	session := ui.discordSvc.GetSession()
	friends := utils.Filter(session.State.Relationships, func(rel *discordgo.Relationship) bool {
		return rel.Type == 1
	})

	userText := user.DisplayName + " | " + user.Username

	bannerTextUpper := fmt.Sprintf("LightningBot %s | %s | Prefix: %s", constants.VERSION, userText, cfg.Selfbot.CommandPrefix)
	bannerTextLower := fmt.Sprintf("%d Commands | %d Servers | %d Friends", numCmds, len(guilds), len(friends))

	upper := ui.colored(bannerTextUpper, 0, 203, 255, false)
	lower := ui.colored(bannerTextLower, 0, 203, 255, false)

	if constants.CONSOLE {
		upper := ui.colored(utils.String.Center(bannerTextUpper), 0, 203, 255, false)
		lower := ui.colored(utils.String.Center(bannerTextLower), 0, 203, 255, false)

		fmt.Println(upper + "\n" + lower + "\n")
	}

	ui.eventBroker.Publish(eventbroker.EventSetConsole, "ALIGN_CENTER"+upper+"\nALIGN_CENTER"+lower+"\n")
}

func (ui *uiService) logNitroSniper(msg *discordgo.Message, elapsed time.Duration, status, code, account string) {
	session := ui.discordSvc.GetSession()
	guildName := utils.Discord.GetGuildName(session, msg.GuildID)
	channelName := utils.Discord.GetChannelName(session, msg.ChannelID)

	cfg := ui.configSvc.GetConfig()

	_ = ui.ShowNotification("Found a Nitro code in %s - #%s!\nStatus: %s\nSent by: %s",
		guildName, channelName, status, msg.Author.String())

	_ = ui.NewConsoleMessage().
		SetTitle("Nitro Sniper").
		SetColor(consoleMsgColor).
		AddInfo("Status", status).
		AddInfo("Author", msg.Author.String()).
		AddInfo("Server", guildName).
		AddInfo("Channel", channelName).
		AddInfo("Code", code).
		AddInfo("Snipe Time", elapsed.String()).
		AddInfo("Message Link", msg.GetLink()).
		AddInfo("On Account", account, account != "").
		SetWebhookURL(cfg.Webhooks.NitroSniper).
		SetImportant(status == "Redeemed Nitro").
		Show()
}

func (ui *uiService) logNitroSniperFake(msg *discordgo.Message, code, account string) {
	session := ui.discordSvc.GetSession()
	guildName := utils.Discord.GetGuildName(session, msg.GuildID)
	channelName := utils.Discord.GetChannelName(session, msg.ChannelID)

	cfg := ui.configSvc.GetConfig()

	_ = ui.ShowNotification("Found a fake Nitro code in %s - #%s from %s",
		guildName, channelName, msg.Author.String())

	_ = ui.NewConsoleMessage().
		SetTitle("Nitro Sniper").
		SetColor(consoleMsgColor).
		AddInfo("Status", "Fake Code").
		AddInfo("Author", msg.Author.String()).
		AddInfo("Server", guildName).
		AddInfo("Channel", channelName).
		AddInfo("Code", code).
		AddInfo("Message Link", msg.GetLink()).
		AddInfo("On Account", account, account != "").
		SetWebhookURL(cfg.Webhooks.NitroSniper).
		Show()
}

func (ui *uiService) logPrivnoteSniper(msg *discordgo.Message, elapsed time.Duration, status, path, account string) {
	session := ui.discordSvc.GetSession()
	guildName := utils.Discord.GetGuildName(session, msg.GuildID)
	channelName := utils.Discord.GetChannelName(session, msg.ChannelID)

	cfg := ui.configSvc.GetConfig()

	_ = ui.ShowNotification("Sniped a Privnote in %s - $%s!\nFrom: %s\nSaved to: %s",
		guildName, channelName, msg.Author.String(), path)

	_ = ui.NewConsoleMessage().
		SetTitle("Privnote Sniper").
		SetColor(consoleMsgColor).
		AddInfo("Status", status).
		AddInfo("Author", msg.Author.String()).
		AddInfo("Server", guildName).
		AddInfo("Channel", channelName).
		AddInfo("Saved At", path).
		AddInfo("Snipe Time", elapsed.String()).
		AddInfo("Message Link", msg.GetLink()).
		AddInfo("On Account", account, account != "").
		SetWebhookURL(cfg.Webhooks.Misc).
		SetImportant(status == "Saved Note").
		Show()
}

func (ui *uiService) logGiveawaySniperFound(msg *discordgo.Message, prize, hostID, endsIn, account string) {
	session := ui.discordSvc.GetSession()
	guildName := utils.Discord.GetGuildName(session, msg.GuildID)
	channelName := utils.Discord.GetChannelName(session, msg.ChannelID)
	var hostName string
	host, err := utils.Discord.GetUser(session, hostID)
	if err == nil {
		hostName = host.DisplayName
	}

	cfg := ui.configSvc.GetConfig()

	notification := fmt.Sprintf("Found a giveaway in %s!", guildName)
	if prize != "" {
		notification += fmt.Sprintf("\nPrize: %s", prize)
	}
	if endsIn != "" {
		notification += fmt.Sprintf("\nEnds in: %s", endsIn)
	}

	_ = ui.ShowNotification("%s", notification)
	_ = ui.NewConsoleMessage().
		SetTitle("Giveaway Sniper").
		SetColor(consoleMsgColor).
		AddInfo("", "Found Giveaway").
		AddInfo("Server", guildName).
		AddInfo("Channel", channelName).
		AddInfo("Bot", msg.Author.String()).
		AddInfo("Prize", prize, prize != "").
		AddInfo("Hosted By", hostName, hostName != "").
		AddInfo("Ends In", endsIn, endsIn != "").
		AddInfo("Message Link", msg.GetLink()).
		AddInfo("On Account", account, account != "").
		SetWebhookURL(cfg.Webhooks.GiveawaySniper).
		Show()
}

func (ui *uiService) logGiveawaySniperWon(msg *discordgo.Message, prize, account string) {
	session := ui.discordSvc.GetSession()
	guildName := utils.Discord.GetGuildName(session, msg.GuildID)
	channelName := utils.Discord.GetChannelName(session, msg.ChannelID)

	cfg := ui.configSvc.GetConfig()

	notification := fmt.Sprintf("You won a giveaway in %s! 🎉", guildName)
	if prize != "" {
		notification += fmt.Sprintf("\nPrize: %s", prize)
	}

	_ = ui.ShowNotification("%s", notification)
	_ = ui.NewConsoleMessage().
		SetTitle("Giveaway Sniper").
		SetColor(consoleMsgColor).
		AddInfo("", "Won Giveaway").
		AddInfo("Server", guildName).
		AddInfo("Channel", channelName).
		AddInfo("Bot", msg.Author.String()).
		AddInfo("Prize", prize, prize != "").
		AddInfo("Message Link", msg.GetLink()).
		AddInfo("On Account", account, account != "").
		SetWebhookURL(cfg.Webhooks.GiveawaySniper).
		SetImportant(true).
		Show()
}

func (ui *uiService) logPokeTwoSniper(msg *discordgo.Message, elapsed time.Duration, confidence float64, pokemon, received, iv, account string) {
	session := ui.discordSvc.GetSession()
	guildName := utils.Discord.GetGuildName(session, msg.GuildID)
	channelName := utils.Discord.GetChannelName(session, msg.ChannelID)

	cfg := ui.configSvc.GetConfig()

	notification := fmt.Sprintf("PokéTwo Sniper caught a %s in %s - $%s! ✨", pokemon, guildName, channelName)
	if iv != "" {
		notification += fmt.Sprintf("\nIV: %s", iv)
	}

	_ = ui.ShowNotification("%s", notification)
	_ = ui.NewConsoleMessage().
		SetTitle("PokéTwo Sniper").
		SetColor(consoleMsgColor).
		AddInfo("", "Caught a "+pokemon).
		AddInfo("Received", received, received != "").
		AddInfo("IV", iv, iv != "").
		AddInfo("AI Confidence", fmt.Sprint(confidence)).
		AddInfo("Detection Time", elapsed.String()).
		AddInfo("Server", guildName).
		AddInfo("Channel", channelName).
		AddInfo("Message Link", msg.GetLink()).
		AddInfo("On Account", account, account != "").
		SetWebhookURL(cfg.Webhooks.Misc).
		SetImportant(true).
		Show()
}

func (ui *uiService) logBanNotifier(guildID string) {
	session := ui.discordSvc.GetSession()
	guildName := utils.Discord.GetGuildName(session, guildID)

	cfg := ui.configSvc.GetConfig()

	_ = ui.ShowNotification("You have been banned from %s!", guildName)
	_ = ui.NewConsoleMessage().
		SetTitle("Ban Notifier").
		SetColor(consoleMsgColor).
		AddInfo("", "You have been banned from "+guildName).
		SetWebhookURL(cfg.Webhooks.Notifiers).
		Show()
}

func (ui *uiService) logDMNotifier(msg *discordgo.Message) {
	cfg := ui.configSvc.GetConfig()

	_ = ui.ShowNotification("New DM from %s", msg.Author.String())
	_ = ui.NewConsoleMessage().
		SetTitle("DM Notifier").
		SetColor(consoleMsgColor).
		AddInfo("", "New DM from "+msg.Author.String()).
		AddInfo("Content", msg.Content).
		AddInfo("Message Link", msg.GetLink()).
		SetWebhookURL(cfg.Webhooks.Notifiers).
		Show()
}

func (ui *uiService) logDMEditNotifier(msg *discordgo.Message, before, after string) {
	cfg := ui.configSvc.GetConfig()

	_ = ui.ShowNotification("%s edited their DM", msg.Author.String())
	_ = ui.NewConsoleMessage().
		SetTitle("DM Edit Notifier").
		SetColor(consoleMsgColor).
		AddInfo("", "Edited DM from "+msg.Author.String()).
		AddInfo("Before", before).
		AddInfo("After", after).
		SetWebhookURL(cfg.Webhooks.Notifiers).
		Show()
}

func (ui *uiService) logDMDeleteNotifier(msg *discordgo.Message, attachments []string) {
	cfg := ui.configSvc.GetConfig()

	_ = ui.ShowNotification("%s deleted their DM", msg.Author.String())
	_ = ui.NewConsoleMessage().
		SetTitle("DM Delete Notifier").
		SetColor(consoleMsgColor).
		AddInfo("", "Deleted DM from "+msg.Author.String()).
		AddInfo("Content", msg.Content).
		AddInfo("Attachments", strings.Join(attachments, ", "), len(attachments) > 0).
		SetWebhookURL(cfg.Webhooks.Notifiers).
		Show()
}

func (ui *uiService) logDMTypingNotifier(username string) {
	cfg := ui.configSvc.GetConfig()

	_ = ui.ShowNotification("%s is typing in your DM...", username)
	_ = ui.NewConsoleMessage().
		SetTitle("DM Typing Notifier").
		SetColor(consoleMsgColor).
		AddInfo("", username+" is typing").
		SetWebhookURL(cfg.Webhooks.Notifiers).
		Show()
}

func (ui *uiService) logMessageDeleteNotifier(msg *discordgo.Message, attachments []string) {
	session := ui.discordSvc.GetSession()
	guildName := utils.Discord.GetGuildName(session, msg.GuildID)
	channelName := utils.Discord.GetChannelName(session, msg.ChannelID)

	cfg := ui.configSvc.GetConfig()

	_ = ui.ShowNotification("%s deleted a message in %s - #%s", msg.Author.String(), guildName, channelName)
	_ = ui.NewConsoleMessage().
		SetTitle("Message Delete Notifier").
		SetColor(consoleMsgColor).
		AddInfo("", "Deleted message from "+msg.Author.String()).
		AddInfo("Content", msg.Content).
		AddInfo("Server", guildName).
		AddInfo("Channel", channelName).
		AddInfo("Attachments", strings.Join(attachments, ", "), len(attachments) > 0).
		SetWebhookURL(cfg.Webhooks.Notifiers).
		Show()
}

func (ui *uiService) logDiscordStaffDetection(msg *discordgo.Message) {
	session := ui.discordSvc.GetSession()
	guildName := utils.Discord.GetGuildName(session, msg.GuildID)
	channelName := utils.Discord.GetChannelName(session, msg.ChannelID)

	cfg := ui.configSvc.GetConfig()

	_ = ui.ShowNotification("⚠️ Discord staff member detected: %s in %s - #%s",
		msg.Author.String(), guildName, channelName)

	_ = ui.NewConsoleMessage().
		SetTitle("Discord Staff Detection").
		SetColor(consoleMsgColor).
		AddInfo("", msg.Author.String()+" is Discord staff").
		AddInfo("Server", guildName).
		AddInfo("Channel", channelName).
		AddInfo("Message Link", msg.GetLink()).
		SetWebhookURL(cfg.Webhooks.StaffDetection).
		Show()
}

func (ui *uiService) logGhostPingDetection(msg *discordgo.Message) {
	session := ui.discordSvc.GetSession()
	guildName := utils.Discord.GetGuildName(session, msg.GuildID)
	channelName := utils.Discord.GetChannelName(session, msg.ChannelID)

	cfg := ui.configSvc.GetConfig()

	_ = ui.ShowNotification("%s ghost pinged you in %s - #%s", msg.Author.String(), guildName, channelName)
	_ = ui.NewConsoleMessage().
		SetTitle("Ghostping Detection").
		SetColor(consoleMsgColor).
		AddInfo("", msg.Author.String()+" has pinged you").
		AddInfo("Server", guildName).
		AddInfo("Channel", channelName).
		SetWebhookURL(cfg.Webhooks.PingDetection).
		Show()
}

func (ui *uiService) logSelfbotDetection(msg *discordgo.Message, selfbot string) {
	session := ui.discordSvc.GetSession()
	guildName := utils.Discord.GetGuildName(session, msg.GuildID)
	channelName := utils.Discord.GetChannelName(session, msg.ChannelID)

	cfg := ui.configSvc.GetConfig()

	_ = ui.ShowNotification("Detected %s using a selfbot (%s) in %s - #%s", msg.Author.String(), selfbot, guildName, channelName)
	_ = ui.NewConsoleMessage().
		SetTitle("Selfbot Detection").
		SetColor(consoleMsgColor).
		AddInfo("", msg.Author.String()+" was detected using a selfbot").
		AddInfo("Server", guildName).
		AddInfo("Channel", channelName).
		AddInfo("Selfbot", selfbot).
		AddInfo("Message Link", msg.GetLink()).
		SetWebhookURL(cfg.Webhooks.SelfbotDetection).
		Show()
}

func (ui *uiService) logTicketNotifier(channel *discordgo.Channel) {
	session := ui.discordSvc.GetSession()
	guildName := utils.Discord.GetGuildName(session, channel.GuildID)
	channelName := utils.Discord.GetChannelName(session, channel.ID)

	cfg := ui.configSvc.GetConfig()

	_ = ui.ShowNotification("New ticket was created in %s - #%s", guildName, channelName)
	_ = ui.NewConsoleMessage().
		SetTitle("Ticket Notifier").
		SetColor(consoleMsgColor).
		AddInfo("", "New ticket was created").
		AddInfo("Server", guildName).
		AddInfo("Channel", channelName).
		AddInfo("Ticket", channelName).
		SetWebhookURL(cfg.Webhooks.Notifiers).
		Show()
}

func (ui *uiService) logKeywordNotifier(msg *discordgo.Message, keyword string) {
	session := ui.discordSvc.GetSession()
	guildName := utils.Discord.GetGuildName(session, msg.GuildID)
	channelName := utils.Discord.GetChannelName(session, msg.ChannelID)

	cfg := ui.configSvc.GetConfig()

	_ = ui.ShowNotification("Keyword Notifier: Keyword sent by %s\nServer: %s\nChannel: %s\nKeyword: %s",
		msg.Author.Username, guildName, channelName, keyword)

	_ = ui.NewConsoleMessage().
		SetTitle("Keyword Notifier").
		SetColor(consoleMsgColor).
		AddInfo("", "Keyword found").
		AddInfo("Server", guildName).
		AddInfo("Channel", channelName).
		AddInfo("Author", msg.Author.String()).
		AddInfo("Keyword", keyword).
		AddInfo("Message Link", msg.GetLink()).
		SetWebhookURL(cfg.Webhooks.Misc).
		SetImportant(true).
		Show()
}

func (ui *uiService) logPhishingDetection(msg *discordgo.Message, domain string) {
	session := ui.discordSvc.GetSession()
	guildName := utils.Discord.GetGuildName(session, msg.GuildID)
	channelName := utils.Discord.GetChannelName(session, msg.ChannelID)

	cfg := ui.configSvc.GetConfig()

	_ = ui.ShowNotification("⚠️ Phishing link detected from %s in %s - #%s\nDomain: %s",
		msg.Author.String(), guildName, channelName, domain)

	_ = ui.NewConsoleMessage().
		SetTitle("Protections").
		SetColor(consoleMsgColor).
		AddInfo("", "Phishing Detected!").
		AddInfo("Domain", domain).
		AddInfo("Author", msg.Author.String()).
		AddInfo("Server", guildName).
		AddInfo("Channel", channelName).
		AddInfo("Message Link", msg.GetLink()).
		SetWebhookURL(cfg.Webhooks.Misc).
		Show()
}

func (ui *uiService) logMassDMBlocked() {
	cfg := ui.configSvc.GetConfig()

	_ = ui.ShowNotification("Mass DM spam blocked! DMs will be enabled again in 5 minutes")
	_ = ui.NewConsoleMessage().
		SetTitle("Protections").
		SetColor(consoleMsgColor).
		AddInfo("", "Mass DM Blocked!").
		AddInfo("", "DMs will be enabled again in 5 minutes.").
		SetWebhookURL(cfg.Webhooks.Misc).
		Show()
}

func (ui *uiService) logFriendRequestSpamBlocked() {
	cfg := ui.configSvc.GetConfig()

	_ = ui.ShowNotification("Friend request spam blocked! Friend requests will be enabled again in 5 minutes")
	_ = ui.NewConsoleMessage().
		SetTitle("Protections").
		SetColor(consoleMsgColor).
		AddInfo("", "Friend Request Spam Blocked!").
		AddInfo("", "Friend Requests will be enabled again in 5 minutes.").
		SetWebhookURL(cfg.Webhooks.Misc).
		Show()
}

func (ui *uiService) logAntiDelete(msg *discordgo.Message) {
	session := ui.discordSvc.GetSession()
	guildName := utils.Discord.GetGuildName(session, msg.GuildID)
	channelName := utils.Discord.GetChannelName(session, msg.ChannelID)

	cfg := ui.configSvc.GetConfig()

	_ = ui.ShowNotification("Prevented %s from deleting a message in %s - #%s", msg.Author.String(), guildName, channelName)
	_ = ui.NewConsoleMessage().
		SetTitle("Protections").
		SetColor(consoleMsgColor).
		AddInfo("", msg.Author.String()+" deleted a message").
		AddInfo("Server", guildName).
		AddInfo("Channel", channelName).
		SetWebhookURL(cfg.Webhooks.Misc).
		Show()
}

func (ui *uiService) logAFKMessage(msg *discordgo.Message, response string) {
	session := ui.discordSvc.GetSession()
	guildName := utils.Discord.GetGuildName(session, msg.GuildID)
	channelName := utils.Discord.GetChannelName(session, msg.ChannelID)

	cfg := ui.configSvc.GetConfig()

	_ = ui.ShowNotification("Sent AFK response to %s in %s - #%s",
		msg.Author.String(), guildName, channelName)

	_ = ui.NewConsoleMessage().
		SetTitle("AFK Responder").
		SetColor(consoleMsgColor).
		AddInfo("", "Sent AFK message").
		AddInfo("Server", guildName).
		AddInfo("Channel", channelName).
		AddInfo("Author", msg.Author.String()).
		AddInfo("Content", msg.Content).
		AddInfo("Response", response).
		AddInfo("Message Link", msg.GetLink()).
		SetWebhookURL(cfg.Webhooks.Misc).
		Show()
}

func (ui *uiService) logMultiTokenSniper(numTokens, numGuilds int) {
	cfg := ui.configSvc.GetConfig()

	notification := fmt.Sprintf("Connected to %d ", numTokens)
	notification += map[bool]string{true: "accounts", false: "account"}[numTokens > 1]
	notification += fmt.Sprintf(" and sniping in %d servers", numGuilds)

	_ = ui.ShowNotification("%s", notification)
	_ = ui.NewConsoleMessage().
		SetTitle("Multi Token Sniper").
		SetColor(consoleMsgColor).
		AddInfo("", notification).
		SetWebhookURL(cfg.Webhooks.Misc).
		Show()
}

func (ui *uiService) logSessionConnected(id, os, client string) {
	cfg := ui.configSvc.GetConfig()

	if !cfg.Events.SessionConnected {
		return
	}

	_ = ui.ShowNotification("New session connected from %s on %s", client, os)
	_ = ui.NewConsoleMessage().
		SetTitle("Session Connected").
		SetColor(consoleMsgColor).
		AddInfo("Session ID", id).
		AddInfo("OS", os).
		AddInfo("Client", client).
		SetWebhookURL(cfg.Webhooks.Sessions).
		Show()
}

func (ui *uiService) logSessionDisconnected(id, os, client string) {
	cfg := ui.configSvc.GetConfig()

	if !cfg.Events.SessionDisconnected {
		return
	}

	_ = ui.ShowNotification("Session disconnected from %s on %s", client, os)
	_ = ui.NewConsoleMessage().
		SetTitle("Session Disconnected").
		SetColor(consoleMsgColor).
		AddInfo("Session ID", id).
		AddInfo("OS", os).
		AddInfo("Client", client).
		SetWebhookURL(cfg.Webhooks.Sessions).
		Show()
}

func (ui *uiService) logSCJobs(numJobs int) {
	if numJobs == 0 {
		return
	}

	cfg := ui.configSvc.GetConfig()

	notification := fmt.Sprintf("Started %d slash command ", numJobs)
	notification += map[bool]string{true: "jobs", false: "job"}[numJobs > 1]

	_ = ui.ShowNotification("%s", notification)
	_ = ui.NewConsoleMessage().
		SetTitle("Slash Commands").
		SetColor(consoleMsgColor).
		AddInfo("", notification).
		SetWebhookURL(cfg.Webhooks.Misc).
		Show()
}

func (ui *uiService) logEventReactions(numReactions int) {
	cfg := ui.configSvc.GetConfig()

	notification := fmt.Sprintf("Set up %d event ", numReactions)
	notification += map[bool]string{true: "reactions", false: "reaction"}[numReactions > 1]

	_ = ui.ShowNotification("%s", notification)
	_ = ui.NewConsoleMessage().
		SetTitle("Event Reactions").
		SetColor(consoleMsgColor).
		AddInfo("", notification).
		SetWebhookURL(cfg.Webhooks.Misc).
		Show()
}

func (ui *uiService) logDangerousCmdBlocked(user *discordgo.User, cmd string) {
	cfg := ui.configSvc.GetConfig()

	_ = ui.ShowNotification("Blocked %s from using dangerous command: %s", user.DisplayName, cmd)
	_ = ui.NewConsoleMessage().
		SetTitle("Sharing").
		SetColor(consoleMsgColor).
		AddInfo("", "Blocked someone from using a dangerous command").
		AddInfo("User", user.String()).
		AddInfo("Command", cmd).
		SetWebhookURL(cfg.Webhooks.Misc).
		Show()
}

func (ui *uiService) logDangerousHostingCmdBlocked(cmd string) {
	cfg := ui.configSvc.GetConfig()

	_ = ui.ShowNotification("Blocked dangerous hosting command: %s", cmd)
	_ = ui.NewConsoleMessage().
		SetTitle("LightningBot Hosting").
		SetColor(consoleMsgColor).
		AddInfo("", "Blocked dangerous hosting command").
		AddInfo("Command", cmd).
		SetWebhookURL(cfg.Webhooks.Misc).
		Show()
}

func (ui *uiService) logReminder(reminder config.Reminder) {
	cfg := ui.configSvc.GetConfig()

	_ = ui.ShowNotification("Reminder: %s", reminder.Content)
	_ = ui.NewConsoleMessage().
		SetTitle("Reminder").
		SetColor(consoleMsgColor).
		AddInfo("Message", reminder.Content).
		SetWebhookURL(cfg.Webhooks.Misc).
		SetImportant(true).
		Show()
}
