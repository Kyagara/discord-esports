package main

import (
	"discord-esports/models"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"os"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var (
	wikiChampions = make(map[string]WikiChampion)
	champions     = make(map[string]models.Champion)
)

func main() {
	if _, err := os.Stat("./champions/normalized"); os.IsNotExist(err) {
		if err := os.Mkdir("./champions/normalized", os.ModePerm); err != nil {
			log.Fatal(err)
		}
	} else {
		err := os.RemoveAll("./champions/normalized")
		if err != nil {
			log.Fatal(err)
		}
		if err := os.Mkdir("./champions/normalized", os.ModePerm); err != nil {
			log.Fatal(err)
		}
	}

	files, err := os.ReadDir("./champions")
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		if f.IsDir() {
			continue
		}

		file, err := os.ReadFile(fmt.Sprintf("./champions/%v", f.Name()))
		if err != nil {
			log.Fatal(err)
		}

		champion := WikiChampion{}
		err = json.Unmarshal(file, &champion)
		if err != nil {
			log.Fatal(err)
		}

		wikiChampions[champion.Key] = champion
	}

	for _, wikiChampion := range wikiChampions {
		champion := models.Champion{
			ID:               wikiChampion.ID,
			Key:              wikiChampion.Key,
			Name:             wikiChampion.Name,
			Icon:             wikiChampion.Icon,
			Lore:             wikiChampion.Lore,
			PatchLastChanged: wikiChampion.PatchLastChanged,
			FullTitle:        fmt.Sprintf("%v, %v", wikiChampion.Name, wikiChampion.Title),
			AttackType:       getCapitalized(wikiChampion.AttackType),
			Roles:            getRoles(wikiChampion.Roles),
			Resource:         getCapitalizedEnum(wikiChampion.Resource),
			AdaptiveType:     getCapitalizedEnum(wikiChampion.AdaptiveType),
			OfficialPage:     fmt.Sprintf("https://www.leagueoflegends.com/en-us/champions/%v/", strings.Replace(wikiChampion.Name, " ", "-", -1)),
			Stats:            getChampionStats(wikiChampion.Stats),
			Spells:           getSpells(wikiChampion.Spells, wikiChampion.Name, wikiChampion.ID),
			WikiPage:         fmt.Sprintf("https://leagueoflegends.fandom.com/wiki/%v/LoL", strings.Replace(wikiChampion.Name, " ", "_", -1)),
		}
		champions[wikiChampion.Name] = champion
	}

	for _, champion := range champions {
		file, err := json.MarshalIndent(champion, "", " ")
		if err != nil {
			log.Fatal(err)
		}

		err = os.WriteFile(fmt.Sprintf("./champions/normalized/%v.json", champion.Key), file, 0644)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func getChampionStats(wikiStats WikiChampionStats) models.ChampionStats {
	return models.ChampionStats{
		Health:          getFlatAndPerLevelString(wikiStats.Health),
		HealthRegen:     getFlatAndPerLevelString(wikiStats.HealthRegen),
		Mana:            getFlatAndPerLevelString(wikiStats.Mana),
		ManaRegen:       getFlatAndPerLevelString(wikiStats.ManaRegen),
		Armor:           getFlatAndPerLevelString(wikiStats.Armor),
		MagicResistance: getFlatAndPerLevelString(wikiStats.MagicResistance),
		MovementSpeed:   getFlatAndPerLevelString(wikiStats.MovementSpeed),
		AttackRange:     getFlatAndPerLevelString(wikiStats.AttackRange),
		AttackSpeed:     getFlatAndPerLevelString(wikiStats.AttackSpeed),
		AttackDamage:    getFlatAndPerLevelString(wikiStats.AttackDamage),
	}
}

func getSpells(wikiSpells WikiChampionSpells, championName string, championID int) models.ChampionSpells {
	spells := models.ChampionSpells{
		Passive: []models.ChampionSpell{},
		Q:       []models.ChampionSpell{},
		W:       []models.ChampionSpell{},
		E:       []models.ChampionSpell{},
		R:       []models.ChampionSpell{},
	}
	for _, spell := range wikiSpells.Passive {
		spells.Passive = append(spells.Passive, getSpellStats(spell, championName, championID, "P"))
	}
	for _, spell := range wikiSpells.Q {
		spells.Q = append(spells.Q, getSpellStats(spell, championName, championID, "Q"))
	}
	for _, spell := range wikiSpells.W {
		spells.W = append(spells.W, getSpellStats(spell, championName, championID, "W"))
	}
	for _, spell := range wikiSpells.E {
		spells.E = append(spells.E, getSpellStats(spell, championName, championID, "E"))
	}
	for _, spell := range wikiSpells.R {
		spells.R = append(spells.R, getSpellStats(spell, championName, championID, "R"))
	}
	return spells
}

func getSpellStats(wikiSpell WikiSpell, championName string, championID int, spellKey string) models.ChampionSpell {
	spellShieldable := getCapitalized(wikiSpell.SpellShieldable)
	switch spellShieldable {
	case "Special":
		spellShieldable = "Yes"
	case "True":
		spellShieldable = "Yes"
	case "False":
		spellShieldable = "No"
	default:
		spellShieldable = "No"
	}

	projectile := getCapitalized(wikiSpell.Projectile)
	switch projectile {
	case "True":
		projectile = "Yes"
	case "False":
		projectile = "No"
	default:
		projectile = "No"
	}

	// For now, ignoring tables from spells notes
	// Notes also need works to allow for bullet points
	notesTemp := strings.Split(wikiSpell.Notes, "\n")
	var notes []string
	tableIndex := -1
	for i, note := range notesTemp {
		if note == "Type" {
			break
		}

		if strings.HasPrefix(note, "The following table") {
			tableIndex = i
			break
		}

		noteTemp := strings.Replace(note, "  ", " ", -1)

		if strings.HasPrefix(noteTemp, " ") {
			noteTemp = strings.Replace(noteTemp, " ", "", 1)
		}

		if noteTemp != "" {
			notes = append(notes, noteTemp)
		}
	}

	if tableIndex != -1 {
		notes = notes[:tableIndex]
	}

	var effects []models.ChampionSpellEffect
	for _, wikiEffect := range wikiSpell.Effects {
		var leveling []models.ChampionSpellLeveling

		for _, wikiLeveling := range wikiEffect.Leveling {
			var modifiers []models.ChampionSpellModifier

			for _, modifier := range wikiLeveling.Modifiers {
				modifiersValues := []string{}
				for _, value := range modifier.Values {
					modifiersValues = append(modifiersValues, fmt.Sprintf("%v", value))
				}

				values := strings.Join(modifiersValues, "/")
				if modifiersValues[0] == modifiersValues[len(modifiersValues)-1] {
					values = modifiersValues[0]
				}

				modifiers = append(modifiers, models.ChampionSpellModifier{Values: values, Unit: modifier.Units[0]})
			}
			leveling = append(leveling, models.ChampionSpellLeveling{Attribute: wikiLeveling.Attribute, Modifiers: modifiers})
		}
		effects = append(effects, models.ChampionSpellEffect{Description: wikiEffect.Description, Leveling: leveling})
	}

	cooldown, affectedByCDR := getSpellCooldown(wikiSpell.Cooldown)

	return models.ChampionSpell{
		Name:            wikiSpell.Name,
		Icon:            wikiSpell.Icon,
		Cooldown:        cooldown,
		WikiPage:        fmt.Sprintf("https://leagueoflegends.fandom.com/wiki/%v/LoL#%v", strings.Replace(championName, " ", "_", -1), strings.Replace(wikiSpell.Name, " ", "_", -1)),
		Video:           fmt.Sprintf("https://d28xe8vt774jo5.cloudfront.net/champion-abilities/%04d/ability_%04d_%v1.webm", championID, championID, spellKey),
		AffectedByCDR:   affectedByCDR,
		Resource:        getCapitalizedEnum(wikiSpell.Resource),
		DamageType:      getCapitalizedEnum(wikiSpell.DamageType),
		Cost:            getSpellCost(wikiSpell.Cost),
		Angle:           getDefaultIntString(wikiSpell.Angle),
		Affects:         getDefaultString(wikiSpell.Affects),
		TargetRange:     getSpellRangeOrRadius(wikiSpell.TargetRange),
		EffectRadius:    getSpellRangeOrRadius(wikiSpell.EffectRadius),
		Speed:           getDefaultIntString(wikiSpell.Speed),
		CastTime:        getDefaultIntString(wikiSpell.CastTime),
		Width:           getDefaultIntString(wikiSpell.Width),
		SpellShieldable: spellShieldable,
		Projectile:      projectile,
		Targeting:       getCapitalized(wikiSpell.Targeting),
		Notes:           notes,
		Effects:         effects,
	}
}

func getFlatAndPerLevelString(stat WikiStat) string {
	field := getDefaultIntString(fmt.Sprintf("%v", stat.Flat))
	if stat.PerLevel > 0 {
		field += fmt.Sprintf(" (+ %v)", stat.PerLevel)
	}
	return field
}

func getCapitalizedEnum(value string) string {
	if value == "" {
		return "None"
	}
	temp := []string{}
	words := strings.Split(value, "_")
	for _, word := range words {
		temp = append(temp, getCapitalized(word))
	}
	return strings.Join(temp, " ")
}

func getRoles(value []string) []string {
	var tags []string
	for _, tag := range value {
		tags = append(tags, getCapitalized(tag))
	}
	if len(tags) == 0 || tags[0] == "" {
		tags = []string{"None"}
	}
	return tags
}

func getCapitalized(value string) string {
	return cases.Title(language.English, cases.NoLower).String(strings.ToLower(value))
}

func getDefaultIntString(value string) string {
	if value != "" {
		if value == "None" {
			return "0"
		}
		return getCapitalized(value)
	} else {
		return "0"
	}
}

func getDefaultString(value string) string {
	if value != "" {
		return getCapitalized(value)
	} else {
		return "None"
	}
}

func getSpellCooldown(wikiCooldown WikiSpellCooldown) (string, string) {
	var cd []string
	for _, modifier := range wikiCooldown.Modifiers {
		for _, value := range modifier.Values {
			if value == math.Trunc(value) {
				cd = append(cd, fmt.Sprintf("%v", value))
			} else {
				cd = append(cd, fmt.Sprintf("%.2f", value))
			}
		}
	}

	if len(cd) == 0 {
		return "0", "No"
	}

	var affectedByCDR string
	if wikiCooldown.AffectedByCDR {
		affectedByCDR = "Yes"
	} else {
		affectedByCDR = "No"
	}

	if cd[0] == cd[len(cd)-1] {
		return cd[0], affectedByCDR
	}

	return strings.Join(cd, "/"), affectedByCDR
}

func getSpellCost(wikiCost WikiSpellCost) string {
	var costs []string
	for _, modifier := range wikiCost.Modifiers {
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
