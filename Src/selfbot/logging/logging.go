package logging

import (
	"fmt"
	"strings"
	"time"

	"github.com/LightningDev1/LB-Selfbot/selfbot/config"
	"github.com/LightningDev1/LB-Selfbot/selfbot/constants"
	"github.com/LightningDev1/LB-Selfbot/selfbot/data"
	"github.com/LightningDev1/LB-Selfbot/selfbot/utils"
	"github.com/LightningDev1/dgc"
	"github.com/LightningDev1/discordgo"
)

type _Logging struct{}

func (_Logging) PrintBanner(console bool, ui bool) {
	if console {
		Console.PrintBanner()
	}

	if ui {
		message := "<image Source=\"assets/Banner.png\" Width=\"650\" Height=\"150\" Centered />"
		UI.Set(message)
	}
}

// Specific logging functions

func (_Logging) CommandUsed(commandName string) {
	UI.Notification("Command " + commandName + " used")

	NewConsoleMessage().
		SetTitle("Command Used").
		SetColor("#C1EFFF").
		AddInfo("Command", commandName).
		Show()
}

func (_Logging) OnReady(router *dgc.Router, event *discordgo.Ready) {
	cfg, err := config.GetConfig()
	if err != nil {
		Console.Error("Error getting config:", err)
		return
	}

	if event.User.Discriminator == "0" {
		UI.Notification(fmt.Sprintf("Connected to %s (%s)", event.User.DisplayName, event.User.Username))
	} else {
		UI.Notification(fmt.Sprintf("Connected to %s", event.User.String()))
	}

	// Count friends
	friends := utils.Filter(event.Relationships, func(rel *discordgo.Relationship) bool {
		return rel.Type == 1
	})

	userText := event.User.Username
	if event.User.DisplayName != "" {
		userText = event.User.DisplayName + " | " + event.User.Username
	} else if event.User.Discriminator != "0" {
		userText = event.User.Username + "#" + event.User.Discriminator
	}

	data.TotalCommands = len(router.Commands)
	data.Guilds = len(event.Guilds)
	data.Friends = len(friends)
	data.AvatarImageKey = strings.ReplaceAll(event.User.AvatarURL(""), "https://cdn.discordapp.com/", "mp:")
	data.UserID = event.User.ID

	bannerTextUpper := fmt.Sprintf("LightningBot %s | %s | Prefix: %s", constants.VERSION, userText, cfg.Selfbot.CommandPrefix)
	bannerTextLower := fmt.Sprintf("%d Commands | %d Servers | %d Friends", len(router.Commands), len(event.Guilds), len(friends))

	// Console

	if constants.CONSOLE {
		upper := Console.Colored(utils.String.Center(bannerTextUpper), 0, 203, 255, false)
		lower := Console.Colored(utils.String.Center(bannerTextLower), 0, 203, 255, false)

		Console.Clear()
		Logging.PrintBanner(true, false)
		fmt.Println(upper + "\n" + lower + "\n")
	}

	// UI

	upper := Console.Colored(bannerTextUpper, 0, 203, 255, false)
	lower := Console.Colored(bannerTextLower, 0, 203, 255, false)

	UI.Set("ALIGN_CENTER" + upper + "\nALIGN_CENTER" + lower + "\n")
}

func (_Logging) NitroSniperFake(session *discordgo.Session, message *discordgo.Message, code, account string) {
	guildName := utils.Discord.GetGuildName(session, message.GuildID)
	channelName := utils.Discord.GetChannelName(session, message.ChannelID)
	UI.Notification(fmt.Sprintf("Nitro Sniper: Fake Code\nServer: %s\nChannel: %s\nAuthor: %s\nCode: %s", guildName, channelName, message.Author.String(), code))

	cfg, _ := config.GetConfig()

	NewConsoleMessage().
		SetTitle("Nitro Sniper").
		SetColor("#C1EFFF").
		AddInfo("Status", "Fake Code").
		AddInfo("Author", message.Author.String()).
		AddInfo("Server", guildName).
		AddInfo("Channel", channelName).
		AddInfo("Code", code).
		AddInfo("Message Link", message.GetLink()).
		AddInfoConditional("On Account", account, account != "").
		SetWebhookURL(cfg.Webhooks.NitroSniper).
		Show()
}

