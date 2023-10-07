package main

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func readyEvent(session *discordgo.Session, ready *discordgo.Ready) {
	client.logger.Info(fmt.Sprintf("Logged in as: %v#%v", session.State.User.Username, session.State.User.Discriminator))

	err := updateAllData()
	if err != nil {
		client.logger.Error(fmt.Sprintf("Error updating data: %v", err))
	}

	err = postAllData()
	if err != nil {
		client.logger.Error(fmt.Sprintf("Error sending data: %v", err))
	}
}
