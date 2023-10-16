package main

import (
	"fmt"
)

type Configuration struct {
	Token           string   `json:"token"`
	GuildID         string   `json:"guild_id"`
	LOLChannel      string   `json:"lol_channel"`
	VALChannel      string   `json:"val_channel"`
	ModRoles        []string `json:"mod_roles"`
	UpdateDateTimer int      `json:"update_data_timer"`
	PostDataTimer   int      `json:"post_data_timer"`
	Commands        struct {
		Esports  bool `json:"esports"`
		Info     bool `json:"info"`
		Champion bool `json:"champion"`
		Spell    bool `json:"spell"`
	} `json:"commands"`
}

func (config *Configuration) loadConfig() error {
	err := loadFile(CONFIG_FILE_PATH, &config)
	if err != nil {
		return fmt.Errorf("error loading config file: %v", err)
	}

	if config.Token == "" {
		return fmt.Errorf("token field not set")
	}

	if config.GuildID == "" {
		return fmt.Errorf("guild_id field not set")
	}

	if config.LOLChannel == "" {
		return fmt.Errorf("lol_channel field not set")
	}

	if config.VALChannel == "" {
		return fmt.Errorf("val_channel field not set")
	}

	if config.UpdateDateTimer < 1800000 {
		return fmt.Errorf("update_data_timer is set too low")
	}

	if config.PostDataTimer < 3600000 {
		return fmt.Errorf("post_data_timer is set too low")
	}

	return nil
}
