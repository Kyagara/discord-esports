package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

type VLRGGUpcomingResponse struct {
	Data VLRGGUpcomingData `json:"data"`
}

type VLRGGUpcomingData struct {
	Segments []VLRGGUpcomingSegments `json:"segments"`
}

type VLRGGUpcomingSegments struct {
	Team1         string `json:"team1"`
	Team2         string `json:"team2"`
	UnixTimestamp string `json:"unix_timestamp"`
	Tournament    string `json:"match_event"`
	Series        string `json:"match_series"`
}

type VALEsportsTournamentSchedule struct {
	TeamA      string
	TeamB      string
	Time       time.Time
	Tournament string
	Series     string
	URL        string
}

func updateVALEsportsData() error {
	http := http.DefaultClient
	req, err := newRequest("https://vlrggapi.vercel.app/match?q=upcoming")
	if err != nil {
		return err
	}

	res, err := http.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	var upcoming VLRGGUpcomingResponse
	err = json.NewDecoder(res.Body).Decode(&upcoming)
	if err != nil {
		return err
	}

	var schedule = make(map[time.Time]map[string][]VALEsportsTournamentSchedule, 10)

	for _, segment := range upcoming.Data.Segments {
		time, err := time.Parse("2006-01-02 15:04:05", segment.UnixTimestamp)
		if err != nil {
			return err
		}

		gameData := VALEsportsTournamentSchedule{
			TeamA:      segment.Team1,
			TeamB:      segment.Team2,
			Tournament: segment.Tournament,
			Series:     segment.Series,
			Time:       time,
		}

		if gameData.TeamA == "TBD" && gameData.TeamB == "TBD" {
			continue
		}

		if schedule[time] == nil {
			schedule[time] = make(map[string][]VALEsportsTournamentSchedule)
		}

		tournament := fmt.Sprintf("%v - %v", segment.Tournament, segment.Series)

		if schedule[time] != nil && schedule[time][tournament] == nil {
			schedule[time][tournament] = make([]VALEsportsTournamentSchedule, 0)
		}

		schedule[time][tournament] = append(schedule[time][tournament], gameData)
	}

	tempSchedule := make(map[string][]VALEsportsTournamentSchedule, 0)

	for time, entries := range schedule {
		if time.After(tomorrow) {
			continue
		}

		for tournament, entryList := range entries {
			for _, item := range entryList {
				if time.Before(now) {
					continue
				}

				if tempSchedule[tournament] == nil {
					tempSchedule[tournament] = make([]VALEsportsTournamentSchedule, 0, 10)
				}

				tempSchedule[tournament] = append(tempSchedule[tournament], item)
			}
		}
	}

	esports.VALSchedule = tempSchedule
	client.logger.Info("Updated VAL esports data.")
	return nil
}

func createVALMessageEmbed() *discordgo.MessageEmbed {
	var fields []*discordgo.MessageEmbedField

	for tournament, games := range esports.VALSchedule {
		output := ""

		if len(games) == 0 {
			continue
		}

		if len(games) > 10 {
			games = games[:10]
		}

		for _, game := range games {
			if game.Time.After(now) {
				output += fmt.Sprintf("%v vs %v, <t:%v:R>.\n", game.TeamA, game.TeamB, game.Time.UnixMilli()/1000)
			}
		}

		if len(strings.TrimSpace(output)) == 0 {
			continue
		}

		fields = append(fields, &discordgo.MessageEmbedField{Name: tournament, Value: output})
	}

	if len(fields) == 0 {
		return &discordgo.MessageEmbed{Title: fmt.Sprintf("Valorant games on %v", tomorrow.Format("2006/01/02")), Color: DISCORD_EMBED_COLOR, Description: "No games found :/"}
	}

	fields = append(fields, &discordgo.MessageEmbedField{Name: "Upcoming matches", Value: "[Check all upcoming matches here](https://www.vlr.gg/matches)"})

	return &discordgo.MessageEmbed{Title: fmt.Sprintf("Valorant games on %v", tomorrow.Format("2006/01/02")), Color: DISCORD_EMBED_COLOR, Fields: fields}
}

func postVALEsportsEmbed() error {
	_, err := client.session.ChannelMessageSendEmbed(client.config.VALChannel, createVALMessageEmbed())
	if err != nil {
		return err
	}

	client.logger.Info("Sent VAL data.")
	return nil
}
