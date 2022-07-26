package main

import (
	"fmt"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"
)

type Config struct {
	Env          string `env:"ENV" envDefault:"dev"`
	AppToken     string `env:"SLACK_APP_TOKEN,required"`
	BotToken     string `env:"SLACK_BOT_TOKEN,required"`
	EmojiChannel string `env:"EMOJI_CHANNEL" envDefault:"#general"`
}

func (c Config) isProdLike() bool {
	return c.Env != "dev"
}

type Bot struct {
	conf         Config
	apiClient    *slack.Client
	socketClient *socketmode.Client
}

func NewBot(config Config) *Bot {
	if !strings.HasPrefix(config.EmojiChannel, "#") {
		config.EmojiChannel = "#" + config.EmojiChannel
	}
	api := slack.New(
		config.BotToken,
		slack.OptionAppLevelToken(config.AppToken),
		slack.OptionDebug(!config.isProdLike()),
		slack.OptionLog(logger{}),
	)

	client := socketmode.New(
		api,
		socketmode.OptionDebug(!config.isProdLike()),
		socketmode.OptionLog(logger{}),
	)
	return &Bot{
		conf:         config,
		apiClient:    api,
		socketClient: client,
	}
}

func (b *Bot) Run() error {
	go b.handleIncomingEvents()
	return b.socketClient.Run()
}

func (b *Bot) handleIncomingEvents() {
	for evt := range b.socketClient.Events {
		switch evt.Type {
		case socketmode.EventTypeConnecting:
			log.Info().Msgf("Connecting to Slack with Socket Mode...")
		case socketmode.EventTypeConnectionError:
			log.Warn().Msgf("Connection failed. Retrying later...")
		case socketmode.EventTypeConnected:
			log.Info().Msgf("Connected to Slack with Socket Mode.")
		case socketmode.EventTypeEventsAPI:
			event, ok := evt.Data.(slackevents.EventsAPIEvent)
			if !ok {
				continue
			}

			go b.handleEventsApiEvent(event)
			b.socketClient.Ack(*evt.Request)
		case socketmode.EventTypeInteractive:
			log.Info().Msgf("socketmode.EventTypeInteractive")
		case socketmode.EventTypeSlashCommand:
			log.Info().Msgf("socketmode.EventTypeSlashCommand")
		default:
			log.Info().Msgf("unknown type: %+v", evt.Type)
		}
	}
}

func (b *Bot) handleEventsApiEvent(apiEvent slackevents.EventsAPIEvent) {
	switch apiEvent.Type {
	case slackevents.CallbackEvent:
		innerEvent := apiEvent.InnerEvent
		switch ev := innerEvent.Data.(type) {
		case *slackevents.AppMentionEvent:
			b.handleAppMentionEvent(ev)
		case *slackevents.EmojiChangedEvent:
			b.handleEmojiChangedEvent(ev)
		}
	default:
		log.Debug().Msgf("unsupported Events API event received: %T", apiEvent.Type)
	}
}

func (b *Bot) handleAppMentionEvent(event *slackevents.AppMentionEvent) {
	// todo fill this out more
	if _, _, err := b.apiClient.PostMessage(
		event.Channel,
		slack.MsgOptionText("Yes, hello.", false),
	); err != nil {
		log.Error().Err(err).Msgf("failed posting message")
	}
}

func (b *Bot) handleEmojiChangedEvent(event *slackevents.EmojiChangedEvent) {
	log.Info().
		Str("event_type", "emoji_changed").
		Str("subtype", event.Subtype).
		Msgf("handling emoji_changed event")

	var message []slack.MsgOption
	switch event.Subtype {
	case "remove":
		if len(event.Names) == 0 {
			return
		}
		names := mapSlice(event.Names, func(t string) string {
			return fmt.Sprintf(":%s:", t)
		})
		// Header Section
		headerSection := slack.NewHeaderBlock(
			slack.NewTextBlockObject(slack.PlainTextType, "Emojis Removed", false, false),
		)
		// main message
		mainMsg := slack.NewTextBlockObject(
			slack.MarkdownType,
			fmt.Sprintf("`%s`", strings.Join(names, " — ")),
			false,
			false,
		)
		fieldsSection := slack.NewSectionBlock(nil, []*slack.TextBlockObject{mainMsg}, nil)

		message = append(message,
			slack.MsgOptionText("Emojis Removed", false),
			slack.MsgOptionBlocks(
				headerSection,
				slack.NewDividerBlock(),
				fieldsSection,
			),
		)
	case "add":
		headerText := slack.NewTextBlockObject(slack.PlainTextType, "New emoji added! :wave:", false, false)
		headerSection := slack.NewHeaderBlock(headerText)

		// main message
		mainMsg := slack.NewTextBlockObject(
			slack.MarkdownType,
			fmt.Sprintf(":%s: — `:%s:`", event.Name, event.Name),
			false,
			false,
		)
		fieldsSection := slack.NewSectionBlock(nil, []*slack.TextBlockObject{mainMsg}, nil)

		message = append(message,
			slack.MsgOptionText("New emoji added! :wave:", false),
			slack.MsgOptionBlocks(
				headerSection,
				slack.NewDividerBlock(),
				fieldsSection,
			),
		)
	default:
		return // shrug
	}

	if _, _, err := b.apiClient.PostMessage(b.conf.EmojiChannel, message...); err != nil {
		log.Error().Err(err).Msgf("failed to publish message to #general")
	}
}

func mapSlice[T any, M any](a []T, f func(T) M) []M {
	n := make([]M, len(a))
	for i, e := range a {
		n[i] = f(e)
	}
	return n
}
