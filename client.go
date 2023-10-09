package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/Kyagara/equinox"
	"github.com/bwmarrin/discordgo"
	"go.uber.org/zap"
)

type Client struct {
	session  *discordgo.Session
	equinox  *equinox.Equinox
	config   *Configuration
	logger   *zap.Logger
	commands map[string]Command
}

type Command struct {
	Interaction *discordgo.ApplicationCommand
	Handler     func(session *discordgo.Session, interaction *discordgo.InteractionCreate)
}

type Configuration struct {
	Token           string   `json:"token"`
	GuildID         string   `json:"guild_id"`
	LOLChannel      string   `json:"lol_channel"`
	VALChannel      string   `json:"val_channel"`
	ModRoles        []string `json:"mod_roles"`
	UpdateDateTimer int      `json:"update_data_timer"`
	PostDataTimer   int      `json:"post_data_timer"`
	Commands        struct {
		LOL      bool `json:"lol"`
		VAL      bool `json:"val"`
		Post     bool `json:"post"`
		Update   bool `json:"update"`
		Info     bool `json:"info"`
		Champion bool `json:"champion"`
	} `json:"commands"`
}

func (config *Configuration) loadConfig() error {
	bytes, err := os.ReadFile("./config.json")
	if err != nil {
		return fmt.Errorf("error reading config file: %v", err)
	}

	err = json.Unmarshal(bytes, &config)
	if err != nil {
		return fmt.Errorf("error parsing config file: %v", err)
	}

	if config.Token == "" {
		return fmt.Errorf("token field not set")
	}

	if config.GuildID == "" {
		return fmt.Errorf("guild_id field not set")
	}

	if config.LOLChannel == "" {
		return fmt.Errorf("lol_channel field not set")
	}

	if config.VALChannel == "" {
		return fmt.Errorf("val_channel field not set")
	}

	if config.UpdateDateTimer < 1800000 {
		return fmt.Errorf("update_data_timer is set too low")
	}

	if config.PostDataTimer < 1800000 {
		return fmt.Errorf("post_data_timer is set too low")
	}

	return nil
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

	equinox, err := equinox.NewClient("")
	if err != nil {
		return nil, fmt.Errorf("error starting equinox client: %v", err)
	}

	session, err := discordgo.New("Bot " + config.Token)
	if err != nil {
		return nil, fmt.Errorf("error creating Discord session: %v", err)
	}

	client := &Client{
		session:  session,
		equinox:  equinox,
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

	if c.config.Commands.LOL {
		c.commands["lol"] = Command{Interaction: &discordgo.ApplicationCommand{
			Name:                     "lol",
			Description:              "Send a list of upcoming League of Legends games.",
			DescriptionLocalizations: &map[discordgo.Locale]string{discordgo.PortugueseBR: "Envia uma lista de próximos jogos de League of Legends."},
		}, Handler: LOLEsportsCommand}
	}

	if c.config.Commands.VAL {
		c.commands["val"] = Command{Interaction: &discordgo.ApplicationCommand{
			Name:                     "val",
			Description:              "Send a list of upcoming Valorant games.",
			DescriptionLocalizations: &map[discordgo.Locale]string{discordgo.PortugueseBR: "Envia uma lista de próximos jogos de Valorant."},
		}, Handler: VALEsportsCommand}
	}

	if c.config.Commands.Update {
		c.commands["update"] = Command{Interaction: &discordgo.ApplicationCommand{
			Name:                     "update",
			Description:              "Force all data to update.",
			DescriptionLocalizations: &map[discordgo.Locale]string{discordgo.PortugueseBR: "Força todos os dados a serem atualizados."},
		}, Handler: UpdateCommand}
	}

	if c.config.Commands.Info {
		c.commands["info"] = Command{Interaction: &discordgo.ApplicationCommand{
			Name:                     "info",
			Description:              "Send information about the bot.",
			DescriptionLocalizations: &map[discordgo.Locale]string{discordgo.PortugueseBR: "Envia informação sobre o bot."},
		}, Handler: InfoCommand}
	}

	if c.config.Commands.Post {
		c.commands["post"] = Command{Interaction: &discordgo.ApplicationCommand{
			Name:                     "post",
			Description:              "Force all data to be sent again.",
			DescriptionLocalizations: &map[discordgo.Locale]string{discordgo.PortugueseBR: "Força todos os dados a serem enviados."},
		}, Handler: PostCommand}
	}

	if c.config.Commands.Champion {
		c.commands["champion"] = Command{Interaction: &discordgo.ApplicationCommand{
			Name:                     "champion",
			Description:              "Get a League of Legends champion stats.",
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

func (c *Client) removeCommands() error {
	c.logger.Info("Removing commands")

	_, err := c.session.ApplicationCommandBulkOverwrite(c.session.State.User.ID, c.config.GuildID, make([]*discordgo.ApplicationCommand, 0))
	if err != nil {
		return fmt.Errorf("error removing guild commands: %v", err)
	}

	c.logger.Info("Removed guild commands.")

	return nil
}

func (c *Client) mainLoop() {
	post := time.NewTicker(time.Duration(c.config.PostDataTimer * int(time.Millisecond)))
	defer post.Stop()

	update := time.NewTicker(time.Duration(c.config.UpdateDateTimer * int(time.Millisecond)))
	defer update.Stop()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	// The first update/send is done in the ready event, there would be a race condition if it was done by the ticker
	firstUpdate := true
	firstSend := true

loop:
	for {
		select {
		case <-stop:
			c.logger.Info("Shutting down.")
			break loop
		case <-update.C:
			if !firstUpdate {
				err := updateAllData()
				if err != nil {
					c.logger.Error(fmt.Sprintf("error updating data: %v", err))
				}
			}

			firstUpdate = false

		case <-post.C:
			if !firstSend {
				err := postAllData()
				if err != nil {
					c.logger.Error(fmt.Sprintf("error sending data: %v", err))
				}
			}

			firstSend = false
		}
	}
}
