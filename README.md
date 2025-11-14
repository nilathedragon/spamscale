# Spam Scale

Spam Scale is a Telegram bot aiming to block spam bots and known scammers by integrating with existing Telegram moderation tools and permission systems. The goal is to provide an easy-to-install and easy-to-configure bot, that is fun to use for both admins and members.

## Features

- Join Captchas using Telegram's "Approve new members" feature
- FAST known scammer blocklist integration
- In-chat reporting, automatically notifying admins
- Configurable entirely via Telegram commands

## Demo

Demo video is coming soon!

## Usage

To use Spam Scale in your Telegram group, you can either use the public instance ([@SpamScaleBot](https://t.me/SpamScaleBot)) or you can self-host the bot if you would like to customize its appearance or prefer having more control of your data and how it is stored. If you would like to self-host please follow [the installation guide](#Installation) first.

Once you add Spam Scale to your group, it will tell you what to do next.

### Join Captcha

If you would like to make use of the join captcha feature, you will need to configure your Telegram group to use the "Approve new members" feature. For public groups, you can do this under "Manage group" -> "Group type". For private groups, you have to configure this setting when creating an invite link. You may need to revoke the current invite link and generate a new one.

You can then configure the type of captcha new members will need to solve. Simply execute the following command in the chat you would like to configure.

```
/setcaptchatype
```

### FAST - Block known scammers

If you would like to prevent known scammers from joining your chat, you can enable the FAST blocklist within your chat. Simply execute the following command in the chat you would like to configure.

```
/setfastblocklistenabled
```

### In-Chat reporting

In case a scammer or spammer does slip through the configured protection mechanisms, members of your chat can alert all administrators with basic moderation privileges (Delete Messages or Restrict Member). To do this, your members will simply execute the following command.

```
/report
```

## Installation

Spam Scale provides a rolling-release docker container that you can use to easily deploy your own instance. It can be fully configured through environment variables or by mounting in a configuration file.

### Quickstart

If you simply want to run Spam Scale with the default configuration using your own bot credentials, you can use the following command.

```bash
docker run -it -e SPAMSCALE_BOT_TOKEN="<your bot token here>" -v ./spamscale:/mnt ghcr.io/nilathedragon/spamscale
```

### Configuration

All configuration can be done through environment variables.

| Environment Variable     | Type    | Description                                 |
| ------------------------ | ------- | ------------------------------------------- |
| SPAMSCALE_BOT_TOKEN      | String  | Telegram Bot Token                          |
| SPAMSCALE_CACHE_DURATION | Integer | How long, in minutes, items remain in cache |

## Roadmap

- Announcement of successful blocks
  - With the option of petting the bot for the good work!
- Better configuration command
  - Inline-keyboard based settings menu instead of individual commands

### Ideas

- Command to check user in chat against known blocklists
- Accept Chat Rules before accepting join request (After captcha)

## Project Maintainers

- [@drafolin](https://www.github.com/drafolin)
- [@nilathedragon](https://www.github.com/nilathedragon)
