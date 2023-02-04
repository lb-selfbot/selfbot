package commands

import (
	"strconv"
	"strings"

	"github.com/LightningDev1/LB-Selfbot/conversion"
	"github.com/LightningDev1/LB-Selfbot/embed"
	"github.com/LightningDev1/LB-Selfbot/utils"
	"github.com/LightningDev1/dgc"
)

type ConversionCommands struct{}

func (c ConversionCommands) RegisterCommands(router *dgc.Router) {
	router.RegisterCmd(&dgc.Command{
		Name:        "convertvelocity",
		Description: "Convert velocity from one unit to another",
		Usage:       "[p]convertvelocity <value> <source unit> <destination unit>",
		Handler:     c.convertvelocity,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "convertarea",
		Description: "Convert area from one unit to another",
		Usage:       "[p]convertarea <value> <source unit> <destination unit>",
		Handler:     c.convertarea,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "convertlength",
		Description: "Convert length from one unit to another",
		Usage:       "[p]convertlength <value> <source unit> <destination unit>",
		Handler:     c.convertlength,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "converttemperature",
		Description: "Convert temperature from one unit to another",
		Usage:       "[p]converttemperature <value> <source unit> <destination unit>",
		Handler:     c.converttemperature,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "converttime",
		Description: "Convert time from one unit to another",
		Usage:       "[p]converttime <value> <source unit> <destination unit>",
		Handler:     c.converttime,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "convertvolume",
		Description: "Convert volume from one unit to another",
		Usage:       "[p]convertvolume <value> <source unit> <destination unit>",
		Handler:     c.convertvolume,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "convertweight",
		Description: "Convert weight from one unit to another",
		Usage:       "[p]convertweight <value> <source unit> <destination unit>",
		Handler:     c.convertweight,
	})
}

func (ConversionCommands) checkIfUnitValid(ctx *dgc.Ctx, unit string, validUnits []string) bool {
	isValid := utils.SliceContains(validUnits, strings.ToLower(unit))
	if isValid {
		return true
	}

	embed.NewEmbed().
		AddField("Error", "Invalid unit").
		AddField("Valid units", strings.Join(validUnits, ", ")).
		Send(ctx)

	return false
}

func (c ConversionCommands) convertvelocity(ctx *dgc.Ctx) {
	value, err := strconv.ParseFloat(utils.Arguments["<value>"], 64)
	if err != nil {
		ctx.RespondText("Invalid value")
		return
	}

	sourceUnit := utils.Arguments["<source unit>"]
	destinationUnit := utils.Arguments["<destination unit>"]

	if !c.checkIfUnitValid(ctx, sourceUnit, conversion.SpeedUnits) || !c.checkIfUnitValid(ctx, destinationUnit, conversion.SpeedUnits) {
		return
	}

	result, err := conversion.ConvertSpeed(sourceUnit, destinationUnit, value)
	if err != nil {
		ctx.RespondText("Error converting value: " + err.Error())
	}

	embed.NewEmbed().
		AddField("Value", strconv.FormatFloat(value, 'g', -1, 64)+" "+sourceUnit).
		AddField("Converted", strconv.FormatFloat(result, 'g', -1, 64)+" "+destinationUnit).
		Send(ctx)
}

func (c ConversionCommands) convertarea(ctx *dgc.Ctx) {
	value, err := strconv.ParseFloat(utils.Arguments["<value>"], 64)
	if err != nil {
		ctx.RespondText("Invalid value")
		return
	}

	sourceUnit := utils.Arguments["<source unit>"]
	destinationUnit := utils.Arguments["<destination unit>"]

	if !c.checkIfUnitValid(ctx, sourceUnit, conversion.AreaUnits) || !c.checkIfUnitValid(ctx, destinationUnit, conversion.AreaUnits) {
		return
	}

	result, err := conversion.ConvertArea(sourceUnit, destinationUnit, value)
	if err != nil {
		ctx.RespondText("Error converting value: " + err.Error())
	}

	embed.NewEmbed().
		AddField("Value", strconv.FormatFloat(value, 'g', -1, 64)+" "+sourceUnit).
		AddField("Converted", strconv.FormatFloat(result, 'g', -1, 64)+" "+destinationUnit).
		Send(ctx)
}

