package main

import (
	"discord-esports/models"
	"fmt"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
)

type SpellEmbeds struct {
	General   discordgo.MessageEmbed
	Modifiers discordgo.MessageEmbed
	Notes     discordgo.MessageEmbed
}

type SpellInfo struct {
	Key      string
	FullName string
	Index    int
}

func SpellCommand(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
	options := interaction.ApplicationCommandData().Options
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
	for _, opt := range options {
		optionMap[opt.Name] = opt
	}

	if interaction.Type == discordgo.InteractionApplicationCommandAutocomplete {
		if champion, ok := optionMap["champion"]; ok {
			if champion.Focused {
				autoCompleteChampionName(session, interaction, champion.StringValue())
				return
			}
		}

		if spell, ok := optionMap["spell"]; ok {
			if spell.Focused {
				autoCompleteSpell(session, interaction, optionMap["champion"].StringValue())
			}
		}

		return
	}

	championKey := optionMap["champion"].StringValue()
	if _, ok := championsNames[championKey]; !ok {
		respondWithMessage(interaction.Interaction, fmt.Sprintf("Champion '%v' not found.", championKey))
		return
	}

	spell := strings.Split(optionMap["spell"].StringValue(), ",")

	if len(spell) != 2 {
		respondWithMessage(interaction.Interaction, "Spell not provided.")
		return
	}

	spellEmbeds, ok := spellsEmbeds[championKey][spell[0]]
	if !ok {
		respondWithMessage(interaction.Interaction, fmt.Sprintf("Spell '%v' not found.", optionMap["spell"].StringValue()))
		return
	}

	spellIndex, err := strconv.Atoi(spell[1])
	if err != nil {
		client.logger.Error(fmt.Sprintf("error converting spell index to int: %v", err))
	}

	if spellIndex < 0 || spellIndex > len(spellEmbeds)-1 {
		respondWithMessage(interaction.Interaction, fmt.Sprintf("Spell '%v' not found.", optionMap["spell"].StringValue()))
		return
	}

	components := []discordgo.MessageComponent{discordgo.ActionsRow{Components: []discordgo.MessageComponent{
		discordgo.Button{Label: "Modifiers", CustomID: fmt.Sprintf("modifiers_%v_%v_%v", championKey, spell[0], spellIndex), Disabled: true},
		discordgo.Button{Label: "Notes", CustomID: fmt.Sprintf("notes_%v_%v_%v", championKey, spell[0], spellIndex), Disabled: true},
	}}}

	var embed *discordgo.MessageEmbed
	spellEmbed := spellEmbeds[spellIndex]

	embedType, ok := optionMap["type"]
	if !ok || embedType.StringValue() == "" {
		embed = &spellEmbed.General
	} else {
		switch embedType.StringValue() {
		case "modifiers":
			embed = &spellEmbed.Modifiers
			components = []discordgo.MessageComponent{}
		case "notes":
			embed = &spellEmbed.Notes
			components = []discordgo.MessageComponent{}
		default:
			embed = &spellEmbed.General
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

func createChampionSpellEmbed(champion *models.Champion, spell *models.ChampionSpell) SpellEmbeds {
	description := fmt.Sprintf("[Wiki](%v) - [Video](%v)\n\n", spell.WikiPage, spell.Video)
	for _, effects := range spell.Effects {
		description += fmt.Sprintf("%v\n\n", effects.Description)
	}

	championAndSkillName := fmt.Sprintf("%v - %v", champion.Name, spell.Name)

	spellEmbed := discordgo.MessageEmbed{
		Title:       championAndSkillName,
		URL:         spell.WikiPage,
		Color:       DISCORD_EMBED_COLOR,
		Description: description,
		Thumbnail:   &discordgo.MessageEmbedThumbnail{URL: spell.Icon},
		Fields: []*discordgo.MessageEmbedField{
			{Name: "Cooldown", Value: spell.Cooldown, Inline: true},
			{Name: "Affected by CDR", Value: spell.AffectedByCDR, Inline: true},
			{Name: "Damage Type", Value: spell.DamageType, Inline: true},

			{Name: "", Value: ""},
			{Name: "Cost", Value: spell.Cost, Inline: true},
			{Name: "Resource", Value: spell.Resource, Inline: true},
			{Name: "Affects", Value: spell.Affects, Inline: true},

			{Name: "", Value: ""},
			{Name: "Targeting", Value: spell.Targeting, Inline: true},
			{Name: "Spell Shieldable", Value: spell.SpellShieldable, Inline: true},
			{Name: "Projecile", Value: spell.Projectile, Inline: true},

			{Name: "", Value: ""},
			{Name: "Angle", Value: spell.Angle, Inline: true},
			{Name: "Effect Radius", Value: spell.EffectRadius, Inline: true},
			{Name: "Target Range", Value: spell.TargetRange, Inline: true},

			{Name: "", Value: ""},
			{Name: "Speed", Value: spell.Speed, Inline: true},
			{Name: "Cast Time", Value: spell.CastTime, Inline: true},
			{Name: "Width", Value: spell.Width, Inline: true},
		}}

	fields := []*discordgo.MessageEmbedField{}
	for _, effect := range spell.Effects {

		for _, leveling := range effect.Leveling {
			modifiers := []string{}
			for _, modifier := range leveling.Modifiers {
				values := modifier.Values
				if modifier.Unit != "" {
					values += fmt.Sprintf(" %v", modifier.Unit)
					modifiers = append(modifiers, values)
				}
			}

			for _, modifier := range modifiers {
				fields = append(fields, &discordgo.MessageEmbedField{Name: leveling.Attribute, Value: modifier, Inline: true})
			}
		}
	}

	modifiersDesc := ""
	if len(fields) == 0 {
		modifiersDesc = "No modifiers available."
	}

	spellModifiersEmbed := discordgo.MessageEmbed{
		Title:       championAndSkillName,
		URL:         spell.WikiPage,
		Color:       DISCORD_EMBED_COLOR,
		Description: modifiersDesc,
		Fields:      fields,
	}

	spellNotes := ""
	if len(spell.Notes) == 0 {
		spellNotes = "No additional notes available."
	} else {
		spellNotes = strings.Join(spell.Notes, "\n\n")
	}

	// Making the notes smaller to avoid hitting the discord embed limit
	if len(spellNotes) > 2000 {
		r := []rune(spellNotes)
		trunc := r[:2000]
		spellNotes = string(trunc) + fmt.Sprintf("... check the [Wiki](%v) page for more details.", spell.WikiPage)
	}

	spellNotesEmbed := discordgo.MessageEmbed{
		Title:       championAndSkillName,
		URL:         spell.WikiPage,
		Color:       DISCORD_EMBED_COLOR,
		Description: spellNotes,
	}

	return SpellEmbeds{General: spellEmbed, Modifiers: spellModifiersEmbed, Notes: spellNotesEmbed}
}

func createSpellInfo(spell *models.ChampionSpell, key string, index int) SpellInfo {
	var spellFullName string
	if key == "P" {
		spellFullName = fmt.Sprintf("Passive - %v", spell.Name)
	} else {
		spellFullName = fmt.Sprintf("%v - %v", key, spell.Name)
	}

	return SpellInfo{Index: index, Key: key, FullName: spellFullName}
}
