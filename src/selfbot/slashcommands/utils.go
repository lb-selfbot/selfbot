package slashcommands

import (
	"crypto/md5"
	"errors"
	"fmt"
	"strconv"

	"github.com/LightningDev1/LB-Selfbot/selfbot/utils"
	"github.com/LightningDev1/dgc"
	"github.com/LightningDev1/discordgo"
)

var errInvalidSubCommand = errors.New("invalid subcommand")

func optionError(name string, err error) error {
	return errors.New("Error in option " + name + ": " + err.Error())
}

func ArgumentsToCommandOptions(options []*discordgo.ApplicationCommandOption, args *dgc.Arguments, offset int) ([]*discordgo.ApplicationCommandOption, error) {
	if len(options) == 0 {
		// No need to do anything because there are no options
		return options, nil
	}

	// Subcommand handling
	isSubCommand := options[0].Type == discordgo.ApplicationCommandOptionSubCommand ||
		options[0].Type == discordgo.ApplicationCommandOptionSubCommandGroup

	if isSubCommand {
		subCommands := utils.Filter(options, func(option *discordgo.ApplicationCommandOption) bool {
			return option.Name == args.Get(0+offset).Raw()
		})

		if len(subCommands) == 0 {
			return nil, errInvalidSubCommand
		}

		subCommand := subCommands[0]

		subOptions, err := ArgumentsToCommandOptions(subCommand.Options, args, offset+1)
		if err != nil {
			return nil, err
		}

		return []*discordgo.ApplicationCommandOption{
			{
				Type:    subCommand.Type,
				Name:    subCommand.Name,
				Options: subOptions,
			},
		}, nil
	}

	// Normal command option handling
	for i, option := range options {
		arg := args.Get(i + offset)

		switch option.Type {
		case discordgo.ApplicationCommandOptionString:
			option.Value = arg.Raw()
		case discordgo.ApplicationCommandOptionInteger, discordgo.ApplicationCommandOptionNumber:
			asInt, err := arg.AsInt()
			if err != nil {
				return nil, optionError(option.Name, err)
			}

			option.Value = asInt
		case discordgo.ApplicationCommandOptionBoolean:
			asBool, err := arg.AsBool()
			if err != nil {
				return nil, optionError(option.Name, err)
			}

			option.Value = asBool
		case discordgo.ApplicationCommandOptionUser, discordgo.ApplicationCommandOptionChannel, discordgo.ApplicationCommandOptionRole, discordgo.ApplicationCommandOptionMentionable:
			option.Value = arg.AsUserMentionID()
		}
	}

	return options, nil
}

func CommandIdentifier(j *Job) string {
	return j.GuildObject.ID + j.ApplicationObject.ID + j.CommandObject.ID
}

func JobIdentifier(j *Job) string {
	id := []byte(j.GuildObject.ID +
		j.ApplicationObject.ID +
		j.ChannelObject.ID +
		j.CommandObject.ID +
		j.Name +
		j.RawArguments +
		strconv.Itoa(j.IntervalSec))

	return fmt.Sprintf("%x", md5.Sum(id))
}
