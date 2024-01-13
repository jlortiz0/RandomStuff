# AgeCounter

A BetterDiscord plugin that allows you to add emotes from servers you aren't in to the emoji picker. This does not allow you to use them; you will need DiscordFreeEmojis. Nitro has not been tested.

## Usage

This plugin requires that you have [BetterDiscord](https://betterdiscord.app) and [BDFDB](https://betterdiscord.app/plugin/BDFDB).

Why did I use BDFDB? Well, I needed a way to add things to context menus and have a dynamic settings menu and it seemed handy. If mwittrien changes that functionality then ¯\\\_(ツ)\_/¯.

Anyway, download the plugin to the plugins folder and enable it in BD's settings. To grab an emote, you must see it in the wild. Right click on it and select "Grab Emote". It should be added to your emote picker below the last server. All emotes from the server will be shown, not just the one you picked.

To remove a server, go into Plugin settings and click the Remove buttong next to the emote you wish to remove from the picker. All emotes from the same server will be removed.

## Known Bugs

This plugin may stop working if you join or leave a server while it is active. Disable and re-enable the plugin to get it working again.

Some users have reported favorites not working properly when grabbed emotes are used. I have had difficulty replicating this issue, so I cannot offer advice.

The plugin directly modifies some Discord functions related to unavailable servers. While this does not appear to have adverse effects at time of writing, this may change at any time.

Grabbed emotes will always appear in a category called "Custom", which will push all lower emote categories down by one. If you frequently use the buttons on the left to scroll to default emotes, this will break that.

The functionality of this plugin on messages with multiple distinct emotes has not been fully tested.

## Fun Fact

I wanted to put the grab button in the popout that appears when you click on an external emoji, but I couldn't figure out how to inject into it because it's an unexported functional component. 

## DISCLAIMER

I am not responsible for your use of this plugin and personally believe that you should not use it. This plugin has bugs which may result in reduced functionality of your Discord client and will not be maintained. Some people may not be comfortable with your use of their private emotes.

Parts of this plugin are derived from [ServerHider](https://github.com/mwittrien/BetterDiscordAddons/blob/master/Plugins/ServerHider).
