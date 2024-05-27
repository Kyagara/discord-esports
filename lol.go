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
	Data LOLEsportsLeagueData `json:"data"`
}

type LOLEsportsLeagueData struct {
	Leagues []LOLEsportsLeague `json:"leagues"`
}

type LOLEsportsLeague struct {
	ID string `json:"id"`
}

type LOLEsportsEventListResponse struct {
	Data LOLEsportsEventData `json:"data"`
}

type LOLEsportsEventData struct {
	Esports LOLEsportsEventEsports `json:"esports"`
}

type LOLEsportsEventEsports struct {
	Events []LOLEsportsEvent `json:"events"`
}

type LOLEsportsEvent struct {
	StartTime string                `json:"startTime"`
	League    LOLEsportsEventLeague `json:"league"`
	Match     LOLEsportsEventMatch  `json:"match"`
}

type LOLEsportsEventLeague struct {
	Slug string `json:"slug"`
	Name string `json:"name"`
}

type LOLEsportsEventMatch struct {
	Teams []LOLEsportsEventMatchTeam `json:"teams"`
}

type LOLEsportsEventMatchTeam struct {
	Code string `json:"code"`
}

type LOLEsportsLeagueSchedule struct {
	TeamA  string
	TeamB  string
	League string
	Time   string
	URL    string
}

func updateLOLEsportsData() error {
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

	tempSchedule := make(map[string][]LOLEsportsLeagueSchedule, 0)

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

				if tempSchedule[league] == nil {
					tempSchedule[league] = make([]LOLEsportsLeagueSchedule, 0)
				}

				tempSchedule[league] = append(tempSchedule[league], item)
			}
		}
	}

	esports.LOLSchedule = tempSchedule
	client.logger.Info("Updated LOL esports data.")
	return nil
}

func createLOLMessageEmbed() *discordgo.MessageEmbed {
	var fields []*discordgo.MessageEmbedField

	for league, games := range esports.LOLSchedule {
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
				if game.TeamA == "TFT" && game.TeamB == "TFT" {
					output += fmt.Sprintf("<t:%v:R>.\n", date.UnixMilli()/1000)
					continue
				}

				output += fmt.Sprintf("%v vs %v, <t:%v:R>.\n", game.TeamA, game.TeamB, date.UnixMilli()/1000)
			}
		}

		if len(strings.TrimSpace(output)) == 0 {
			continue
		}

		fields = append(fields, &discordgo.MessageEmbedField{Name: league, Value: output + fmt.Sprintf("[More about %v here](%v)\n", league, getLOLUrlByLeague(games[0]))})
	}

	if len(fields) == 0 {
		return &discordgo.MessageEmbed{Title: fmt.Sprintf("League games on %v", tomorrow.Format("2006/01/02")), Color: DISCORD_EMBED_COLOR, Description: "No games found :/"}
	}

	return &discordgo.MessageEmbed{Title: fmt.Sprintf("League games on %v", tomorrow.Format("2006/01/02")), Color: DISCORD_EMBED_COLOR, Fields: fields}
}

func getLOLUrlByLeague(leagueName LOLEsportsLeagueSchedule) string {
	return "https://lolesports.com/schedule?leagues=" + leagueName.URL
}

func postLOLEsportsEmbed() error {
	_, err := client.session.ChannelMessageSendEmbed(client.config.LOLChannel, createLOLMessageEmbed())
	if err != nil {
		return err
	}

	client.logger.Info("Sent LOL data.")
	return nil
}
