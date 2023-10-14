## discord-esports

This bot periodically sends information about upcoming League of Legends and Valorant professional games, its ~~a carbon copy~~ inspired by the [BottyMcBotface](https://github.com/Querijn/BottyMcBotface) bot.

This bot uses the Lolesports api, [League of Legends Wiki](https://leagueoflegends.fandom.com/wiki/League_of_Legends_Wiki), [vlrggapi](https://github.com/axsddlr/vlrggapi), [cdragon](https://github.com/CommunityDragon/) and ddragon.

## Slash Commands

- esports - Get a list of upcoming matches from a game. User needs one of the mod roles set in `config.json` to use the update option.
- info - Send information about the bot, includes link for this page, the last time data was updated and posted and all commands.
- champion - Sends stats for a champion, includes links for a wiki page and LoLalytics page.
- spell - Sends information about a spell, includes links for its wiki page, video, modifiers and notes.

## Todo

- `item` command.
- Disable button when pressed by a user.

## Setup

### lolstaticdata

This bot requires the data provided from [lolstaticdata](https://github.com/meraki-analytics/lolstaticdata). For now, only champion data is needed, so just can `python lolstaticdata.champions`.

After copying the champions data folder to the root of the project, a champion should have a path like `./champions/Aatrox.json`.

Run `go run ./normalize`, this will create a folder with a path for a champion like `./champions/normalized/Aatrox.json`.

### Bot Settings

The bot needs the following permissions:

- Send Messages
- Embed Links

This invite link will have everything the bot needs, replace BOT_ID with the ID of your bot:

https://discord.com/api/oauth2/authorize?client_id=BOT_ID&permissions=18432&scope=bot%20applications.commands

### Config

Edit the `config.json.example` and rename it to just `config.json`.

If mod_roles is empty, anyone will be able to use the `update` option from the `esports` commands.

### Running

After building the bot with `go build .`, run the bot with the flag `-register`, this will register all commands to the guild specified in the config file, if you want to remove all commands use the `-remove` flag.

## Disclaimer

discord-esports isn't endorsed by Riot Games and doesn't reflect the views or opinions of Riot Games or anyone officially involved in producing or managing Riot Games properties. Riot Games, and all associated properties are trademarks or registered trademarks of Riot Games, Inc.
