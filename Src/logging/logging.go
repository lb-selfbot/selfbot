package logging

import (
	"container/list"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/LightningDev1/LB-Selfbot/config"
	"github.com/LightningDev1/LB-Selfbot/constants"
	"github.com/LightningDev1/LB-Selfbot/ipc"
	"github.com/LightningDev1/LB-Selfbot/utils"
	"github.com/LightningDev1/LB-Selfbot/websocket"
	"github.com/LightningDev1/dgc"
	"github.com/LightningDev1/discordgo"
	"github.com/common-nighthawk/go-figure"
	"github.com/gen2brain/beeep"
)

// UI Logging

func getGuildName(session *discordgo.Session, guildID string) string {
	guild, _ := utils.Discord.GetGuild(session, guildID)
	guildName := "No Server"
	if guild != nil {
		guildName = guild.Name
	}
	return guildName
}

func getChannelName(session *discordgo.Session, channelID string) string {
	channel, _ := utils.Discord.GetChannel(session, channelID)
	channelName := "No Channel"
	if channel != nil {
		channelName = channel.Name
	}
	return channelName
}

func getUserName(session *discordgo.Session, userID string) string {
	arg := &dgc.Argument{}
	userName, _ := utils.Discord.TryGetUsernameFromMention(session, arg.SetRaw(userID))
	return utils.Misc.Or(userName, "Unknown")
}

func ConsoleAppend(text string) {
	ipc.IPC.ConsoleAppend(text)
	websocket.Send(websocket.WsMessage{
		Op: "CONSOLE_APPEND",
		Data: websocket.WsData{
			Text: text,
		},
	})
}

func ConsoleClear() {
	ipc.IPC.ConsoleClear()
	websocket.Send(websocket.WsMessage{
		Op:   "CONSOLE_CLEAR",
		Data: websocket.WsData{},
	})
}

type ConsoleMessage struct {
	Title string
	Color string
	Info  *list.List
}

func (c *ConsoleMessage) SetTitle(title string) *ConsoleMessage {
	c.Title = title
	return c
}

func (c *ConsoleMessage) SetColor(color string) *ConsoleMessage {
	c.Color = color
	return c
}

func (c *ConsoleMessage) AddInfo(key, value string) *ConsoleMessage {
	c.Info.PushBack([]string{key, value})
	return c
}

func (c *ConsoleMessage) AddInfoConditional(key, value string, condition bool) *ConsoleMessage {
	if condition {
		c.AddInfo(key, value)
	}
	return c
}

func (c *ConsoleMessage) Show() {
	if constants.Console {
		r, g, b := utils.Misc.HexToRGB(c.Color)
		timeFormat := UI.ColoredConsole(time.Now().Local().Format("15:04:05"), 193, 239, 255, false)
		message := UI.ColoredConsole("[", 238, 238, 238, false) + timeFormat + UI.ColoredConsole("] ", 238, 238, 238, false) + UI.ColoredConsole(c.Title, r, g, b, false) + "\n"

		for e := c.Info.Front(); e != nil; e = e.Next() {
			info := e.Value.([]string)
			if info[0] != "" {
				message += UI.ColoredConsole("  "+info[0]+": ", 238, 238, 238, false) + UI.ColoredConsole(info[1], 193, 239, 255, false) + "\n"
			} else {
				message += UI.ColoredConsole("  "+info[1], 193, 239, 255, false) + "\n"
			}
		}

		fmt.Println(message)

	}

	timeFormat := UI.Colored(time.Now().Local().Format("15:04:05"), "#C1EFFF")
	message := UI.Colored("[", "#eeeeee") + timeFormat + UI.Colored("] ", "#eeeeee") + UI.Colored(c.Title, c.Color) + "\n"

	for e := c.Info.Front(); e != nil; e = e.Next() {
		info := e.Value.([]string)
		if info[0] != "" {
			message += UI.Colored("  "+info[0]+": ", "#eeeeee") + UI.Colored(info[1], "#C1EFFF") + "\n"
		} else {
			message += UI.Colored("  "+info[1], "#C1EFFF") + "\n"
		}
	}

	ConsoleAppend(message)
}

func NewConsoleMessage() *ConsoleMessage {
	return &ConsoleMessage{
		Info: list.New(),
	}
}

type UILogging struct{}

func (UILogging) Notification(message string) {
	cfg, _ := config.GetConfig()
	if cfg.Notifications.Enabled {
		go beeep.Notify("LightningBot", message, filepath.Join(utils.File.GetMainFolder(), "Logo.png"))
	}
}

func (UILogging) Colored(text string, color string) string {
	return fmt.Sprintf("<text Color=\"%s\">%s</text>", utils.Misc.Or(color, "#ffffff"), text)
}