func (c ConversionCommands) convertlength(ctx *dgc.Ctx) {
	value, err := strconv.ParseFloat(utils.Arguments["<value>"], 64)
	if err != nil {
		ctx.RespondText("Invalid value")
		return
	}

	sourceUnit := utils.Arguments["<source unit>"]
	destinationUnit := utils.Arguments["<destination unit>"]

	if !c.checkIfUnitValid(ctx, sourceUnit, conversion.LengthUnits) || !c.checkIfUnitValid(ctx, destinationUnit, conversion.LengthUnits) {
		return
	}

	result, err := conversion.ConvertLength(sourceUnit, destinationUnit, value)
	if err != nil {
		ctx.RespondText("Error converting value: " + err.Error())
	}

	embed.NewEmbed().
		AddField("Value", strconv.FormatFloat(value, 'g', -1, 64)+" "+sourceUnit).
		AddField("Converted", strconv.FormatFloat(result, 'g', -1, 64)+" "+destinationUnit).
		Send(ctx)
}

func (c ConversionCommands) converttemperature(ctx *dgc.Ctx) {
	value, err := strconv.ParseFloat(utils.Arguments["<value>"], 64)
	if err != nil {
		ctx.RespondText("Invalid value")
		return
	}

	sourceUnit := utils.Arguments["<source unit>"]
	destinationUnit := utils.Arguments["<destination unit>"]

	if !c.checkIfUnitValid(ctx, sourceUnit, conversion.TemperatureUnits) || !c.checkIfUnitValid(ctx, destinationUnit, conversion.TemperatureUnits) {
		return
	}

	result, err := conversion.ConvertTemperature(sourceUnit, destinationUnit, value)
	if err != nil {
		ctx.RespondText("Error converting value: " + err.Error())
	}

	embed.NewEmbed().
		AddField("Value", strconv.FormatFloat(value, 'g', -1, 64)+" "+sourceUnit).
		AddField("Converted", strconv.FormatFloat(result, 'g', -1, 64)+" "+destinationUnit).
		Send(ctx)
}

func (c ConversionCommands) converttime(ctx *dgc.Ctx) {
	value, err := strconv.ParseFloat(utils.Arguments["<value>"], 64)
	if err != nil {
		ctx.RespondText("Invalid value")
		return
	}

	sourceUnit := utils.Arguments["<source unit>"]
	destinationUnit := utils.Arguments["<destination unit>"]

	if !c.checkIfUnitValid(ctx, sourceUnit, conversion.TimeUnits) || !c.checkIfUnitValid(ctx, destinationUnit, conversion.TimeUnits) {
		return
	}

	result, err := conversion.ConvertTime(sourceUnit, destinationUnit, value)
	if err != nil {
		ctx.RespondText("Error converting value: " + err.Error())
	}

	embed.NewEmbed().
		AddField("Value", strconv.FormatFloat(value, 'g', -1, 64)+" "+sourceUnit).
		AddField("Converted", strconv.FormatFloat(result, 'g', -1, 64)+" "+destinationUnit).
		Send(ctx)
}

func (c ConversionCommands) convertvolume(ctx *dgc.Ctx) {
	value, err := strconv.ParseFloat(utils.Arguments["<value>"], 64)
	if err != nil {
		ctx.RespondText("Invalid value")
		return
	}

	sourceUnit := utils.Arguments["<source unit>"]
	destinationUnit := utils.Arguments["<destination unit>"]

	if !c.checkIfUnitValid(ctx, sourceUnit, conversion.VolumeUnits) || !c.checkIfUnitValid(ctx, destinationUnit, conversion.VolumeUnits) {
		return
	}

	result, err := conversion.ConvertVolume(sourceUnit, destinationUnit, value)
	if err != nil {
		ctx.RespondText("Error converting value: " + err.Error())
	}

	embed.NewEmbed().
		AddField("Value", strconv.FormatFloat(value, 'g', -1, 64)+" "+sourceUnit).
		AddField("Converted", strconv.FormatFloat(result, 'g', -1, 64)+" "+destinationUnit).
		Send(ctx)
}

func (c ConversionCommands) convertweight(ctx *dgc.Ctx) {
	value, err := strconv.ParseFloat(utils.Arguments["<value>"], 64)
	if err != nil {
		ctx.RespondText("Invalid value")
		return
	}

	sourceUnit := utils.Arguments["<source unit>"]
	destinationUnit := utils.Arguments["<destination unit>"]

	if !c.checkIfUnitValid(ctx, sourceUnit, conversion.WeightUnits) || !c.checkIfUnitValid(ctx, destinationUnit, conversion.WeightUnits) {
		return
	}

	result, err := conversion.ConvertWeight(sourceUnit, destinationUnit, value)
	if err != nil {
		ctx.RespondText("Error converting value: " + err.Error())
	}

	embed.NewEmbed().
		AddField("Value", strconv.FormatFloat(value, 'g', -1, 64)+" "+sourceUnit).
		AddField("Converted", strconv.FormatFloat(result, 'g', -1, 64)+" "+destinationUnit).
		Send(ctx)
}
