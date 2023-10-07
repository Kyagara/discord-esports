package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

type LOLEsportsLeagueResponse struct {
	Data struct {
		Leagues []struct {
			ID              string `json:"id"`
			Slug            string `json:"slug"`
			Name            string `json:"name"`
			Region          string `json:"region"`
			Image           string `json:"image"`
			Priority        int    `json:"priority"`
			DisplayPriority struct {
				Position int    `json:"position"`
				Status   string `json:"status"`
			} `json:"displayPriority"`
		} `json:"leagues"`
	} `json:"data"`
}

type LOLEsportsEventListResponse struct {
	Data struct {
		Esports struct {
			Events []struct {
				StartTime string `json:"startTime"`
				Match     struct {
					Teams []struct {
						Code  string `json:"code"`
						Image string `json:"image"`
					} `json:"teams"`
				} `json:"match"`
				League struct {
					ID   string `json:"id"`
					Slug string `json:"slug"`
					Name string `json:"name"`
				} `json:"league"`
			} `json:"events"`
		} `json:"esports"`
	} `json:"data"`
}

type LOLEsportsLeagueSchedule struct {
	League string
	URL    string
	Time   string
	TeamA  string
	TeamB  string
}

func LOLEsportsCommand(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
	respondWithEmbed(interaction.Interaction, []*discordgo.MessageEmbed{createLOLMessageEmbed()})
}

func updateLOLData() error {
	client.logger.Info("Updating LOL data.")
	lolSchedule = make(map[string][]LOLEsportsLeagueSchedule)
	http := http.DefaultClient

	req, err := newRequest("https://esports-api.lolesports.com/persisted/gw/getLeagues?hl=en-US")
	if err != nil {
		return err
	}
	req.Header.Add("x-api-key", "0TvQnueqKa5mxJntVWt0w4LpLfEkrV1Ta8rQBb9Z")

	res, err := http.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	var leagueList LOLEsportsLeagueResponse

	err = json.NewDecoder(res.Body).Decode(&leagueList)
	if err != nil {
		return err
	}

	var leagueIds []string

	for _, league := range leagueList.Data.Leagues {
		leagueIds = append(leagueIds, league.ID)
	}

	req, err = newRequest("https://esports-api.lolesports.com/persisted/gw/getEventList?hl=en-US&leagueId=" + strings.Join(leagueIds, ","))
	if err != nil {
		return err
	}

	req.Header.Add("x-api-key", "0TvQnueqKa5mxJntVWt0w4LpLfEkrV1Ta8rQBb9Z")

	res, err = http.Do(req)
	if err != nil {
		return err
	}

	var events LOLEsportsEventListResponse

	err = json.NewDecoder(res.Body).Decode(&events)
	if err != nil {
		return err
	}

	var schedule = make(map[string]map[string][]LOLEsportsLeagueSchedule)

	for _, event := range events.Data.Esports.Events {
		gameData := LOLEsportsLeagueSchedule{
			League: event.League.Slug,
			URL:    event.League.Slug,
			Time:   event.StartTime,
			TeamA:  event.Match.Teams[0].Code,
			TeamB:  event.Match.Teams[1].Code,
		}

		dateSplit := strings.Split(event.StartTime, "-")
		realDate := fmt.Sprintf("%v %v %v", dateSplit[0], dateSplit[1], strings.Split(dateSplit[2], "T")[0])

		if schedule[realDate] == nil {
			schedule[realDate] = make(map[string][]LOLEsportsLeagueSchedule)
		}

		if schedule[realDate] != nil && schedule[realDate][event.League.Name] == nil {
			schedule[realDate][event.League.Name] = make([]LOLEsportsLeagueSchedule, 0)
		}

		schedule[realDate][event.League.Name] = append(schedule[realDate][event.League.Name], gameData)
	}

	for dateKey, entries := range schedule {
		parsedTime, err := time.Parse("2006 01 02", dateKey)
		if err != nil {
			client.logger.Error(fmt.Sprintf("Error parsing date key: %v", err))
			continue
		}

		if parsedTime.After(tomorrow) {
			continue
		}

		for league, entryList := range entries {
			for _, item := range entryList {
				parsedTime, err = time.Parse("2006-01-02T15:04:05Z", item.Time)
				if err != nil {
					client.logger.Error(fmt.Sprintf("Error parsing entry time: %v", err))
					continue
				}

				if parsedTime.Before(now) {
					continue
				}

				if lolSchedule[league] == nil {
					lolSchedule[league] = make([]LOLEsportsLeagueSchedule, 0)
				}

				lolSchedule[league] = append(lolSchedule[league], item)
			}
		}
	}

	client.logger.Info("Updated LOL data.")
	return nil
}

func createLOLMessageEmbed() *discordgo.MessageEmbed {
	var fields []*discordgo.MessageEmbedField

	for league, games := range lolSchedule {
		output := ""

		if len(games) == 0 {
			continue
		}

		if len(games) > 10 {
			games = games[:10]
		}

		for _, game := range games {
			date, err := time.Parse("2006-01-02T15:04:05Z", game.Time)
			if err != nil {
				client.logger.Error(fmt.Sprintf("Error parsing game time: %v", err))
				continue
			}

			if date.After(now) {
				output += fmt.Sprintf("%v vs %v, <t:%v:R>.\n", game.TeamA, game.TeamB, date.UnixMilli()/1000)
			}
		}

		if len(strings.TrimSpace(output)) == 0 {
			continue
		}

		fields = append(fields, &discordgo.MessageEmbedField{Name: league, Value: output + fmt.Sprintf("[More about %v here](%v)\n", league, getLOLUrlByLeague(games[0]))})
	}

	if len(fields) == 0 {
		client.logger.Info("No LOL games found.")
		return &discordgo.MessageEmbed{Title: fmt.Sprintf("League games on %v", tomorrow.Format("2006/01/02")), Color: embedColor, Description: "No games found :/"}
	}

	return &discordgo.MessageEmbed{Title: fmt.Sprintf("League games on %v", tomorrow.Format("2006/01/02")), Color: embedColor, Fields: fields}
}

func getLOLUrlByLeague(leagueName LOLEsportsLeagueSchedule) string {
	return "https://lolesports.com/schedule?leagues=" + leagueName.URL
}

func sendLOLEmbed() error {
	client.logger.Info("Sending LOL data.")

	_, err := client.session.ChannelMessageSendEmbed(client.config.LOLChannel, createLOLMessageEmbed())
	if err != nil {
		return err
	}

	client.logger.Info("Sent LOL data.")
	return nil
}
