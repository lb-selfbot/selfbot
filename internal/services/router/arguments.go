package router

import (
	"errors"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/lb-selfbot/discordgo"
)

var (
	// ErrArgumentEmpty is returned when an argument is empty.
	ErrArgumentEmpty = errors.New("argument is empty")

	// ErrNotInGuild is returned when the command is not executed in a guild.
	ErrNotInGuild = errors.New("command is not executed in a guild")

	// ErrMemberNotFound is returned when a member is not found.
	ErrMemberNotFound = errors.New("member not found")

	// ErrRoleNotFound is returned when a role is not found.
	ErrRoleNotFound = errors.New("role not found")

	// ErrChannelNotFound is returned when a channel is not found.
	ErrChannelNotFound = errors.New("channel not found")
)

type ArgumentFlag int

const (
	// ArgumentFlagFallbackDefault returns the default value if the argument is not found.
	ArgumentFlagFallbackDefault ArgumentFlag = iota
)

var (
	// regexArguments defines the regex the argument string must match.
	regexArguments = regexp.MustCompile(`[^\s"']+|"([^"\\]|\\.)*"|'([^'\\]|\\.)*'`)

	// regexUserMention defines the regex a user mention must match.
	regexUserMention = regexp.MustCompile(`<@!?(\d+)>`)

	// regexRoleMention defines the regex a role mention must match.
	regexRoleMention = regexp.MustCompile(`<@&(\d+)>`)

	// regexChannelMention defines the regex a channel mention must match.
	regexChannelMention = regexp.MustCompile(`<#(\d+)>`)
)

// Arguments represents command arguments.
type Arguments struct {
	ctx *Ctx

	raw       string
	arguments []*Argument
}

