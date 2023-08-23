package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

type BotConfiguration struct {
	Token                string   `json:"token"`
	GuildID              string   `json:"guild_id"`
	LOLChannel           string   `json:"lol_channel"`
	VALChannel           string   `json:"val_channel"`
	ModRoles             []string `json:"mod_roles"`
	DataUpdateTimeout    int      `json:"data_update_timeout"`
	SendToChannelTimeout int      `json:"send_to_channel_timeout"`
}

func loadConfig() error {
	bytes, err := os.ReadFile("./config.json")
	if err != nil {
		return fmt.Errorf("error reading config file: %v", err)
	}

	err = json.Unmarshal(bytes, &config)
	if err != nil {
		return fmt.Errorf("error parsing config file: %v", err)
	}

	if len(config.ModRoles) == 0 {
		log.Print("You have not set any mod_roles, anyone will be able to use the post and update commands.")
	}

	return nil
}
