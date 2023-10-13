package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func SpellCommand(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
	options := interaction.ApplicationCommandData().Options
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
	for _, opt := range options {
		optionMap[opt.Name] = opt
	}

	if interaction.Type == discordgo.InteractionApplicationCommandAutocomplete {
		champion, ok := optionMap["champion"]
		if ok && champion.Focused {
			autoCompleteChampionName(session, interaction, champion.StringValue())
			return
		}

		spell, ok := optionMap["spell"]
		if ok && spell.Focused {
			autoCompleteSpell(session, interaction, optionMap["champion"].StringValue())
		}

		return
	}

	championKey := optionMap["champion"].StringValue()

	spell := strings.Split(optionMap["spell"].StringValue(), ",")
	spellIndex, err := strconv.Atoi(spell[1])
	if err != nil {
		client.logger.Error(fmt.Sprintf("error converting spell index to int: %v", err))
	}

	embed := &discordgo.MessageEmbed{}
	components := []discordgo.MessageComponent{discordgo.ActionsRow{Components: []discordgo.MessageComponent{
		discordgo.Button{Label: "Modifiers", CustomID: fmt.Sprintf("modifiers_%v_%v_%v", championKey, spell[0], spellIndex)},
		discordgo.Button{Label: "Notes", CustomID: fmt.Sprintf("notes_%v_%v_%v", championKey, spell[0], spellIndex)},
	}}}

	embedType, ok := optionMap["type"]
	if !ok || embedType.StringValue() == "" {
		embed = &spellsEmbeds[championKey][spell[0]][spellIndex].General
	} else {
		switch embedType.StringValue() {
		case "modifiers":
			embed = &spellsEmbeds[championKey][spell[0]][spellIndex].Modifiers
			components = []discordgo.MessageComponent{}
		case "notes":
			embed = &spellsEmbeds[championKey][spell[0]][spellIndex].Notes
			components = []discordgo.MessageComponent{}
		default:
			embed = &spellsEmbeds[championKey][spell[0]][spellIndex].General
		}
	}

	err = client.session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds:     []*discordgo.MessageEmbed{embed},
			Components: components,
		},
	})

	if err != nil {
		client.logger.Error(fmt.Sprintf("Error responding with embed: %v", err))
	}
}