func (_Logging) NitroSniper(session *discordgo.Session, message *discordgo.Message, elapsed time.Duration, status, code, account string) {
	guildName := utils.Discord.GetGuildName(session, message.GuildID)
	channelName := utils.Discord.GetChannelName(session, message.ChannelID)
	UI.Notification(fmt.Sprintf("Nitro Sniper: %s\nServer: %s\nChannel: %s\nAuthor: %s\nCode: %s", status, guildName, channelName, message.Author.String(), code))
	cfg, _ := config.GetConfig()
	NewConsoleMessage().
		SetTitle("Nitro Sniper").
		SetColor("#C1EFFF").
		AddInfo("Status", status).
		AddInfo("Author", message.Author.String()).
		AddInfo("Server", guildName).
		AddInfo("Channel", channelName).
		AddInfo("Code", code).
		AddInfo("Snipe Time", elapsed.String()).
		AddInfo("Message Link", message.GetLink()).
		AddInfoConditional("On Account", account, account != "").
		SetWebhookURL(cfg.Webhooks.NitroSniper).
		Show()
}

func (_Logging) PrivnoteSniper(session *discordgo.Session, message *discordgo.Message, elapsed time.Duration, status, path, account string) {
	guildName := utils.Discord.GetGuildName(session, message.GuildID)
	channelName := utils.Discord.GetChannelName(session, message.ChannelID)
	UI.Notification(fmt.Sprintf("Privnote Sniper: %s\nServer: %s\nChannel: %s\nAuthor: %s\nSaved At: %s", status, guildName, channelName, message.Author.String(), path))
	cfg, _ := config.GetConfig()
	NewConsoleMessage().
		SetTitle("Privnote Sniper").
		SetColor("#C1EFFF").
		AddInfo("Status", status).
		AddInfo("Author", message.Author.String()).
		AddInfo("Server", guildName).
		AddInfo("Channel", channelName).
		AddInfo("Saved At", path).
		AddInfo("Snipe Time", elapsed.String()).
		AddInfo("Message Link", message.GetLink()).
		AddInfoConditional("On Account", account, account != "").
		SetWebhookURL(cfg.Webhooks.Misc).
		Show()
}

func (_Logging) BanNotifier(session *discordgo.Session, guildID string) {
	guildName := utils.Discord.GetGuildName(session, guildID)
	UI.Notification("Ban Notifier: You have been banned from " + guildName)
	cfg, _ := config.GetConfig()
	NewConsoleMessage().
		SetTitle("Ban Notifier").
		SetColor("#C1EFFF").
		AddInfo("", "You have been banned from "+guildName).
		SetWebhookURL(cfg.Webhooks.Notifiers).
		Show()
}

func (_Logging) DMNotifier(message *discordgo.Message) {
	UI.Notification("DM Notifier: New DM from " + message.Author.String())
	cfg, _ := config.GetConfig()
	NewConsoleMessage().
		SetTitle("DM Notifier").
		SetColor("#C1EFFF").
		AddInfo("", "New DM from "+message.Author.String()).
		AddInfo("Content", message.Content).
		AddInfo("Message Link", message.GetLink()).
		SetWebhookURL(cfg.Webhooks.Notifiers).
		Show()
}

func (_Logging) DMEditNotifier(message *discordgo.Message, before, after string) {
	UI.Notification("DM Edit Notifier: Edited DM from " + message.Author.String())
	cfg, _ := config.GetConfig()
	NewConsoleMessage().
		SetTitle("DM Edit Notifier").
		SetColor("#C1EFFF").
		AddInfo("", "Edited DM from "+message.Author.String()).
		AddInfo("Before", before).
		AddInfo("After", after).
		SetWebhookURL(cfg.Webhooks.Notifiers).
		Show()
}

