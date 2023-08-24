package main

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

func readyEvent(session *discordgo.Session, ready *discordgo.Ready) {
	log.Printf("Logged in as: %v#%v", session.State.User.Username, session.State.User.Discriminator)

	err := updateAllData()
	if err != nil {
		log.Printf("Error updating data: %v", err)
	}

	err = postAllData()
	if err != nil {
		log.Printf("Error sending data: %v", err)
	}
}
