# Dashy
## Read-only log for Discord servers

Dashy is a read-only Discord client that lives in the terminal. Because of how the Discord API is restricted, it must be 'installed' as a bot in every server you want to use it in.

### Setup
Log in to the Discord developer portal. Create a new bot. The only permission you need to enable is the server members intent. Add the bot to the server (guild) you want to view messages in. Compile the go file. Enter the guild ID number and the bot token as command line arguments like so:

`./dashy -guild GUILD_ID -token BOT_TOKEN`