func (_Logging) DMDeleteNotifier(message *discordgo.Message, attachments []string) {
	UI.Notification("DM Delete Notifier: Deleted DM from " + message.Author.String())
	cfg, _ := config.GetConfig()
	c := NewConsoleMessage().
		SetTitle("DM Delete Notifier").
		SetColor("#C1EFFF").
		AddInfo("", "Deleted DM from "+message.Author.String()).
		AddInfo("Content", message.Content)

	if len(attachments) > 0 {
		c.AddInfo("Attachments", strings.Join(attachments, ", "))
	}

	c.SetWebhookURL(cfg.Webhooks.Notifiers).Show()
}

func (_Logging) DMTypingNotifier(userName string) {
	UI.Notification("DM Typing Notifier: " + userName + " is typing")
	cfg, _ := config.GetConfig()
	NewConsoleMessage().
		SetTitle("DM Typing Notifier").
		SetColor("#C1EFFF").
		AddInfo("", userName+" is typing").
		SetWebhookURL(cfg.Webhooks.Notifiers).
		Show()
}

func (_Logging) MessageDeleteNotifier(session *discordgo.Session, message *discordgo.Message, attachments []string) {
	guildName := utils.Discord.GetGuildName(session, message.GuildID)
	channelName := utils.Discord.GetChannelName(session, message.ChannelID)
	UI.Notification("Message Delete Notifier: Deleted message from " + message.Author.String())
	cfg, _ := config.GetConfig()
	c := NewConsoleMessage().
		SetTitle("Message Delete Notifier").
		SetColor("#C1EFFF").
		AddInfo("", "Deleted message from "+message.Author.String()).
		AddInfo("Content", message.Content).
		AddInfo("Server", guildName).
		AddInfo("Channel", channelName)

	if len(attachments) > 0 {
		c.AddInfo("Attachments", strings.Join(attachments, ", "))
	}

	c.SetWebhookURL(cfg.Webhooks.Notifiers).Show()
}

func (_Logging) DiscordStaffDetection(session *discordgo.Session, message *discordgo.Message) {
	guildName := utils.Discord.GetGuildName(session, message.GuildID)
	channelName := utils.Discord.GetChannelName(session, message.ChannelID)
	UI.Notification("Discord Staff Detection: " + message.Author.String() + " is Discord staff")
	cfg, _ := config.GetConfig()
	NewConsoleMessage().
		SetTitle("Discord Staff Detection").
		SetColor("#C1EFFF").
		AddInfo("", message.Author.String()+" is Discord staff").
		AddInfo("Server", guildName).
		AddInfo("Channel", channelName).
		AddInfo("Message Link", message.GetLink()).
		SetWebhookURL(cfg.Webhooks.StaffDetection).
		Show()
}

func (_Logging) GhostPingDetection(session *discordgo.Session, message *discordgo.Message) {
	guildName := utils.Discord.GetGuildName(session, message.GuildID)
	channelName := utils.Discord.GetChannelName(session, message.ChannelID)
	UI.Notification("Ghostping Detection: " + message.Author.String() + " has pinged you")
	cfg, _ := config.GetConfig()
	NewConsoleMessage().
		SetTitle("Ghostping Detection").
		SetColor("#C1EFFF").
		AddInfo("", message.Author.String()+" has pinged you").
		AddInfo("Server", guildName).
		AddInfo("Channel", channelName).
		SetWebhookURL(cfg.Webhooks.PingDetection).
		Show()
}

