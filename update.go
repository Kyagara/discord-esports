package main

import (
	"fmt"
	"time"
)

func updateAllData() error {
	now = time.Now()
	tomorrow = time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, now.Location())

	err := updateLOLData()
	if err != nil {
		return fmt.Errorf("error updating LOL data: %v", err)
	}

	err = updateVALData()
	if err != nil {
		return fmt.Errorf("error updating VAL data: %v", err)
	}

	lastUpdate = time.Now()
	return nil
}
