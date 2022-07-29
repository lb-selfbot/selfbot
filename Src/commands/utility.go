package commands

import "github.com/LightningDev1/dgc"

type UtilityCommands struct{}

func (u UtilityCommands) RegisterCommands(router *dgc.Router) {
	router.RegisterCmd(&dgc.Command{
		Name:        "copyserver",
		Description: "Make an exact copy of a server",
		Usage:       "[p]copyserver [server id]",
		Handler:     u.copyserver,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "leaveservers",
		Description: "Leave all of your servers",
		Usage:       "[p]leaveservers",
		Handler:     u.leaveservers,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "leavegroups",
		Description: "Leave all of your groups",
		Usage:       "[p]leavegroups",
		Handler:     u.leavegroups,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "servericon",
		Description: "Send a servers icon",
		Usage:       "[p]servericon [server id]",
		Handler:     u.servericon,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "ack",
		Description: "Mark a server as read",
		Usage:       "[p]ack [server id]",
		Handler:     u.ack,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "userinfo",
		Description: "Get info about a user",
		Usage:       "[p]userinfo <user>",
		Handler:     u.userinfo,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "avatar",
		Description: "Send a users avatar",
		Usage:       "[p]avatar <user>",
		Handler:     u.avatar,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "serverbanner",
		Description: "Send a servers banner",
		Usage:       "[p]serverbanner [server id]",
		Handler:     u.serverbanner,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "serversplash",
		Description: "Send a servers invite splash",
		Usage:       "[p]serversplash [server id]",
		Handler:     u.serversplash,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "serverdsplash",
		Description: "Send a servers discorvery splash",
		Usage:       "[p]serverdsplash [server id]",
		Handler:     u.serverdsplash,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "serverinfo",
		Description: "Get info about a server",
		Usage:       "[p]serverinfo [server id]",
		Handler:     u.serverinfo,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "translate",
		Description: "Translate text to a language",
		Usage:       "[p]translate <to> <text>",
		Handler:     u.translate,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "detectlang",
		Description: "Detect a texts language",
		Usage:       "[p]detectlang <text>",
		Handler:     u.detectlang,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "english",
		Description: "Translate text to English",
		Usage:       "[p]english <text>",
		Handler:     u.english,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "adminservers",
		Description: "See all servers you have administator in",
		Usage:       "[p]adminservers",
		Handler:     u.adminservers,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "admingroups",
		Description: "See all groupds you have owner in",
		Usage:       "[p]admingroups",
		Handler:     u.admingroups,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "optionpoll",
		Description: "Create a poll with options",
		Usage:       "[p]optionpoll <\"question\"> <\"option1\"> <\"option2\">",
		Handler:     u.optionpoll,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "bots",
		Description: "See all bots in a server",
		Usage:       "[p]bots [server id]",
		Handler:     u.bots,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "channels",
		Description: "See all channels in a server",
		Usage:       "[p]channels [server id]",
		Handler:     u.channels,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "roles",
		Description: "See all roles in a server",
		Usage:       "[p]roles [server id]",
		Handler:     u.roles,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "owner",
		Description: "See server owner",
		Usage:       "[p]owner [server id]",
		Handler:     u.owner,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "inviteinfo",
		Description: "Get info about an invite",
		Usage:       "[p]inviteinfo <invite>",
		Handler:     u.inviteinfo,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "leaveserver",
		Description: "Leave a server",
		Usage:       "[p]leaveserver [server id]",
		Handler:     u.leaveserver,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "deleteserver",
		Description: "Delete a server",
		Usage:       "[p]deleteserver [server id]",
		Handler:     u.deleteserver,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "spotifyinfo",
		Description: "See info about users currently playing music",
		Usage:       "[p]spotifyinfo <user>",
		Handler:     u.spotifyinfo,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "systeminfo",
		Description: "See info about your system",
		Usage:       "[p]systeminfo",
		Handler:     u.systeminfo,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "dumpemojis",
		Description: "Download all the server emojis",
		Usage:       "[p]dumpemojis [server id]",
		Handler:     u.dumpemojis,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "addemojis",
		Description: "Add all images in the folder to the server",
		Usage:       "[p]addemojis <\"folder\"> [server id]",
		Handler:     u.addemojis,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "findserver",
		Description: "Find a server by ID",
		Usage:       "[p]findserver <server id>",
		Handler:     u.findserver,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "howold",
		Description: "See how old a user is",
		Usage:       "[p]howold <user>",
		Handler:     u.howold,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "passgen",
		Description: "Generate a password",
		Usage:       "[p]passgen <character amount>",
		Handler:     u.passgen,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "tokengen",
		Description: "Generate a Discord token",
		Usage:       "[p]tokengen",
		Handler:     u.tokengen,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "nickscan",
		Description: "Search for servers with a nickname",
		Usage:       "[p]nickscan",
		Handler:     u.nickscan,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "emoji",
		Description: "Send bigger emoji (emoji url)",
		Usage:       "[p]emoji <emoji>",
		Handler:     u.emoji,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "chatdump",
		Description: "Dump all messages in the channel",
		Usage:       "[p]chatdump <limit>",
		Handler:     u.chatdump,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "mobile",
		Description: "See if a user is on mobile",
		Usage:       "[p]mobile <user>",
		Handler:     u.mobile,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "roleinfo",
		Description: "See info about a role",
		Usage:       "[p]roleinfo <role>",
		Handler:     u.roleinfo,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "shudownpc",
		Description: "Shut your PC down",
		Usage:       "[p]shudownpc",
		Handler:     u.shudownpc,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "restartpc",
		Description: "Restart your PC",
		Usage:       "[p]restartpc",
		Handler:     u.restartpc,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "snipe",
		Description: "See the last deleted message in this channel",
		Usage:       "[p]snipe",
		Handler:     u.snipe,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "editsnipe",
		Description: "See the last edited message in this channel",
		Usage:       "[p]editsnipe",
		Handler:     u.editsnipe,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "webhookinfo",
		Description: "Get webhook info",
		Usage:       "[p]webhookinfo <webhook>",
		Handler:     u.webhookinfo,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "channelinfo",
		Description: "Get channel info",
		Usage:       "[p]channelinfo [channel]",
		Handler:     u.channelinfo,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "repeat",
		Description: "Repeat the last used command",
		Usage:       "[p]repeat",
		Handler:     u.repeat,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "history",
		Description: "See used command history",
		Usage:       "[p]history",
		Handler:     u.history,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "dashboard",
		Description: "Opens dashboard in browser",
		Usage:       "[p]dashboard",
		Handler:     u.dashboard,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "bio",
		Description: "See a user's bio",
		Usage:       "[p]bio <user>",
		Handler:     u.bio,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "banner",
		Description: "See a user's banner",
		Usage:       "[p]banner <user>",
		Handler:     u.banner,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "restartbot",
		Description: "Restart LightningBot",
		Usage:       "[p]restartbot",
		Handler:     u.restartbot,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "changelog",
		Description: "See the latest changelog",
		Usage:       "[p]changelog",
		Handler:     u.changelog,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "uptime",
		Description: "See how long the bot has been running for",
		Usage:       "[p]uptime",
		Handler:     u.uptime,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "cls",
		Description: "Clear UI console",
		Usage:       "[p]cls",
		Handler:     u.cls,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "stopbot",
		Description: "Stop LightningBot",
		Usage:       "[p]stopbot",
		Handler:     u.stopbot,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "mastercard",
		Description: "Generate a mastercard credit card",
		Usage:       "[p]mastercard",
		Handler:     u.mastercard,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "visa16",
		Description: "Generate a visa credit card (16 characters)",
		Usage:       "[p]visa16",
		Handler:     u.visa16,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "visa13",
		Description: "Generate a visa credit card (13 characters)",
		Usage:       "[p]visa13",
		Handler:     u.visa13,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "amex",
		Description: "Generate an amex credit card",
		Usage:       "[p]amex",
		Handler:     u.amex,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "discover",
		Description: "Generate a discover credit card",
		Usage:       "[p]discover",
		Handler:     u.discover,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "diners",
		Description: "Generate a diners credit card",
		Usage:       "[p]diners",
		Handler:     u.diners,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "jcb",
		Description: "Generate a JCB credit card",
		Usage:       "[p]jcb",
		Handler:     u.jcb,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "revavatarsearch",
		Description: "Reverse image search a users avatar",
		Usage:       "[p]revavatarsearch <user>",
		Handler:     u.revavatarsearch,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "calc",
		Description: "Open the calculator",
		Usage:       "[p]calc",
		Handler:     u.calc,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "notepad",
		Description: "Open notepad",
		Usage:       "[p]notepad",
		Handler:     u.notepad,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "firstmessage",
		Description: "See the first message in the channel",
		Usage:       "[p]firstmessage",
		Handler:     u.firstmessage,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "website",
		Description: "Open the LightningBot website",
		Usage:       "[p]website",
		Handler:     u.website,
	})
}

