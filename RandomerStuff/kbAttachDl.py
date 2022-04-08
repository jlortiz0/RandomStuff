#!/usr/bin/python3

import os
import json
import time
import mimetypes
import subprocess

def api_call(msg):
    # Don't open a new terminal window on Windows
    return subprocess.run(("keybase", "chat", "api", "-m", msg), stdout=subprocess.PIPE, creationflags=subprocess.CREATE_NO_WINDOW if os.name == 'nt' else 0).stdout

teams = {}
# I haven't found a good way to load teams by themselves, so this just checks all subscribed channels for team membership
# This does have the benefit of only loading subscribed channels, in case you're in a hypercategorized team
for c in json.loads(api_call('{"method": "list"}'))["result"]["conversations"]:
    # Keybase sometimes doesn't consider #general to have a topic, so we have to account for that to avoid KeyError
    # This does have the side effect of making it impossible to download from #general, so if there's a team that uses it, sorry!
    if c["channel"]["members_type"] == "team" and "topic_name" in c["channel"]:
        if c["channel"]["name"] not in teams:
            teams[c["channel"]["name"]] = []
        teams[c["channel"]["name"]].append(c["channel"]["topic_name"])
for x in teams.values():
    x.sort()

def menu(options, header=""):
    if not options:
        return 0
    while True:
        if os.name == 'nt':
            os.system('cls')
        else:
            os.system('clear')
        if header:
            print(header)
        for i,x in enumerate(options, start=1):
            print(str(i) + ". "+x)
        n = input("\nSelect an option: ")
        if n and n.isnumeric():
            n = int(n)
            if 0 < n < len(options)+1:
                return n-1

while True:
    teamL = list(teams.keys())
    team = menu(teamL+["Exit"], "Select team to download from")
    if team == len(teamL):
        break
    team = teamL[team]
    while True:
        channel = menu(teams[team]+["Back"], "Select channel to download from")
        if channel == len(teams[team]):
            break
        print("\nDownloading message list for #"+teams[team][channel]+" in "+team+"...")
        # Error handling: 0
        msgList = json.loads(api_call('{"method": "read", "params": {"options": {"channel": {"name": "'+team+'", "members_type": "team", "topic_name": "'+teams[team][channel]+'"}, "peek": true}}}'))["result"]["messages"]
        idList = {}
        size = 0
        for x in msgList:
            if x["msg"]["content"]["type"] == "attachment":
                idList[x["msg"]["id"]] = x["msg"]["content"]["attachment"]["object"]["mimeType"]
                size += x["msg"]["content"]["attachment"]["object"]["size"]
        del msgList
        if idList:
            dir = os.path.abspath(team+os.sep+teams[team][channel])
            if size > 1073741824:
                sizeSt = str(round(size/1073741824, 2))+" GiB"
            elif size > 1048576:
                sizeSt = str(round(size/1048576, 2))+" MiB"
            else:
                sizeSt = str(round(size/1024, 2))+" KiB"
            print("Identified "+str(len(idList))+" attachments to download, totalling "+sizeSt+". Files will be saved in "+dir)
            # Yes, I added this after downloading a couple gigs off the wrong channel.
            if input("If this is okay, type 'yes' and press enter: ").lower() == "yes":
                os.makedirs(dir, exist_ok=True)
                sizeSoFar = 0
                for c,x in enumerate(idList):
                    print("Downloading file "+str(c+1)+" of "+str(len(idList))+" ("+str(round(sizeSoFar/size*100, 1))+"%)...\r", end='', flush=True)
                    ext = mimetypes.guess_extension(idList[x])
                    if not ext:
                        ext = '.'+idList[x].split('/')[1]
                    # .jpe? For shame
                    if ext == '.jpe':
                        ext = '.jpg'
                    if not os.path.exists(dir+os.sep+str(x)+ext):
                        api_call('{"method": "download", "params": {"options": {"channel": {"name": "'+team+'", "members_type": "team", "topic_name": "'+teams[team][channel]+'"}, "message_id": '+str(x)+', "output": "'+dir+os.sep+str(x)+ext+'"}}}')
                    try:
                        sizeSoFar += os.path.getsize(dir+os.sep+str(x)+ext)
                    except OSError:
                        pass
                    counter += 1
                print("All files downloaded to "+dir+" sucessfully!")
                input("Press enter to continue.")
                # This does allow for buffering of downloads
                # It's a feature
            else:
                print("Operation cancelled.")
                time.sleep(1.5)
        else:
            print("This channel has no attachments.")
            input("Press enter to go back to the channel list.")
