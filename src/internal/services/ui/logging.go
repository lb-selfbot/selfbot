package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/LightningDev1/LB-Selfbot/selfbot/config"
	"github.com/LightningDev1/LB-Selfbot/selfbot/utils"
	"github.com/LightningDev1/discordgo"
)

func (ui *UIServiceImpl) LogCmdUsed(cmdName string) {
	_ = ui.ShowNotification("You used the %s command.", cmdName)
	_ = ui.NewConsoleMessage().
		SetTitle("Command Used").
		SetColor("#C1EFFF").
		AddInfo("Command", cmdName).
		Show()
}

func (ui *UIServiceImpl) LogReady(event *discordgo.Ready, numCmds int) {
	cfg, err := config.GetConfig()
	if err != nil {
		ui.LogError("Error getting config: %v", err)
		return
	}

	_ = ui.ShowNotification("Successfully connected as %s", event.User.DisplayName)

	// Count friends
	friends := utils.Filter(event.Relationships, func(rel *discordgo.Relationship) bool {
		return rel.Type == 1
	})

	userText := event.User.DisplayName + " | " + event.User.Username

	bannerTextUpper := fmt.Sprintf("LightningBot %s | %s | Prefix: %s", "VERSION", userText, cfg.Selfbot.CommandPrefix)
	bannerTextLower := fmt.Sprintf("%d Commands | %d Servers | %d Friends", numCmds, len(event.Guilds), len(friends))

	upper := ui.ANSIColored(bannerTextUpper, 0, 203, 255, false)
	lower := ui.ANSIColored(bannerTextLower, 0, 203, 255, false)

	ui.SetConsole("ALIGN_CENTER" + upper + "\nALIGN_CENTER" + lower + "\n")
}

func (ui *UIServiceImpl) LogNitroSniper(msg *discordgo.Message, elapsed time.Duration, status, code, account string) {
	session := ui.discordSvc.GetSession()
	guildName := utils.Discord.GetGuildName(session, msg.GuildID)
	channelName := utils.Discord.GetChannelName(session, msg.ChannelID)

	cfg, err := config.GetConfig()
	if err != nil {
		return
	}

	_ = ui.ShowNotification("Found a Nitro code in %s - #%s!\nStatus: %s\nSent by: %s",
		guildName, channelName, status, msg.Author.String())

	_ = ui.NewConsoleMessage().
		SetTitle("Nitro Sniper").
		SetColor("#C1EFFF").
		AddInfo("Status", status).
		AddInfo("Author", msg.Author.String()).
		AddInfo("Server", guildName).
		AddInfo("Channel", channelName).
		AddInfo("Code", code).
		AddInfo("Snipe Time", elapsed.String()).
		AddInfo("Message Link", msg.GetLink()).
		AddInfoConditional("On Account", account, account != "").
		SetWebhookURL(cfg.Webhooks.NitroSniper).
		SetImportant(status == "Redeemed Nitro").
		Show()
}

func (ui *UIServiceImpl) LogNitroSniperFake(msg *discordgo.Message, code, account string) {
	session := ui.discordSvc.GetSession()
	guildName := utils.Discord.GetGuildName(session, msg.GuildID)
	channelName := utils.Discord.GetChannelName(session, msg.ChannelID)

	cfg, err := config.GetConfig()
	if err != nil {
		return
	}

	_ = ui.ShowNotification("Found a fake Nitro code in %s - #%s from %s",
		guildName, channelName, msg.Author.String())

	_ = ui.NewConsoleMessage().
		SetTitle("Nitro Sniper").
		SetColor("#C1EFFF").
		AddInfo("Status", "Fake Code").
		AddInfo("Author", msg.Author.String()).
		AddInfo("Server", guildName).
		AddInfo("Channel", channelName).
		AddInfo("Code", code).
		AddInfo("Message Link", msg.GetLink()).
		AddInfoConditional("On Account", account, account != "").
		SetWebhookURL(cfg.Webhooks.NitroSniper).
		Show()
}

