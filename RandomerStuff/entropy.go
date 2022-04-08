package main

import (
	"math"
    "os"
    "os/signal"
    "fmt"
    "syscall"
    "time"

	"github.com/bwmarrin/discordgo"
)

const GUILD_ID = "822542519669227611"
const USER_ID = "210556673188823041"

// Portions of this program derived from example programs created by Prof. Darrel Long (darrellong) and Eugene Chou (eugenechou)

func countEntrop(self *discordgo.Session, userID string, guildID string) (float64, error) {
    chans, err := self.GuildChannels(guildID)
    if err != nil {
        return math.NaN(), err
    }
    me, err := self.User("@me")
    if err != nil {
        return math.NaN(), err
    }
    myID := me.ID
    var chars uint64
    var entropy float64
    // fmt.Println(len(chans))
    for _, channel := range chans {
        if channel.Type != discordgo.ChannelTypeGuildText {
            continue
        }
        // fmt.Println(channel.Name)
        perms, _ := self.State.UserChannelPermissions(myID, channel.ID)
        if perms & discordgo.PermissionReadMessages == 0 {
            continue
        }
        var lastMsg string
        for {
            toProc, err := self.ChannelMessages(channel.ID, 100, lastMsg, "", "")
            if err != nil {
                return math.NaN(), err
            }
            if len(toProc) == 0 {
                break
            }
            lastMsg = toProc[len(toProc) - 1].ID
            for _, v := range toProc {
                if v.Type != discordgo.MessageTypeDefault && v.Type != discordgo.MessageTypeReply {
                    continue
                }
                if v.Author.ID != userID {
                    continue
                }
                chars += uint64(len(v.Content))
                var counts [256]uint16
                for i := range v.Content {
                    counts[i] += 1
                }
                var msgEn float64
                msgSize := float64(len(v.Content))
                for _, c := range counts {
                    temp := float64(c) / msgSize
                    if temp > 0 {
                        msgEn += temp * math.Log2(temp)
                    }
                }
                entropy += -msgEn * msgSize
            }
        }
    }
	return entropy / float64(chars), err
}

var sc chan os.Signal

func main() {
	f, err := os.Open("key.txt")
	if err != nil {
		panic(err)
	}
	strBytes := make([]byte, 64)
	c, err := f.Read(strBytes)
	f.Close()
	if err != nil {
		panic(err)
	}

	client, err := discordgo.New("Bot " + string(strBytes[:c]))
	if err != nil {
		panic(err)
	}

	client.AddHandlerOnce(ready)
	// client.Identify.Intents = discordgo.IntentsGuilds | discordgo.IntentsGuildMessages
    client.Identify.Intents = discordgo.IntentsAllWithoutPrivileged
	err = client.Open()
	if err != nil {
		panic(err)
	}

	sc = make(chan os.Signal)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
    <-sc
	client.Close()
}

func ready(self *discordgo.Session, _ *discordgo.Ready) {
	time.Sleep(5 * time.Millisecond)
    fmt.Println(countEntrop(self, USER_ID, GUILD_ID))
    sc<-nil
}
