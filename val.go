package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

type VLRGGUpcomingResponse struct {
	Data struct {
		Status   int `json:"status"`
		Segments []struct {
			Team1          string `json:"team1"`
			Team2          string `json:"team2"`
			Flag1          string `json:"flag1"`
			Flag2          string `json:"flag2"`
			Score1         string `json:"score1"`
			Score2         string `json:"score2"`
			TimeUntilMatch string `json:"time_until_match"`
			RoundInfo      string `json:"round_info"`
			TournamentName string `json:"tournament_name"`
			UnixTimestamp  int    `json:"unix_timestamp"`
			MatchPage      string `json:"match_page"`
		} `json:"segments"`
	} `json:"data"`
}

type VALEsportsTournamentSchedule struct {
	Tournament string
	URL        string
	RoundInfo  string
	Time       int
	TeamA      string
	TeamB      string
}

func updateVALEsportsData() error {
	http := http.DefaultClient
	req, err := newRequest("https://vlrggapi.vercel.app/match/upcoming_index")
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

	var schedule = make(map[string]map[string][]VALEsportsTournamentSchedule)

	for _, segment := range upcoming.Data.Segments {
		gameData := VALEsportsTournamentSchedule{
			Tournament: segment.TournamentName,
			RoundInfo:  segment.RoundInfo,
			Time:       segment.UnixTimestamp,
			TeamA:      segment.Team1,
			URL:        segment.MatchPage,
			TeamB:      segment.Team2,
		}

		if gameData.TeamA == "TBD" && gameData.TeamB == "TBD" {
			continue
		}

		date := time.Unix(int64(gameData.Time), 0)
		realDate := fmt.Sprintf("%v %v %v", date.Year(), int(date.Month()), date.Day())

		if schedule[realDate] == nil {
			schedule[realDate] = make(map[string][]VALEsportsTournamentSchedule)
		}

		tournament := fmt.Sprintf("%v - %v", segment.RoundInfo, segment.TournamentName)

		if schedule[realDate] != nil && schedule[realDate][tournament] == nil {
			schedule[realDate][tournament] = make([]VALEsportsTournamentSchedule, 0)
		}

		schedule[realDate][tournament] = append(schedule[realDate][tournament], gameData)
	}

	tempSchedule := make(map[string][]VALEsportsTournamentSchedule, 0)

	for dateKey, entries := range schedule {
		parsedTime, err := time.Parse("2006 1 2", dateKey)
		if err != nil {
			client.logger.Error(fmt.Sprintf("Error parsing date key: %v", err))
			continue
		}

		if parsedTime.After(tomorrow) {
			continue
		}

		for tournament, entryList := range entries {
			for _, item := range entryList {
				parsedTime := time.Unix(int64(item.Time), 0)
				if parsedTime.Before(now) {
					continue
				}

				if tempSchedule[tournament] == nil {
					tempSchedule[tournament] = make([]VALEsportsTournamentSchedule, 0)
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
			i, err := strconv.ParseInt(fmt.Sprintf("%v", game.Time), 10, 64)
			if err != nil {
				client.logger.Error(fmt.Sprintf("Error parsing game time: %v", err))
				continue
			}

			date := time.Unix(i, 0)
			if date.After(now) {
				output += fmt.Sprintf("%v vs %v, <t:%v:R>.\n", game.TeamA, game.TeamB, date.UnixMilli()/1000)
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