func (ui *UIServiceImpl) LogPrivnoteSniper(msg *discordgo.Message, elapsed time.Duration, status, path, account string) {
	session := ui.discordSvc.GetSession()
	guildName := utils.Discord.GetGuildName(session, msg.GuildID)
	channelName := utils.Discord.GetChannelName(session, msg.ChannelID)

	cfg, err := config.GetConfig()
	if err != nil {
		return
	}

	_ = ui.ShowNotification("Sniped a Privnote in %s - $%s!\nFrom: %s\nSaved to: %s",
		guildName, channelName, msg.Author.String(), path)

	_ = ui.NewConsoleMessage().
		SetTitle("Privnote Sniper").
		SetColor("#C1EFFF").
		AddInfo("Status", status).
		AddInfo("Author", msg.Author.String()).
		AddInfo("Server", guildName).
		AddInfo("Channel", channelName).
		AddInfo("Saved At", path).
		AddInfo("Snipe Time", elapsed.String()).
		AddInfo("Message Link", msg.GetLink()).
		AddInfoConditional("On Account", account, account != "").
		SetWebhookURL(cfg.Webhooks.Misc).
		SetImportant(status == "Saved Note").
		Show()
}

func (ui *UIServiceImpl) LogGiveawaySniperFound(msg *discordgo.Message, prize, host, endsIn, account string) {
	session := ui.discordSvc.GetSession()
	guildName := utils.Discord.GetGuildName(session, msg.GuildID)
	channelName := utils.Discord.GetChannelName(session, msg.ChannelID)
	hostName := utils.Discord.GetUserName(session, host)
	if hostName == "Unknown" {
		hostName = host
	}

	cfg, err := config.GetConfig()
	if err != nil {
		return
	}

	notification := fmt.Sprintf("Found a giveaway in %s!", guildName)
	if prize != "" {
		notification += fmt.Sprintf("\nPrize: %s", prize)
	}
	if endsIn != "" {
		notification += fmt.Sprintf("\nEnds in: %s", endsIn)
	}

	_ = ui.ShowNotification(notification)
	_ = ui.NewConsoleMessage().
		SetTitle("Giveaway Sniper").
		SetColor("#C1EFFF").
		AddInfo("", "Found Giveaway").
		AddInfo("Server", guildName).
		AddInfo("Channel", channelName).
		AddInfo("Bot", msg.Author.String()).
		AddInfoConditional("Prize", prize, prize != "").
		AddInfoConditional("Hosted By", hostName, hostName != "").
		AddInfoConditional("Ends In", endsIn, endsIn != "").
		AddInfo("Message Link", msg.GetLink()).
		AddInfoConditional("On Account", account, account != "").
		SetWebhookURL(cfg.Webhooks.GiveawaySniper).
		Show()
}

func (ui *UIServiceImpl) LogGiveawaySniperWon(msg *discordgo.Message, prize, account string) {
	session := ui.discordSvc.GetSession()
	guildName := utils.Discord.GetGuildName(session, msg.GuildID)
	channelName := utils.Discord.GetChannelName(session, msg.ChannelID)

	cfg, err := config.GetConfig()
	if err != nil {
		return
	}

	notification := fmt.Sprintf("You won a giveaway in %s! 🎉", guildName)
	if prize != "" {
		notification += fmt.Sprintf("\nPrize: %s", prize)
	}

	_ = ui.ShowNotification(notification)
	_ = ui.NewConsoleMessage().
		SetTitle("Giveaway Sniper").
		SetColor("#C1EFFF").
		AddInfo("", "Won Giveaway").
		AddInfo("Server", guildName).
		AddInfo("Channel", channelName).
		AddInfo("Bot", msg.Author.String()).
		AddInfoConditional("Prize", prize, prize != "").
		AddInfo("Message Link", msg.GetLink()).
		AddInfoConditional("On Account", account, account != "").
		SetWebhookURL(cfg.Webhooks.GiveawaySniper).
		SetImportant(true).
		Show()
}

func (ui *UIServiceImpl) LogPokeTwoSniper(msg *discordgo.Message, elapsed time.Duration, confidence float64, pokemon, received, iv, account string) {
	session := ui.discordSvc.GetSession()
	guildName := utils.Discord.GetGuildName(session, msg.GuildID)
	channelName := utils.Discord.GetChannelName(session, msg.ChannelID)

	cfg, err := config.GetConfig()
	if err != nil {
		return
	}

	notification := fmt.Sprintf("PokéTwo Sniper caught a %s in %s - $%s! ✨", pokemon, guildName, channelName)
	if iv != "" {
		notification += fmt.Sprintf("\nIV: %s", iv)
	}

	_ = ui.ShowNotification(notification)
	_ = ui.NewConsoleMessage().
		SetTitle("PokéTwo Sniper").
		SetColor("#C1EFFF").
		AddInfo("", "Caught a "+pokemon).
		AddInfoConditional("Received", received, received != "").
		AddInfoConditional("IV", iv, iv != "").
		AddInfo("AI Confidence", fmt.Sprint(confidence)).
		AddInfo("Detection Time", elapsed.String()).
		AddInfo("Server", guildName).
		AddInfo("Channel", channelName).
		AddInfo("Message Link", msg.GetLink()).
		AddInfoConditional("On Account", account, account != "").
		SetWebhookURL(cfg.Webhooks.Misc).
		SetImportant(true).
		Show()
}

