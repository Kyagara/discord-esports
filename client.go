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
	session *discordgo.Session
	equinox *equinox.Equinox
	config  *Configuration
	logger  *zap.Logger
}

type Configuration struct {
	Token           string   `json:"token"`
	GuildID         string   `json:"guild_id"`
	LOLChannel      string   `json:"lol_channel"`
	VALChannel      string   `json:"val_channel"`
	ModRoles        []string `json:"mod_roles"`
	UpdateDateTimer int      `json:"update_data_timer"`
	PostDataTimer   int      `json:"post_data_timer"`
}

func (c *Configuration) loadConfig() error {
	bytes, err := os.ReadFile("./config.json")
	if err != nil {
		return fmt.Errorf("error reading config file: %v", err)
	}

	err = json.Unmarshal(bytes, &c)
	if err != nil {
		return fmt.Errorf("error parsing config file: %v", err)
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
		session: session,
		equinox: equinox,
		config:  &config,
		logger:  logger,
	}

	logger.Info("Client successfully created.")
	return client, nil
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

	registeredCommands := make([]*discordgo.ApplicationCommand, len(commands))

	for i, cmd := range commands {
		command, err := c.session.ApplicationCommandCreate(c.session.State.User.ID, c.config.GuildID, cmd)
		if err != nil {
			return fmt.Errorf("error registering '%v' command: %v", cmd.Name, err)
		}

		registeredCommands[i] = command
		c.logger.Info(fmt.Sprintf("Registered '%v' command.", command.Name))
	}

	c.logger.Info("Commands registered.")
	return nil
}

func (c *Client) removeCommands() error {
	c.logger.Info("Removing commands")

	registeredCommands, err := c.session.ApplicationCommands(c.session.State.User.ID, c.config.GuildID)
	if err != nil {
		return fmt.Errorf("error fetching registered commands: %v", err)
	}

	for _, cmd := range registeredCommands {
		err := c.session.ApplicationCommandDelete(c.session.State.User.ID, c.config.GuildID, cmd.ID)
		if err != nil {
			return fmt.Errorf("error deleting '%v' command: %v", cmd.Name, err)
		}

		c.logger.Info(fmt.Sprintf("Deleted '%v' command.", cmd.Name))
	}

	registeredCommands, err = c.session.ApplicationCommands(c.session.State.User.ID, "")
	if err != nil {
		return fmt.Errorf("error fetching registered global commands: %v", err)
	}

	for _, cmd := range registeredCommands {
		err := c.session.ApplicationCommandDelete(c.session.State.User.ID, "", cmd.ID)
		if err != nil {
			return fmt.Errorf("error deleting global '%v' command: %v", cmd.Name, err)
		}

		c.logger.Info(fmt.Sprintf("Deleted global '%v' command.", cmd.Name))
	}

	c.logger.Info("Commands removed.")
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
