package bot

import (
	"context"

	"tgbot-skeleton/internal/ai"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
)

// Handler handles Telegram updates
type Handler struct {
	bot      *tgbotapi.BotAPI
	logger   *zap.Logger
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
		h.sendMessage(chatID, "🤖 Привет! Я AI-бот, готовый помочь вам с любыми вопросами.\n\nПросто напишите мне сообщение, и я отвечу с помощью искусственного интеллекта.\n\nИспользуй /help для получения списка команд.")
	case "help":
		h.sendMessage(chatID, "📋 Доступные команды:\n\n/start - Начать работу с ботом\n/help - Показать это сообщение\n/status - Показать статус бота\n\n💬 Просто напишите любое сообщение, и я отвечу с помощью AI!")
	case "status":
		h.sendMessage(chatID, "✅ AI-бот работает нормально!\n\nВерсия: 1.0.0\nСтатус: Активен\nAI: Подключен")
	default:
		h.sendMessage(chatID, "❓ Неизвестная команда. Используй /help для получения списка команд.")
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

	// Send AI response
	h.sendMessage(chatID, response)
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