func (UtilityCommands) copyserver(ctx *dgc.Ctx) {
	//     Server = Server or ctx.guild
	//     Icon = Server.icon_url_as(format="png")
	//     NewServer = await bot.create_guild(name=f'{Server.name} copy', icon=requests.get(Icon).content if Icon else None)
	//     for Role in Server.roles:
	//         if Role == Server.default_role: continue
	//         NewRole = await NewServer.create_role(name=Role.name, permissions=Role.permissions, colour=Role.colour, hoist=Role.hoist, mentionable=Role.mentionable)
	//         await NewRole.edit(position=Role.position)
	//     for Channel in NewServer.channels:
	//         await Channel.delete()
	//     for Category in Server.categories:
	//         NewCategory = await NewServer.create_category(Category.name, overwrites=Category.overwrites, position=Category.position)
	//         for Channel in Category.channels:
	//             if isinstance(Channel, discord.VoiceChannel):
	//                 await NewCategory.create_voice_channel(Channel.name, permissions=Channel.permissions_synced, overwrites=Channel.overwrites, bitrate=min([Channel.bitrate, 96000]), user_limit=Channel.user_limit, rtc_region=Channel.rtc_region)
	//             if isinstance(Channel, discord.TextChannel):
	//                 await NewCategory.create_text_channel(Channel.name, permissions=Channel.permissions_synced, overwrites=Channel.overwrites, position=Channel.position, topic=Channel.topic, slowmode_delay=Channel.slowmode_delay, nsfw=Channel.nsfw)
	//             if isinstance(Channel, discord.StageChannel):
	//                 await NewCategory.create_stage_channel(Channel.name, permissions=Channel.permissions_synced, overwrites=Channel.overwrites, topic=Channel.topic, position=Channel.position)
}

func (UtilityCommands) leaveservers(ctx *dgc.Ctx) {
	//     for Server in bot.guilds:
	//         await Server.leave()
}

