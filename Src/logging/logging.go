package logging

import (
	"container/list"
	"fmt"
	"strings"
	"time"

	"github.com/LightningDev1/LB-Selfbot/config"
	"github.com/LightningDev1/LB-Selfbot/constants"
	"github.com/LightningDev1/LB-Selfbot/ipc"
	"github.com/LightningDev1/LB-Selfbot/utils"
	"github.com/LightningDev1/LB-Selfbot/websocket"
	"github.com/LightningDev1/dgc"
	"github.com/LightningDev1/discordgo"
	toast "github.com/electricbubble/go-toast"
)

func Debug(v ...any) {
	t := time.Now().Local().Format("2006/01/02 15:04:05")
	fmt.Printf("[\x1b[34;1m%s\x1b[0m] [\x1b[35;1mDEBUG\x1b[0m] ", t)
	fmt.Println(v...)
}

func Info(v ...any) {
	t := time.Now().Local().Format("2006/01/02 15:04:05")
	fmt.Printf("[\x1b[34;1m%s\x1b[0m] [\x1b[32;1mINFO\x1b[0m] ", t)
	fmt.Println(v...)
}

func Warn(v ...any) {
	t := time.Now().Local().Format("2006-01-02 15:04:05")
	fmt.Printf("[\x1b[34;1m%s\x1b[0m] [\x1b[33;1mWARNING\x1b[0m] ", t)
	fmt.Println(v...)
}

func Error(v ...any) {
	t := time.Now().Local().Format("2006-01-02 15:04:05")
	fmt.Printf("[\x1b[34;1m%s\x1b[0m] [\x1b[31;1mERROR\x1b[0m] ", t)
	fmt.Println(v...)
}

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
		Op: "CONSOLE_CLEAR",
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
		go toast.Push(message,
			toast.WithTitle("LightningBot"),
			toast.WithAppID("LightningBot"),
			toast.WithIcon(utils.File.GetMainFolder()+"/Logo.png"),
		)
	}
}

func (UILogging) Colored(text string, color string) string {
	return fmt.Sprintf("<text Color=\"%s\">%s</text>", utils.Misc.Or(color, "#ffffff"), text)
}

func (UILogging) Centered(text string, color string) string {
	return fmt.Sprintf("<text Color=\"%s\" Centered>%s</text>", utils.Misc.Or(color, "#ffffff"), text)
}

func (UILogging) ConsoleMessage(title string, color string, v map[string]string) {
	timeFormat := UI.Colored(time.Now().Local().Format("15:04:05"), "#C1EFFF")
	message := UI.Colored("[", "#eeeeee") + timeFormat + UI.Colored("] ", "#eeeeee") + UI.Colored(title, color) + "\n"
	for key, value := range v {
		if key != "" && key != "_" {
			message += UI.Colored("  "+key+": ", "#eeeeee") + UI.Colored(value, "#C1EFFF") + "\n"
		} else {
			message += UI.Colored("  "+value, "#C1EFFF") + "\n"
		}
	}
	ConsoleAppend(message)
}

func (UILogging) CommandUsed(commandName string) {
	UI.Notification("Command " + commandName + " used")
	UI.ConsoleMessage("Command Used", "#C1EFFF", map[string]string{"Command": commandName})
}

func (UILogging) OnReady(router *dgc.Router, event *discordgo.Ready) {
	cfg, _ := config.GetConfig()
	friends := 0
	for _, rel := range event.Relationships {
		if rel.Type == 1 {
			friends++
		}
	}
	ConsoleClear()
	UI.Notification("Connected to " + event.User.String())
	message := "<image Source=\"assets/Banner.png\" Width=\"650\" Height=\"150\" Centered />\n" +
		UI.Centered("LightningBot "+constants.VersionString+" | "+event.User.String()+" | Prefix: "+cfg.Selfbot.CommandPrefix, "#C1EFFF") + "\n" +
		UI.Centered(fmt.Sprint(len(router.Commands))+" Commands | "+fmt.Sprint(len(event.Guilds))+" Servers | "+fmt.Sprint(friends)+" Friends", "#C1EFFF") + "\n"
	ConsoleAppend(message)
}

