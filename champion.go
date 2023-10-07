package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/Kyagara/equinox/clients/ddragon"
	"github.com/bwmarrin/discordgo"
)

func ChampionCommand(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
	options := interaction.ApplicationCommandData().Options

	notUpToDate := time.Since(ddVersionUpdated) > 4*time.Minute

	if len(championsNames) == 0 || notUpToDate {
		ddVersion, err := client.equinox.DDragon.Version.Latest()
		if err != nil {
			respondWithError(interaction.Interaction, err)
			return
		}

		ddVersionUpdated = time.Now()

		champions, err = client.equinox.DDragon.Champion.AllChampions(ddVersion, ddragon.EnUS)
		if err != nil {
			respondWithError(interaction.Interaction, err)
			return
		}

		championsNames = []string{}
		for _, c := range champions {
			championsNames = append(championsNames, c.ID)
		}
	}

	if interaction.Type == discordgo.InteractionApplicationCommandAutocomplete {
		var filteredNames []string
		for _, c := range championsNames {
			if strings.HasPrefix(strings.ToLower(c), strings.ToLower(options[0].StringValue())) {
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
			client.logger.Error(fmt.Sprintf("Error sending autocomplete: %v", err))
		}

		return
	}

	champion, ok := champions[options[0].StringValue()]

	if !ok {
		respondWithError(interaction.Interaction, fmt.Errorf("champion '%s' not found", options[0].StringValue()))
		return
	}

	data, err := client.equinox.CDragon.Champion.ByName(ddVersion, options[0].StringValue())
	if err != nil {
		respondWithError(interaction.Interaction, err)
		return
	}

	cds := []string{"Q - %v/%v/%v/%v/%v", "W - %v/%v/%v/%v/%v", "E - %v/%v/%v/%v/%v", "R - %v/%v/%v"}

	for i := 0; i < 3; i++ {
		cd := data.Spells[i].CooldownCoefficients
		cds[i] = fmt.Sprintf(cds[i], cd[0], cd[1], cd[2], cd[3], cd[4])
	}

	cd := data.Spells[3].CooldownCoefficients
	cds[3] = fmt.Sprintf(cds[3], cd[0], cd[1], cd[2])

	embed := &discordgo.MessageEmbed{
		Title: fmt.Sprintf("%s, %s", champion.Name, champion.Title),
		URL:   fmt.Sprintf("https://www.leagueoflegends.com/en-us/champions/%s/", champion.ID),
		Color: embedColor,
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: fmt.Sprintf("https://ddragon.leagueoflegends.com/cdn/%s/img/champion/%s.png", ddVersion, champion.ID),
		},
		Description: fmt.Sprintf("[Wiki](https://leagueoflegends.fandom.com/wiki/%s/LoL) - [LoLalytics](https://lolalytics.com/lol/%s/build/)", champion.ID, strings.ToLower(champion.ID)),
		Fields: []*discordgo.MessageEmbedField{
			{Name: "HP | Regen", Value: fmt.Sprintf("``%v | %v per lvl``\n``%v | %v per lvl``", champion.Stats.HP, champion.Stats.HPPerLevel, champion.Stats.HPRegen, champion.Stats.HPRegenPerLevel), Inline: true},
			{Name: "MP | Regen", Value: fmt.Sprintf("``%v | %v per lvl``\n``%v | %v per lvl``", champion.Stats.MP, champion.Stats.MPPerLevel, champion.Stats.MPRegen, champion.Stats.MPRegenPerLevel), Inline: true},
			{Name: "Armor | MR", Value: fmt.Sprintf("``%v | %v per lvl``\n``%v | %v per lvl``", champion.Stats.Armor, champion.Stats.ArmorPerLevel, champion.Stats.SpellBlock, champion.Stats.SpellBlockPerLevel), Inline: true},

			{Name: "Attack damage", Value: fmt.Sprintf("``%v | %v per lvl``", champion.Stats.AttackDamage, champion.Stats.AttackDamagePerLevel), Inline: true},
			{Name: "Attack speed", Value: fmt.Sprintf("``%v | %v per lvl``", champion.Stats.AttackSpeed, champion.Stats.AttackSpeedPerLevel), Inline: true},
			{Name: "Crit", Value: fmt.Sprintf("``%v | %v per lvl``", champion.Stats.Crit, champion.Stats.CritPerLevel), Inline: true},

			{Name: "Cooldowns", Value: fmt.Sprintf("``%v\n%v\n%v\n%v``", cds[0], cds[1], cds[2], cds[3])},

			{Name: "Range", Value: fmt.Sprintf("``%v``", champion.Stats.AttackRange), Inline: true},
			{Name: "Movement", Value: fmt.Sprintf("``%v``", champion.Stats.MovementSpeed), Inline: true},
			{Name: "Resource", Value: fmt.Sprintf("``%v``", champion.Partype), Inline: true},
		},
		Footer: &discordgo.MessageEmbedFooter{Text: strings.Join(champion.Tags, ", ")},
	}

	respondWithEmbed(interaction.Interaction, []*discordgo.MessageEmbed{embed})
}
