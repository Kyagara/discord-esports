package main

import (
	"fmt"
	"time"
)

func postEsportsData() {
	err := postLOLEsportsEmbed()
	if err != nil {
		client.logger.Error(fmt.Sprintf("error sending LOL embed: %v", err))
	}

	err = postVALEsportsEmbed()
	if err != nil {
		client.logger.Error(fmt.Sprintf("error sending VAL embed: %v", err))
	}

	esports.LastPostTimestamp = time.Now()
}
