package telegram

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

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

// sendMarkdownV2 sends a message using MarkdownV2 parse mode
func (c *Client) sendMarkdownV2(chatID int64, text string) error {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = tgbotapi.ModeMarkdownV2
	_, err := c.bot.Send(msg)
	return err
}

// DeleteMessage удаляет сообщение
func (c *Client) DeleteMessage(chatID int64, messageID int) error {
	del := tgbotapi.NewDeleteMessage(chatID, messageID)
	_, err := c.bot.Request(del)
	return err
}

// SendMediaGroup отправляет группу фото (файлы с диска). Возвращает ID первого сообщения.
func (c *Client) SendMediaGroup(chatID int64, filePaths []string) error {
	if len(filePaths) == 0 {
		return nil
	}
	var media []interface{}
	for _, p := range filePaths {
		photo := tgbotapi.NewInputMediaPhoto(tgbotapi.FilePath(p))
		media = append(media, photo)
	}
	mg := tgbotapi.NewMediaGroup(chatID, media)
	_, err := c.bot.SendMediaGroup(mg)
	return err
}

// SendMessageWithKeyboard sends a new message (returns sent message)
func (c *Client) SendMessageWithKeyboardFull(chatID int64, text string, buttons [][]tgbotapi.InlineKeyboardButton) (tgbotapi.Message, error) {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = tgbotapi.ModeMarkdown
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(buttons...)
	return c.bot.Send(msg)
}

// DownloadPhoto загружает файл из Telegram по fileID и сохраняет в destPath.
func (c *Client) DownloadPhoto(fileID, destPath string) error {
	file, err := c.bot.GetFile(tgbotapi.FileConfig{FileID: fileID})
	if err != nil {
		return fmt.Errorf("DownloadPhoto GetFile: %w", err)
	}

	return c.downloadFile(file, destPath)
}

// GetFileExt возвращает расширение файла из Telegram по fileID (напр. ".jpg", ".png").
func (c *Client) GetFileExt(fileID string) (string, error) {
	file, err := c.bot.GetFile(tgbotapi.FileConfig{FileID: fileID})
	if err != nil {
		return "", fmt.Errorf("GetFileExt: %w", err)
	}
	return strings.ToLower(filepath.Ext(file.FilePath)), nil
}

func (c *Client) downloadFile(file tgbotapi.File, destPath string) error {
	url := file.Link(c.bot.Token)

	resp, err := http.Get(url) //nolint:gosec
	if err != nil {
		return fmt.Errorf("downloadFile http.Get: %w", err)
	}
	defer resp.Body.Close()

	if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
		return fmt.Errorf("downloadFile mkdir: %w", err)
	}
	f, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("downloadFile create: %w", err)
	}
	defer f.Close()
	_, err = io.Copy(f, resp.Body)
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
	wh, err := tgbotapi.NewWebhook(webhookURL)
	if err != nil {
		return fmt.Errorf("NewWebhook: %w", err)
	}
	_, err = c.bot.Request(wh)
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
