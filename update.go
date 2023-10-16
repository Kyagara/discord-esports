package main

import (
	"fmt"
	"time"
)

func updateEsportsData() {
	now = time.Now()
	tomorrow = time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, now.Location())

	err := updateLOLEsportsData()
	if err != nil {
		client.logger.Error(fmt.Sprintf("error updating LOL data: %v", err))
	}

	err = updateVALEsportsData()
	if err != nil {
		client.logger.Error(fmt.Sprintf("error updating VAL data: %v", err))
	}

	esports.LastUpdateTimestamp = time.Now()
	saveEsportsFile()
}
