package logging

import (
	"fmt"
	"strings"
	"time"

	"github.com/LightningDev1/LB-Selfbot/config"
	"github.com/LightningDev1/LB-Selfbot/constants"
	"github.com/LightningDev1/LB-Selfbot/utils"
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
		UI.Append(message)
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

	// Count friends
	friends := 0
	for _, rel := range event.Relationships {
		if rel.Type == 1 {
			friends++
		}
	}

	UI.Notification(fmt.Sprintf("Connected to %s", event.User.String()))

	bannerTextUpper := fmt.Sprintf("LightningBot %s | %s | Prefix: %s", constants.VERSION, event.User.String(), cfg.Selfbot.CommandPrefix)
	bannerTextLower := fmt.Sprintf("%d Commands | %d Servers | %d Friends", len(router.Commands), len(event.Guilds), friends)

	if constants.CONSOLE {
		Console.Clear()

		Logging.PrintBanner(true, false)

		message := Console.Colored(
			utils.String.Center(bannerTextUpper)+"\n"+utils.String.Center(bannerTextLower), 0, 203, 255, false,
		)

		fmt.Println(message)
	}

	message := "<image Source=\"assets/Banner.png\" Width=\"650\" Height=\"150\" Centered />\n" +
		UI.Centered(bannerTextUpper, "#C1EFFF") + "\n" +
		UI.Centered(bannerTextLower, "#C1EFFF") + "\n"

	UI.Clear()
	UI.Append(message)
}

func (_Logging) NitroSniperFake(session *discordgo.Session, message *discordgo.Message, code, account string) {
	guildName := utils.Discord.GetGuildName(session, message.GuildID)
	channelName := utils.Discord.GetChannelName(session, message.ChannelID)
	UI.Notification(fmt.Sprintf("Nitro Sniper: Fake Code\nServer: %s\nChannel: %s\nAuthor: %s\nCode: %s", guildName, channelName, message.Author.String(), code))
	NewConsoleMessage().
		SetTitle("Nitro Sniper").
		SetColor("#C1EFFF").
		AddInfo("Status", "Fake Code").
		AddInfo("Author", message.Author.String()).
		AddInfo("Server", guildName).
		AddInfo("Channel", channelName).
		AddInfo("Code", code).
		AddInfoConditional("On Account", account, account != "").
		Show()
}

func (_Logging) NitroSniper(session *discordgo.Session, message *discordgo.Message, elapsed time.Duration, status, code, account string) {
	guildName := utils.Discord.GetGuildName(session, message.GuildID)
	channelName := utils.Discord.GetChannelName(session, message.ChannelID)
	UI.Notification(fmt.Sprintf("Nitro Sniper: %s\nServer: %s\nChannel: %s\nAuthor: %s\nCode: %s", status, guildName, channelName, message.Author.String(), code))
	NewConsoleMessage().
		SetTitle("Nitro Sniper").
		SetColor("#C1EFFF").
		AddInfo("Status", status).
		AddInfo("Author", message.Author.String()).
		AddInfo("Server", guildName).
		AddInfo("Channel", channelName).
		AddInfo("Code", code).
		AddInfo("Snipe Time", elapsed.String()).
		AddInfoConditional("On Account", account, account != "").
		Show()
}

func (_Logging) BanNotifier(session *discordgo.Session, guildID string) {
	guildName := utils.Discord.GetGuildName(session, guildID)
	UI.Notification("Ban Notifier: You have been banned from " + guildName)
	NewConsoleMessage().
		SetTitle("Ban Notifier").
		SetColor("#C1EFFF").
		AddInfo("", "You have been banned from "+guildName).
		Show()
}

func (_Logging) DMNotifier(message *discordgo.Message) {
	UI.Notification("DM Notifier: New DM from " + message.Author.String())
	NewConsoleMessage().
		SetTitle("DM Notifier").
		SetColor("#C1EFFF").
		AddInfo("", "New DM from "+message.Author.String()).
		Show()
}

func (_Logging) DMEditNotifier(message *discordgo.Message) {
	UI.Notification("DM Edit Notifier: Edited DM from " + message.Author.String())
	NewConsoleMessage().
		SetTitle("DM Edit Notifier").
		SetColor("#C1EFFF").
		AddInfo("", "Edited DM from "+message.Author.String()).
		Show()
}