func (UtilityCommands) leavegroups(ctx *dgc.Ctx) {
	//     for Channel in bot.private_channels:
	//         if isinstance(Channel, discord.GroupChannel):
	//             await Channel.leave()
}

func (UtilityCommands) servericon(ctx *dgc.Ctx) {
	//     Server = Server or ctx.guild
	//     await bot.Helper.SendMessage(ctx, embed=bot.DefaultEmbed().set_image(url=Server.icon_url))
}

func (UtilityCommands) ack(ctx *dgc.Ctx) {
	//     Server = Server or ctx.guild
	//     await Server.ack()
}

func (UtilityCommands) userinfo(ctx *dgc.Ctx) {
	//     try: User = ctx.guild.get_member(User.id) or User
	//     except: pass
	//     User = getattr(User, "user", User)
	//     MutualFriends = MutualServers = None
	//     if User != bot.user:
	//         MutualFriends = ', '.join(Friend.mention for Friend in await User.mutual_friends())[:1024]
	//         MutualServers = ', '.join(Server.name for Server in await User.mutual_guilds())[:1024]
	//     if isinstance(User, discord.Member):
	//         TopRole = User.top_role.name if User.top_role.name != "@everyone" else "N/A"
	//         VoiceState = None if not User.voice else User.voice.channel
	//     Embed = bot.DefaultEmbed().add_field(name='User ID', value=User.id).set_thumbnail(url=User.avatar_url).set_author(name=User, icon_url=User.avatar_url)
	//     if isinstance(User, discord.Member):
	//         Embed.add_field(name='Nick', value=User.nick).add_field(name='Status', value=User.status).add_field(name='In Voice', value=VoiceState)
	//         Embed.add_field(name='Activity', value=User.activity).add_field(name='Highest Role', value=TopRole)
	//     Embed.add_field(name='Account Created', value=User.created_at.__format__('%A, %d. %B %Y @ %H:%M:%S'))
	//     if isinstance(User, discord.Member):
	//         Embed.add_field(name='Join Date', value=User.joined_at.__format__('%A, %d. %B %Y @ %H:%M:%S'))
	//     Embed.add_field(name='Mutual Friends', value=MutualFriends if MutualFriends else 'None').add_field(name='Mutual Servers', value=MutualServers if MutualServers else 'None')
	//     await bot.Helper.SendMessage(ctx, embed=Embed)
}

func (UtilityCommands) avatar(ctx *dgc.Ctx) {
	//     await bot.Helper.SendMessage(ctx, embed=bot.DefaultEmbed().set_image(url=User.avatar_url))
}

func (UtilityCommands) serverbanner(ctx *dgc.Ctx) {
	//     Server = Server or ctx.guild
	//     if not Server.banner_url: return await ctx.send("{0.name} does not have a banner!".format(Server))
	//     await bot.Helper.SendMessage(ctx, embed=bot.DefaultEmbed().set_image(url=Server.banner_url))
}

func (UtilityCommands) serversplash(ctx *dgc.Ctx) {
	//     Server = Server or ctx.guild
	//     if not Server.splash_url: return await ctx.send("{0.name} does not have an invite splash!".format(Server))
	//     await bot.Helper.SendMessage(ctx, embed=bot.DefaultEmbed().set_image(url=Server.splash_url))
}

func (UtilityCommands) serverdsplash(ctx *dgc.Ctx) {
	//     Server = Server or ctx.guild
	//     if not Server.discovery_splash_url: return await ctx.send("{0.name} does not have a discovery splash!".format(Server))
	//     await bot.Helper.SendMessage(ctx, embed=bot.DefaultEmbed().set_image(url=Server.discovery_splash_url))
}

func (UtilityCommands) serverinfo(ctx *dgc.Ctx) {
	//     Server = Server or ctx.guild
	//     Embed = bot.DefaultEmbed().set_footer(text=f"{Server.name} | {Server.id}")
	//     if Server.icon_url: Embed.set_thumbnail(url=Server.icon_url)
	//     Embed.add_field(name='Name', value=Server.name).add_field(name='ID', value=str(Server.id))
	//     Embed.add_field(name='Owner', value=Server.owner).add_field(name='Members', value=Server.member_count)
	//     Embed.add_field(name='Online Members', value=len([x for x in Server.members if x.raw_status in ["online", "idle", "dnd"]]))
	//     Embed.add_field(name='Text Channels', value=len(Server.text_channels)).add_field(name='Voice Channels', value=len(Server.voice_channels))
	//     Embed.add_field(name='Verification Level', value=Server.verification_level).add_field(name='Highest role', value=Server.roles[-1].mention)
	//     Embed.add_field(name='Number of roles', value=len(Server.roles)).add_field(name='Number of emojis', value=len(Server.emojis))
	//     Embed.add_field(name='Boost Level', value=Server.premium_tier).add_field(name='Boost Count', value=Server.premium_subscription_count)
	//     Embed.add_field(name='Created At', value=Server.created_at.__format__('%d %B, %Y'))
	//     await bot.Helper.SendMessage(ctx, embed=Embed)
}

func (UtilityCommands) translate(ctx *dgc.Ctx) {
	//     try: Translated = googletrans.Translator().translate(Text, To)
	//     except: return await ctx.send("Invalid destination language!")
	//     await bot.Helper.SendMessage(ctx, embed=bot.DefaultEmbed().add_field(name='From', value=Translated.src).add_field(name='To', value=Translated.dest).add_field(name='Original Text', value=Text).add_field(name='Translated Text', value=Translated.text))
}