func (UILogging) NitroSniperFake(session *discordgo.Session, message *discordgo.Message, code, account string) {
	guildName := getGuildName(session, message.GuildID)
	channelName := getChannelName(session, message.ChannelID)
	UI.Notification(fmt.Sprintf("Nitro Sniper: Fake Code\nServer: %s\nChannel: %s\nAuthor: %s\nCode: %s", guildName, channelName, message.Author.String(), code))
	v := map[string]string{
		"Status":  "Fake Code",
		"Author":  message.Author.String(),
		"Server":  guildName,
		"Channel": channelName,
		"Code":    code,
	}
	if account != "" {
		v["On Account"] = account
	}
	UI.ConsoleMessage("Nitro Sniper", "#C1EFFF", v)
}

func (UILogging) NitroSniper(session *discordgo.Session, message *discordgo.Message, elapsed time.Duration, status, code, account string) {
	guildName := getGuildName(session, message.GuildID)
	channelName := getChannelName(session, message.ChannelID)
	UI.Notification(fmt.Sprintf("Nitro Sniper: %s\nServer: %s\nChannel: %s\nAuthor: %s\nCode: %s", status, guildName, channelName, message.Author.String(), code))
	v := map[string]string{
		"Status":     status,
		"Author":     message.Author.String(),
		"Server":     guildName,
		"Channel":    channelName,
		"Code":       code,
		"Snipe Time": elapsed.String(),
	}
	if account != "" {
		v["On Account"] = account
	}
	UI.ConsoleMessage("Nitro Sniper", "#C1EFFF", v)
}

func (UILogging) BanNotifier(session *discordgo.Session, guildID string) {
	guildName := getGuildName(session, guildID)
	UI.Notification("Ban Notifier: You have been banned from " + guildName)
	UI.ConsoleMessage("Ban Notifier", "#C1EFFF", map[string]string{"": "You have been banned from " + guildName})
}

func (UILogging) DMNotifier(message *discordgo.Message) {
	UI.Notification("DM Notifier: New DM from " + message.Author.String())
	UI.ConsoleMessage("DM Notifier", "#C1EFFF", map[string]string{"": "New DM from " + message.Author.String()})
}

func (UILogging) DMEditNotifier(message *discordgo.Message) {
	UI.Notification("DM Edit Notifier: Edited DM from " + message.Author.String())
	UI.ConsoleMessage("DM Edit Notifier", "#C1EFFF", map[string]string{"": "Edited DM from " + message.Author.String()})
}

func (UILogging) DMDeleteNotifier(message *discordgo.Message) {
	UI.Notification("DM Delete Notifier: Deleted DM from " + message.Author.String())
	UI.ConsoleMessage("DM Delete Notifier", "#C1EFFF", map[string]string{"": "Deleted DM from " + message.Author.String()})
}

func (UILogging) DMTypingNotifier(userName string) {
	UI.Notification("DM Typing Notifier: " + userName + " is typing")
	UI.ConsoleMessage("DM Typing Notifier", "#C1EFFF", map[string]string{"": userName + " is typing"})
}

func (UILogging) MessageDeleteNotifier(session *discordgo.Session, message *discordgo.Message) {
	guildName := getGuildName(session, message.GuildID)
	channelName := getChannelName(session, message.ChannelID)
	UI.Notification("Message Delete Notifier: Deleted message from " + message.Author.String())
	UI.ConsoleMessage("Message Delete Notifier", "#C1EFFF", map[string]string{
		"":        "Deleted message from " + message.Author.String(),
		"Content": message.Content,
		"Server":  guildName,
		"Channel": channelName,
	})
}

