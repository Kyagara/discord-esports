package main

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

func readyEvent(session *discordgo.Session, ready *discordgo.Ready) {
	log.Printf("Logged in as: %v#%v", session.State.User.Username, session.State.User.Discriminator)

	// Since the bot just started, its not an issue to use log.Fatal so the problem can be checked early

	err := updateAllData()
	if err != nil {
		log.Fatal(err)
	}

	err = postAllData()
	if err != nil {
		log.Fatal(err)
	}
}