func ParseArguments(ctx *Ctx, raw string) *Arguments {
	matches := regexArguments.FindAllString(raw, -1)
	arguments := make([]*Argument, 0, len(matches))

	for _, match := range matches {
		if match == "" {
			continue
		}

		// Remove quotes from the start and end of the string
		first, last := match[0], match[len(match)-1]
		if (first == '"' && last == '"') || (first == '\'' && last == '\'') {
			match = match[1 : len(match)-1]
		}

		// Remove escape characters
		match = strings.ReplaceAll(match, `\"`, `"`)
		match = strings.ReplaceAll(match, `\'`, `'`)
		match = strings.ReplaceAll(match, `\\`, `\`)

		arguments = append(arguments, &Argument{
			ctx: ctx,
			raw: match,
		})
	}

	return &Arguments{
		ctx:       ctx,
		raw:       raw,
		arguments: arguments,
	}
}

// Raw returns the raw arguments.
func (a *Arguments) Raw() string {
	return a.raw
}

// RawFrom returns the raw arguments from the given index.
func (a *Arguments) RawFrom(index int) string {
	if a.Amount() <= index {
		return ""
	}

	args := a.arguments[index:]

	var sb strings.Builder
	for _, arg := range args {
		sb.WriteString(arg.raw)
		sb.WriteString(" ")
	}

	return strings.TrimSpace(sb.String())
}

// Amount returns the amount of arguments.
func (a *Arguments) Amount() int {
	return len(a.arguments)
}

// Get returns the argument at the given index.
func (a *Arguments) Get(index int) *Argument {
	if a.Amount() <= index {
		return &Argument{ctx: a.ctx, raw: ""}
	}
	return a.arguments[index]
}

// Argument represents a command argument.
type Argument struct {
	ctx *Ctx

	raw string
}

// Raw returns the raw argument.
func (a *Argument) Raw() string {
	return a.raw
}

// Int returns the argument as an int.
func (a *Argument) Int() (int, error) {
	return strconv.Atoi(a.raw)
}

// Bool returns the argument as a bool.
func (a *Argument) Bool() (bool, error) {
	switch arg := strings.ToLower(a.raw); arg {
	case "on":
		return true, nil
	case "off":
		return false, nil
	default:
		return strconv.ParseBool(arg)
	}
}

// MustBool returns the argument as a bool or panics if it fails.
func (a *Argument) MustBool() bool {
	b, err := a.Bool()
	if err != nil {
		panic(err)
	}

	return b
}

// MustInt returns the argument as an int or panics if it fails.
func (a *Argument) MustInt() int {
	i, err := a.Int()
	if err != nil {
		panic(err)
	}

	return i
}

// Int64 returns the argument as an int64.
func (a *Argument) Int64() (int64, error) {
	return strconv.ParseInt(a.raw, 10, 64)
}

// MustInt64 returns the argument as an int64 or panics if it fails.
func (a *Argument) MustInt64() int64 {
	i, err := a.Int64()
	if err != nil {
		panic(err)
	}

	return i
}

// Float64 returns the argument as a float64.
func (a *Argument) Float64() (float64, error) {
	return strconv.ParseFloat(a.raw, 64)
}

// MustFloat64 returns the argument as a float64 or panics if it fails.
func (a *Argument) MustFloat64() float64 {
	f, err := a.Float64()
	if err != nil {
		panic(err)
	}

	return f
}

// Duration returns the argument as a duration.
func (a *Argument) Duration() (time.Duration, error) {
	d, err := time.ParseDuration(a.raw)
	if err == nil {
		return d, nil
	}

	// Fallback to seconds
	return time.ParseDuration(a.raw + "s")
}

// MustDuration returns the argument as a duration or panics if it fails.
func (a *Argument) MustDuration() time.Duration {
	d, err := a.Duration()
	if err != nil {
		panic(err)
	}

	return d
}

// UserMentionID returns the ID of the user mentioned in the argument
// or the raw argument if no user is mentioned.
func (a *Argument) UserMentionID() string {
	if matches := regexUserMention.FindStringSubmatch(a.raw); len(matches) > 1 {
		return matches[1]
	}

	return a.raw
}

// RoleMentionID returns the ID of the role mentioned in the argument
// or the raw argument if no role is mentioned.
func (a *Argument) RoleMentionID() string {
	if matches := regexRoleMention.FindStringSubmatch(a.raw); len(matches) > 1 {
		return matches[1]
	}

	return a.raw
}

// ChannelMentionID returns the ID of the channel mentioned in the argument
// or the raw argument if no channel is mentioned.
func (a *Argument) ChannelMentionID() string {
	if matches := regexChannelMention.FindStringSubmatch(a.raw); len(matches) > 1 {
		return matches[1]
	}

	return a.raw
}

// User returns the argument as a user.
func (a *Argument) User(flags ...ArgumentFlag) (*discordgo.User, error) {
	userID := a.UserMentionID()
	if userID == "" {
		if slices.Contains(flags, ArgumentFlagFallbackDefault) {
			// Default to author instead of session user as the
			// command could have been run by a shared user
			userID = a.ctx.Message.Author.ID
		} else {
			return nil, ErrArgumentEmpty
		}
	}

	// Try to get the user from the state
	if user, err := a.ctx.Session.State.GetUser(userID); err == nil {
		return user, nil
	}

	// Fallback to API
	return a.ctx.Session.User(userID)
}

// MustUser returns the argument as a user or panics if it fails.
func (a *Argument) MustUser(flags ...ArgumentFlag) *discordgo.User {
	user, err := a.User(flags...)
	if err != nil {
		panic(err)
	}

	return user
}

// Member returns the argument as a guild member.
func (a *Argument) Member(flags ...ArgumentFlag) (*discordgo.Member, error) {
	if a.ctx.Message.GuildID == "" {
		return nil, ErrNotInGuild
	}

	userID := a.UserMentionID()
	if userID == "" {
		if slices.Contains(flags, ArgumentFlagFallbackDefault) {
			// Default to author instead of session user as the
			// command could have been run by a shared user
			userID = a.ctx.Message.Author.ID
		} else {
			return nil, ErrArgumentEmpty
		}
	}

	// Try to get the member from the state
	if member, err := a.ctx.Session.State.Member(a.ctx.Message.GuildID, userID); err == nil {
		return member, nil
	}

	// Fallback to QueryGuildMembers gateway operation
	userIDInt, err := strconv.Atoi(userID)
	if err != nil {
		return nil, err
	}

	queryParams := discordgo.NewQueryGuildMembersParams(a.ctx.Message.GuildID)
	queryParams.UserIDs = []int{userIDInt}

	members, err := a.ctx.Session.QueryGuildMembers(queryParams)
	if err != nil {
		return nil, err
	}

	if len(members) == 0 {
		return nil, nil
	}

	return members[0], nil
}

// MustMember returns the argument as a guild member or panics if it fails.
func (a *Argument) MustMember(flags ...ArgumentFlag) *discordgo.Member {
	member, err := a.Member(flags...)
	if err != nil {
		panic(err)
	}

	return member
}

// Role returns the argument as a role.
func (a *Argument) Role(flags ...ArgumentFlag) (*discordgo.Role, error) {
	if a.ctx.Message.GuildID == "" {
		return nil, ErrNotInGuild
	}

	roleID := a.RoleMentionID()
	if roleID == "" {
		if slices.Contains(flags, ArgumentFlagFallbackDefault) {
			// roleID = guildID will default to @everyone
			roleID = a.ctx.Message.GuildID
		} else {
			return nil, ErrArgumentEmpty
		}
	}

	// Try to get the role from the state
	if role, err := a.ctx.Session.State.Role(a.ctx.Message.GuildID, roleID); err == nil {
		return role, nil
	}

	// Fallback to API
	roles, err := a.ctx.Session.GuildRoles(a.ctx.Message.GuildID)
	if err != nil {
		return nil, err
	}

	for _, role := range roles {
		if role.ID == roleID {
			return role, nil
		}
	}

	return nil, ErrRoleNotFound
}

// MustRole returns the argument as a role or panics if it fails.
func (a *Argument) MustRole(flags ...ArgumentFlag) *discordgo.Role {
	role, err := a.Role(flags...)
	if err != nil {
		panic(err)
	}

	return role
}

// Channel returns the argument as a channel.
func (a *Argument) Channel(flags ...ArgumentFlag) (*discordgo.Channel, error) {
	channelID := a.ChannelMentionID()
	if channelID == "" {
		if slices.Contains(flags, ArgumentFlagFallbackDefault) {
			channelID = a.ctx.Message.ChannelID
		} else {
			return nil, ErrArgumentEmpty
		}
	}

	// Try to get the channel from the state
	if channel, err := a.ctx.Session.State.Channel(channelID); err == nil {
		return channel, nil
	}

	// Fallback to API
	return a.ctx.Session.Channel(channelID)
}

// MustChannel returns the argument as a channel or panics if it fails.
func (a *Argument) MustChannel(flags ...ArgumentFlag) *discordgo.Channel {
	channel, err := a.Channel(flags...)
	if err != nil {
		panic(err)
	}

	return channel
}

// Guild returns the argument as a guild.
func (a *Argument) Guild(flags ...ArgumentFlag) (*discordgo.Guild, error) {
	guildID := a.raw
	if guildID == "" {
		if slices.Contains(flags, ArgumentFlagFallbackDefault) && a.ctx.Message.GuildID != "" {
			guildID = a.ctx.Message.GuildID
		} else {
			return nil, ErrArgumentEmpty
		}
	}

	// Try to get the guild from the state
	if guild, err := a.ctx.Session.State.Guild(guildID); err == nil {
		return guild, nil
	}

	// Fallback to API
	return a.ctx.Session.Guild(guildID)
}

// MustGuild returns the argument as a guild or panics if it fails.
func (a *Argument) MustGuild(flags ...ArgumentFlag) *discordgo.Guild {
	guild, err := a.Guild(flags...)
	if err != nil {
		panic(err)
	}

	return guild
}
