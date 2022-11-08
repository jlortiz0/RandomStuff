# AgeCounter

A BetterDiscord plugin that lists the account age of everyone* in a server.

\*in servers with 250 or more members, only online users will be counted

## Usage

This plugin requires that you have [BetterDiscord](https://betterdiscord.app) and [BDFDB](https://betterdiscord.app/plugin/BDFDB).

Why did I use BDFDB? Well, I needed a way to add things to context menus and it seemed handy. If mwittrien changes that functionality then ¯\\\_(ツ)\_/¯.

Anyway, download the plugin to the plugins folder and enable it in BD's settings. To generate the list, right click on a server and select "Count Ages" You will be asked where to save the output file.

The output is in CSV format. The first line is a comment. Some CSV readers ignore this, some don't, if yours doesn't then just delete the line in a text editor or something. Timestamps are in `dd-mm-yyyy` format.

## DISCLAIMER

All information given by this plugin is available through Discord's interfaces, albeit in a different form; this is merely an aggregation tool. This plugin does not perform any API calls.

Parts of this plugin are derived from [PersonalPins](https://github.com/mwittrien/BetterDiscordAddons/tree/master/Plugins/PersonalPins).