func (_Logging) SelfbotDetection(session *discordgo.Session, message *discordgo.Message, selfbot string) {
	guildName := utils.Discord.GetGuildName(session, message.GuildID)
	channelName := utils.Discord.GetChannelName(session, message.ChannelID)
	UI.Notification("Selfbot Detection: " + message.Author.String() + " was detected using a selfbot")
	cfg, _ := config.GetConfig()
	NewConsoleMessage().
		SetTitle("Selfbot Detection").
		SetColor("#C1EFFF").
		AddInfo("", message.Author.String()+" was detected using a selfbot").
		AddInfo("Server", guildName).
		AddInfo("Channel", channelName).
		AddInfo("Selfbot", selfbot).
		AddInfo("Message Link", message.GetLink()).
		SetWebhookURL(cfg.Webhooks.SelfbotDetection).
		Show()
}

func (_Logging) TicketNotifier(session *discordgo.Session, channel *discordgo.Channel) {
	guildName := utils.Discord.GetGuildName(session, channel.GuildID)
	channelName := utils.Discord.GetChannelName(session, channel.ID)
	cfg, _ := config.GetConfig()
	NewConsoleMessage().
		SetTitle("Ticket Notifier").
		SetColor("#C1EFFF").
		AddInfo("", "New ticket was created").
		AddInfo("Server", guildName).
		AddInfo("Channel", channelName).
		AddInfo("Ticket", channelName).
		SetWebhookURL(cfg.Webhooks.Notifiers).
		Show()
}

func (_Logging) PhishingDetected(session *discordgo.Session, message *discordgo.Message, domain string) {
	guildName := utils.Discord.GetGuildName(session, message.GuildID)
	channelName := utils.Discord.GetChannelName(session, message.ChannelID)
	UI.Notification(fmt.Sprintf("Protections: Phishing Detected!\nDomain: %s\nAuthor: %s", domain, message.Author.String()))
	cfg, _ := config.GetConfig()
	NewConsoleMessage().
		SetTitle("Protections").
		SetColor("#C1EFFF").
		AddInfo("", "Phishing Detected!").
		AddInfo("Domain", domain).
		AddInfo("Author", message.Author.String()).
		AddInfo("Server", guildName).
		AddInfo("Channel", channelName).
		AddInfo("Message Link", message.GetLink()).
		SetWebhookURL(cfg.Webhooks.Misc).
		Show()
}

func (_Logging) MassDMBlocked() {
	UI.Notification("Protections: Mass DM Blocked!\nDMs will be enabled again in 5 minutes.")
	cfg, _ := config.GetConfig()
	NewConsoleMessage().
		SetTitle("Protections").
		SetColor("#C1EFFF").
		AddInfo("", "Mass DM Blocked!").
		AddInfo("", "DMs will be enabled again in 5 minutes.").
		SetWebhookURL(cfg.Webhooks.Misc).
		Show()
}

func (_Logging) FriendRequestSpamBlocked() {
	UI.Notification("Protections: Friend Request Spam Blocked!\nFriend Requests will be enabled again in 5 minutes.")
	cfg, _ := config.GetConfig()
	NewConsoleMessage().
		SetTitle("Protections").
		SetColor("#C1EFFF").
		AddInfo("", "Friend Request Spam Blocked!").
		AddInfo("", "Friend Requests will be enabled again in 5 minutes.").
		SetWebhookURL(cfg.Webhooks.Misc).
		Show()
}

func (_Logging) AntiDelete(session *discordgo.Session, message *discordgo.Message) {
	guildName := utils.Discord.GetGuildName(session, message.GuildID)
	channelName := utils.Discord.GetChannelName(session, message.ChannelID)
	UI.Notification("Anti Delete: " + message.Author.String() + " deleted a message in " + guildName)
	cfg, _ := config.GetConfig()
	NewConsoleMessage().
		SetTitle("Protections").
		SetColor("#C1EFFF").
		AddInfo("", message.Author.String()+" deleted a message").
		AddInfo("Server", guildName).
		AddInfo("Channel", channelName).
		SetWebhookURL(cfg.Webhooks.Misc).
		Show()
}

