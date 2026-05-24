package telegram

import (
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Client struct {
	bot *tgbotapi.BotAPI
}

func New(token string) (*Client, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, fmt.Errorf("failed to create bot: %w", err)
	}

	bot.Debug = false
	log.Printf("Telegram bot authorized on account %s", bot.Self.UserName)

	return &Client{bot: bot}, nil
}

// SendMessage sends a text message to a user
func (c *Client) SendMessage(chatID int64, text string) error {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = tgbotapi.ModeMarkdown

	_, err := c.bot.Send(msg)
	return err
}

// SendMessageWithKeyboard sends a message with inline keyboard buttons
func (c *Client) SendMessageWithKeyboard(chatID int64, text string, buttons [][]tgbotapi.InlineKeyboardButton) error {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = tgbotapi.ModeMarkdown
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(buttons...)

	_, err := c.bot.Send(msg)
	return err
}

// AnswerCallbackQuery answers inline button clicks
func (c *Client) AnswerCallbackQuery(callbackQueryID, text string, alert bool) error {
	callback := tgbotapi.NewCallback(callbackQueryID, text)
	callback.ShowAlert = alert

	_, err := c.bot.Request(callback)
	return err
}

// EditMessage edits an existing message
func (c *Client) EditMessage(chatID int64, messageID int, text string, buttons [][]tgbotapi.InlineKeyboardButton) error {
	edit := tgbotapi.NewEditMessageText(chatID, messageID, text)
	edit.ParseMode = tgbotapi.ModeMarkdown
	if len(buttons) > 0 {
		edit.ReplyMarkup = &tgbotapi.InlineKeyboardMarkup{InlineKeyboard: buttons}
	}

	_, err := c.bot.Send(edit)
	return err
}

// GetUpdates returns a channel with updates from Telegram
func (c *Client) GetUpdates(timeout int) (tgbotapi.UpdatesChannel, error) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = timeout

	return c.bot.GetUpdatesChan(u), nil
}

// SetWebhook sets webhook URL for receiving updates
func (c *Client) SetWebhook(webhookURL string) error {
	_, err := c.bot.Request(tgbotapi.NewSetWebhook(webhookURL))
	return err
}

// ProcessUpdate handles incoming Telegram update
func (c *Client) ProcessUpdate(update *tgbotapi.Update) {
	if update == nil {
		return
	}

	if update.Message != nil {
		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
	} else if update.CallbackQuery != nil {
		log.Printf("[%s] Callback: %s", update.CallbackQuery.From.UserName, update.CallbackQuery.Data)
	}
}
