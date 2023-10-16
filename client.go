package main

import (
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/bwmarrin/discordgo"
	"go.uber.org/zap"
)

type Client struct {
	session  *discordgo.Session
	config   *Configuration
	logger   *zap.Logger
	commands map[string]Command
}

type Command struct {
	Interaction *discordgo.ApplicationCommand
	Handler     func(session *discordgo.Session, interaction *discordgo.InteractionCreate)
}

func newClient() (*Client, error) {
	logger, err := zap.NewProduction()
	if err != nil {
		return nil, fmt.Errorf("error starting logger: %v", err)
	}

	config := Configuration{}
	err = config.loadConfig()
	if err != nil {
		return nil, fmt.Errorf("error loading configuration file: %v", err)
	}

	if len(config.ModRoles) == 0 {
		logger.Info("You have not set any mod_roles, anyone will be able to use the post and update commands, this can be abused.")
	}

	s, err := discordgo.New("Bot " + config.Token)
	if err != nil {
		return nil, fmt.Errorf("error creating Discord session: %v", err)
	}

	s.State.MaxMessageCount = 0
	s.StateEnabled = false
	s.State.TrackPresences = false
	s.State.TrackRoles = false
	s.State.TrackEmojis = false
	s.State.TrackMembers = false
	s.State.TrackVoice = false
	s.State.TrackChannels = false

	client := &Client{
		session:  s,
		config:   &config,
		logger:   logger,
		commands: map[string]Command{},
	}

	client.loadEnabledCommands()

	logger.Info("Client successfully created.")
	return client, nil
}

func (c *Client) loadEnabledCommands() {
	client.commands = make(map[string]Command)

	if c.config.Commands.Esports {
		c.commands["esports"] = Command{Interaction: &discordgo.ApplicationCommand{
			Name:                     "esports",
			Description:              "Send a list of upcoming games from the provided game.",
			DescriptionLocalizations: &map[discordgo.Locale]string{discordgo.PortugueseBR: "Envia uma lista de próximas partidas do jogo escolhido."},
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:                     "game",
					Required:                 true,
					Type:                     discordgo.ApplicationCommandOptionString,
					Description:              "Force an update of the upcoming matches.",
					DescriptionLocalizations: map[discordgo.Locale]string{discordgo.PortugueseBR: "Força uma atualização da lista de partidas."},
					Choices: []*discordgo.ApplicationCommandOptionChoice{
						{Name: "LOL", Value: "lol"},
						{Name: "VAL", Value: "val"},
					},
				},
				{
					Name:                     "update",
					Type:                     discordgo.ApplicationCommandOptionBoolean,
					Description:              "Force an update of the upcoming matches.",
					DescriptionLocalizations: map[discordgo.Locale]string{discordgo.PortugueseBR: "Força uma atualização da lista de partidas."},
					Choices: []*discordgo.ApplicationCommandOptionChoice{
						{Name: "Yes", Value: true},
					},
				},
			}}, Handler: EsportsCommand}
	}

	if c.config.Commands.Info {
		c.commands["info"] = Command{Interaction: &discordgo.ApplicationCommand{
			Name:                     "info",
			Description:              "Send information about the bot.",
			DescriptionLocalizations: &map[discordgo.Locale]string{discordgo.PortugueseBR: "Envia informação sobre o bot."},
		}, Handler: InfoCommand}
	}

	if c.config.Commands.Champion {
		c.commands["champion"] = Command{Interaction: &discordgo.ApplicationCommand{
			Name:                     "champion",
			Description:              "Get information about a League of Legends champion.",
			DescriptionLocalizations: &map[discordgo.Locale]string{discordgo.PortugueseBR: "Envia informações de um champion de League of Legends."},
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:                     "champion",
					Required:                 true,
					Autocomplete:             true,
					Type:                     discordgo.ApplicationCommandOptionString,
					Description:              "Champion name.",
					DescriptionLocalizations: map[discordgo.Locale]string{discordgo.PortugueseBR: "Nome do champion."},
				},
			},
		}, Handler: ChampionCommand}
	}

	if c.config.Commands.Spell {
		c.commands["spell"] = Command{Interaction: &discordgo.ApplicationCommand{
			Name:                     "spell",
			Description:              "Get spells from a League of Legends champion.",
			DescriptionLocalizations: &map[discordgo.Locale]string{discordgo.PortugueseBR: "Envia informações de spells de um de champion League of Legends."},
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:                     "champion",
					Description:              "Champion name.",
					DescriptionLocalizations: map[discordgo.Locale]string{discordgo.PortugueseBR: "Nome do champion."},
					Required:                 true,
					Autocomplete:             true,
					Type:                     discordgo.ApplicationCommandOptionString,
				},
				{
					Name:                     "spell",
					NameLocalizations:        map[discordgo.Locale]string{discordgo.PortugueseBR: "habilidade"},
					Description:              "Spell name.",
					DescriptionLocalizations: map[discordgo.Locale]string{discordgo.PortugueseBR: "Nome da habilidade."},
					Required:                 true,
					Autocomplete:             true,
					Type:                     discordgo.ApplicationCommandOptionString,
				},
				{
					Name:                     "type",
					NameLocalizations:        map[discordgo.Locale]string{discordgo.PortugueseBR: "tipo"},
					Description:              "Choose the type of information you want from the spell.",
					DescriptionLocalizations: map[discordgo.Locale]string{discordgo.PortugueseBR: "Escolha o tipo de informação que você quer dessa habilidade."},
					Type:                     discordgo.ApplicationCommandOptionString,
					Choices: []*discordgo.ApplicationCommandOptionChoice{
						{Name: "Modifiers", Value: "modifiers", NameLocalizations: map[discordgo.Locale]string{discordgo.PortugueseBR: "Modificadores"}},
						{Name: "Notes", Value: "notes", NameLocalizations: map[discordgo.Locale]string{discordgo.PortugueseBR: "Notas"}},
					},
				},
			},
		}, Handler: SpellCommand}
	}
}