func (_Logging) DMDeleteNotifier(message *discordgo.Message) {
	UI.Notification("DM Delete Notifier: Deleted DM from " + message.Author.String())
	NewConsoleMessage().
		SetTitle("DM Delete Notifier").
		SetColor("#C1EFFF").
		AddInfo("", "Deleted DM from "+message.Author.String()).
		AddInfo("Content", message.Content).
		Show()
}

func (_Logging) DMTypingNotifier(userName string) {
	UI.Notification("DM Typing Notifier: " + userName + " is typing")
	NewConsoleMessage().
		SetTitle("DM Typing Notifier").
		SetColor("#C1EFFF").
		AddInfo("", userName+" is typing").
		Show()
}

func (_Logging) MessageDeleteNotifier(session *discordgo.Session, message *discordgo.Message) {
	guildName := utils.Discord.GetGuildName(session, message.GuildID)
	channelName := utils.Discord.GetChannelName(session, message.ChannelID)
	UI.Notification("Message Delete Notifier: Deleted message from " + message.Author.String())
	NewConsoleMessage().
		SetTitle("Message Delete Notifier").
		SetColor("#C1EFFF").
		AddInfo("", "Deleted message from "+message.Author.String()).
		AddInfo("Content", message.Content).
		AddInfo("Server", guildName).
		AddInfo("Channel", channelName).
		Show()
}

func (_Logging) DiscordStaffDetection(session *discordgo.Session, message *discordgo.Message) {
	guildName := utils.Discord.GetGuildName(session, message.GuildID)
	channelName := utils.Discord.GetChannelName(session, message.ChannelID)
	UI.Notification("Discord Staff Detection: " + message.Author.String() + " is Discord staff")
	NewConsoleMessage().
		SetTitle("Discord Staff Detection").
		SetColor("#C1EFFF").
		AddInfo("", message.Author.String()+" is Discord staff").
		AddInfo("Server", guildName).
		AddInfo("Channel", channelName).
		Show()
}

func (_Logging) GhostPingDetection(session *discordgo.Session, message *discordgo.Message) {
	guildName := utils.Discord.GetGuildName(session, message.GuildID)
	channelName := utils.Discord.GetChannelName(session, message.ChannelID)
	UI.Notification("Ghostping Detection: " + message.Author.String() + " has pinged you")
	NewConsoleMessage().
		SetTitle("Ghostping Detection").
		SetColor("#C1EFFF").
		AddInfo("", message.Author.String()+" has pinged you").
		AddInfo("Server", guildName).
		AddInfo("Channel", channelName).
		Show()
}

func (_Logging) SelfbotDetection(session *discordgo.Session, message *discordgo.Message, selfbot string) {
	guildName := utils.Discord.GetGuildName(session, message.GuildID)
	channelName := utils.Discord.GetChannelName(session, message.ChannelID)
	UI.Notification("Selfbot Detection: " + message.Author.String() + " was detected using a selfbot")
	NewConsoleMessage().
		SetTitle("Selfbot Detection").
		SetColor("#C1EFFF").
		AddInfo("", message.Author.String()+" was detected using a selfbot").
		AddInfo("Server", guildName).
		AddInfo("Channel", channelName).
		AddInfo("Selfbot", selfbot).
		Show()
}

func (_Logging) TicketNotifier(session *discordgo.Session, channel *discordgo.Channel) {
	guildName := utils.Discord.GetGuildName(session, channel.GuildID)
	channelName := utils.Discord.GetChannelName(session, channel.ID)
	NewConsoleMessage().
		SetTitle("Ticket Notifier").
		SetColor("#C1EFFF").
		AddInfo("", "New ticket was created").
		AddInfo("Server", guildName).
		AddInfo("Channel", channelName).
		AddInfo("Ticket", channelName).
		Show()
}

