#!/usr/bin/python3

import os
import csv
import subprocess
import music_tag

infile = csv.reader(open("songs.csv"))
next(infile)
for row in infile:
    url, title, artist, track, disc = row
    if not os.path.exists(title + ".m4a"):
        track = int(track)
        disc = int(disc)
        subprocess.run(("yt-dlp", url, "-f", "m4a", "-o", title + ".%(ext)s"), check = True)
        f = music_tag.load_file(title + ".m4a")
        f['title'] = title
        f['artist'] = artist
        f['tracknumber'] = track
        f['discnumber'] = disc
        f['compilation'] = True
        f.save()
        print(title)

input()
