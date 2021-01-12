package main

import (
	"flag"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/slack-go/slack"
	"strings"
)

var (
	githubAPIToken string
	slackAPIToken string
)
func init() {
	flag.StringVar(&slackAPIToken, "slack-token", "", "Slack API Token")
	flag.StringVar(&githubAPIToken, "github-token", "", "Github API Token")
	flag.Parse()
}
func main() {
	api := slack.New(
		slackAPIToken,
		slack.OptionDebug(false),
		//slack.OptionLog(log.New(os.Stdout, "slack-bot: ", log.Lshortfile|log.LstdFlags)),
	)

	rtm := api.NewRTM()
	go rtm.ManageConnection()

	for msg := range rtm.IncomingEvents {
		//fmt.Print("Event Received: ")
		switch ev := msg.Data.(type) {
		case *slack.HelloEvent:
			// Ignore hello

		case *slack.ConnectedEvent:
			fmt.Println("Infos:", ev.Info)
			fmt.Println("Connection counter:", ev.ConnectionCount)
			// Replace C2147483705 with your Channel ID
			rtm.SendMessage(rtm.NewOutgoingMessage("Hello world", "C01DLLG9JQK"))

		case *slack.MessageEvent:
			fmt.Printf("Message: %v\n", ev)
			if ev.SubType != "bot_message" && strings.Contains(ev.Text,"<@W01BZ4H40EM>"){
				messages, _, _, err := api.GetConversationReplies(&slack.GetConversationRepliesParameters{ChannelID: ev.Channel, Timestamp: ev.ThreadTimestamp})
				if err != nil {
					fmt.Printf("error occured: %v", err)
				}
				fmt.Printf("All the messages: ")
				spew.Dump(messages[0])
				fmt.Printf("single message: ")
				spew.Dump(ev)
				api.PostMessage(ev.Channel,
					slack.MsgOptionText("Hello!", false),
					slack.MsgOptionTS(ev.Timestamp))
			}

		case *slack.PresenceChangeEvent:
			fmt.Printf("Presence Change: %v\n", ev)

		//case *slack.LatencyReport:
		//	fmt.Printf("Current latency: %v\n", ev.Value)

		//case *slack.DesktopNotificationEvent:
		//	fmt.Printf("Desktop Notification: %v\n", ev)


		case *slack.RTMError:
			fmt.Printf("Error: %s\n", ev.Error())

		case *slack.InvalidAuthEvent:
			fmt.Printf("Invalid credentials")
			return

		default:

			// Ignore other events..
			// fmt.Printf("Unexpected: %v\n", msg.Data)
		}
	}
}