func (_Logging) PhisingDetected(session *discordgo.Session, message *discordgo.Message, domain string) {
	guildName := utils.Discord.GetGuildName(session, message.GuildID)
	channelName := utils.Discord.GetChannelName(session, message.ChannelID)
	UI.Notification(fmt.Sprintf("Protections: Phising Detected!\nDomain: %s\nAuthor: %s", domain, message.Author.String()))
	NewConsoleMessage().
		SetTitle("Protections").
		SetColor("#C1EFFF").
		AddInfo("", "Phising Detected!").
		AddInfo("Domain", domain).
		AddInfo("Author", message.Author.String()).
		AddInfo("Server", guildName).
		AddInfo("Channel", channelName).
		Show()
}

func (_Logging) MassDMBlocked() {
	UI.Notification("Protections: Mass DM Blocked!\nDMs will be enabled again in 5 minutes.")
	NewConsoleMessage().
		SetTitle("Protections").
		SetColor("#C1EFFF").
		AddInfo("", "Mass DM Blocked!").
		AddInfo("", "DMs will be enabled again in 5 minutes.").
		Show()
}

func (_Logging) FriendRequestSpamBlocked() {
	UI.Notification("Protections: Friend Request Spam Blocked!\nFriend Requests will be enabled again in 5 minutes.")
	NewConsoleMessage().
		SetTitle("Protections").
		SetColor("#C1EFFF").
		AddInfo("", "Friend Request Spam Blocked!").
		AddInfo("", "Friend Requests will be enabled again in 5 minutes.").
		Show()
}

func (_Logging) AntiDelete(session *discordgo.Session, message *discordgo.Message) {
	guildName := utils.Discord.GetGuildName(session, message.GuildID)
	channelName := utils.Discord.GetChannelName(session, message.ChannelID)
	UI.Notification("Anti Delete: " + message.Author.String() + " deleted a message in " + guildName)
	NewConsoleMessage().
		SetTitle("Protections").
		SetColor("#C1EFFF").
		AddInfo("", message.Author.String()+" deleted a message").
		AddInfo("Server", guildName).
		AddInfo("Channel", channelName).
		Show()
}

func (_Logging) GiveawayFound(session *discordgo.Session, message *discordgo.Message, prize, host, endsIn, account string) {
	guildName := utils.Discord.GetGuildName(session, message.GuildID)
	channelName := utils.Discord.GetChannelName(session, message.ChannelID)
	hostName := utils.Discord.GetUserName(session, host)
	if hostName == "Unknown" {
		hostName = host
	}
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
		AddInfoConditional("On Account", account, account != "").
		Show()
}

func (_Logging) GiveawayWon(session *discordgo.Session, message *discordgo.Message, prize, account string) {
	guildName := utils.Discord.GetGuildName(session, message.GuildID)
	channelName := utils.Discord.GetChannelName(session, message.ChannelID)
	UI.Notification("Giveaway Sniper: Won a giveaway in " + guildName + "\nPrize: " + prize)
	NewConsoleMessage().
		SetTitle("Giveaway Sniper").
		SetColor("#C1EFFF").
		AddInfo("", "Won Giveaway").
		AddInfo("Server", guildName).
		AddInfo("Channel", channelName).
		AddInfo("Bot", message.Author.String()).
		AddInfoConditional("Prize", prize, prize != "").
		AddInfoConditional("On Account", account, account != "").
		Show()
}

func (_Logging) MultiTokenSniperConnected(amount int) {
	info := fmt.Sprintf("%d Account Connected", amount)
	if amount > 1 {
		info = strings.ReplaceAll(info, "Account", "Accounts")
	}
	
	UI.Notification("Multi Token Sniper: " + info)
	NewConsoleMessage().
		SetTitle("Multi Token Sniper").
		SetColor("#C1EFFF").
		AddInfo("", info).
		Show()
}

func (_Logging) SessionConnected(id, os, client string) {
	cfg, err := config.GetConfig()
	if err != nil {
		return
	}

	if !cfg.Events.SessionDisconnected {
		return
	}

	UI.Notification("Session Connected\nOS: " + os + "\nClient: " + client)
	NewConsoleMessage().
		SetTitle("Session Connected").
		SetColor("#C1EFFF").
		AddInfo("Session ID", id).
		AddInfo("OS", os).
		AddInfo("Client", client).
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
		Show()
}

var Logging = _Logging{}
