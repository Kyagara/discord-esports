## discord-esports

This bot periodically sends information about upcoming League of Legends and Valorant professional games, its ~~a carbon copy~~ inspired by the [BottyMcBotface](https://github.com/Querijn/BottyMcBotface) bot. It runs in a small Discord server, so it needs work to properly function in a big server.

## Slash Commands

- esports - Get a list of upcoming matches from a game. User needs one of the mod roles set in `./data/config.json` to use the update option.
- info - Send information about the bot, includes link for this page, the last time data was updated and posted and all commands.
- champion - Sends stats for a champion, includes links for a wiki page and LoLalytics page.
- spell - Sends information about a spell, includes links for its wiki page, video, modifiers and notes.

## Todo

- `item` command.
- Cache requests to esports apis in a file to avoid requests when starting up unless its past the post time. (Partially complete)
- Disable button when pressed by a user. (To avoid errors for now, buttons are disabled)
- Fix edge cases with user input, for example giving a invalid champion/spell or not giving one at all

## What data does this bot provide and how?

### Champions and spells (stats, modifiers and notes)

Data from [DataDragon](https://developer.riotgames.com/docs/lol#data-dragon) provided by riot can be sometimes... interesting, take Miss Fortune for example, her ultimate, `Bullet Time`, has apparently a range of 25000, the same range as Ashe's ultimate, a global spell. [CommunityDragon](https://github.com/CommunityDragon/)'s alternative to the DataDragon champion endpoint provides more data, but runs into the same issue in the case of Miss Fortune's ultimate. This leaves the [League of Legends Wiki](https://leagueoflegends.fandom.com/wiki/League_of_Legends_Wiki) as a pretty good source of information.

Using the [lolstaticdata](https://github.com/meraki-analytics/lolstaticdata) tool to gather data from the wiki and normalizing it for this bot's use case (numbers can be string since it will be thrown into Discord Embeds anyway) and the CommunityDragon's CDN to get links for some static data, we can get some pretty good information, nothing ground breaking or that "discards" a visit to the wiki (which is why we keep links to the wiki everywhere).

This data of course comes with its own problems, since its being constantly updated by users, it can have errors, not be formatted properly causing issues when gathering the data and/or normalizing it, but most of the time, its very complete. And since its not an endpoint which can be cached and revalidate after some time, it has to be manually updated, there exists some ideas to allow for updating the data without downtime, but its out of scope for now.

### Esports

For League of Legends, this bot uses the unofficial [LolEsports](https://lolesports.com/) api ([documentation](https://vickz84259.github.io/lolesports-api-docs)), and the unofficial api for [VLR.gg](https://www.vlr.gg/), [vlrggapi](https://github.com/axsddlr/vlrggapi) for Valorant.

## Setup

### lolstaticdata

> You can disable the spell and champion command in the config file and skip this step if you don't need this command.

This bot requires the data provided from [lolstaticdata](https://github.com/meraki-analytics/lolstaticdata). For now, only champion data is needed, so you can just run `python -m lolstaticdata.champions`.

After copying the champions data folder to the root of the project, a champion should have a path like `./data/champions/Aatrox.json`.

Run `go run ./normalize`, this will create a folder with a path for a champion like `./data/champions/normalized/Aatrox.json`.

### Bot Settings

The bot needs the ability of creating slash comamnds and the following permissions:

- Send Messages
- Embed Links

This invite link will have everything the bot needs, replace BOT_ID with the ID of your bot:

https://discord.com/api/oauth2/authorize?client_id=BOT_ID&permissions=18432&scope=bot%20applications.commands

### Config

Duplicate the `./data/config.json.example` and rename the copy to just `config.json`.

You can disable any command, disabling the spell AND champion command will skip loading wiki data.

If mod_roles is empty, anyone will be able to use the `update` option from the `esports` commands, this command makes a request to an esports api, this can be abused.

### Running

Now run the bot with the flag `-register` once, this will **overwrite all guild commands** to the specified guild in the config file, after that you can run the bot without specifying any flag.

## Updating

When a new League of Legends patch is out, just repeat the process of gathering and normalizing the data, please try to not constantly update the data, when a new champion is out for example, check if the wiki has enough information on it to justify gathering all data again.

## Disclaimer

discord-esports isn't endorsed by Riot Games and doesn't reflect the views or opinions of Riot Games or anyone officially involved in producing or managing Riot Games properties. Riot Games, and all associated properties are trademarks or registered trademarks of Riot Games, Inc.
