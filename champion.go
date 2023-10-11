package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/Kyagara/equinox/clients/ddragon"
	"github.com/bwmarrin/discordgo"
)

var spellKeysDict = map[int]string{
	-1: "P",
	0:  "Q",
	1:  "W",
	2:  "E",
	3:  "R",
}

func ChampionCommand(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
	options := interaction.ApplicationCommandData().Options

	notUpToDate := time.Since(ddVersionUpdated) > time.Duration(4*time.Minute)

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

		championsNames = make(map[string]string)
		for _, c := range champions {
			championsNames[c.Name] = c.ID
		}
	}

	if interaction.Type == discordgo.InteractionApplicationCommandAutocomplete {
		filteredNames := make(map[string]string)
		for id, name := range championsNames {
			if strings.HasPrefix(strings.ToLower(name), strings.ToLower(options[0].StringValue())) {
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
			client.logger.Error(fmt.Sprintf("Error sending autocomplete: %v", err))
		}

		return
	}

	champion, ok := champions[options[0].StringValue()]
	if !ok {
		respondWithError(interaction.Interaction, fmt.Errorf("champion '%s' not found", options[0].StringValue()))
		return
	}

	data, err := client.equinox.DDragon.Champion.ByName(ddVersion, ddragon.EnUS, options[0].StringValue())
	if err != nil {
		respondWithError(interaction.Interaction, err)
		return
	}

	if options[1] != nil {
		optionValue := options[1].IntValue()

		var name string
		var description string
		var image string
		rangeBurn := "0"
		cooldownBurn := "0"
		costBurn := "0"
		key, err := strconv.Atoi(data.Key)
		if err != nil {
			client.logger.Error(fmt.Sprintf("error converting championg key to int: %v", err))
		}

		if optionValue == -1 {
			name = data.Passive.Name
			description = data.Passive.Description
			image = data.Passive.Image.Full
		} else {
			name = data.Spells[optionValue].Name
			description = data.Spells[optionValue].Description
			image = data.Spells[optionValue].Image.Full
			rangeBurn = data.Spells[optionValue].RangeBurn
			cooldownBurn = data.Spells[optionValue].CooldownBurn
			costBurn = data.Spells[optionValue].CostBurn
		}

		description = strings.Replace(description, "<br><br>", "\n", -1)
		description = strings.Replace(description, "<spellName>", "`", -1)
		description = strings.Replace(description, "</spellName>", "`", -1)

		r := regexp.MustCompile(`<.*?>`)
		description = r.ReplaceAllString(description, "")

		video := fmt.Sprintf("https://d28xe8vt774jo5.cloudfront.net/champion-abilities/%04d/ability_%04d_%v1.webm", key, key, spellKeysDict[int(optionValue)])

		embed := &discordgo.MessageEmbed{
			Title:       name,
			URL:         video,
			Color:       embedColor,
			Description: description,
			Thumbnail:   &discordgo.MessageEmbedThumbnail{URL: fmt.Sprintf("https://ddragon.leagueoflegends.com/cdn/%v/img/spell/%v", ddVersion, image)},
			Fields: []*discordgo.MessageEmbedField{
				{Name: "Links", Value: fmt.Sprintf("[Wiki](https://leagueoflegends.fandom.com/wiki/%v/LoL#%v) - [Video](%v)", strings.Replace(champion.Name, " ", "_", -1), strings.Replace(name, " ", "", -1), video)},
				{Name: "Range", Value: fmt.Sprintf("``%v``", rangeBurn), Inline: true},
				{Name: "Cost", Value: fmt.Sprintf("``%v``", costBurn), Inline: true},
				{Name: "Cooldown", Value: fmt.Sprintf("``%v``", cooldownBurn), Inline: true},
			}}

		respondWithEmbed(interaction.Interaction, []*discordgo.MessageEmbed{embed})
		return
	}

	resource := champion.Partype
	if champion.Partype == "" {
		resource = "None"
	}

	fields := []*discordgo.MessageEmbedField{
		{Name: "HP | Regen", Value: fmt.Sprintf("``%v (+ %v)``\n``%v (+ %v)``", champion.Stats.HP, champion.Stats.HPPerLevel, champion.Stats.HPRegen, champion.Stats.HPRegenPerLevel), Inline: true},
	}

	if champion.Stats.MP != 0 && champion.Stats.MPPerLevel != 0 {
		fields = append(fields, &discordgo.MessageEmbedField{Name: "MP | Regen", Value: fmt.Sprintf("``%v (+ %v)``\n``%v (+ %v)``", champion.Stats.MP, champion.Stats.MPPerLevel, champion.Stats.MPRegen, champion.Stats.MPRegenPerLevel), Inline: true})
	}

	fields = append(fields, &discordgo.MessageEmbedField{Name: "Armor | MR", Value: fmt.Sprintf("``%v (+ %v)``\n``%v (+ %v)``", champion.Stats.Armor, champion.Stats.ArmorPerLevel, champion.Stats.SpellBlock, champion.Stats.SpellBlockPerLevel), Inline: true})

	fields = append(fields, &discordgo.MessageEmbedField{Name: "", Value: ""})
	fields = append(fields, &discordgo.MessageEmbedField{Name: "Attack Damage", Value: fmt.Sprintf("``%v (+ %v)``", champion.Stats.AttackDamage, champion.Stats.AttackDamagePerLevel), Inline: true})

	if champion.Stats.AttackSpeed != 0 && champion.Stats.AttackSpeedPerLevel != 0 {
		fields = append(fields, &discordgo.MessageEmbedField{Name: "Attack Speed", Value: fmt.Sprintf("``%v (+ %v)``", champion.Stats.AttackSpeed, champion.Stats.AttackSpeedPerLevel), Inline: true})
	}

	if champion.Stats.Crit != 0 && champion.Stats.CritPerLevel != 0 {
		fields = append(fields, &discordgo.MessageEmbedField{Name: "Crit", Value: fmt.Sprintf("``%v (+ %v)``", champion.Stats.Crit, champion.Stats.CritPerLevel), Inline: true})
	}

	cds := fmt.Sprintf("``Q - %v\nW - %v\nE - %v\nR - %v``", data.Spells[0].CooldownBurn, data.Spells[1].CooldownBurn, data.Spells[2].CooldownBurn, data.Spells[3].CooldownBurn)
	costs := fmt.Sprintf("``Q - %v\nW - %v\nE - %v\nR - %v``", data.Spells[0].CostBurn, data.Spells[1].CostBurn, data.Spells[2].CostBurn, data.Spells[3].CostBurn)
	ranges := fmt.Sprintf("``Q - %v\nW - %v\nE - %v\nR - %v``", data.Spells[0].RangeBurn, data.Spells[1].RangeBurn, data.Spells[2].RangeBurn, data.Spells[3].RangeBurn)

	fields = append(fields, &discordgo.MessageEmbedField{Name: "", Value: ""})
	fields = append(fields, &discordgo.MessageEmbedField{Name: "Cooldown", Value: cds, Inline: true})
	fields = append(fields, &discordgo.MessageEmbedField{Name: "Spell Cost", Value: costs, Inline: true})
	fields = append(fields, &discordgo.MessageEmbedField{Name: "Spell Range", Value: ranges, Inline: true})

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
		Description: fmt.Sprintf("[Wiki](https://leagueoflegends.fandom.com/wiki/%s/LoL) - [LoLalytics](https://lolalytics.com/lol/%s/build/)", strings.Replace(champion.Name, " ", "_", 1), strings.ToLower(champion.ID)),
		Fields:      fields,
		Footer:      &discordgo.MessageEmbedFooter{Text: strings.Join(tags, ", ")},
	}

	respondWithEmbed(interaction.Interaction, []*discordgo.MessageEmbed{embed})
}