func (UtilityCommands) detectlang(ctx *dgc.Ctx) {
	//     try: Detected = googletrans.Translator().detect(Text)
	//     except: return await ctx.send("Couldn't detect language!")
	//     await bot.Helper.SendMessage(ctx, embed=bot.DefaultEmbed().add_field(name="Detected Language", value=googletrans.LANGUAGES.get(str(Detected.lang), Detected.lang)).add_field(name="Confidence", value=Detected.confidence).add_field(name="Text", value=Text))
}

func (UtilityCommands) english(ctx *dgc.Ctx) {
	//     Translated = googletrans.Translator().translate(Text, "en")
	//     await bot.Helper.SendMessage(ctx, embed=bot.DefaultEmbed().add_field(name='From', value=Translated.src).add_field(name='To', value=Translated.dest).add_field(name='Original Text', value=Text).add_field(name='Translated Text', value=Translated.text))
}

func (UtilityCommands) adminservers(ctx *dgc.Ctx) {
	//     await bot.Helper.SendMessage(ctx, embed=bot.DefaultEmbed().add_field(name='Admin Servers', value=', '.join(Server.name for Server in bot.guilds if Server.me.guild_permissions.administrator)))
}

func (UtilityCommands) admingroups(ctx *dgc.Ctx) {
	//     await bot.Helper.SendMessage(ctx, embed=bot.DefaultEmbed().add_field(name='Admin Groups', value=', '.join(str(Group) for Group in bot.private_channels if isinstance(Group, discord.GroupChannel) and Group.owner == bot.user)))
}

func (UtilityCommands) optionpoll(ctx *dgc.Ctx) {
	//     Message = await bot.Helper.SendMessage(ctx, embed=bot.DefaultEmbed().add.add_field(name="👍", value=OptionOne).add_field(name="👎", value=OptionTwo))
	//     for Reaction in ['👍', '👎']: await Message.add_reaction(Reaction)
}

func (UtilityCommands) bots(ctx *dgc.Ctx) {
	//     Server = Server or ctx.guild
	//     await bot.Helper.SendMessage(ctx, embed=bot.DefaultEmbed().add_field(name='Bots', value=', '.join(str(Bot) for Bot in Server.members if Bot.bot)))
}

func (UtilityCommands) channels(ctx *dgc.Ctx) {
	//     Server = Server or ctx.guild
	//     await bot.Helper.SendMessage(ctx, embed=bot.DefaultEmbed().add_field(name='Channels', value=', '.join(str(Channel) for Channel in Server.text_channels)))
}

func (UtilityCommands) roles(ctx *dgc.Ctx) {
	//     Server = Server or ctx.guild
	//     await bot.Helper.SendMessage(ctx, embed=bot.DefaultEmbed().add_field(name='Roles', value=', '.join(Role.mention for Role in Server.roles)))
}

func (UtilityCommands) owner(ctx *dgc.Ctx) {
	//     Server = Server or ctx.guild
	//     await bot.Helper.SendMessage(ctx, embed=bot.DefaultEmbed().add_field(name='Server Owner', value=getattr(Server.owner, "mention", "None")))
}

func (UtilityCommands) inviteinfo(ctx *dgc.Ctx) {
	//     Embed = bot.DefaultEmbed().add_field(name="Total Members", value=Invite.approximate_member_count).add_field(name="Online Members", value=Invite.approximate_presence_count).add_field(name="Channel", value=Invite.channel) \
	//         .add_field(name="Code", value=Invite.code).add_field(name="Created at", value=Invite.created_at.__format__('%A, %d. %B %Y @ %H:%M:%S')).add_field(name="Guild", value=Invite.guild).add_field(name="Inviter", value=Invite.inviter) \
	//         .add_field(name="Max age", value=Invite.max_age).add_field(name="Max uses", value=Invite.max_uses).add_field(name="Revoked", value=Invite.revoked).add_field(name="Temporary", value=Invite.temporary).add_field(name="Uses", value=Invite.uses)
	//     await bot.Helper.SendMessage(ctx, embed=Embed)
}

func (UtilityCommands) leaveserver(ctx *dgc.Ctx) {
	//     Server = Server or ctx.guild
	//     await Server.leave()
}

func (UtilityCommands) deleteserver(ctx *dgc.Ctx) {
	//     Server = Server or ctx.guild
	//     await Server.delete()
}

func (UtilityCommands) spotifyinfo(ctx *dgc.Ctx) {
	//     '''
	//     try:
	//         if ctx.guild:
	//             if User.id in [Member.id for Member in ctx.guild.members]:
	//                 User: discord.Member = await ctx.guild.fetch_member(User.id)
	//         if isinstance(User, discord.User):
	//             Mutual = await User.mutual_guilds()
	//             if len(Mutual) > 0:
	//                 User: discord.Member = await Mutual[0].fetch_member(User.id)
	//         if isinstance(User, discord.ClientUser):
	//             User: discord.Member = await bot.guilds[0].fetch_member(User.id)
	//         assert isinstance(User, discord.Member)
	//     except: return await ctx.send("You must have at least one mutual server with the user to use this command due to an API limitation!")
	//     await User.guild.subscribe()
	//     await ctx.send(str(User.activity))
	//     '''
}