func (c *Client) connect() error {
	c.session.AddHandler(interactionsEvent)
	c.session.AddHandler(readyEvent)

	c.session.Identify.Intents = discordgo.IntentMessageContent | discordgo.IntentsGuildMessages
	client = c

	err := c.session.Open()
	if err != nil {
		return fmt.Errorf("error opening Discord session: %v", err)
	}

	return nil
}

func (c *Client) registerCommands() error {
	c.logger.Info("Registering commands.")

	var cmds []*discordgo.ApplicationCommand
	for _, c := range c.commands {
		cmds = append(cmds, c.Interaction)
	}

	cmds, err := c.session.ApplicationCommandBulkOverwrite(c.session.State.User.ID, c.config.GuildID, cmds)
	if err != nil {
		return fmt.Errorf("error registering guild commands: %v", err)
	}

	c.logger.Info(fmt.Sprintf("Registered %v guild commands.", len(cmds)))
	return nil
}

func (c *Client) mainLoop() {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	if !client.config.Commands.Esports {
		<-stop
		c.logger.Info("Shutting down.")
		return
	}

	var post *time.Ticker
	var update *time.Ticker

	// Booleans to avoid a race condition with tickers as both of them immediately start
	firstUpdate := true
	firstPost := true

	resumeUpdate := false
	resumePost := false

	// If first start up
	if esports.LastPostTimestamp.IsZero() || esports.LastUpdateTimestamp.IsZero() {
		updateEsportsData()
		postEsportsData()

		update = getTimer(c.config.UpdateDateTimer)
		post = getTimer(c.config.PostDataTimer)
	} else {
		nextUpdateTime := esports.LastUpdateTimestamp.Add(time.Duration(c.config.UpdateDateTimer * int(time.Millisecond)))
		pastTime, leftForNextUpdate := isPastTime(nextUpdateTime, c.config.UpdateDateTimer)

		if pastTime {
			client.logger.Info("Past time to update, updating now.")
			updateEsportsData()
			update = getTimer(c.config.UpdateDateTimer)
		} else {
			client.logger.Info(fmt.Sprintf("Time to update still valid, waiting for next update in %v.", leftForNextUpdate.String()))
			update = time.NewTicker(leftForNextUpdate)
			resumeUpdate = true
		}

		nextPostTime := esports.LastPostTimestamp.Add(time.Duration(c.config.PostDataTimer * int(time.Millisecond)))
		pastTime, leftForNextPost := isPastTime(nextPostTime, c.config.PostDataTimer)

		if pastTime {
			client.logger.Info("Past time to post, posting now.")
			postEsportsData()
			post = getTimer(c.config.PostDataTimer)
		} else {
			client.logger.Info(fmt.Sprintf("Time to post still valid, waiting for next post in %v.", leftForNextPost.String()))
			post = time.NewTicker(leftForNextPost)
			resumePost = true
		}
	}

	defer update.Stop()
	defer post.Stop()

	for {
		select {
		case <-stop:
			c.logger.Info("Shutting down.")
			saveEsportsFile()
			return

		case <-update.C:
			if resumeUpdate {
				updateEsportsData()
				client.logger.Info("Resetting the Update ticker to specified duration.")
				update.Reset(time.Duration(c.config.UpdateDateTimer * int(time.Millisecond)))
				resumeUpdate = false
				firstUpdate = false
				continue
			}

			if !firstUpdate {
				updateEsportsData()
				firstUpdate = false
			}

		case <-post.C:
			if resumePost {
				postEsportsData()
				client.logger.Info("Resetting the Post ticker to specified duration.")
				post.Reset(time.Duration(c.config.PostDataTimer * int(time.Millisecond)))
				resumePost = false
				firstPost = false
				continue
			}

			if !firstPost {
				postEsportsData()
				firstPost = false
			}
		}
	}
}

func getTimer(duration int) *time.Ticker {
	return time.NewTicker(time.Duration(duration * int(time.Millisecond)))
}

func isPastTime(nextTime time.Time, timerConfiguration int) (bool, time.Duration) {
	timeLeftForNextTick := nextTime.Sub(started)

	// If time left for the next update/post is negative,
	// its way past the time to update/post
	if timeLeftForNextTick < 0 {
		return true, 0
	}

	// Return the time left for the next tick
	timerDuration := time.Duration(timerConfiguration * int(time.Millisecond))
	return timeLeftForNextTick > timerDuration, timeLeftForNextTick
}
