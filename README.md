## discord-esports

This bot periodically sends information about upcoming League of Legends and Valorant professional games, its ~~a carbon copy~~ inspired by the [BottyMcBotface](https://github.com/Querijn/BottyMcBotface) bot.

This bot uses the Lolesports api, [vlrggapi](https://github.com/axsddlr/vlrggapi), [cdragon](https://github.com/CommunityDragon/) and ddragon.

## Slash Commands

- lol - Send upcoming league games.
- val - Send upcoming val games.
- update - Force all data to update. User needs mod roles set in `config.json`.
- post - Force all posts to be sent again. User needs mod roles set in `config.json`.
- info - Send information about the bot, includes link for this page, the last time data was updated and posted and all commands.
- champion - Sends stats for a champion, includes links for a wiki and LoLalytics page.

## Setup

### Bot Settings

The bot needs the following permissions:

- Send Messages
- Embed Links

This invite link will have everything the bot needs, replace BOT_ID with the ID of your bot:

https://discord.com/api/oauth2/authorize?client_id=BOT_ID&permissions=18432&scope=bot%20applications.commands

### Config

Edit the `config.json.example` and rename it to just `config.json`.

If a guild_id is not provided, the bot will register commands as global commands.

If mod_roles is empty, anyone will be able to use the `post` and `update` commands.

### Running

After building the bot with `go build .`, run the bot with the flag `-register`, this will register all commands to the guild specified in the config file, if you want to remove all commands (including global ones) use the `-remove` flag.

## Disclaimer

discord-esports isn't endorsed by Riot Games and doesn't reflect the views or opinions of Riot Games or anyone officially involved in producing or managing Riot Games properties. Riot Games, and all associated properties are trademarks or registered trademarks of Riot Games, Inc.
