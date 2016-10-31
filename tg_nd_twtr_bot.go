package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/coreos/pkg/flagutil"
	"github.com/dghubble/go-twitter/twitter"
	"gopkg.in/telegram-bot-api.v4"
	"github.com/dghubble/oauth1"
)

func main() {
	flags := flag.NewFlagSet("user-auth", flag.ExitOnError)
	consumerKey := flags.String("consumer-key", "", "Twitter Consumer Key")
	consumerSecret := flags.String("consumer-secret", "", "Twitter Consumer Secret")
	accessToken := flags.String("access-token", "", "Twitter Access Token")
	accessSecret := flags.String("access-secret", "", "Twitter Access Secret")
	tgtoken := flags.String("tg-token", "", "Telegram token")
	tg_chat_id := flags.Int64("tg-chat-id", 0, "Telegram token")

	flags.Parse(os.Args[1:])
	flagutil.SetFlagsFromEnv(flags, "TWITTER_BOT")

	if *consumerKey == "" || *consumerSecret == "" || *accessToken == "" || *accessSecret == "" {
		log.Fatal("Consumer key/secret and Access token/secret required")
	}

	config := oauth1.NewConfig(*consumerKey, *consumerSecret)
	token := oauth1.NewToken(*accessToken, *accessSecret)
	// OAuth1 http.Client will automatically authorize Requests
	httpClient := config.Client(oauth1.NoContext, token)

	// Twitter Client
	client := twitter.NewClient(httpClient)

	tgBot, err := tgbotapi.NewBotAPI(*tgtoken)
	if err != nil {
		log.Panic(err)
	}

	log.Printf("Authorized on telegram account %s", tgBot.Self.UserName)


	// Convenience Demux demultiplexed stream messages
	demux := twitter.NewSwitchDemux()
	demux.Tweet = func(tweet *twitter.Tweet) {
		log.Print(tweet)
		msg := tweet.User.Name + " "
		msg += "@" + tweet.User.ScreenName
		msg += "\n\n"
		msg += tweet.Text

		tgBot.Send(tgbotapi.NewMessage(*tg_chat_id, msg))
	}
	demux.DM = func(dm *twitter.DirectMessage) {
		log.Print(dm)
		msg := "[DM]"
		msg += dm.Sender.Name + " "
		msg += "@" + dm.Sender.ScreenName
		msg += "\n\n"
		msg += dm.Text

		tgBot.Send(tgbotapi.NewMessage(*tg_chat_id, msg))
	}
	demux.Event = func(event *twitter.Event) {
		log.Print(event)
		msg := event.Source.Name
		msg += "\n"
		msg += event.Event

		tgBot.Send(tgbotapi.NewMessage(*tg_chat_id, msg))
	}

	fmt.Println("Starting Stream...")

	// USER (quick test: auth'd user likes a tweet -> event)
	userParams := &twitter.StreamUserParams{
		StallWarnings: twitter.Bool(true),
		With:          "followings",
	}
	stream, err := client.Streams.User(userParams)
	if err != nil {
		log.Fatal(err)
	}

	// Receive messages until stopped or stream quits
	go demux.HandleChan(stream.Messages)

	// Wait for SIGINT and SIGTERM (HIT CTRL-C)
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	log.Println(<-ch)

	fmt.Println("Stopping Stream...")
	stream.Stop()
}