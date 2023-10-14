package main

import (
	"discord-esports/models"
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

type ChampionEmbeds struct {
	General discordgo.MessageEmbed
}

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
	champion := championsEmbeds[championKey]

	err := client.session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{&champion.General},
		},
	})

	if err != nil {
		client.logger.Error(fmt.Sprintf("Error responding with embed: %v", err))
	}
}

func createChampionEmbed(champion *models.Champion) ChampionEmbeds {
	championEmbed := discordgo.MessageEmbed{
		Title: champion.FullTitle,
		URL:   champion.OfficialPage,
		Color: embedColor,
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: champion.Icon,
		},
		Description: fmt.Sprintf("[Wiki](%s) - [LoLalytics](https://lolalytics.com/lol/%s/build/)\n\n%v", champion.WikiPage, strings.ToLower(champion.Key), champion.Lore),
		Fields: []*discordgo.MessageEmbedField{
			{Name: "HP | Regen", Value: fmt.Sprintf("%v\n%v", champion.Stats.Health, champion.Stats.HealthRegen), Inline: true},
			{Name: "MP | Regen", Value: fmt.Sprintf("%v\n%v", champion.Stats.Mana, champion.Stats.ManaRegen), Inline: true},
			{Name: "Armor | MR", Value: fmt.Sprintf("%v\n%v", champion.Stats.Armor, champion.Stats.MagicResistance), Inline: true},
			{Name: "", Value: ""},
			{Name: "Attack Range", Value: fmt.Sprintf("%v", champion.Stats.AttackRange), Inline: true},
			{Name: "Attack Damage", Value: fmt.Sprintf("%v", champion.Stats.AttackDamage), Inline: true},
			{Name: "Attack Speed", Value: fmt.Sprintf("%v", champion.Stats.AttackSpeed), Inline: true},
			{Name: "", Value: ""},
			{Name: "Movement", Value: fmt.Sprintf("%v", champion.Stats.MovementSpeed), Inline: true},
			{Name: "Adaptive Type", Value: fmt.Sprintf("%v", champion.AdaptiveType), Inline: true},
			{Name: "Resource", Value: champion.Resource, Inline: true},
			{Name: "", Value: ""},
			{Name: "Attack Type", Value: champion.AttackType, Inline: true},
			{Name: "Patch last changed", Value: champion.PatchLastChanged, Inline: true},
		},
		Footer: &discordgo.MessageEmbedFooter{Text: strings.Join(champion.Roles, ", ")},
	}

	return ChampionEmbeds{General: championEmbed}
}
