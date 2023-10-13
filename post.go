package main

import (
	"fmt"
	"time"
)

func postAllData() error {
	err := sendLOLEmbed()
	if err != nil {
		return fmt.Errorf("error sending LOL embed: %v", err)
	}

	err = sendVALEmbed()
	if err != nil {
		return fmt.Errorf("error sending VAL embed: %v", err)
	}

	lastPost = time.Now()
	return nil
}