func (ui *UIServiceImpl) LogBanNotifier(guildID string) {
	session := ui.discordSvc.GetSession()
	guildName := utils.Discord.GetGuildName(session, guildID)

	cfg, err := config.GetConfig()
	if err != nil {
		return
	}

	_ = ui.ShowNotification("You have been banned from %s!", guildName)
	_ = ui.NewConsoleMessage().
		SetTitle("Ban Notifier").
		SetColor("#C1EFFF").
		AddInfo("", "You have been banned from "+guildName).
		SetWebhookURL(cfg.Webhooks.Notifiers).
		Show()
}

func (ui *UIServiceImpl) LogDMNotifier(msg *discordgo.Message) {
	cfg, err := config.GetConfig()
	if err != nil {
		return
	}

	_ = ui.ShowNotification("New DM from %s", msg.Author.String())
	_ = ui.NewConsoleMessage().
		SetTitle("DM Notifier").
		SetColor("#C1EFFF").
		AddInfo("", "New DM from "+msg.Author.String()).
		AddInfo("Content", msg.Content).
		AddInfo("Message Link", msg.GetLink()).
		SetWebhookURL(cfg.Webhooks.Notifiers).
		Show()
}

func (ui *UIServiceImpl) LogDMEditNotifier(msg *discordgo.Message, before, after string) {
	cfg, err := config.GetConfig()
	if err != nil {
		return
	}

	_ = ui.ShowNotification("%s edited their DM", msg.Author.String())
	_ = ui.NewConsoleMessage().
		SetTitle("DM Edit Notifier").
		SetColor("#C1EFFF").
		AddInfo("", "Edited DM from "+msg.Author.String()).
		AddInfo("Before", before).
		AddInfo("After", after).
		SetWebhookURL(cfg.Webhooks.Notifiers).
		Show()
}

func (ui *UIServiceImpl) LogDMDeleteNotifier(msg *discordgo.Message, attachments []string) {
	cfg, err := config.GetConfig()
	if err != nil {
		return
	}

	_ = ui.ShowNotification("%s deleted their DM", msg.Author.String())
	_ = ui.NewConsoleMessage().
		SetTitle("DM Delete Notifier").
		SetColor("#C1EFFF").
		AddInfo("", "Deleted DM from "+msg.Author.String()).
		AddInfo("Content", msg.Content).
		AddInfoConditional("Attachments", strings.Join(attachments, ", "), len(attachments) > 0).
		SetWebhookURL(cfg.Webhooks.Notifiers).
		Show()
}

func (ui *UIServiceImpl) LogDMTypingNotifier(username string) {
	cfg, err := config.GetConfig()
	if err != nil {
		return
	}

	_ = ui.ShowNotification("%s is typing in your DM...", username)
	_ = ui.NewConsoleMessage().
		SetTitle("DM Typing Notifier").
		SetColor("#C1EFFF").
		AddInfo("", username+" is typing").
		SetWebhookURL(cfg.Webhooks.Notifiers).
		Show()
}

func (ui *UIServiceImpl) LogMessageDeleteNotifier(msg *discordgo.Message, attachments []string) {
	session := ui.discordSvc.GetSession()
	guildName := utils.Discord.GetGuildName(session, msg.GuildID)
	channelName := utils.Discord.GetChannelName(session, msg.ChannelID)

	cfg, err := config.GetConfig()
	if err != nil {
		return
	}

	_ = ui.ShowNotification("%s deleted a message in %s - #%s", msg.Author.String(), guildName, channelName)
	_ = ui.NewConsoleMessage().
		SetTitle("Message Delete Notifier").
		SetColor("#C1EFFF").
		AddInfo("", "Deleted message from "+msg.Author.String()).
		AddInfo("Content", msg.Content).
		AddInfo("Server", guildName).
		AddInfo("Channel", channelName).
		AddInfoConditional("Attachments", strings.Join(attachments, ", "), len(attachments) > 0).
		SetWebhookURL(cfg.Webhooks.Notifiers).
		Show()
}

