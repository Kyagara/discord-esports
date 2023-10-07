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

func VALEsportsCommand(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
	respondWithEmbed(interaction.Interaction, []*discordgo.MessageEmbed{createVALMessageEmbed()})
}

func updateVALData() error {
	client.logger.Info("Updating VAL data.")
	valSchedule = make(map[string][]VALEsportsTournamentSchedule)
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
			TeamB:      segment.Team2,
		}

		i, err := strconv.ParseInt(fmt.Sprintf("%v", gameData.Time), 10, 64)
		if err != nil {
			client.logger.Error(fmt.Sprintf("Error parsing unix date: %v", err))
			continue
		}

		date := time.Unix(i, 0)
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
				i, err := strconv.ParseInt(fmt.Sprintf("%v", item.Time), 10, 64)
				if err != nil {
					client.logger.Error(fmt.Sprintf("Error parsing entry time: %v", err))
					continue
				}

				parsedTime := time.Unix(i, 0)

				if parsedTime.Before(now) {
					continue
				}

				if valSchedule[tournament] == nil {
					valSchedule[tournament] = make([]VALEsportsTournamentSchedule, 0)
				}

				valSchedule[tournament] = append(valSchedule[tournament], item)
			}
		}
	}

	client.logger.Info("Updated VAL data.")
	return nil
}

func createVALMessageEmbed() *discordgo.MessageEmbed {
	var fields []*discordgo.MessageEmbedField

	for tournament, games := range valSchedule {
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
		client.logger.Info("No VAL games found.")
		return &discordgo.MessageEmbed{Title: fmt.Sprintf("Valorant games on %v", tomorrow.Format("2006/01/02")), Color: embedColor, Description: "No games found :/"}
	}

	fields = append(fields, &discordgo.MessageEmbedField{Name: "Upcoming matches", Value: "[Check all upcoming matches here](https://www.vlr.gg/matches)"})

	return &discordgo.MessageEmbed{Title: fmt.Sprintf("Valorant games on %v", tomorrow.Format("2006/01/02")), Color: embedColor, Fields: fields}
}

func sendVALEmbed() error {
	client.logger.Info("Sending VAL data.")

	_, err := client.session.ChannelMessageSendEmbed(client.config.VALChannel, createVALMessageEmbed())
	if err != nil {
		return err
	}

	client.logger.Info("Sent VAL data.")
	return nil
}
