package main

import (
	"flag"
	"fmt"
	"github.com/aquilax/truncate"
	"github.com/bwmarrin/discordgo"
	"log"
	"sort"
	"time"
)

// Bot parameters
var (
	GuildID  = flag.String("guild", "", "Test guild ID. If not passed - bot registers commands globally")
	BotToken = flag.String("token", "", "Bot access token")
)

var s *discordgo.Session

func init() { flag.Parse() }

func init() {
	var err error
	s, err = discordgo.New("Bot " + *BotToken)
	if err != nil {
		log.Fatalf("Invalid bot parameters: %v", err)
	}
}

type CacheMessage struct {
	Message   string
	Sent      time.Time
	ID        string
	ChannelID string
}

func main() {
	fmt.Println("Loading...")

	err := s.Open()
	defer s.Close()

	channels, err := s.GuildChannels(*GuildID)
	if err != nil {
		fmt.Println(err)
	}

	// init struct
	var lastMessages []CacheMessage

	//init memory
	lastInChannel := make(map[string]string)

	//init member map
	memberMap := make(map[string]string)

	//get current time zone
	t := time.Now()
	_, zoneOffset := t.Zone()

	for {
		lastMessages = nil
		// load messages into batch, key by time
		for _, channel := range channels {

			channelName := truncate.Truncate("#"+channel.Name, 11, "", truncate.PositionEnd)

			if channel.Type == discordgo.ChannelTypeGuildText {
				messages, err := s.ChannelMessages(channel.ID, 20, "", lastInChannel[channel.ID], "") //, "", afterMessageId, ""
				if err != nil {
					continue //if we error out here, means we can't get all messages. carry on
				}

				for _, message := range messages {
					messageContent := message.Content
					messageTime := message.Timestamp.Add(time.Second * time.Duration(zoneOffset))

					if memberMap[message.Author.ID] == "" {
						member, err := s.GuildMember(*GuildID, message.Author.ID)
						if err == nil {
							memberMap[message.Author.ID] = member.Nick
						} else {
							memberMap[message.Author.ID] = message.Author.Username
						}

						memberMap[message.Author.ID] = truncate.Truncate(memberMap[message.Author.ID], 9, "", truncate.PositionEnd)
					}
					messageAuthor := memberMap[message.Author.ID]

					//break down into string and format
					formattedMessageString := fmt.Sprintf("(%s) - %11s - %9s: %s", messageTime.Format("03:04:05PM"), channelName, messageAuthor, messageContent)

					//place into hash table with date so that it reads date - message
					messageToSend := CacheMessage{
						Message:   formattedMessageString,
						Sent:      messageTime,
						ID:        message.ID,
						ChannelID: channel.ID,
					}

					lastMessages = append(lastMessages, messageToSend)
				}
			}
		}

		//sort
		sort.Slice(lastMessages, func(i, j int) bool {
			return lastMessages[i].Sent.Before(lastMessages[j].Sent)
		})

		//write key batch
		for _, lastMessage := range lastMessages {
			fmt.Println(lastMessage.Message)
			lastInChannel[lastMessage.ChannelID] = lastMessage.ID // cache last message id
		}

		//wait and load again
		time.Sleep(time.Second)
	}
}