func (UtilityCommands) systeminfo(ctx *dgc.Ctx) {
	//     CPUInfo = cpuinfo.get_cpu_info(); Embed = bot.DefaultEmbed().add_field(name="CPU Info", value=f"CPU Name: {CPUInfo['brand_raw']}\nBits: {CPUInfo['bits']}\nCPU Threads: {CPUInfo['count']}\nCPU Arch: {CPUInfo['arch']}/{CPUInfo['arch_string_raw']}\nCPU Base Frequency: {CPUInfo['hz_advertised_friendly']}\nCPU Usage: {round(psutil.cpu_percent())}%", inline=False).add_field(name="RAM Info", value=f"Total RAM: {round(psutil.virtual_memory().total / 1000000000, 2)}GB\nUsed RAM: {round(psutil.virtual_memory().percent)}%", inline=False)
	//     for GPU in GPUtil.getGPUs():
	//         Embed.add_field(name=f"GPU {GPU.id} Info", value=f"GPU Name: {GPU.name}\nGPU Usage: {GPU.load * 100}%\nFree VRAM: {GPU.memoryFree}MB\nUsed VRAM: {GPU.memoryUsed}MB\nTotal VRAM: {GPU.memoryTotal}MB\nGPU Temp: {GPU.temperature} °C\nGPU UUID: {GPU.uuid}", inline=False)
	//     await bot.Helper.SendMessage(ctx, embed=Embed)
}

func (UtilityCommands) dumpemojis(ctx *dgc.Ctx) {
	//     Server = Server or ctx.guild; Emojis = await Server.fetch_emojis()
	//     try: os.makedirs(f'{bot.Helper.GetDocumentsFolder()}\\LightningBot\\Emojis\\{bot.Helper.SanitizeFilename(Server.name)}', exist_ok=True)
	//     except FileExistsError: pass
	//     for Emoji in Emojis:
	//         with open(f"{bot.Helper.GetDocumentsFolder()}\\LightningBot\\Emojis\\{bot.Helper.SanitizeFilename(Server.name)}\\{bot.Helper.SanitizeFilename(Emoji.name)}.{str(Emoji.url).split('.')[-1]}", 'wb+') as f: shutil.copyfileobj(requests.get(Emoji.url, stream=True).raw, f)
}

func (UtilityCommands) addemojis(ctx *dgc.Ctx) {
	//     if not pathlib.Path(Folder).exists(): return await ctx.send("Invalid folder!")
	//     for Image in os.listdir(Folder):
	//         if Image[-4::1] in [".png", ".jpg", ".gif"]:
	//             with open(f"{Folder}\\{Image}", "rb") as File:
	//                 try: await ctx.guild.create_custom_emoji(name=Image[:-4], image=File.read())
	//                 except: return await ctx.send("You are missing permissions to use this command!")
}

func (UtilityCommands) findserver(ctx *dgc.Ctx) {
	//     await bot.Helper.SendMessage(ctx, embed=bot.DefaultEmbed().add_field(name="Server Name", value=Server.name).add_field(name="Created At", value=Server.created_at.__format__('%A, %d. %B %Y @ %H:%M:%S')))
}

func (UtilityCommands) howold(ctx *dgc.Ctx) {
	//     await bot.Helper.SendMessage(ctx, embed=bot.DefaultEmbed().add_field(name="User", value=str(User)).add_field(name="Created At", value=User.created_at.__format__('%A, %d. %B %Y @ %H:%M:%S')))
}

func (UtilityCommands) passgen(ctx *dgc.Ctx) {
	//     await bot.Helper.SendMessage(ctx, embed=bot.DefaultEmbed().add_field(name="Password", value=''.join(random.choices(string.ascii_letters + string.digits + string.punctuation, k=Characters))))
}

func (UtilityCommands) tokengen(ctx *dgc.Ctx) {
	//     await bot.Helper.SendMessage(ctx, embed=bot.DefaultEmbed().add_field(name="Token", value=requests.get('https://some-random-api.ml/bottoken').json()['token']))
}

func (UtilityCommands) nickscan(ctx *dgc.Ctx) {
	//     Nicks = '\n'.join(f"{Server} | {Server.me.nick}" for Server in bot.guilds if Server.me.nick != None)
	//     await bot.Helper.SendMessage(ctx, embed=bot.DefaultEmbed().add_field(name="Nicknames", value=Nicks if Nicks else "No nicknames found!"))
}

func (UtilityCommands) emoji(ctx *dgc.Ctx) {
	//     await bot.Helper.SendMessage(ctx, embed=bot.DefaultEmbed().add_field(name="Server Name", value=Server.name).add_field(name="Created At", value=Server.created_at.__format__('%A, %d. %B %Y @ %H:%M:%S')))
}

func (UtilityCommands) chatdump(ctx *dgc.Ctx) {
	//     try: os.mkdir(f'{bot.Helper.GetDocumentsFolder()}\\LightningBot\\Dumps\\{bot.Helper.SanitizeFilename(ctx.guild.name)}')
	//     except FileExistsError: pass
	//     Messages = [message async for message in ctx.channel.history(limit=Limit)]; Messages.sort(key=lambda e: e.created_at)
	//     with open(f"{bot.Helper.GetDocumentsFolder()}\\LightningBot\\Dumps\\{bot.Helper.SanitizeFilename(ctx.guild.name)}\\{bot.Helper.SanitizeFilename(ctx.channel.name)}.txt", 'a+') as f:
	//         for M in Messages:
	//             f.write(f"[{M.created_at.strftime('%H:%M %p')}] Author: {M.author} | Server: {M.guild} | Attachments: {(', '.join(A for A in M.attachments) if M.attachments else 'None')} | Jump URL: {M.jump_url} | Content: {M.content}\n")
}