func (UILogging) DiscordStaffDetection(session *discordgo.Session, message *discordgo.Message) {
	guildName := getGuildName(session, message.GuildID)
	channelName := getChannelName(session, message.ChannelID)
	UI.Notification("Discord Staff Detection: " + message.Author.String() + " is Discord staff")
	UI.ConsoleMessage("Discord Staff Detection", "#C1EFFF", map[string]string{
		"":        message.Author.String() + " is Discord staff",
		"Server":  guildName,
		"Channel": channelName,
	})
}

func (UILogging) GhostPingDetection(session *discordgo.Session, message *discordgo.Message) {
	guildName := getGuildName(session, message.GuildID)
	channelName := getChannelName(session, message.ChannelID)
	UI.Notification("Ghostping Detection: " + message.Author.String() + " has pinged you")
	UI.ConsoleMessage("Ghostping Detection", "#C1EFFF", map[string]string{
		"":        message.Author.String() + " has pinged you",
		"Server":  guildName,
		"Channel": channelName,
	})
}

func (UILogging) SelfbotDetection(session *discordgo.Session, message *discordgo.Message, selfbot string) {
	guildName := getGuildName(session, message.GuildID)
	channelName := getChannelName(session, message.ChannelID)
	UI.Notification("Selfbot Detection: " + message.Author.String() + " was detected using a selfbot")
	UI.ConsoleMessage("Selfbot Detection", "#C1EFFF", map[string]string{
		"":        message.Author.String() + " was detected using a selfbot",
		"Server":  guildName,
		"Channel": channelName,
		"Selfbot": selfbot,
	})
}

func (UILogging) TicketNotifier(session *discordgo.Session, channel *discordgo.Channel) {
	guildName := getGuildName(session, channel.GuildID)
	channelName := getChannelName(session, channel.ID)
	UI.Notification("Ticket Notifier: New ticket " + channelName + " was created in " + guildName)
	UI.ConsoleMessage("Ticket Notifier", "#C1EFFF", map[string]string{
		"":        "New ticket was created",
		"Server":  guildName,
		"Channel": channelName,
		"Ticket":  channelName,
	})
}

func (UILogging) PhisingDetected(session *discordgo.Session, message *discordgo.Message, domain string) {
	guildName := getGuildName(session, message.GuildID)
	channelName := getChannelName(session, message.ChannelID)
	UI.Notification(fmt.Sprintf("Protections: Phising Detected!\nDomain: %s\nAuthor: %s", domain, message.Author.String()))
	UI.ConsoleMessage("Protections", "#C1EFFF", map[string]string{
		"":        "Phising Detected!",
		"Domain":  domain,
		"Author":  message.Author.String(),
		"Server":  guildName,
		"Channel": channelName,
	})
}

func (UILogging) MassDMBlocked() {
	UI.Notification("Protections: Mass DM Blocked!\nDMs will be enabled again in 5 minutes.")
	UI.ConsoleMessage("Protections", "#C1EFFF", map[string]string{
		"":  "Mass DM Blocked!",
		"_": "DMs will be enabled again in 5 minutes.",
	})
}

func (UILogging) FriendRequestSpamBlocked() {
	UI.Notification("Protections: Friend Request Spam Blocked!\nFriend Requests will be enabled again in 5 minutes.")
	UI.ConsoleMessage("Protections", "#C1EFFF", map[string]string{
		"":  "Friend Request Spam Blocked!",
		"_": "Friend Requests will be enabled again in 5 minutes.",
	})
}

func (UILogging) AntiDelete(session *discordgo.Session, message *discordgo.Message) {
	guildName := getGuildName(session, message.GuildID)
	channelName := getChannelName(session, message.ChannelID)
	UI.Notification("Anti Delete: " + message.Author.String() + " deleted a message in " + guildName)
	UI.ConsoleMessage("Protections", "#C1EFFF", map[string]string{
		"":        message.Author.String() + " deleted a message.",
		"Server":  guildName,
		"Channel": channelName,
	})
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
