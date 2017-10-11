# discord-warframe-orokin-alerts
Notifies Discord channels of Orokin alerts

## Usage:

* first compile (`go build -o dwf-bot`)
* then run (`./dwf-bot -token <TOKEN>`)
  * TOKEN is your discord bot application token

The bot will periodically (each minute) poll for alerts that contain the Orokin Catalyst or Orokin Reactor blueprint reward. When it finds one, it will broadcast a message to your Warframe channels. That is, any channel with the word "warframe", "orokin", or "potato" in the name.