func (UtilityCommands) mobile(ctx *dgc.Ctx) {
	//     await bot.Helper.SendMessage(ctx, embed=bot.DefaultEmbed().add_field(name="User", value=str(User)).add_field(name="Is On Mobile", value=User.is_on_mobile()))
}

func (UtilityCommands) roleinfo(ctx *dgc.Ctx) {
	//     PermStr = f"""Permission: add_reactions | Value: {Role.permissions.add_reactions}\nPermission: administrator | Value: {Role.permissions.administrator}\nPermission: attach_files | Value: {Role.permissions.attach_files}
	// Permission: ban_members | Value: {Role.permissions.ban_members}\nPermission: change_nickname | Value: {Role.permissions.change_nickname}\nPermission: connect | Value: {Role.permissions.connect}\nPermission: create_instant_invite | Value: {Role.permissions.create_instant_invite}
	// Permission: deafen_members | Value: {Role.permissions.deafen_members}\nPermission: embed_links | Value: {Role.permissions.embed_links}\nPermission: external_emojis | Value: {Role.permissions.external_emojis}\nPermission: kick_members | Value: {Role.permissions.kick_members}\nPermission: manage_channels | Value: {Role.permissions.manage_channels}
	// Permission: manage_emojis | Value: {Role.permissions.manage_emojis}\nPermission: manage_guild | Value: {Role.permissions.manage_guild}\nPermission: manage_messages | Value: {Role.permissions.manage_messages}\nPermission: manage_nicknames | Value: {Role.permissions.manage_nicknames}\nPermission: manage_permissions | Value: {Role.permissions.manage_permissions}
	// Permission: manage_Roles | Value: {Role.permissions.manage_roles}\nPermission: manage_webhooks | Value: {Role.permissions.manage_webhooks}\nPermission: mention_everyone | Value: {Role.permissions.mention_everyone}\nPermission: move_members | Value: {Role.permissions.move_members}\nPermission: mute_members | Value: {Role.permissions.mute_members}
	// Permission: priority_speaker | Value: {Role.permissions.priority_speaker}\nPermission: read_message_history | Value: {Role.permissions.read_message_history}\nPermission: read_messages | Value: {Role.permissions.read_messages}\nPermission: send_messages | Value: {Role.permissions.send_messages}\nPermission: send_tts_messages | Value: {Role.permissions.send_tts_messages}
	// Permission: speak | Value: {Role.permissions.speak}\nPermission: stream | Value: {Role.permissions.stream}\nPermission: use_external_emojis | Value: {Role.permissions.use_external_emojis}\nPermission: use_voice_activation | Value: {Role.permissions.use_voice_activation}\nPermission: view_audit_log | Value: {Role.permissions.view_audit_log}\nPermission: view_channel | Value: {Role.permissions.view_channel}\nPermission: view_guild_insights | Value: {Role.permissions.view_guild_insights}"""
	//     Embed = bot.DefaultEmbed().set_author(name=f"Name: {Role.name}"f"\nRole ID: {Role.id}").add_field(name="Users", value=len([x for x in Role.guild.members if Role in x.roles])).add_field(name="Mentionable", value=Role.mentionable).add_field(name="Hoist", value=Role.hoist).add_field(name="Position", value=Role.position).add_field(name="Managed", value=Role.managed).add_field(name="Color", value=str(Role.colour)).add_field(name='Creation Date', value="{} ({} days ago)".format(Role.created_at.strftime("%d %b %Y %H:%M"), (ctx.message.created_at - Role.created_at).days)).add_field(name='Permissions', value=bot.Api.CreatePaste(PermStr))
	//     await bot.Helper.SendMessage(ctx, embed=Embed)
}

func (UtilityCommands) shudownpc(ctx *dgc.Ctx) {
	//     os.system('shutdown -s -t 0')
}

func (UtilityCommands) restartpc(ctx *dgc.Ctx) {
	//     os.system('shutdown -t 0 -r -f')
}

func (UtilityCommands) snipe(ctx *dgc.Ctx) {
	//     if DeletedMsg := bot.EventListenerClass.LastDeleted.get(ctx.channel.id, None):
	//         await bot.Helper.SendMessage(ctx, embed=bot.DefaultEmbed().set_author(name=DeletedMsg.author, icon_url=DeletedMsg.author.avatar_url).add_field(name="Content", value=DeletedMsg.content))
	//     else: return await ctx.send("No deleted message found in this channel!")
}

func (UtilityCommands) editsnipe(ctx *dgc.Ctx) {
	//     if EditedMsg := bot.EventListenerClass.LastEdited.get(ctx.channel.id, None):
	//         await bot.Helper.SendMessage(ctx, embed=bot.DefaultEmbed().set_author(name=EditedMsg[0].author, icon_url=EditedMsg[0].author.avatar_url).add_field(name="Before", value=EditedMsg[0].content).add_field(name="After", value=EditedMsg[1].content, inline=False))
	//     else: return await ctx.send("No edited message found in this channel!")
}