func (UILogging) ColoredConsole(text string, r, g, b int, bold bool) string {
	boldArg := ";1"
	if !bold {
		boldArg = ""
	}
	return fmt.Sprintf("\033[38;2;%d;%d;%d%sm%s\033[0m", r, g, b, boldArg, text)
}

func (UILogging) Centered(text string, color string) string {
	return fmt.Sprintf("<text Color=\"%s\" Centered>%s</text>", utils.Misc.Or(color, "#ffffff"), text)
}

func (UILogging) CommandUsed(commandName string) {
	UI.Notification("Command " + commandName + " used")
	NewConsoleMessage().
		SetTitle("Command Used").
		SetColor("#C1EFFF").
		AddInfo("Command", commandName).
		Show()
}

func (UILogging) OnReady(router *dgc.Router, event *discordgo.Ready) {
	cfg, _ := config.GetConfig()
	friends := 0
	for _, rel := range event.Relationships {
		if rel.Type == 1 {
			friends++
		}
	}
	UI.Notification("Connected to " + event.User.String())

	bannerTextUpper := "LightningBot " + constants.VersionString + " | " + event.User.String() + " | Prefix: " + cfg.Selfbot.CommandPrefix
	bannerTextLower := fmt.Sprint(len(router.Commands)) + " Commands | " + fmt.Sprint(len(event.Guilds)) + " Servers | " + fmt.Sprint(friends) + " Friends"

	if constants.Console {
		// Clear console
		fmt.Print("\033[H\033[2J")

		banner := utils.String.CenterSlices(figure.NewFigure("LightningBot", "stop", true).Slicify())

		message := UI.ColoredConsole(banner, 0, 203, 255, true) + "\n\n"

		message += UI.ColoredConsole(
			utils.String.CenterSlices([]string{bannerTextUpper})+"\n"+utils.String.CenterSlices([]string{bannerTextLower}), 0, 203, 255, false,
		)

		fmt.Println(message)

	}

	message := "<image Source=\"assets/Banner.png\" Width=\"650\" Height=\"150\" Centered />\n" +
		UI.Centered("LightningBot "+constants.VersionString+" | "+event.User.String()+" | Prefix: "+cfg.Selfbot.CommandPrefix, "#C1EFFF") + "\n" +
		UI.Centered(fmt.Sprint(len(router.Commands))+" Commands | "+fmt.Sprint(len(event.Guilds))+" Servers | "+fmt.Sprint(friends)+" Friends", "#C1EFFF") + "\n"

	ConsoleClear()
	ConsoleAppend(message)
}

