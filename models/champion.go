package models

type Champion struct {
	ID               int      `json:"id"`
	Key              string   `json:"key"`
	Name             string   `json:"name"`
	FullTitle        string   `json:"fullTitle"`
	Icon             string   `json:"icon"`
	Lore             string   `json:"lore"`
	Resource         string   `json:"resource"`
	AttackType       string   `json:"attackType"`
	AdaptiveType     string   `json:"adaptiveType"`
	Roles            []string `json:"roles"`
	PatchLastChanged string   `json:"patchLastChanged"`
	OfficialPage     string   `json:"officialPage"`
	WikiPage         string   `json:"wikiPage"`
	// ChampionStats also include per level number when available
	Stats  ChampionStats  `json:"stats"`
	Spells ChampionSpells `json:"spells"`
}

type ChampionStats struct {
	Health          string `json:"health"`
	HealthRegen     string `json:"healthRegen"`
	Mana            string `json:"mana"`
	ManaRegen       string `json:"manaRegen"`
	Armor           string `json:"armor"`
	MagicResistance string `json:"magicResistance"`
	AttackDamage    string `json:"attackDamage"`
	MovementSpeed   string `json:"movespeed"`
	AttackSpeed     string `json:"attackSpeed"`
	AttackRange     string `json:"attackRange"`
}

type ChampionSpells struct {
	Passive []ChampionSpell `json:"P"`
	Q       []ChampionSpell `json:"Q"`
	W       []ChampionSpell `json:"W"`
	E       []ChampionSpell `json:"E"`
	R       []ChampionSpell `json:"R"`
}

type ChampionSpell struct {
	Name            string                `json:"name"`
	Icon            string                `json:"icon"`
	Video           string                `json:"video"`
	WikiPage        string                `json:"wikiPage"`
	Effects         []ChampionSpellEffect `json:"effects"`
	AffectedByCDR   string                `json:"affectedByCDR"`
	Cost            string                `json:"cost"`
	Cooldown        string                `json:"cooldown"`
	Targeting       string                `json:"targeting"`
	Affects         string                `json:"affects"`
	SpellShieldable string                `json:"spellshieldable"`
	Resource        string                `json:"resource"`
	DamageType      string                `json:"damageType"`
	Projectile      string                `json:"projectile"`
	Notes           []string              `json:"notes"`
	Speed           string                `json:"speed"`
	Width           string                `json:"width"`
	Angle           string                `json:"angle"`
	CastTime        string                `json:"castTime"`
	EffectRadius    string                `json:"effectRadius"`
	TargetRange     string                `json:"targetRange"`
}

type ChampionSpellEffect struct {
	Description string                  `json:"description"`
	Leveling    []ChampionSpellLeveling `json:"leveling"`
}

type ChampionSpellLeveling struct {
	Attribute string                  `json:"attribute"`
	Modifiers []ChampionSpellModifier `json:"modifiers"`
}

type ChampionSpellModifier struct {
	Values string `json:"values"`
	Unit   string `json:"unit"`
}