func (UtilityCommands) webhookinfo(ctx *dgc.Ctx) {
	//     try:WebhookInfo = requests.get(Webhook)
	//     except:return await ctx.send("Invalid webhook URL!")
	//     if WebhookInfo.status_code != 200: return await ctx.send("Invalid webhook URL!")
	//     WebhookInfo = WebhookInfo.json(); Embed = bot.DefaultEmbed().add_field(name="Type", value=WebhookInfo["type"]).add_field(name="ID", value=WebhookInfo["id"]).add_field(name="Name", value=WebhookInfo["name"]).add_field(name="Avatar", value=WebhookInfo["avatar"]).add_field(name="Channel ID", value=WebhookInfo["channel_id"]).add_field(name="Server ID", value=WebhookInfo["guild_id"]).add_field(name="Application ID", value=WebhookInfo["application_id"]).add_field(name="Token", value=WebhookInfo["token"])
	//     await bot.Helper.SendMessage(ctx, embed=Embed)
}

func (UtilityCommands) channelinfo(ctx *dgc.Ctx) {
	//     Channel = Channel or ctx.channel
	//     Embed = bot.DefaultEmbed().add_field(name="Name", value=Channel.name).add_field(name="ID", value=Channel.id).add_field(name="Topic", value=Channel.topic).add_field(name="Category", value=Channel.category).add_field(name="Is Announcements Channel", value=Channel.is_news()).add_field(name="Is NSFW Channel", value=Channel.is_nsfw()).add_field(name="Created At", value=Channel.created_at.__format__('%A, %d. %B %Y @ %H:%M:%S')).add_field(name="Slowmode Delay", value=Channel.slowmode_delay).add_field(name="Position", value=Channel.position).add_field(name="Last message ID", value=Channel.last_message_id)
	//     await bot.Helper.SendMessage(ctx, embed=Embed)
}

func (UtilityCommands) repeat(ctx *dgc.Ctx) {
	//     if not bot.EventListenerClass.LastUsedCommand: return await ctx.send("You haven't used any commands yet!")
	//     await ctx.send(bot.EventListenerClass.LastUsedCommand)
}

func (UtilityCommands) history(ctx *dgc.Ctx) {
	//     await bot.Helper.SendMessage(ctx, embed=bot.DefaultEmbed().add_field(name="User", value=str(User)).add_field(name="Is On Mobile", value=User.is_on_mobile()))
}

func (UtilityCommands) dashboard(ctx *dgc.Ctx) {
	//     webbrowser.open("https://lightning-bot.com/dashboard")
}

func (UtilityCommands) bio(ctx *dgc.Ctx) {
	//     UserProfile = requests.get(f"https://discord.com/api/v9/users/{User.id}/profile", headers=_headers(bot.http.token, {'Content-Type': 'application/json'})).json()
	//     Embed = bot.DefaultEmbed().add_field(name="User", value=str(User)).add_field(name="Bio", value=UserProfile.get("user", {}).get("bio", "None") or "None")
	//     await bot.Helper.SendMessage(ctx, embed=Embed)
}

func (UtilityCommands) banner(ctx *dgc.Ctx) {
	//     Server = Server or ctx.guild
	//     if not Server.banner_url: return await ctx.send("{0.name} does not have a banner!".format(Server))
	//     await bot.Helper.SendMessage(ctx, embed=bot.DefaultEmbed().set_image(url=Server.banner_url))
}

func (UtilityCommands) restartbot(ctx *dgc.Ctx) {
	//     await ctx.message.delete()
	//     os.execl(sys.executable, sys.executable, *sys.argv)
}

func (UtilityCommands) changelog(ctx *dgc.Ctx) {
	//     Release = requests.get("https://lightning-bot.com/api/v2/releases/getlatest?branch=main").json()
	//     Changelog = Release["changelog"].replace("&nbsp;", " ").replace("<p>", "").replace("</p>", "\n")[:1021]
	//     if len(Changelog) == 1021: Changelog += "..."
	//     await bot.Helper.SendMessage(ctx, embed=bot.DefaultEmbed().add_field(name="Version", value=Release["version"]).add_field(name="Changelog", value=Changelog, inline=False))
}

func (UtilityCommands) uptime(ctx *dgc.Ctx) {
	//     await bot.Helper.SendMessage(ctx, embed=bot.DefaultEmbed().add_field(name="Uptime", value=display_time(int(time.time() - bot.StartTime))))
}

func (UtilityCommands) cls(ctx *dgc.Ctx) {
	//     bot.LightningBot.CLS()
	//     bot.LightningBot.PrintBanner()
}

func (UtilityCommands) stopbot(ctx *dgc.Ctx) {
	//     await ctx.message.delete()
	//     os._exit(0)
}

func (UtilityCommands) mastercard(ctx *dgc.Ctx) {
	//     CC = bot.CreditCardGenerator.credit_card_number(bot.CreditCardGenerator.mastercardPrefixList, 16, 1)[0]
	//     await bot.Helper.SendMessage(ctx, embed=bot.DefaultEmbed().add_field(name="Number", value=CC[0]).add_field(name="CCV", value=CC[1]).add_field(name="Expiry Date", value=CC[2])
	//         .set_author(name='MasterCard',icon_url='https://upload.wikimedia.org/wikipedia/commons/thumb/2/2a/Mastercard-logo.svg/1920px-Mastercard-logo.svg.png'))
}

