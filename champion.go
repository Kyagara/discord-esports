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

	cds := []string{"%v/%v/%v/%v/%v", "%v/%v/%v/%v/%v", "%v/%v/%v/%v/%v", "%v/%v/%v"}

	for i := 0; i < 3; i++ {
		cd := data.Spells[i].CooldownCoefficients
		if cd[0] == cd[4] {
			cds[i] = fmt.Sprint(cd[0])
			continue
		}
		cds[i] = fmt.Sprintf(cds[i], cd[0], cd[1], cd[2], cd[3], cd[4])
	}

	cd := data.Spells[3].CooldownCoefficients
	if cd[0] == cd[2] {
		cds[3] = fmt.Sprint(cd[0])
	} else {
		cds[3] = fmt.Sprintf(cds[3], cd[0], cd[1], cd[2])
	}

	spellRanges := []string{"%v/%v/%v/%v/%v", "%v/%v/%v/%v/%v", "%v/%v/%v/%v/%v", "%v/%v/%v"}

	for i := 0; i < 3; i++ {
		spellRange := data.Spells[i].Range
		if spellRange[0] == spellRange[4] {
			spellRanges[i] = fmt.Sprint(spellRange[0])
			continue
		}
		spellRanges[i] = fmt.Sprintf(spellRanges[i], spellRange[0], spellRange[1], spellRange[2], spellRange[3], spellRange[4])
	}

	spellRange := data.Spells[3].Range
	if spellRange[0] == spellRange[2] {
		spellRanges[3] = fmt.Sprint(spellRange[0])
	} else {
		spellRanges[3] = fmt.Sprintf(spellRanges[3], spellRange[0], spellRange[1], spellRange[2])
	}

	spellCosts := []string{"%v/%v/%v/%v/%v", "%v/%v/%v/%v/%v", "%v/%v/%v/%v/%v", "%v/%v/%v"}

	for i := 0; i < 3; i++ {
		spellCost := data.Spells[i].CostCoefficients
		if spellCost[0] == spellCost[4] {
			spellCosts[i] = fmt.Sprint(spellCost[0])
			continue
		}
		spellCosts[i] = fmt.Sprintf(spellCosts[i], spellCost[0], spellCost[1], spellCost[2], spellCost[3], spellCost[4])
	}

	spellCost := data.Spells[3].CostCoefficients
	if spellCost[0] == spellCost[2] {
		spellCosts[3] = fmt.Sprint(spellCost[0])
	} else {
		spellCosts[3] = fmt.Sprintf(spellCosts[3], spellCost[0], spellCost[1], spellCost[2])
	}

	resource := champion.Partype
	if champion.Partype == "" {
		resource = "None"
	}

	fields := []*discordgo.MessageEmbedField{
		{Name: "HP | Regen", Value: fmt.Sprintf("``%v | %v per lvl``\n``%v | %v per lvl``", champion.Stats.HP, champion.Stats.HPPerLevel, champion.Stats.HPRegen, champion.Stats.HPRegenPerLevel), Inline: true},
	}

	if champion.Stats.MP != 0 && champion.Stats.MPPerLevel != 0 {
		fields = append(fields, &discordgo.MessageEmbedField{Name: "MP | Regen", Value: fmt.Sprintf("``%v | %v per lvl\n%v | %v per lvl``", champion.Stats.MP, champion.Stats.MPPerLevel, champion.Stats.MPRegen, champion.Stats.MPRegenPerLevel), Inline: true})
	}

	fields = append(fields, &discordgo.MessageEmbedField{Name: "Armor | MR", Value: fmt.Sprintf("``%v | %v per lvl\n%v | %v per lvl``", champion.Stats.Armor, champion.Stats.ArmorPerLevel, champion.Stats.SpellBlock, champion.Stats.SpellBlockPerLevel), Inline: true})

	fields = append(fields, &discordgo.MessageEmbedField{Name: "", Value: ""})
	fields = append(fields, &discordgo.MessageEmbedField{Name: "Attack Damage", Value: fmt.Sprintf("``%v | %v per lvl``", champion.Stats.AttackDamage, champion.Stats.AttackDamagePerLevel), Inline: true})

	if champion.Stats.AttackSpeed != 0 && champion.Stats.AttackSpeedPerLevel != 0 {
		fields = append(fields, &discordgo.MessageEmbedField{Name: "Attack Speed", Value: fmt.Sprintf("``%v | %v per lvl``", champion.Stats.AttackSpeed, champion.Stats.AttackSpeedPerLevel), Inline: true})
	}

	if champion.Stats.Crit != 0 && champion.Stats.CritPerLevel != 0 {
		fields = append(fields, &discordgo.MessageEmbedField{Name: "Crit", Value: fmt.Sprintf("``%v | %v per lvl``", champion.Stats.Crit, champion.Stats.CritPerLevel), Inline: true})
	}

	fields = append(fields, &discordgo.MessageEmbedField{Name: "", Value: ""})
	fields = append(fields, &discordgo.MessageEmbedField{Name: "Cooldown", Value: fmt.Sprintf("``Q - %v\nW - %v\nE - %v\nR - %v``", cds[0], cds[1], cds[2], cds[3]), Inline: true})
	fields = append(fields, &discordgo.MessageEmbedField{Name: "Spell Cost", Value: fmt.Sprintf("``Q - %v\nW - %v\nE - %v\nR - %v``", spellCosts[0], spellCosts[1], spellCosts[2], spellCosts[3]), Inline: true})
	fields = append(fields, &discordgo.MessageEmbedField{Name: "Spell Range", Value: fmt.Sprintf("``Q - %v\nW - %v\nE - %v\nR - %v``", spellRanges[0], spellRanges[1], spellRanges[2], spellRanges[3]), Inline: true})

	fields = append(fields, &discordgo.MessageEmbedField{Name: "", Value: ""})
	fields = append(fields, &discordgo.MessageEmbedField{Name: "Range", Value: fmt.Sprintf("``%v``", champion.Stats.AttackRange), Inline: true})
	fields = append(fields, &discordgo.MessageEmbedField{Name: "Movement", Value: fmt.Sprintf("``%v``", champion.Stats.MovementSpeed), Inline: true})
	fields = append(fields, &discordgo.MessageEmbedField{Name: "Resource", Value: fmt.Sprintf("``%v``", resource), Inline: true})

	tags := champion.Tags
	if len(tags) == 0 || tags[0] == "" {
		tags = []string{"None"}
	}

	embed := &discordgo.MessageEmbed{
		Title: fmt.Sprintf("%s, %s", champion.Name, champion.Title),
		URL:   fmt.Sprintf("https://www.leagueoflegends.com/en-us/champions/%s/", champion.ID),
		Color: embedColor,
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: fmt.Sprintf("https://ddragon.leagueoflegends.com/cdn/%s/img/champion/%s.png", ddVersion, champion.ID),
		},
		Description: fmt.Sprintf("[Wiki](https://leagueoflegends.fandom.com/wiki/%s/LoL) - [LoLalytics](https://lolalytics.com/lol/%s/build/)", champion.ID, strings.ToLower(champion.ID)),
		Fields:      fields,
		Footer:      &discordgo.MessageEmbedFooter{Text: strings.Join(tags, ", ")},
	}

	respondWithEmbed(interaction.Interaction, []*discordgo.MessageEmbed{embed})
}
