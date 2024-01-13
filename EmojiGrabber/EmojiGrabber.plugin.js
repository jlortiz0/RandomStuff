/**
 * @name EmojiGrabber
 * @displayName EmojiGrabber
 * @description Add emojis from servers you aren't in to the picker.
 * @author jlortiz
 * @authorId 210556673188823041
 * @version 0.4.4
 */
/*@cc_on
@if (@_jscript)

var shell = WScript.CreateObject("WScript.Shell");
shell.Popup("It looks like you've mistakenly tried to run me directly. That's not how you install plugins. \n(So don't do that!)", 0, "I'm a plugin for BetterDiscord", 0x30);

@else@*/
module.exports = (_ => {
    const changeLog = {
        
    };

    return !window.BDFDB_Global || (!window.BDFDB_Global.loaded && !window.BDFDB_Global.started) ? class {
        constructor (meta) {for (let key in meta) this[key] = meta[key];}
        getName () {return this.name;}
        getAuthor () {return this.author;}
        getVersion () {return this.version;}
        getDescription () {return `The Library Plugin needed for ${this.name} is missing. Open the Plugin Settings to download it. \n\n${this.description}`;}
        
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
                BdApi.showConfirmationModal("Library Missing", `The Library Plugin needed for ${this.name} is missing. Please click "Download Now" to install it.`, {
                    confirmText: "Download Now",
                    cancelText: "Cancel",
                    onCancel: _ => {delete window.BDFDB_Global.downloadModal;},
                    onConfirm: _ => {
                        delete window.BDFDB_Global.downloadModal;
                        this.downloadLibrary();
                    }
                });
            }
            if (!window.BDFDB_Global.pluginQueue.includes(this.name)) window.BDFDB_Global.pluginQueue.push(this.name);
        }
        start () {this.load();}
        stop () {}
        getSettingsPanel () {
            let template = document.createElement("template");
            template.innerHTML = `<div style="color: var(--header-primary); font-size: 16px; font-weight: 300; white-space: pre; line-height: 22px;">The Library Plugin needed for ${this.name} is missing.\nPlease click <a style="font-weight: 500;">Download Now</a> to install it.</div>`;
            template.content.firstElementChild.querySelector("a").addEventListener("click", this.downloadLibrary);
            return template.content.firstElementChild;
        }
    } : (([Plugin, BDFDB]) => {
        let requests, endpoints, persist, EmojiUtils;
            
        return class EmojiGrabber extends Plugin {
            onLoad() {
                this.defaults = {};
            }

            grabEmote(e, save = true) {
                const plugin = this;
                requests.get({
                    url: endpoints.EMOJI_SOURCE_DATA(e)
                }).then(function (result) {
                    const guilds = EmojiUtils.getGuilds();
                    if (guilds.length == 0) {
                        console.log("Can't get constructor!");
                        return;
                    }
                    const fetched = result.body.guild;
                    if (guilds[fetched.id]) {
                        if (save) BDFDB.NotificationUtils.toast("Already have this emote!", {type: "success"});
                        return;
                    }
                    const val = Object.values(guilds)[0];
                    const newG = new val.constructor(fetched.id, val._userId, fetched.emojis);
                    guilds[fetched.id] = newG;
                    persist[fetched.id] = {
                        id: fetched.id,
                        name: fetched.name,
                        emote: e,
                    };
                    if (save) {
                        BDFDB.DataUtils.save(persist, plugin, "emotes");
                        BDFDB.NotificationUtils.toast("Grabbed " + fetched.name, {type: "success"});
                    }
                }, console.log);
            }

            onStart() {
                requests = BDFDB.ModuleUtils.findByProperties("V8APIError");
                endpoints = BDFDB.ModuleUtils.findByProperties("Endpoints").Endpoints;
                EmojiUtils = BDFDB.ModuleUtils.findByProperties("hasUsableEmojiInAnyGuild");
                persist = {};
                const loading = BDFDB.DataUtils.load(this, "emotes");
                Object.values(loading).forEach(x => this.grabEmote(x.emote, false));
                const fgi = BDFDB.ModuleUtils.findByProperties("getFlattenedGuildIds");
                fgi.getFlattenedGuildIds = () => Object.keys(EmojiUtils.getGuilds());
                // HACK: Fucking AWFUL fix to emotes. This will break something.
                const ugs = BDFDB.ModuleUtils.findByProperties("totalUnavailableGuilds");
                Object.defineProperty(ugs, "totalUnavailableGuilds", {
                    get() {
                        return 1;
                    }
                });
            }

            onStop() {
                BDFDB.DataUtils.save(persist, this, "emotes");
            }

            getSettingsPanel (collapseStates = {}) {
                let settingsPanel;
                return settingsPanel = BDFDB.PluginUtils.createSettingsPanel(this, {
                    collapseStates: collapseStates,
                    children: _ => {
                        let settingsItems = [];
                        settingsItems.push(BDFDB.ReactUtils.createElement(BDFDB.LibraryComponents.SettingsItem, {
                            type: "Button",
                            label: "Manage Emotes",
                            onClick: _ => this.showEmoteModal(),
                            children: BDFDB.LanguageUtils.LanguageStrings.EDIT
                        }));
                        
                        return settingsItems;
                    }
                });
            }

            showEmoteModal () {
                let i = 0;

                BDFDB.ModalUtils.open(this, {
                    size: "MEDIUM",
                    header: "Emoji Grabber",
                    subHeader: "",
                    contentClassName: BDFDB.disCN.listscroller,
                    children: Object.values(persist).map(emote => {
                    const j = i;
                    i++;
                    return [
                        (i > 0) && BDFDB.ReactUtils.createElement(BDFDB.LibraryComponents.FormComponents.FormDivider, {
                            className: BDFDB.disCNS.margintop4 + BDFDB.disCN.marginbottom4
                        }),
                        BDFDB.ReactUtils.createElement(BDFDB.LibraryComponents.ListRow, {
                            prefix: BDFDB.ReactUtils.createElement(BDFDB.LibraryComponents.Emoji, {
                                className: BDFDB.DOMUtils.formatClassName(BDFDB.disCN.listavatar, BDFDB.disCN.marginleft8),
                                emojiId: emote.emote,
                                emojiName: ":" + emote.name + ":",
                                size: "jumbo"
                            }),
                            label: BDFDB.ReactUtils.createElement(BDFDB.LibraryComponents.TextScroller, {
                                children: emote.name
                            }),
                            suffix: BDFDB.ReactUtils.createElement(BDFDB.LibraryComponents.Button, {
                                color: BDFDB.LibraryComponents.Button.Colors.RED,
                                onClick: _ => {
                                    document.getElementById("grabbed-emote-li-" + j).remove();
                                    delete persist[emote.id];
                                    delete EmojiUtils.getGuilds()[emote.id];
                                    BDFDB.DataUtils.save(persist, this, "emotes");
                                    BDFDB.NotificationUtils.toast("Removed " + emote.name, {type: "danger"});
                                },
                                children: BDFDB.LanguageUtils.LanguageStrings.DELETE
                            }),
                            id: "grabbed-emote-li-" + j,
                        })
                    ];
                }).flat(10).filter(n => n),
                    buttons: [{
                        contents: BDFDB.LanguageUtils.LanguageStrings.OKAY,
                        color: "BRAND",
                        close: true
                    }]
                });
            }

            onMessageContextMenu (e) {
                let grab = false;
                if (!e.instance.props.message || !e.instance.props.channel || !e.instance.props.target) return;
                const target = e.instance.props.target.tagName == "img" || e.instance.props.target;
                if (target.tagName == "IMG" && target.complete && target.naturalHeight && BDFDB.DOMUtils.containsClass(target, BDFDB.disCN.emojiold, "emote", false)) {
                    grab = true;
                } else {
                    const reaction = BDFDB.DOMUtils.getParent(BDFDB.dotCN.messagereaction, target);
                    if (reaction) {
                        const emoji = reaction.querySelector(BDFDB.dotCN.emojiold);
                        if (emoji) grab = true;
                    }
                }
                if (grab) {
                    let [children, index] = BDFDB.ContextMenuUtils.findItem(e.returnvalue, {id: "reply"});
                    children.splice(index > -1 ? index + 1 : 0, 0, BDFDB.ContextMenuUtils.createItem(BDFDB.LibraryComponents.MenuItems.MenuItem, {
                            label: "Grab Emote",
                            id: BDFDB.ContextMenuUtils.createItemId(this.name, "grabEmote"),
                            action: _ => this.grabEmote(target.getAttribute("data-id")),
                            icon: _ => BDFDB.ReactUtils.createElement(BDFDB.LibraryComponents.Emoji, {
                                className: "emoji",
                                emojiId: target.getAttribute("data-id"),
                                emojiName: target.getAttribute("aria-label"),
                                size: "small"
                            }),
                        })
                    );
                }
            }
        };
    })(window.BDFDB_Global.PluginUtils.buildPlugin(changeLog));
})();
/*@end@*/

if (module.exports.default) {
    module.exports = module.exports.default;
}
if (typeof(module.exports) !== "function") {
    module.exports = eval("EmojiGrabber");
}