package main

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func ChampionCommand(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
	options := interaction.ApplicationCommandData().Options
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
	for _, opt := range options {
		optionMap[opt.Name] = opt
	}

	if interaction.Type == discordgo.InteractionApplicationCommandAutocomplete {
		champion, ok := optionMap["champion"]
		if ok && champion.Focused {
			autoCompleteChampionName(session, interaction, champion.StringValue())
		}

		return
	}

	championKey := optionMap["champion"].StringValue()

	embed := &discordgo.MessageEmbed{}
	components := []discordgo.MessageComponent{discordgo.ActionsRow{Components: []discordgo.MessageComponent{
		discordgo.Button{Label: "Spells", CustomID: fmt.Sprintf("spells_%v", championKey)},
		discordgo.Button{Label: "Skins", CustomID: fmt.Sprintf("skins_%v", championKey)},
	}}}

	champion := championsEmbeds[championKey]

	embedType, ok := optionMap["type"]
	if !ok || embedType.StringValue() == "" {
		embed = &champion.General
	} else {
		switch embedType.StringValue() {
		case "spells":
			embed = &champion.Spells
			components = []discordgo.MessageComponent{}
		case "skins":
			embed = &champion.Skins
			components = []discordgo.MessageComponent{}
		default:
			embed = &champion.General
		}
	}

	err := client.session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
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

func createChampionEmbed(champion *WikiChampion) ChampionEmbeds {
	resource := champion.Resource
	if champion.Resource == "" {
		resource = "None"
	}

	/*
		cds := fmt.Sprintf("``Q - %v\nW - %v\nE - %v\nR - %v``", champion.Spells.[0].CooldownBurn, data.Spells[1].CooldownBurn, data.Spells[2].CooldownBurn, data.Spells[3].CooldownBurn)
		costs := fmt.Sprintf("``Q - %v\nW - %v\nE - %v\nR - %v``", data.Spells[0].CostBurn, data.Spells[1].CostBurn, data.Spells[2].CostBurn, data.Spells[3].CostBurn)
		ranges := fmt.Sprintf("``Q - %v\nW - %v\nE - %v\nR - %v``", data.Spells[0].RangeBurn, data.Spells[1].RangeBurn, data.Spells[2].RangeBurn, data.Spells[3].RangeBurn)
	*/ /*

		fields = append(fields, &discordgo.MessageEmbedField{Name: "", Value: ""})
		fields = append(fields, &discordgo.MessageEmbedField{Name: "Cooldown", Value: cds, Inline: true})
		fields = append(fields, &discordgo.MessageEmbedField{Name: "Spell Cost", Value: costs, Inline: true})
		fields = append(fields, &discordgo.MessageEmbedField{Name: "Spell Range", Value: ranges, Inline: true}) */

	tags := champion.Roles
	if len(tags) == 0 || tags[0] == "" {
		tags = []string{"None"}
	}

	championEmbed := discordgo.MessageEmbed{
		Title: fmt.Sprintf("%s, %s", champion.Name, champion.Title),
		URL:   fmt.Sprintf("https://www.leagueoflegends.com/en-us/champions/%s/", champion.Key),
		Color: embedColor,
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: champion.Icon,
		},
		Description: fmt.Sprintf("[Wiki](https://leagueoflegends.fandom.com/wiki/%s/LoL) - [LoLalytics](https://lolalytics.com/lol/%s/build/)\n\n%v", strings.Replace(champion.Name, " ", "_", 1), strings.ToLower(champion.Key), champion.Lore),
		Fields: []*discordgo.MessageEmbedField{
			{Name: "HP | Regen", Value: fmt.Sprintf("%v (+ %v)\n%v (+ %v)", champion.Stats.Health.Flat, champion.Stats.Health.PerLevel, champion.Stats.HealthRegen.Flat, champion.Stats.HealthRegen.PerLevel), Inline: true},
			{Name: "MP | Regen", Value: fmt.Sprintf("%v (+ %v)\n%v (+ %v)", champion.Stats.Mana.Flat, champion.Stats.Mana.PerLevel, champion.Stats.ManaRegen.Flat, champion.Stats.ManaRegen.PerLevel), Inline: true},
			{Name: "Armor | MR", Value: fmt.Sprintf("%v (+ %v)\n%v (+ %v)", champion.Stats.Armor.Flat, champion.Stats.Armor.PerLevel, champion.Stats.MagicResistance.Flat, champion.Stats.MagicResistance.PerLevel), Inline: true},
			{Name: "", Value: ""},
			{Name: "Attack Range", Value: fmt.Sprintf("%v", champion.Stats.AttackRange.Flat), Inline: true},
			{Name: "Attack Damage", Value: fmt.Sprintf("%v (+ %v)", champion.Stats.AttackDamage.Flat, champion.Stats.AttackDamage.PerLevel), Inline: true},
			{Name: "Attack Speed", Value: fmt.Sprintf("%v (+ %v)", champion.Stats.AttackSpeed.Flat, champion.Stats.AttackSpeed.PerLevel), Inline: true},
			{Name: "", Value: ""},
			{Name: "Movement", Value: fmt.Sprintf("%v", champion.Stats.MovementSpeed.Flat), Inline: true},
			{Name: "Adaptive Type", Value: fmt.Sprintf("%v", champion.AdaptiveType), Inline: true},
			{Name: "Resource", Value: resource, Inline: true},
			{Name: "", Value: ""},
			{Name: "Patch last changed", Value: champion.PatchLastChanged, Inline: true},
		},
		Footer: &discordgo.MessageEmbedFooter{Text: strings.Join(tags, ", ")},
	}

	return ChampionEmbeds{General: championEmbed}
}
