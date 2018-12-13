# About

Provides a simple way to recieve notifications of new messages from Discord through Telegram. Once a message is recieved on Discord the username of the sender is sent to Telegram through the Telegram Bot API. Before use of this program, please read the disclaimer [outlined below](#disclaimer). 

The intended audience of this program is of users of AOSP (or one of its derivatives such as LineageOS) who do not have Google Play services. This is because apps like Discord rely on GCM or [FCM](https://firebase.google.com/docs/cloud-messaging/) to deliver push notifications, which is a service provided by Google Play services.

## Similar projects:

- [slack-to-telegram](https://github.com/dan-v/slack-to-telegram)

# Building

To install the dependancies:

```
go get gopkg.in/bwmarrin/discordgo.v0
go get gopkg.in/go-telegram-bot-api/telegram-bot-api.v4
```

To compile:

```
go build main.go
```

# Configuration

The `conf.json` configuration file requires three entries (see `conf.json.example`).

`discordToken` is the session token for Discord. This can be generated either interactively by running this program (do not attempt if two factor authentication is enabled on your account) and entering your username and password, or it may be obtained by [logging into Discord](https://discordapp.com/login) on your browser and finding `token` in `Local Storage` (`SHIFT`+`F9` on Firefox).

`telegramToken` is the token for the Telegram bot, which is provided by creating a bot by following the steps provided in the [Telegram bot API documentation](https://core.telegram.org/bots#3-how-do-i-create-a-bot).

[//]: # (TODO: Make this process interactive, similar to the generation of the Telegram token.)
`telegramChatID` refers to the Telegram user account, for which the ID can be obtained by messaging the Telegram bot `@get_id_bot`.

# Usage

```
./main
```

# Disclaimer

The use of a user account interacting with Discord's API, opposed to that of a bot, is implicitly against Discord's terms of service. Therefore use of this program may result in the **termination of your account**. I take no responsibility in the event this may occur. However, as I have no knowledge of this occuring and that the program interacts to a very minimal extent with the API, the percieved risk is low.