func (UILogging) NitroSniperFake(session *discordgo.Session, message *discordgo.Message, code, account string) {
	guildName := getGuildName(session, message.GuildID)
	channelName := getChannelName(session, message.ChannelID)
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

func (UILogging) NitroSniper(session *discordgo.Session, message *discordgo.Message, elapsed time.Duration, status, code, account string) {
	guildName := getGuildName(session, message.GuildID)
	channelName := getChannelName(session, message.ChannelID)
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

func (UILogging) BanNotifier(session *discordgo.Session, guildID string) {
	guildName := getGuildName(session, guildID)
	UI.Notification("Ban Notifier: You have been banned from " + guildName)
	NewConsoleMessage().
		SetTitle("Ban Notifier").
		SetColor("#C1EFFF").
		AddInfo("", "You have been banned from "+guildName).
		Show()
}

func (UILogging) DMNotifier(message *discordgo.Message) {
	UI.Notification("DM Notifier: New DM from " + message.Author.String())
	NewConsoleMessage().
		SetTitle("DM Notifier").
		SetColor("#C1EFFF").
		AddInfo("", "New DM from "+message.Author.String()).
		Show()
}

func (UILogging) DMEditNotifier(message *discordgo.Message) {
	UI.Notification("DM Edit Notifier: Edited DM from " + message.Author.String())
	NewConsoleMessage().
		SetTitle("DM Edit Notifier").
		SetColor("#C1EFFF").
		AddInfo("", "Edited DM from "+message.Author.String()).
		Show()
}

func (UILogging) DMDeleteNotifier(message *discordgo.Message) {
	UI.Notification("DM Delete Notifier: Deleted DM from " + message.Author.String())
	NewConsoleMessage().
		SetTitle("DM Delete Notifier").
		SetColor("#C1EFFF").
		AddInfo("", "Deleted DM from "+message.Author.String()).
		AddInfo("Content", message.Content).
		Show()
}

func (UILogging) DMTypingNotifier(userName string) {
	UI.Notification("DM Typing Notifier: " + userName + " is typing")
	NewConsoleMessage().
		SetTitle("DM Typing Notifier").
		SetColor("#C1EFFF").
		AddInfo("", userName+" is typing").
		Show()
}

func (UILogging) MessageDeleteNotifier(session *discordgo.Session, message *discordgo.Message) {
	guildName := getGuildName(session, message.GuildID)
	channelName := getChannelName(session, message.ChannelID)
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

func (UILogging) DiscordStaffDetection(session *discordgo.Session, message *discordgo.Message) {
	guildName := getGuildName(session, message.GuildID)
	channelName := getChannelName(session, message.ChannelID)
	UI.Notification("Discord Staff Detection: " + message.Author.String() + " is Discord staff")
	NewConsoleMessage().
		SetTitle("Discord Staff Detection").
		SetColor("#C1EFFF").
		AddInfo("", message.Author.String()+" is Discord staff").
		AddInfo("Server", guildName).
		AddInfo("Channel", channelName).
		Show()
}

func (UILogging) GhostPingDetection(session *discordgo.Session, message *discordgo.Message) {
	guildName := getGuildName(session, message.GuildID)
	channelName := getChannelName(session, message.ChannelID)
	UI.Notification("Ghostping Detection: " + message.Author.String() + " has pinged you")
	NewConsoleMessage().
		SetTitle("Ghostping Detection").
		SetColor("#C1EFFF").
		AddInfo("", message.Author.String()+" has pinged you").
		AddInfo("Server", guildName).
		AddInfo("Channel", channelName).
		Show()
}

func (UILogging) SelfbotDetection(session *discordgo.Session, message *discordgo.Message, selfbot string) {
	guildName := getGuildName(session, message.GuildID)
	channelName := getChannelName(session, message.ChannelID)
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

func (UILogging) TicketNotifier(session *discordgo.Session, channel *discordgo.Channel) {
	guildName := getGuildName(session, channel.GuildID)
	channelName := getChannelName(session, channel.ID)
	NewConsoleMessage().
		SetTitle("Ticket Notifier").
		SetColor("#C1EFFF").
		AddInfo("", "New ticket was created").
		AddInfo("Server", guildName).
		AddInfo("Channel", channelName).
		AddInfo("Ticket", channelName).
		Show()
}

func (UILogging) PhisingDetected(session *discordgo.Session, message *discordgo.Message, domain string) {
	guildName := getGuildName(session, message.GuildID)
	channelName := getChannelName(session, message.ChannelID)
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

func (UILogging) MassDMBlocked() {
	UI.Notification("Protections: Mass DM Blocked!\nDMs will be enabled again in 5 minutes.")
	NewConsoleMessage().
		SetTitle("Protections").
		SetColor("#C1EFFF").
		AddInfo("", "Mass DM Blocked!").
		AddInfo("", "DMs will be enabled again in 5 minutes.").
		Show()
}

func (UILogging) FriendRequestSpamBlocked() {
	UI.Notification("Protections: Friend Request Spam Blocked!\nFriend Requests will be enabled again in 5 minutes.")
	NewConsoleMessage().
		SetTitle("Protections").
		SetColor("#C1EFFF").
		AddInfo("", "Friend Request Spam Blocked!").
		AddInfo("", "Friend Requests will be enabled again in 5 minutes.").
		Show()
}

func (UILogging) AntiDelete(session *discordgo.Session, message *discordgo.Message) {
	guildName := getGuildName(session, message.GuildID)
	channelName := getChannelName(session, message.ChannelID)
	UI.Notification("Anti Delete: " + message.Author.String() + " deleted a message in " + guildName)
	NewConsoleMessage().
		SetTitle("Protections").
		SetColor("#C1EFFF").
		AddInfo("", message.Author.String()+" deleted a message").
		AddInfo("Server", guildName).
		AddInfo("Channel", channelName).
		Show()
}

func (UILogging) GiveawayFound(session *discordgo.Session, message *discordgo.Message, prize, host, endsIn, account string) {
	guildName := getGuildName(session, message.GuildID)
	channelName := getChannelName(session, message.ChannelID)
	hostName := getUserName(session, host)
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

func (UILogging) GiveawayWon(session *discordgo.Session, message *discordgo.Message, prize, account string) {
	guildName := getGuildName(session, message.GuildID)
	channelName := getChannelName(session, message.ChannelID)
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

func (UILogging) MultiTokenSniperConnected(amount int) {
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

func (UILogging) SessionConnected(id, os, client string) {
	UI.Notification("Session Connected\nOS: " + os + "\nClient: " + client)
	NewConsoleMessage().
		SetTitle("Session Connected").
		SetColor("#C1EFFF").
		AddInfo("Session ID", id).
		AddInfo("OS", os).
		AddInfo("Client", client).
		Show()
}

func (UILogging) SessionDisconnected(id, os, client string) {
	UI.Notification("Session Disconnected\nOS: " + os + "\nClient: " + client)
	NewConsoleMessage().
		SetTitle("Session Disconnected").
		SetColor("#C1EFFF").
		AddInfo("Session ID", id).
		AddInfo("OS", os).
		AddInfo("Client", client).
		Show()
}

var UI = UILogging{}
