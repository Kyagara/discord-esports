package main

import (
	"fmt"
	"strconv"
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
		if ok := optionMap["champion"].Focused; ok {
			filteredNames := make(map[string]string)
			for id, name := range championsNames {
				if strings.HasPrefix(strings.ToLower(name), strings.ToLower(optionMap["champion"].StringValue())) {
					filteredNames[name] = id
				}

				if len(filteredNames) == 20 {
					break
				}
			}

			var choices []*discordgo.ApplicationCommandOptionChoice
			for id, name := range filteredNames {
				choices = append(choices, &discordgo.ApplicationCommandOptionChoice{Name: name, Value: id})
			}

			err := session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionApplicationCommandAutocompleteResult,
				Data: &discordgo.InteractionResponseData{
					Choices: choices,
				},
			})

			if err != nil {
				client.logger.Error(fmt.Sprintf("Error sending champion autocomplete: %v", err))
			}

			return
		}

		if ok := optionMap["spell"].Focused; ok {
			spells := spellsInfo[optionMap["champion"].StringValue()]

			var choices []*discordgo.ApplicationCommandOptionChoice
			for _, spell := range spells {
				choices = append(choices, &discordgo.ApplicationCommandOptionChoice{Name: spell.FullName, Value: fmt.Sprintf("%v,%v", spell.Key, spell.Index)})
			}

			err := session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionApplicationCommandAutocompleteResult,
				Data: &discordgo.InteractionResponseData{
					Choices: choices,
				},
			})

			if err != nil {
				client.logger.Error(fmt.Sprintf("Error sending spell autocomplete: %v", err))
			}
		}

		return
	}

	key := optionMap["champion"].StringValue()

	if ok := optionMap["spell"].StringValue() != ""; ok {
		spell := strings.Split(optionMap["spell"].StringValue(), ",")
		spellIndex, err := strconv.Atoi(spell[1])
		if err != nil {
			client.logger.Error(fmt.Sprintf("error converting spell index to int: %v", err))
		}

		err = client.session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{&spellsEmbeds[key][spell[0]][spellIndex].General},
				Components: []discordgo.MessageComponent{discordgo.ActionsRow{Components: []discordgo.MessageComponent{
					discordgo.Button{Label: "Modifiers", CustomID: fmt.Sprintf("modifiers_%v_%v_%v", key, spell[0], spellIndex)},
					discordgo.Button{Label: "Notes", CustomID: fmt.Sprintf("notes_%v_%v_%v", key, spell[0], spellIndex)},
				}}},
			},
		})

		if err != nil {
			client.logger.Error(fmt.Sprintf("Error responding with embed: %v", err))
		}

		return
	}

	embed := championsEmbeds[key].General
	err := client.session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{&embed},
			Components: []discordgo.MessageComponent{discordgo.ActionsRow{Components: []discordgo.MessageComponent{
				discordgo.Button{Label: "Spells", CustomID: fmt.Sprintf("spells_%v", key)},
				discordgo.Button{Label: "Skins", CustomID: fmt.Sprintf("skins_%v", key)},
			}}},
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

	fields := []*discordgo.MessageEmbedField{
		{Name: "HP | Regen", Value: fmt.Sprintf("``%v (+ %v)``\n``%v (+ %v)``", champion.Stats.Health.Flat, champion.Stats.Health.PerLevel, champion.Stats.HealthRegen.Flat, champion.Stats.HealthRegen.PerLevel), Inline: true},
	}

	if champion.Stats.Mana.Flat != 0 && champion.Stats.Mana.PerLevel != 0 {
		fields = append(fields, &discordgo.MessageEmbedField{Name: "MP | Regen", Value: fmt.Sprintf("``%v (+ %v)``\n``%v (+ %v)``", champion.Stats.Mana.Flat, champion.Stats.Mana.PerLevel, champion.Stats.ManaRegen.Flat, champion.Stats.ManaRegen.PerLevel), Inline: true})
	}

	fields = append(fields, &discordgo.MessageEmbedField{Name: "Armor | MR", Value: fmt.Sprintf("``%v (+ %v)``\n``%v (+ %v)``", champion.Stats.Armor.Flat, champion.Stats.Armor.PerLevel, champion.Stats.MagicResistance.Flat, champion.Stats.MagicResistance.PerLevel), Inline: true})

	fields = append(fields, &discordgo.MessageEmbedField{Name: "", Value: ""})
	fields = append(fields, &discordgo.MessageEmbedField{Name: "Attack Damage", Value: fmt.Sprintf("``%v (+ %v)``", champion.Stats.AttackDamage.Flat, champion.Stats.AttackDamage.PerLevel), Inline: true})

	if champion.Stats.AttackSpeed.Flat != 0 && champion.Stats.AttackSpeed.PerLevel != 0 {
		fields = append(fields, &discordgo.MessageEmbedField{Name: "Attack Speed", Value: fmt.Sprintf("``%v (+ %v)``", champion.Stats.AttackSpeed.Flat, champion.Stats.AttackSpeed.PerLevel), Inline: true})
	}

	/* 	if champion.Stats.Crit != 0 && champion.Stats.CritPerLevel != 0 {
		fields = append(fields, &discordgo.MessageEmbedField{Name: "Crit", Value: fmt.Sprintf("``%v (+ %v)``", champion.Stats.CriticalStrikeDamage, champion.Stats.CritPerLevel), Inline: true})
	} */

	/*
		cds := fmt.Sprintf("``Q - %v\nW - %v\nE - %v\nR - %v``", champion.Spells.[0].CooldownBurn, data.Spells[1].CooldownBurn, data.Spells[2].CooldownBurn, data.Spells[3].CooldownBurn)
		costs := fmt.Sprintf("``Q - %v\nW - %v\nE - %v\nR - %v``", data.Spells[0].CostBurn, data.Spells[1].CostBurn, data.Spells[2].CostBurn, data.Spells[3].CostBurn)
		ranges := fmt.Sprintf("``Q - %v\nW - %v\nE - %v\nR - %v``", data.Spells[0].RangeBurn, data.Spells[1].RangeBurn, data.Spells[2].RangeBurn, data.Spells[3].RangeBurn)
	*/ /*

		fields = append(fields, &discordgo.MessageEmbedField{Name: "", Value: ""})
		fields = append(fields, &discordgo.MessageEmbedField{Name: "Cooldown", Value: cds, Inline: true})
		fields = append(fields, &discordgo.MessageEmbedField{Name: "Spell Cost", Value: costs, Inline: true})
		fields = append(fields, &discordgo.MessageEmbedField{Name: "Spell Range", Value: ranges, Inline: true}) */

	fields = append(fields, &discordgo.MessageEmbedField{Name: "", Value: ""})
	fields = append(fields, &discordgo.MessageEmbedField{Name: "Range", Value: fmt.Sprintf("``%v``", champion.Stats.AttackRange.Flat), Inline: true})
	fields = append(fields, &discordgo.MessageEmbedField{Name: "Movement", Value: fmt.Sprintf("``%v``", champion.Stats.MovementSpeed.Flat), Inline: true})
	fields = append(fields, &discordgo.MessageEmbedField{Name: "Resource", Value: resource, Inline: true})

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
		Description: fmt.Sprintf("[Wiki](https://leagueoflegends.fandom.com/wiki/%s/LoL) - [LoLalytics](https://lolalytics.com/lol/%s/build/)", strings.Replace(champion.Name, " ", "_", 1), strings.ToLower(champion.Key)),
		Fields:      fields,
		Footer:      &discordgo.MessageEmbedFooter{Text: strings.Join(tags, ", ")},
	}

	return ChampionEmbeds{General: championEmbed}
}
