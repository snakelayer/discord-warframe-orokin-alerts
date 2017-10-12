# discord-warframe-orokin-alerts
Notifies Discord channels of Orokin alerts.

Uses data from [WFCD/warframe-worldstate-data](https://github.com/WFCD/warframe-worldstate-data).

## Usage:

### Adding a bot to your server

* Create: Go to [Discord - My Apps](https://discordapp.com/developers/applications/me) and create a new app. Turn it into an app bot user. Note the **Client ID** under **App Details**.
* Invite: Visit `https://discordapp.com/oauth2/authorize?scope=bot&permissions=6144&client_id=ClientID` to invite the bot to your server, replacing **ClientID** with the value from the previous step. The bot needs `SEND_MESSAGES` and optionally `SEND_TTS_MESSAGES` permissions on your server, which the link should grant.

### Running the bot

* first compile (`go build -o dwf-bot`)
* then run (`./dwf-bot -token <TOKEN>`)
  * TOKEN is from the app bot user that you created earlier.

The bot will periodically (each minute) poll for alerts that contain the Orokin Catalyst or Orokin Reactor blueprint reward. When it finds one, it will broadcast a message to your Warframe channels. That is, any channel with the word "warframe", "orokin", or "potato" in the name.
