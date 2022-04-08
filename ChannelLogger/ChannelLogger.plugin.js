/**
 * @name ChannelLogger
 * @author jlortiz
 * @version 1.1.2
 * @description Lets you download messages from a channel
 */

module.exports = (_ => {
    const config = {
        "info": {
            "name": "ChannelLogger",
            "author": "jlortiz",
            "version": "1.1.2",
            "description": "Lets you download messages from a channel"
        }
    };

    return !window.BDFDB_Global || (!window.BDFDB_Global.loaded && !window.BDFDB_Global.started) ? class {
        getName () {return config.info.name;}
        getAuthor () {return config.info.author;}
        getVersion () {return config.info.version;}
        getDescription () {return `The Library Plugin needed for ${config.info.name} is missing. Open the Plugin Settings to download it. \n\n${config.info.description}`;}
        
        downloadLibrary () {
            require("request").get("https://mwittrien.github.io/BetterDiscordAddons/Library/0BDFDB.plugin.js", (e, r, b) => {
                if (!e && b && r.statusCode == 200) require("fs").writeFile(require("path").join(BdApi.Plugins.folder, "0BDFDB.plugin.js"), b, _ => BdApi.showToast("Finished downloading BDFDB Library", {type: "success"}));
                else BdApi.alert("Error", "Could not download BDFDB Library Plugin. Try again later or download it manually from GitHub: https://mwittrien.github.io/downloader/?library");
            });
        }
        
        load () {
            if (!window.BDFDB_Global || !Array.isArray(window.BDFDB_Global.pluginQueue)) window.BDFDB_Global = Object.assign({}, window.BDFDB_Global, {pluginQueue: []});
            if (!window.BDFDB_Global.downloadModal) {
                window.BDFDB_Global.downloadModal = true;
                BdApi.showConfirmationModal("Library Missing", `The Library Plugin needed for ${config.info.name} is missing. Please click "Download Now" to install it.`, {
                    confirmText: "Download Now",
                    cancelText: "Cancel",
                    onCancel: _ => {delete window.BDFDB_Global.downloadModal;},
                    onConfirm: _ => {
                        delete window.BDFDB_Global.downloadModal;
                        this.downloadLibrary();
                    }
                });
            }
            if (!window.BDFDB_Global.pluginQueue.includes(config.info.name)) window.BDFDB_Global.pluginQueue.push(config.info.name);
        }
        start () {this.load();}
        stop () {}
        getSettingsPanel () {
            let template = document.createElement("template");
            template.innerHTML = `<div style="color: var(--header-primary); font-size: 16px; font-weight: 300; white-space: pre; line-height: 22px;">The Library Plugin needed for ${config.info.name} is missing.\nPlease click <a style="font-weight: 500;">Download Now</a> to install it.</div>`;
            template.content.firstElementChild.querySelector("a").addEventListener("click", this.downloadLibrary);
            return template.content.firstElementChild;
        }
    } : (([Plugin, BDFDB]) => {
        return class ChannelLogger extends Plugin {
            onLoad() {
                // this.patchedModules = {};
                return;
            }

            onStart() {
                BDFDB.PatchUtils.forceAllUpdates(this);
            }

            onStop() {
                BDFDB.PatchUtils.forceAllUpdates(this);
            }

            onChannelContextMenu(e) {
                if (e.instance.props.channel && e.instance.props.channel.type == BDFDB.DiscordConstants.ChannelTypes.GUILD_TEXT) {
                    let [children, index] = BDFDB.ContextMenuUtils.findItem(e.returnvalue, { id: "channel-notifications", group: true });
                    children.splice(index > -1 ? index + 1 : 0, 0, BDFDB.ContextMenuUtils.createItem(BDFDB.LibraryComponents.MenuItems.MenuGroup, {
                        children: BDFDB.ContextMenuUtils.createItem(BDFDB.LibraryComponents.MenuItems.MenuItem, {
                            label: "Log Messages",
                            id: BDFDB.ContextMenuUtils.createItemId(this.name, "logger"),
                            action: _ => this.downloadMessages(e.instance.props.channel)
                        })
                    }));
                }
            }

            onMessageContextMenu(e) {
                if (e.instance.props.message && e.instance.props.channel) {
                    let [children, index] = BDFDB.ContextMenuUtils.findItem(e.returnvalue, {id: "reply"});
                    children.splice(index > -1 ? index + 1 : 0, 0, BDFDB.ContextMenuUtils.createItem(BDFDB.LibraryComponents.MenuItems.MenuGroup, {
                        children: BDFDB.ContextMenuUtils.createItem(BDFDB.LibraryComponents.MenuItems.MenuItem, {
                            label: "Log After",
                            id: BDFDB.ContextMenuUtils.createItemId(this.name, "logger"),
                            action: _ => this.downloadMessages(e.instance.props.channel, e.instance.props.message.id)
                        })
                    }));
                }
            }

            downloadMessages(channel, nextMsg = "0") {
                const requests = BDFDB.ModuleUtils.findByProperties("V8APIError");
                const endpoints = BDFDB.ModuleUtils.findByProperties("Endpoints").Endpoints;
                let fd = "Discord Text Archive created on ";
                fd += BDFDB.DiscordObjects.Timestamp().format("MMM D YYYY HH:mm");
                fd += " by jlortiz's ChannelLogger\n";
                let nicks = new Map();
                let dayOfYear = 0;
                const promiseExceptHelper = function (result) {
                    if (result.status == 429) {
                        new Promise(r => setTimeout(r, 2000)).then(function() {
                            let promise = requests.get({
                                url: endpoints.MESSAGES(channel.id) + "?after=" + nextMsg + "&limit=100"
                            });
                            promise.then(promiseHelper);
                        });
                    } else {
                        throw result;
                    }
                }
                const promiseHelper = function (result) {
                    if (result.body.length == 0) {
                        let hrefURL = window.URL.createObjectURL(new Blob([fd]));
                        let tempLink = document.createElement("a");
                        tempLink.href = hrefURL;
                        tempLink.download = channel.name+"-"+BDFDB.DiscordObjects.Timestamp.now()+".txt";
                        tempLink.click();
                        window.URL.revokeObjectURL(hrefURL);
                        return;
                    }
                    for (let i = result.body.length - 1; i >= 0; i--) {
                        let v = result.body[i];
                        if (v.type != BDFDB.DiscordConstants.MessageTypes.DEFAULT && v.type != BDFDB.DiscordConstants.MessageTypes.REPLY) {
                            continue;
                        }
                        if (!nicks.has(v.author.id)) {
                            nicks.set(v.author.id, v.author.username);
                            if (v.member && v.member.nick) {
                                nicks.set(v.author.id, v.member.nick)
                            }
                        }
                        let t = BDFDB.DiscordObjects.Timestamp(v.timestamp).local();
                        if (t.dayOfYear() != dayOfYear) {
                            fd += t.format("\\[MMM D\\]") + "\n";
                            dayOfYear = t.dayOfYear();
                        }
                        fd += t.format("\\[HH:mm\\] <");
                        fd += nicks.get(v.author.id) + "> ";
                        let content = v.content;
                        for (const mentioned of v.mentions) {
                            if (nicks.has(mentioned.id)) {
                                content = content.replaceAll("<@" + mentioned.id + ">", "@" + nicks.get(mentioned.id));
                                content = content.replaceAll("<@!" + mentioned.id + ">", "@" + nicks.get(mentioned.id));
                            } else {
                                content = content.replaceAll("<@" + mentioned.id + ">", "@" + mentioned.username);
                                content = content.replaceAll("<@!" + mentioned.id + ">", "@" + mentioned.username);
                            }
                        }
                        fd += content;
                        if (v.pinned) {
                            fd += "\n - Pinned";
                        }
                        for (const attach of v.attachments) {
                            fd += "\n - Attachment: " + attach.url;
                        }
                        for (const embed of v.embeds) {
                            if (embed.video) {
                                fd += "\n - Video: " + embed.video.url;
                            } else if (embed.image) {
                                fd += "\n - Image: " + embed.image.url;
                            } else {
                                fd += "\n - Embed: ";
                                if (embed.url) {
                                    fd += embed.url;
                                } else {
                                    fd += embed.title + "(" + embed.description + ")";
                                }
                            }
                        }
                        fd += "\n";
                    }
                    nextMsg = result.body[0].id;
                    let promise = requests.get({
                        url: endpoints.MESSAGES(channel.id) + "?after=" + nextMsg + "&limit=100"
                    });
                    promise.then(promiseHelper, promiseExceptHelper);
                }
                if (nextMsg != "0") {
                    nextMsg = (nextMsg - 1).toString()
                }
                let promise = requests.get({
                    url: endpoints.MESSAGES(channel.id) + "?after="+nextMsg+"&limit=100"
                });
                promise.then(promiseHelper, promiseExceptHelper);
            }
        };
    })(window.BDFDB_Global.PluginUtils.buildPlugin(config));
})();
