package bot

import (
	"context"

	"tgbot-skeleton/internal/ai"
	"tgbot-skeleton/internal/utils"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
)

// Handler handles Telegram updates
type Handler struct {
	bot       *tgbotapi.BotAPI
	logger    *zap.Logger
	aiService *ai.Service
}

// NewHandler creates a new handler
func NewHandler(bot *tgbotapi.BotAPI, logger *zap.Logger, aiService *ai.Service) *Handler {
	return &Handler{
		bot:       bot,
		logger:    logger,
		aiService: aiService,
	}
}

// HandleUpdate handles incoming Telegram updates
func (h *Handler) HandleUpdate(ctx context.Context, update tgbotapi.Update) {
	if update.Message == nil {
		return
	}

	message := update.Message
	h.logger.Info("received message",
		zap.Int64("chat_id", message.Chat.ID),
		zap.String("text", message.Text),
		zap.String("username", message.From.UserName),
	)

	// Handle commands
	if message.IsCommand() {
		h.handleCommand(ctx, message)
		return
	}

	// Handle regular messages
	h.handleMessage(ctx, message)
}

// handleCommand handles bot commands
func (h *Handler) handleCommand(ctx context.Context, message *tgbotapi.Message) {
	command := message.Command()
	chatID := message.Chat.ID

	h.logger.Info("handling command", zap.String("command", command))

	switch command {
	case "start":
		h.sendMessage(chatID, "🇬🇧 Привет! Я ваш персональный преподаватель английского языка!\n\n📝 Что я умею:\n• Исправлять ошибки в английском тексте\n• Переводить с русского на английский\n• Создавать карточки слов с объяснениями\n\n💡 Как пользоваться:\n• Отправьте английский текст → получите исправления\n• Отправьте русский текст → получите перевод\n• Отправьте одно слово → получите карточку слова\n\nИспользуйте /help для подробной информации.")
	case "help":
		h.sendMessage(chatID, "📚 Помощь по использованию бота-преподавателя английского:\n\n🔤 **Одно слово** → Карточка слова с:\n• Частью речи\n• Транскрипцией IPA\n• Определением на русском\n• Примерами использования\n• Формами неправильных глаголов\n\n📝 **Английский текст** → Исправления:\n• Поиск ошибок (орфография, грамматика, пунктуация)\n• Подробные объяснения\n• Исправленная версия\n\n🇷🇺 **Русский текст** → Перевод:\n• Естественный перевод на английский\n• Анализ сложных фраз\n• Сохранение тона и стиля\n\n💬 Просто отправьте текст или слово - я сразу помогу!")
	default:
		h.sendMessage(chatID, "❓ Неизвестная команда. Используйте /help для получения информации о возможностях бота.")
	}
}

// handleMessage handles regular text messages
func (h *Handler) handleMessage(ctx context.Context, message *tgbotapi.Message) {
	chatID := message.Chat.ID
	text := message.Text

	if text == "" {
		h.sendMessage(chatID, "Пожалуйста, отправьте текстовое сообщение.")
		return
	}

	h.logger.Info("processing user message",
		zap.Int64("chat_id", chatID),
		zap.String("text", text),
	)

	// Send typing indicator
	h.sendTyping(chatID)

	// Get AI response
	response, err := h.aiService.GenerateResponse(ctx, text)
	if err != nil {
		h.logger.Error("failed to get AI response", zap.Error(err))
		h.sendMessage(chatID, "Извините, произошла ошибка при обработке вашего сообщения. Попробуйте еще раз.")
		return
	}

	// Convert Markdown to Telegram format and send AI response
	telegramResponse := utils.ConvertMarkdownToTelegram(response)
	h.sendMessage(chatID, telegramResponse)
}

// sendMessage sends a message to the specified chat
func (h *Handler) sendMessage(chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = tgbotapi.ModeMarkdown

	if _, err := h.bot.Send(msg); err != nil {
		h.logger.Error("failed to send message", zap.Error(err))
	}
}

// sendTyping sends a typing indicator to the specified chat
func (h *Handler) sendTyping(chatID int64) {
	action := tgbotapi.NewChatAction(chatID, tgbotapi.ChatTyping)
	if _, err := h.bot.Request(action); err != nil {
		h.logger.Error("failed to send typing indicator", zap.Error(err))
	}
}
