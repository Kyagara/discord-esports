package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Kyagara/equinox/clients/data_dragon"
	"github.com/bwmarrin/discordgo"
)

func ChampionCommand(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
	options := interaction.ApplicationCommandData().Options

	if interaction.Type == discordgo.InteractionApplicationCommandAutocomplete {
		if champions == nil || time.Since(versionUpdated) > 4*time.Minute {
			ddVersion, err := dd.DataDragon.Version.Latest()
			if err != nil {
				respondWithErrorEmbed(interaction.Interaction, err)
				return
			}

			champions, err = dd.DataDragon.Champion.AllChampions(ddVersion, data_dragon.EnUS)
			if err != nil {
				respondWithErrorEmbed(interaction.Interaction, err)
				return
			}
		}

		var filteredNames []string
		for _, c := range championsNames {
			if strings.HasPrefix(c, options[0].StringValue()) {
				filteredNames = append(filteredNames, c)
			}
		}

		if len(filteredNames) > 20 {
			filteredNames = filteredNames[:20]
		}

		var choices []*discordgo.ApplicationCommandOptionChoice

		for _, n := range filteredNames {
			choices = append(choices, &discordgo.ApplicationCommandOptionChoice{Name: n, Value: n})
		}

		err := session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionApplicationCommandAutocompleteResult,
			Data: &discordgo.InteractionResponseData{
				Choices: choices,
			},
		})

		if err != nil {
			log.Printf("Error sending autocomplete: %v", err)
		}

		return
	}

	ddVersion, err := dd.DataDragon.Version.Latest()
	if err != nil {
		respondWithErrorEmbed(interaction.Interaction, err)
		return
	}

	versionUpdated = time.Now()

	champions, err = dd.DataDragon.Champion.AllChampions(ddVersion, data_dragon.EnUS)
	if err != nil {
		respondWithErrorEmbed(interaction.Interaction, err)
		return
	}

	champion := champions[options[0].StringValue()]

	if champion == nil {
		respondWithErrorEmbed(interaction.Interaction, fmt.Errorf("champion '%s' not found", options[0].StringValue()))
		return
	}

	embed := &discordgo.MessageEmbed{
		Title: fmt.Sprintf("%s, %s", champion.Name, champion.Title),
		URL:   fmt.Sprintf("https://www.leagueoflegends.com/en-us/champions/%s/", champion.ID),
		Color: embedColor,
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: fmt.Sprintf("https://ddragon.leagueoflegends.com/cdn/%s/img/champion/%s.png", ddVersion, champion.ID),
		},
		Description: fmt.Sprintf("[Wiki](https://leagueoflegends.fandom.com/wiki/%s/LoL) - [LoLalytics](https://lolalytics.com/lol/%s/build/)", champion.ID, strings.ToLower(champion.ID)),
		Fields: []*discordgo.MessageEmbedField{
			{Name: "HP | Regen", Value: fmt.Sprintf("``%v - %v per lvl``\n``%v - %v per lvl``", champion.Stats.HP, champion.Stats.HPPerLevel, champion.Stats.HPRegen, champion.Stats.HPRegenPerLevel), Inline: true},
			{Name: "MP | Regen", Value: fmt.Sprintf("``%v - %v per lvl``\n``%v - %v per lvl``", champion.Stats.MP, champion.Stats.MPPerLevel, champion.Stats.MPRegen, champion.Stats.MPRegenPerLevel), Inline: true},
			{Name: "Armor | MR", Value: fmt.Sprintf("``%v - %v per lvl``\n``%v - %v per lvl``", champion.Stats.Armor, champion.Stats.ArmorPerLevel, champion.Stats.SpellBlock, champion.Stats.SpellBlockPerLevel), Inline: true},

			{Name: "Attack damage", Value: fmt.Sprintf("``%v - %v per lvl``", champion.Stats.AttackDamage, champion.Stats.AttackDamagePerLevel), Inline: true},
			{Name: "Attack speed", Value: fmt.Sprintf("``%v - %v per lvl``", champion.Stats.AttackSpeedOffset, champion.Stats.AttackSpeedPerLevel), Inline: true},
			{Name: "Crit", Value: fmt.Sprintf("``%v - %v per lvl``", champion.Stats.Crit, champion.Stats.CritPerLevel), Inline: true},

			{Name: "Range", Value: fmt.Sprintf("``%v``", champion.Stats.AttackRange), Inline: true},
			{Name: "Movement", Value: fmt.Sprintf("``%v``", champion.Stats.MovementSpeed), Inline: true},
		},
		Footer: &discordgo.MessageEmbedFooter{Text: strings.Join(champion.Tags, ", ")},
	}

	session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed},
		},
	})
}
