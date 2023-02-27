# OpenAI Discord Bot

This bot integrates with OpenAI's API to provide features such as question completion and image generation

## Developing

This repo uses a `devcontainer` configuration, which will configure a development environment for you automatically using Docker. You can also develop locally.

You will need to set the following environment variables:

* `DISCORD_TOKEN` - Token used to interact with Discord
* `CHATGPT_URL` - Base URL for the OpenAI API
* `OPEN_AI_API_KEY` - API key for OpenAI
* `GUILD_ID` - ID of the server you wish to run the bot in