func (_Logging) GiveawayFound(session *discordgo.Session, message *discordgo.Message, prize, host, endsIn, account string) {
	guildName := utils.Discord.GetGuildName(session, message.GuildID)
	channelName := utils.Discord.GetChannelName(session, message.ChannelID)
	hostName := utils.Discord.GetUserName(session, host)
	if hostName == "Unknown" {
		hostName = host
	}
	cfg, _ := config.GetConfig()
	UI.Notification("Giveaway Sniper: Found a giveaway in " + guildName + "\nPrize: " + prize)
	NewConsoleMessage().
		SetTitle("Giveaway Sniper").
		SetColor("#C1EFFF").
		AddInfo("", "Found Giveaway").
		AddInfo("Server", guildName).
		AddInfo("Channel", channelName).
		AddInfo("Bot", message.Author.String()).
		AddInfoConditional("Prize", prize, prize != "").
		AddInfoConditional("Hosted By", hostName, hostName != "").
		AddInfoConditional("Ends In", endsIn, endsIn != "").
		AddInfo("Message Link", message.GetLink()).
		AddInfoConditional("On Account", account, account != "").
		SetWebhookURL(cfg.Webhooks.GiveawaySniper).
		Show()
}

func (_Logging) GiveawayWon(session *discordgo.Session, message *discordgo.Message, prize, account string) {
	guildName := utils.Discord.GetGuildName(session, message.GuildID)
	channelName := utils.Discord.GetChannelName(session, message.ChannelID)
	UI.Notification("Giveaway Sniper: Won a giveaway in " + guildName + "\nPrize: " + prize)
	cfg, _ := config.GetConfig()
	NewConsoleMessage().
		SetTitle("Giveaway Sniper").
		SetColor("#C1EFFF").
		AddInfo("", "Won Giveaway").
		AddInfo("Server", guildName).
		AddInfo("Channel", channelName).
		AddInfo("Bot", message.Author.String()).
		AddInfoConditional("Prize", prize, prize != "").
		AddInfo("Message Link", message.GetLink()).
		AddInfoConditional("On Account", account, account != "").
		SetWebhookURL(cfg.Webhooks.GiveawaySniper).
		Show()
}

func (_Logging) MultiTokenSniperConnected(amount, guilds int) {
	info := fmt.Sprintf("%d Account Connected, sniping in %d servers", amount, guilds)
	if amount > 1 {
		info = strings.ReplaceAll(info, "Account", "Accounts")
	}

	cfg, _ := config.GetConfig()

	UI.Notification("Multi Token Sniper: " + info)
	NewConsoleMessage().
		SetTitle("Multi Token Sniper").
		SetColor("#C1EFFF").
		AddInfo("", info).
		SetWebhookURL(cfg.Webhooks.Misc).
		Show()
}

func (_Logging) SessionConnected(id, os, client string) {
	cfg, err := config.GetConfig()
	if err != nil {
		return
	}

	if !cfg.Events.SessionConnected {
		return
	}

	UI.Notification("Session Connected\nOS: " + os + "\nClient: " + client)
	NewConsoleMessage().
		SetTitle("Session Connected").
		SetColor("#C1EFFF").
		AddInfo("Session ID", id).
		AddInfo("OS", os).
		AddInfo("Client", client).
		SetWebhookURL(cfg.Webhooks.Sessions).
		Show()
}

func (_Logging) SessionDisconnected(id, os, client string) {
	cfg, err := config.GetConfig()
	if err != nil {
		return
	}

	if !cfg.Events.SessionDisconnected {
		return
	}

	UI.Notification("Session Disconnected\nOS: " + os + "\nClient: " + client)
	NewConsoleMessage().
		SetTitle("Session Disconnected").
		SetColor("#C1EFFF").
		AddInfo("Session ID", id).
		AddInfo("OS", os).
		AddInfo("Client", client).
		SetWebhookURL(cfg.Webhooks.Sessions).
		Show()
}