func (ui *UIServiceImpl) LogDiscordStaffDetection(msg *discordgo.Message) {
	session := ui.discordSvc.GetSession()
	guildName := utils.Discord.GetGuildName(session, msg.GuildID)
	channelName := utils.Discord.GetChannelName(session, msg.ChannelID)

	cfg, err := config.GetConfig()
	if err != nil {
		return
	}

	_ = ui.ShowNotification("⚠️ Discord staff member detected: %s in %s - #%s",
		msg.Author.String(), guildName, channelName)

	_ = ui.NewConsoleMessage().
		SetTitle("Discord Staff Detection").
		SetColor("#C1EFFF").
		AddInfo("", msg.Author.String()+" is Discord staff").
		AddInfo("Server", guildName).
		AddInfo("Channel", channelName).
		AddInfo("Message Link", msg.GetLink()).
		SetWebhookURL(cfg.Webhooks.StaffDetection).
		Show()
}

func (ui *UIServiceImpl) LogGhostPingDetection(msg *discordgo.Message) {
	session := ui.discordSvc.GetSession()
	guildName := utils.Discord.GetGuildName(session, msg.GuildID)
	channelName := utils.Discord.GetChannelName(session, msg.ChannelID)

	cfg, err := config.GetConfig()
	if err != nil {
		return
	}

	_ = ui.ShowNotification("%s ghost pinged you in %s - #%s", msg.Author.String(), guildName, channelName)
	_ = ui.NewConsoleMessage().
		SetTitle("Ghostping Detection").
		SetColor("#C1EFFF").
		AddInfo("", msg.Author.String()+" has pinged you").
		AddInfo("Server", guildName).
		AddInfo("Channel", channelName).
		SetWebhookURL(cfg.Webhooks.PingDetection).
		Show()
}

func (ui *UIServiceImpl) LogSelfbotDetection(msg *discordgo.Message, selfbot string) {
	session := ui.discordSvc.GetSession()
	guildName := utils.Discord.GetGuildName(session, msg.GuildID)
	channelName := utils.Discord.GetChannelName(session, msg.ChannelID)

	cfg, err := config.GetConfig()
	if err != nil {
		return
	}

	_ = ui.ShowNotification("Detected %s using a selfbot (%s) in %s - #%s", msg.Author.String(), selfbot, guildName, channelName)
	_ = ui.NewConsoleMessage().
		SetTitle("Selfbot Detection").
		SetColor("#C1EFFF").
		AddInfo("", msg.Author.String()+" was detected using a selfbot").
		AddInfo("Server", guildName).
		AddInfo("Channel", channelName).
		AddInfo("Selfbot", selfbot).
		AddInfo("Message Link", msg.GetLink()).
		SetWebhookURL(cfg.Webhooks.SelfbotDetection).
		Show()
}

func (ui *UIServiceImpl) LogTicketNotifier(channel *discordgo.Channel) {
	session := ui.discordSvc.GetSession()
	guildName := utils.Discord.GetGuildName(session, channel.GuildID)
	channelName := utils.Discord.GetChannelName(session, channel.ID)

	cfg, err := config.GetConfig()
	if err != nil {
		return
	}

	_ = ui.ShowNotification("New ticket was created in %s - #%s", guildName, channelName)
	_ = ui.NewConsoleMessage().
		SetTitle("Ticket Notifier").
		SetColor("#C1EFFF").
		AddInfo("", "New ticket was created").
		AddInfo("Server", guildName).
		AddInfo("Channel", channelName).
		AddInfo("Ticket", channelName).
		SetWebhookURL(cfg.Webhooks.Notifiers).
		Show()
}