func createChampionSpellEmbed(champion *WikiChampion, spell *WikiSpell, key string) SpellEmbeds {
	championID := champion.ID
	video := fmt.Sprintf("https://d28xe8vt774jo5.cloudfront.net/champion-abilities/%04d/ability_%04d_%v1.webm", championID, championID, key)

	var damageType string
	switch spell.DamageType {
	case "PHYSICAL_DAMAGE":
		damageType = "Physic"
	case "MAGIC_DAMAGE":
		damageType = "Magic"
	case "TRUE_DAMAGE":
		damageType = "True"
	case "PURE_DAMAGE":
		damageType = "Pure"
	case "MIXED_DAMAGE":
		damageType = "Mixed"
	case "OTHER_DAMAGE":
		damageType = "Other"
	default:
		damageType = "None"
	}

	cooldownBurn, affectedByCDR := getSpellCooldownBurn(&spell.Cooldown)

	var affectedStr string
	if affectedByCDR {
		affectedStr = "Yes"
	} else {
		affectedStr = "No"
	}

	var shieldable string
	if spell.SpellShieldable != "" {
		shieldable = cases.Title(language.English, cases.NoLower).String(strings.ToLower(spell.SpellShieldable))
	} else {
		shieldable = "No"
	}

	var projectile string
	if spell.Projectile != "" {
		projectile = cases.Title(language.English, cases.NoLower).String(strings.ToLower(spell.Projectile))
	} else {
		projectile = "False"
	}

	costBurn := getSpellCostBurn(&spell.Cost)
	angle := getDefaultIntString(spell.Angle)
	resource := getDefaultString(spell.Resource)
	affects := getDefaultString(spell.Affects)
	targetRange := getSpellRangeOrRadius(spell.TargetRange)
	effectRadius := getSpellRangeOrRadius(spell.EffectRadius)

	speed := getDefaultIntString(strings.Replace(spell.Speed, " ", "", -1))
	castTime := getDefaultIntString(spell.CastTime)
	width := getDefaultIntString(spell.Width)

	description := fmt.Sprintf("[Wiki](https://leagueoflegends.fandom.com/wiki/%v/LoL#%v) - [Video](%v)\n\n", strings.Replace(champion.Name, " ", "_", -1), strings.Replace(spell.Name, " ", "_", -1), video)
	for _, effects := range spell.Effects {
		description += fmt.Sprintf("%v\n\n", effects.Description)
	}

	championAndSkillName := fmt.Sprintf("%v - %v", champion.Name, spell.Name)

	spellEmbed := discordgo.MessageEmbed{
		Title:       championAndSkillName,
		URL:         video,
		Color:       embedColor,
		Description: description,
		Thumbnail:   &discordgo.MessageEmbedThumbnail{URL: spell.Icon},
		Fields: []*discordgo.MessageEmbedField{
			{Name: "Cooldown", Value: cooldownBurn, Inline: true},
			{Name: "Affected by CDR", Value: affectedStr, Inline: true},
			{Name: "Damage Type", Value: damageType, Inline: true},

			{Name: "", Value: ""},
			{Name: "Cost", Value: costBurn, Inline: true},
			{Name: "Resource", Value: resource, Inline: true},
			{Name: "Affects", Value: affects, Inline: true},

			{Name: "", Value: ""},
			{Name: "Targeting", Value: spell.Targeting, Inline: true},
			{Name: "Spell Shieldable", Value: shieldable, Inline: true},
			{Name: "Projecile", Value: projectile, Inline: true},

			{Name: "", Value: ""},
			{Name: "Angle", Value: angle, Inline: true},
			{Name: "Effect Radius", Value: effectRadius, Inline: true},
			{Name: "Target Range", Value: targetRange, Inline: true},

			{Name: "", Value: ""},
			{Name: "Speed", Value: speed, Inline: true},
			{Name: "Cast Time", Value: castTime, Inline: true},
			{Name: "Width", Value: width, Inline: true},
		}}

	fields := []*discordgo.MessageEmbedField{}
	for _, effect := range spell.Effects {
		for _, leveling := range effect.Leveling {
			values := []string{}
			for _, modifier := range leveling.Modifiers {
				for _, element := range modifier.Values {
					values = append(values, fmt.Sprintf("%v", element))
				}
			}

			fields = append(fields, &discordgo.MessageEmbedField{Name: leveling.Attribute, Value: fmt.Sprintf("%v %v", strings.Join(values, "/"), leveling.Modifiers[0].Units[0]), Inline: true})
		}
	}

	wikiLink := fmt.Sprintf("https://leagueoflegends.fandom.com/wiki/%v/LoL#%v", strings.Replace(champion.Name, " ", "_", -1), strings.Replace(spell.Name, " ", "_", -1))

	spellModifiersEmbed := discordgo.MessageEmbed{
		Title:  championAndSkillName,
		URL:    wikiLink,
		Color:  embedColor,
		Fields: fields,
	}

	spellNotesEmbed := discordgo.MessageEmbed{
		Title:       championAndSkillName,
		URL:         wikiLink,
		Color:       embedColor,
		Description: spell.Notes,
	}

	return SpellEmbeds{General: spellEmbed, Modifiers: spellModifiersEmbed, Notes: spellNotesEmbed}
}

func createSpellInfo(spell *WikiSpell, key string, index int) SpellInfo {
	var spellFullName string
	if key == "P" {
		spellFullName = fmt.Sprintf("Passive - %v", spell.Name)
	} else {
		spellFullName = fmt.Sprintf("%v - %v", key, spell.Name)
	}

	return SpellInfo{Index: index, Key: key, FullName: spellFullName}
}

func getSpellCooldownBurn(cooldown *Cooldown) (string, bool) {
	if cooldown == nil {
		return "0", false
	}

	var cd []string
	for _, modifier := range cooldown.Modifiers {
		for _, value := range modifier.Values {
			cd = append(cd, fmt.Sprintf("%v", value))
		}
	}

	if len(cd) == 0 {
		return "0", false
	}

	if cd[0] == cd[len(cd)-1] {
		return cd[0], cooldown.AffectedByCDR
	}

	return strings.Join(cd, "/"), cooldown.AffectedByCDR
}

func getSpellCostBurn(cost *Cost) string {
	if cost == nil {
		return "0"
	}

	var costs []string
	for _, modifier := range cost.Modifiers {
		for _, value := range modifier.Values {
			costs = append(costs, fmt.Sprintf("%v", value))
		}
	}

	if len(costs) == 0 {
		return "0"
	}

	if costs[0] == costs[len(costs)-1] {
		return costs[0]
	}

	return strings.Join(costs, "/")
}

func getSpellRangeOrRadius(rangeOrRadius string) string {
	if rangeOrRadius == "" {
		return "0"
	}

	if strings.Contains(rangeOrRadius, "based on level") {
		burn := strings.Replace(rangeOrRadius, " ", "/", -1)
		return strings.Replace(burn, "/(basedonlevel)", " (based on level)", -1)
	}

	return strings.Replace(rangeOrRadius, " / ", "/", -1)
}