func (_Logging) SlashCommandJobs(amount int) {
	cfg, err := config.GetConfig()
	if err != nil {
		return
	}

	info := fmt.Sprintf("Started %d slash command job", amount)

	if amount > 1 {
		info += "s"
	}

	UI.Notification(info)
	NewConsoleMessage().
		SetTitle("Slash Commands").
		SetColor("#C1EFFF").
		AddInfo("", info).
		SetWebhookURL(cfg.Webhooks.Misc).
		Show()
}

func (_Logging) EventReactions(amount int) {
	cfg, err := config.GetConfig()
	if err != nil {
		return
	}

	info := fmt.Sprintf("Initialized %d event reaction", amount)

	if amount > 1 {
		info += "s"
	}

	UI.Notification(info)
	NewConsoleMessage().
		SetTitle("Event Reactions").
		SetColor("#C1EFFF").
		AddInfo("", info).
		SetWebhookURL(cfg.Webhooks.Misc).
		Show()
}

func (_Logging) BlockedDangerousCommand(user *discordgo.User, command string) {
	cfg, err := config.GetConfig()
	if err != nil {
		return
	}

	info := fmt.Sprintf("Blocked %s from using a dangerous command (%s)", user.DisplayName, command)

	UI.Notification(info)
	NewConsoleMessage().
		SetTitle("Sharing").
		SetColor("#C1EFFF").
		AddInfo("", info).
		SetWebhookURL(cfg.Webhooks.Misc).
		Show()
}

func (_Logging) BlockedHostingDangerousCommand(command string) {
	cfg, err := config.GetConfig()
	if err != nil {
		return
	}

	info := fmt.Sprintf("Blocked usage of dangerous command %s", command)

	NewConsoleMessage().
		SetTitle("LightningBot Hosting").
		SetColor("#C1EFFF").
		AddInfo("", info).
		SetWebhookURL(cfg.Webhooks.Misc).
		Show()
}

func (_Logging) KeywordNotifier(session *discordgo.Session, message *discordgo.Message, keyword string) {
	cfg, err := config.GetConfig()
	if err != nil {
		return
	}

	guildName := utils.Discord.GetGuildName(session, message.GuildID)
	channelName := utils.Discord.GetChannelName(session, message.ChannelID)

	UI.Notification("Keyword Notifier: Keyword sent by " + message.Author.Username + "\nServer: " + guildName + "\nChannel: " + channelName + "\nKeyword: " + keyword)
	NewConsoleMessage().
		SetTitle("Keyword Notifier").
		SetColor("#C1EFFF").
		AddInfo("", "Keyword found").
		AddInfo("Server", guildName).
		AddInfo("Channel", channelName).
		AddInfo("Author", message.Author.String()).
		AddInfo("Keyword", keyword).
		AddInfo("Message Link", message.GetLink()).
		SetWebhookURL(cfg.Webhooks.Misc).
		Show()
}

func (_Logging) PokeTwoSniper(session *discordgo.Session, message *discordgo.Message, pokemon, received, iv string, confidence float64, elapsed time.Duration, account string) {
	cfg, err := config.GetConfig()
	if err != nil {
		return
	}

	guildName := utils.Discord.GetGuildName(session, message.GuildID)
	channelName := utils.Discord.GetChannelName(session, message.ChannelID)

	UI.Notification("PokéTwo Sniper: Caught a " + pokemon + " in " + guildName + "\nChannel: " + channelName)
	NewConsoleMessage().
		SetTitle("PokéTwo Sniper").
		SetColor("#C1EFFF").
		AddInfo("", "Caught a "+pokemon).
		AddInfoConditional("Received", received, received != "").
		AddInfoConditional("IV", iv, iv != "").
		AddInfo("AI Confidence", fmt.Sprint(confidence)).
		AddInfo("Detection Time", elapsed.String()).
		AddInfo("Server", guildName).
		AddInfo("Channel", channelName).
		AddInfo("Message Link", message.GetLink()).
		AddInfoConditional("On Account", account, account != "").
		SetWebhookURL(cfg.Webhooks.Misc).
		Show()
}

var Logging = _Logging{}