func (ui *UIServiceImpl) LogKeywordNotifier(msg *discordgo.Message, keyword string) {
	session := ui.discordSvc.GetSession()
	guildName := utils.Discord.GetGuildName(session, msg.GuildID)
	channelName := utils.Discord.GetChannelName(session, msg.ChannelID)

	cfg, err := config.GetConfig()
	if err != nil {
		return
	}

	_ = ui.ShowNotification("Keyword Notifier: Keyword sent by %s\nServer: %s\nChannel: %s\nKeyword: %s",
		msg.Author.Username, guildName, channelName, keyword)

	_ = ui.NewConsoleMessage().
		SetTitle("Keyword Notifier").
		SetColor("#C1EFFF").
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

func (ui *UIServiceImpl) LogPhishingDetection(msg *discordgo.Message, domain string) {
	session := ui.discordSvc.GetSession()
	guildName := utils.Discord.GetGuildName(session, msg.GuildID)
	channelName := utils.Discord.GetChannelName(session, msg.ChannelID)

	cfg, err := config.GetConfig()
	if err != nil {
		return
	}

	_ = ui.ShowNotification("⚠️ Phishing link detected from %s in %s - #%s\nDomain: %s",
		msg.Author.String(), guildName, channelName, domain)

	_ = ui.NewConsoleMessage().
		SetTitle("Protections").
		SetColor("#C1EFFF").
		AddInfo("", "Phishing Detected!").
		AddInfo("Domain", domain).
		AddInfo("Author", msg.Author.String()).
		AddInfo("Server", guildName).
		AddInfo("Channel", channelName).
		AddInfo("Message Link", msg.GetLink()).
		SetWebhookURL(cfg.Webhooks.Misc).
		Show()
}

func (ui *UIServiceImpl) LogMassDMBlocked() {
	cfg, err := config.GetConfig()
	if err != nil {
		return
	}

	_ = ui.ShowNotification("Mass DM spam blocked! DMs will be enabled again in 5 minutes")
	_ = ui.NewConsoleMessage().
		SetTitle("Protections").
		SetColor("#C1EFFF").
		AddInfo("", "Mass DM Blocked!").
		AddInfo("", "DMs will be enabled again in 5 minutes.").
		SetWebhookURL(cfg.Webhooks.Misc).
		Show()
}

func (ui *UIServiceImpl) LogFriendRequestSpamBlocked() {
	cfg, err := config.GetConfig()
	if err != nil {
		return
	}

	_ = ui.ShowNotification("Friend request spam blocked! Friend requests will be enabled again in 5 minutes")
	_ = ui.NewConsoleMessage().
		SetTitle("Protections").
		SetColor("#C1EFFF").
		AddInfo("", "Friend Request Spam Blocked!").
		AddInfo("", "Friend Requests will be enabled again in 5 minutes.").
		SetWebhookURL(cfg.Webhooks.Misc).
		Show()
}

func (ui *UIServiceImpl) LogAntiDelete(msg *discordgo.Message) {
	session := ui.discordSvc.GetSession()
	guildName := utils.Discord.GetGuildName(session, msg.GuildID)
	channelName := utils.Discord.GetChannelName(session, msg.ChannelID)

	cfg, err := config.GetConfig()
	if err != nil {
		return
	}

	_ = ui.ShowNotification("Prevented %s from deleting a message in %s - #%s", msg.Author.String(), guildName, channelName)
	_ = ui.NewConsoleMessage().
		SetTitle("Protections").
		SetColor("#C1EFFF").
		AddInfo("", msg.Author.String()+" deleted a message").
		AddInfo("Server", guildName).
		AddInfo("Channel", channelName).
		SetWebhookURL(cfg.Webhooks.Misc).
		Show()
}

func (ui *UIServiceImpl) LogAFKMessage(msg *discordgo.Message, response string) {
	session := ui.discordSvc.GetSession()
	guildName := utils.Discord.GetGuildName(session, msg.GuildID)
	channelName := utils.Discord.GetChannelName(session, msg.ChannelID)

	cfg, err := config.GetConfig()
	if err != nil {
		return
	}

	_ = ui.ShowNotification("Sent AFK response to %s in %s - #%s",
		msg.Author.String(), guildName, channelName)

	_ = ui.NewConsoleMessage().
		SetTitle("AFK Responder").
		SetColor("#C1EFFF").
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

func (ui *UIServiceImpl) LogMultiTokenSniper(numTokens, numGuilds int) {
	cfg, err := config.GetConfig()
	if err != nil {
		return
	}

	notification := fmt.Sprintf("Connected to %d ", numTokens)
	notification += map[bool]string{true: "accounts", false: "account"}[numTokens > 1]
	notification += fmt.Sprintf(" and sniping in %d servers", numGuilds)

	_ = ui.ShowNotification(notification)
	_ = ui.NewConsoleMessage().
		SetTitle("Multi Token Sniper").
		SetColor("#C1EFFF").
		AddInfo("", notification).
		SetWebhookURL(cfg.Webhooks.Misc).
		Show()
}

func (ui *UIServiceImpl) LogSessionConnected(id, os, client string) {
	cfg, err := config.GetConfig()
	if err != nil {
		return
	}

	if !cfg.Events.SessionConnected {
		return
	}

	_ = ui.ShowNotification("New session connected from %s on %s", client, os)
	_ = ui.NewConsoleMessage().
		SetTitle("Session Connected").
		SetColor("#C1EFFF").
		AddInfo("Session ID", id).
		AddInfo("OS", os).
		AddInfo("Client", client).
		SetWebhookURL(cfg.Webhooks.Sessions).
		Show()
}

func (ui *UIServiceImpl) LogSessionDisconnected(id, os, client string) {
	cfg, err := config.GetConfig()
	if err != nil {
		return
	}

	if !cfg.Events.SessionDisconnected {
		return
	}

	_ = ui.ShowNotification("Session disconnected from %s on %s", client, os)
	_ = ui.NewConsoleMessage().
		SetTitle("Session Disconnected").
		SetColor("#C1EFFF").
		AddInfo("Session ID", id).
		AddInfo("OS", os).
		AddInfo("Client", client).
		SetWebhookURL(cfg.Webhooks.Sessions).
		Show()
}

func (ui *UIServiceImpl) LogSCJobs(numJobs int) {
	if numJobs == 0 {
		return
	}

	cfg, err := config.GetConfig()
	if err != nil {
		return
	}

	notification := fmt.Sprintf("Started %d slash command ", numJobs)
	notification += map[bool]string{true: "jobs", false: "job"}[numJobs > 1]

	_ = ui.ShowNotification(notification)
	_ = ui.NewConsoleMessage().
		SetTitle("Slash Commands").
		SetColor("#C1EFFF").
		AddInfo("", notification).
		SetWebhookURL(cfg.Webhooks.Misc).
		Show()
}

func (ui *UIServiceImpl) LogEventReactions(numReactions int) {
	cfg, err := config.GetConfig()
	if err != nil {
		return
	}

	notification := fmt.Sprintf("Set up %d event ", numReactions)
	notification += map[bool]string{true: "reactions", false: "reaction"}[numReactions > 1]

	_ = ui.ShowNotification(notification)
	_ = ui.NewConsoleMessage().
		SetTitle("Event Reactions").
		SetColor("#C1EFFF").
		AddInfo("", notification).
		SetWebhookURL(cfg.Webhooks.Misc).
		Show()
}

func (ui *UIServiceImpl) LogDangerousCmdBlocked(user *discordgo.User, cmd string) {
	cfg, err := config.GetConfig()
	if err != nil {
		return
	}

	notification := fmt.Sprintf("Blocked %s from using dangerous command: %s", user.DisplayName, cmd)

	_ = ui.ShowNotification(notification)
	_ = ui.NewConsoleMessage().
		SetTitle("Sharing").
		SetColor("#C1EFFF").
		AddInfo("", notification).
		SetWebhookURL(cfg.Webhooks.Misc).
		Show()
}

func (ui *UIServiceImpl) LogDangerousHostingCmdBlocked(cmd string) {
	cfg, err := config.GetConfig()
	if err != nil {
		return
	}

	notification := fmt.Sprintf("Blocked dangerous hosting command: %s", cmd)

	_ = ui.ShowNotification(notification)
	_ = ui.NewConsoleMessage().
		SetTitle("LightningBot Hosting").
		SetColor("#C1EFFF").
		AddInfo("", notification).
		SetWebhookURL(cfg.Webhooks.Misc).
		Show()
}

func (ui *UIServiceImpl) LogReminder(reminder config.Reminder) {
	cfg, err := config.GetConfig()
	if err != nil {
		return
	}

	_ = ui.ShowNotification("Reminder: %s", reminder.Content)
	_ = ui.NewConsoleMessage().
		SetTitle("Reminder").
		SetColor("#C1EFFF").
		AddInfo("Message", reminder.Content).
		SetWebhookURL(cfg.Webhooks.Misc).
		SetImportant(true).
		Show()
}
