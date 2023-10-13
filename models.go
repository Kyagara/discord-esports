package main

import "github.com/bwmarrin/discordgo"

type SpellInfo struct {
	Index    int
	Key      string
	FullName string
}

type SpellEmbeds struct {
	General   discordgo.MessageEmbed
	Modifiers discordgo.MessageEmbed
	Notes     discordgo.MessageEmbed
}

type ChampionEmbeds struct {
	General discordgo.MessageEmbed
}

type WikiChampion struct {
	ID           int    `json:"id"`
	Key          string `json:"key"`
	Name         string `json:"name"`
	Title        string `json:"title"`
	Icon         string `json:"icon"`
	Resource     string `json:"resource"`
	AttackType   string `json:"attackType"`
	AdaptiveType string `json:"adaptiveType"`
	Stats        struct {
		Health                       Stat `json:"health"`
		HealthRegen                  Stat `json:"healthRegen"`
		Mana                         Stat `json:"mana"`
		ManaRegen                    Stat `json:"manaRegen"`
		Armor                        Stat `json:"armor"`
		MagicResistance              Stat `json:"magicResistance"`
		AttackDamage                 Stat `json:"attackDamage"`
		MovementSpeed                Stat `json:"movespeed"`
		CriticalStrikeDamage         Stat `json:"criticalStrikeDamage"`
		CriticalStrikeDamageModifier Stat `json:"criticalStrikeDamageModifier"`
		AttackSpeed                  Stat `json:"attackSpeed"`
		AttackRange                  Stat `json:"attackRange"`
	} `json:"stats"`
	Roles            []string `json:"roles"`
	AttributeRatings struct {
		Damage          int `json:"damage"`
		Toughness       int `json:"toughness"`
		Control         int `json:"control"`
		Mobility        int `json:"mobility"`
		Utility         int `json:"utility"`
		AbilityReliance int `json:"abilityReliance"`
		Difficulty      int `json:"difficulty"`
	} `json:"attributeRatings"`
	Spells struct {
		Passive []WikiSpell `json:"P"`
		Q       []WikiSpell `json:"Q"`
		W       []WikiSpell `json:"W"`
		E       []WikiSpell `json:"E"`
		R       []WikiSpell `json:"R"`
	} `json:"abilities"`
	ReleaseDate      string `json:"releaseDate"`
	ReleasePatch     string `json:"releasePatch"`
	PatchLastChanged string `json:"patchLastChanged"`
	Price            struct {
		BlueEssence int `json:"blueEssence"`
		RP          int `json:"rp"`
		SaleRP      int `json:"saleRp"`
	} `json:"price"`
	Lore  string      `json:"lore"`
	Skins []WikiSkins `json:"skins"`
}

type WikiSkins struct {
	Name         string `json:"name"`
	ID           int    `json:"id"`
	IsBase       bool   `json:"isBase"`
	Availability string `json:"availability"`
	FormatName   string `json:"formatName"`
	LootEligible bool   `json:"lootEligible"`
	Cost         any    `json:"cost"`
	Sale         any    `json:"sale"`
	Distribution string `json:"distribution"`
	Rarity       string `json:"rarity"`
	Chromas      []struct {
		Name         string   `json:"name"`
		ID           int      `json:"id"`
		ChromaPath   string   `json:"chromaPath"`
		Colors       []string `json:"colors"`
		Descriptions []struct {
			Description string `json:"description"`
			Region      string `json:"region"`
		} `json:"descriptions"`
		Rarities []struct {
			Rarity int    `json:"rarity"`
			Region string `json:"region"`
		} `json:"rarities"`
	} `json:"chromas"`
	Lore                  string   `json:"lore"`
	Release               string   `json:"release"`
	Set                   []string `json:"set"`
	SplashPath            string   `json:"splashPath"`
	UncenteredSplashPath  string   `json:"uncenteredSplashPath"`
	TilePath              string   `json:"tilePath"`
	LoadScreenPath        string   `json:"loadScreenPath"`
	LoadScreenVintagePath string   `json:"loadScreenVintagePath"`
	NewEffects            bool     `json:"newEffects"`
	NewAnimations         bool     `json:"newAnimations"`
	NewRecall             bool     `json:"newRecall"`
	NewVoice              bool     `json:"newVoice"`
	NewQuotes             bool     `json:"newQuotes"`
	VoiceActor            []string `json:"voiceActor"`
	SplashArtist          []string `json:"splashArtist"`
}

type Stat struct {
	Flat            float64 `json:"flat"`
	Percent         float64 `json:"percent"`
	PerLevel        float64 `json:"perLevel"`
	PercentPerLevel float64 `json:"percentPerLevel"`
}

type Effect struct {
	Description string     `json:"description"`
	Leveling    []Leveling `json:"leveling"`
}

type Leveling struct {
	Attribute string     `json:"attribute"`
	Modifiers []Modifier `json:"modifiers"`
}

type Modifier struct {
	Values []float64 `json:"values"`
	Units  []string  `json:"units"`
}

type Cost struct {
	Modifiers []Modifier `json:"modifiers"`
}

type Cooldown struct {
	Modifiers     []Modifier `json:"modifiers"`
	AffectedByCDR bool       `json:"affectedByCdr"`
}

type WikiSpell struct {
	Name             string    `json:"name"`
	Icon             string    `json:"icon"`
	Effects          []Effect  `json:"effects"`
	Cost             Cost      `json:"cost"`
	Cooldown         Cooldown  `json:"cooldown"`
	Targeting        string    `json:"targeting"`
	Affects          string    `json:"affects"`
	SpellShieldable  string    `json:"spellshieldable"`
	Resource         string    `json:"resource"`
	DamageType       string    `json:"damageType"`
	SpellEffects     string    `json:"spellEffects"`
	Projectile       string    `json:"projectile"`
	OnHitEffects     string    `json:"onHitEffects"`
	Occurrence       string    `json:"occurrence"`
	Notes            string    `json:"notes"`
	Blurb            string    `json:"blurb"`
	MissileSpeed     string    `json:"missileSpeed"`
	RechargeRate     []float64 `json:"rechargeRate"`
	CollisionRadius  string    `json:"collisionRadius"`
	TetherRadius     string    `json:"tetherRadius"`
	OnTargetCDStatic string    `json:"onTargetCdStatic"`
	InnerRadius      string    `json:"innerRadius"`
	Speed            string    `json:"speed"`
	Width            string    `json:"width"`
	Angle            string    `json:"angle"`
	CastTime         string    `json:"castTime"`
	EffectRadius     string    `json:"effectRadius"`
	TargetRange      string    `json:"targetRange"`
}