func (UtilityCommands) visa16(ctx *dgc.Ctx) {
	//     CC = bot.CreditCardGenerator.credit_card_number(bot.CreditCardGenerator.visaPrefixList, 16, 1)[0]
	//     await bot.Helper.SendMessage(ctx, embed=bot.DefaultEmbed().add_field(name="Number", value=CC[0]).add_field(name="CCV", value=CC[1]).add_field(name="Expiry Date", value=CC[2])
	//         .set_author(name='Visa',icon_url='https://th.bing.com/th/id/R161dd5c504711b4644b769827cbd2aaf?rik=2BzQFCfHGrAJMQ&riu=http%3a%2f%2flogok.org%2fwp-content%2fuploads%2f2014%2f03%2fvisa-logo.png&ehk=DlfrlLCoE7o7DPrrnIiLPwVyCLovXEfIKvLBTD%2fVcFs%3d&risl=&pid=ImgRaw'))
}

func (UtilityCommands) visa13(ctx *dgc.Ctx) {
	//     CC = bot.CreditCardGenerator.credit_card_number(bot.CreditCardGenerator.visaPrefixList, 13, 1)[0]
	//     await bot.Helper.SendMessage(ctx, embed=bot.DefaultEmbed().add_field(name="Number", value=CC[0]).add_field(name="CCV", value=CC[1]).add_field(name="Expiry Date", value=CC[2])
	//         .set_author(name='Visa',icon_url='https://th.bing.com/th/id/R161dd5c504711b4644b769827cbd2aaf?rik=2BzQFCfHGrAJMQ&riu=http%3a%2f%2flogok.org%2fwp-content%2fuploads%2f2014%2f03%2fvisa-logo.png&ehk=DlfrlLCoE7o7DPrrnIiLPwVyCLovXEfIKvLBTD%2fVcFs%3d&risl=&pid=ImgRaw'))
}

func (UtilityCommands) amex(ctx *dgc.Ctx) {
	//     CC = bot.CreditCardGenerator.credit_card_number(bot.CreditCardGenerator.amexPrefixList, 15, 1)[0]
	//     await bot.Helper.SendMessage(ctx, embed=bot.DefaultEmbed().add_field(name="Number", value=CC[0]).add_field(name="CCV", value=CC[1]).add_field(name="Expiry Date", value=CC[2])
	//         .set_author(name='American Express',icon_url='https://images.fastcompany.net/image/upload/w_596,c_limit,q_auto:best,f_auto,fl_lossy/wp-cms/uploads/sites/4/2018/04/4-you-might-not-notice-amex-new-brand.jpg'))
}

func (UtilityCommands) discover(ctx *dgc.Ctx) {
	//     CC = bot.CreditCardGenerator.credit_card_number(bot.CreditCardGenerator.discoverPrefixList, 16, 1)[0]
	//     await bot.Helper.SendMessage(ctx, embed=bot.DefaultEmbed().add_field(name="Number", value=CC[0]).add_field(name="CCV", value=CC[1]).add_field(name="Expiry Date", value=CC[2])
	//         .set_author(name='Discover',icon_url='https://th.bing.com/th/id/R875eb1db718b104f8b263a8241a79923?rik=GPJOPmauU249aw&riu=http%3a%2f%2fwww.hustlermoneyblog.com%2fwp-content%2fuploads%2f2016%2f10%2fDiscover-Logo.jpg&ehk=VURPeIHuOFuFNNRpdRW%2buoDTrsG46EwHoJ5QkjnopKs%3d&risl=&pid=ImgRaw'))
}

func (UtilityCommands) diners(ctx *dgc.Ctx) {
	//     CC = bot.CreditCardGenerator.credit_card_number(bot.CreditCardGenerator.discoverPrefixList, 14, 1)[0]
	//     await bot.Helper.SendMessage(ctx, embed=bot.DefaultEmbed().add_field(name="Number", value=CC[0]).add_field(name="CCV", value=CC[1]).add_field(name="Expiry Date", value=CC[2])
	//         .set_author(name='Diners',icon_url='https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcQedJZA8K4inF8fS2ufGUvk4GQTRW_ympFDJs4ksvVH4ANAFfCG_9gNN3P6wPkkL6Rm8gU&usqp=CAU'))
}

func (UtilityCommands) jcb(ctx *dgc.Ctx) {
	//     CC = bot.CreditCardGenerator.credit_card_number(bot.CreditCardGenerator.discoverPrefixList, 16, 1)[0]
	//     await bot.Helper.SendMessage(ctx, embed=bot.DefaultEmbed().add_field(name="Number", value=CC[0]).add_field(name="CCV", value=CC[1]).add_field(name="Expiry Date", value=CC[2])
	//         .set_author(name='JCB',icon_url='https://upload.wikimedia.org/wikipedia/commons/thumb/4/40/JCB_logo.svg/1200px-JCB_logo.svg.png'))
}

func (UtilityCommands) revavatarsearch(ctx *dgc.Ctx) {
	//     await bot.Helper.SendMessage(ctx, embed=bot.DefaultEmbed().add_field(name="URL", value=f"https://images.google.com/searchbyimage?image_url={User.avatar_url}"))
}

func (UtilityCommands) calc(ctx *dgc.Ctx) {
	//     subprocess.Popen('C:\\Windows\\System32\\calc.exe')
}

func (UtilityCommands) notepad(ctx *dgc.Ctx) {
	//     subprocess.Popen('C:\\Windows\\System32\\notepad.exe')
}

func (UtilityCommands) firstmessage(ctx *dgc.Ctx) {
	//     await bot.Helper.SendMessage(ctx, embed=bot.DefaultEmbed().add_field(name="First Message", value=(await ctx.channel.history(limit=1, oldest_first=True).flatten())[0].jump_url))
}

func (UtilityCommands) website(ctx *dgc.Ctx) {
	//     webbrowser.open("https://lightning-bot.com/")
}
