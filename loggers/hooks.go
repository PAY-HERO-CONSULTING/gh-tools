package logger

import (
	"context"
	"strings"
	"sync"
	"time"

	"github.com/nikoksr/notify"
	"github.com/nikoksr/notify/service/telegram"
	"github.com/rs/zerolog"
)

// var wg sync.WaitGroup

type TelegramHook struct {
	wg     *sync.WaitGroup
	config *HookConfig
}

func NewTelegramHook(wg *sync.WaitGroup, config *HookConfig) *TelegramHook {

	return &TelegramHook{
		wg:     wg,
		config: config,
	}
}

func (t *TelegramHook) Run(
	e *zerolog.Event,
	level zerolog.Level,
	message string,
) {
	// TODO: Switch to strategy pattern
	if level >= zerolog.WarnLevel {
		t.wg.Add(1)
		go func() {
			_ = notifyTelegram(t.config.Service, message, t.config.ErrorBotToken, t.config.ChannelID)
			t.wg.Done()
		}()
	} else if strings.Contains(message, docsMesaage) {
		t.wg.Add(1)
		go func() {
			_ = notifyTelegram(t.config.Service, message, t.config.ErrorBotToken, t.config.DocsChatID)
			t.wg.Done()
		}()

	} else {
		t.wg.Add(1)
		go func() {
			_ = notifyTelegram(t.config.Service, message, t.config.ErrorBotToken, t.config.UniversalChannelID)
			t.wg.Done()
		}()
	}
}

func notifyTelegram(
	title,
	msg,
	botToken string,
	chatID int64,
) error {

	log.Print("notifying telegram\n")

	telegramService, err := telegram.New(
		botToken,
	)
	if err != nil {
		log.Print(err.Error())
		return err
	}

	telegramService.AddReceivers(chatID)

	notifier := notify.New()

	notifier.UseServices(telegramService)

	ctx, cancel := context.WithTimeout(
		context.Background(),
		30*time.Second,
	)

	defer cancel()

	return notifier.Send(ctx, title, msg)
}
