# AI Telegram Bot

A simple AI-powered Telegram bot written in Go that integrates with OpenAI-compatible API providers (like OpenRouter).

## Features

- AI-powered responses using OpenAI-compatible APIs
- Support for OpenRouter and other providers
- **Automatic Markdown to Telegram formatting** - converts AI responses to proper Telegram format
- Long polling and webhook support
- Structured logging with Zap
- Configuration management with Viper
- Graceful shutdown
- Health check endpoint
- Docker support

## Requirements

- Go 1.23+
- Docker and Docker Compose (optional)
- Telegram Bot Token
- AI Provider API Key (OpenRouter, OpenAI, etc.)

## Installation and Setup

### 1. Cloning and Setup

```bash
git clone git@github.com:positron48/universal-ai-bot.git
cd ai-bot
make setup
```

### 2. Bot Configuration

1. Create a bot via [@BotFather](https://t.me/botfather)
2. Get the bot token
3. Get an API key from an AI provider (e.g., [OpenRouter](https://openrouter.ai/))
4. Edit the `.env` file:

```bash
# Required settings
TELEGRAM_TOKEN=your_bot_token_here
AI_URL=https://openrouter.ai/api/v1
AI_API_KEY=your_openrouter_api_key
AI_MODEL=gpt-3.5-turbo
AI_PROMPT=You are a helpful AI assistant. Please respond to the user's message in a helpful and informative way.

# Optional settings
TELEGRAM_DEBUG=false
LOG_LEVEL=info
SERVER_ADDRESS=:8080
```

### 3. Running

#### Local Development

```bash
# Install dependencies
make tidy

# Run in development mode
make dev

# Or build and run
make build
make run
```

#### Docker

```bash
# Build and run with Docker Compose
make docker-build
make docker-run

# View logs
make docker-logs

# Stop
make docker-stop
```

## How it Works

1. User sends any text message to the bot
2. Bot sends a typing indicator
3. Bot forwards the message to the AI provider with the configured system prompt
4. AI provider processes the request and returns a response
5. Bot converts Markdown formatting to Telegram format
6. Bot sends the formatted AI response back to the user

## Supported AI Providers

This bot works with any OpenAI-compatible API, including:

- **OpenRouter** - Access to multiple AI models through one API
- **OpenAI** - Direct OpenAI API access
- **Anthropic** - Claude models (if using OpenRouter)
- **Google** - Gemini models (if using OpenRouter)
- **Local models** - Any self-hosted OpenAI-compatible API

## Commands

The bot supports the following commands:

- `/start` - Start the bot and get welcome message
- `/help` - Show help message with available commands
- `/status` - Show bot status and AI connection info

## Configuration

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `TELEGRAM_TOKEN` | Telegram bot token | **Required** |
| `AI_URL` | AI provider API URL | **Required** |
| `AI_API_KEY` | AI provider API key | **Required** |
| `AI_MODEL` | AI model to use | `gpt-3.5-turbo` |
| `AI_PROMPT` | System prompt for AI | **Required** (or use `AI_PROMPT_FILE`) |
| `AI_PROMPT_FILE` | Path to file containing system prompt | Alternative to `AI_PROMPT` |
| `TELEGRAM_DEBUG` | Debug mode | `false` |
| `TELEGRAM_UPDATES_TIMEOUT` | Updates timeout | `30` |
| `TELEGRAM_WEBHOOK_ENABLE` | Enable webhook | `false` |
| `TELEGRAM_WEBHOOK_DOMAIN` | Webhook domain | - |
| `TELEGRAM_WEBHOOK_PATH` | Webhook path | `/webhook` |
| `SERVER_ADDRESS` | Server address | `:8080` |
| `LOG_LEVEL` | Logging level | `info` |

### Example Configuration for OpenRouter

```env
TELEGRAM_TOKEN=your_telegram_bot_token
AI_URL=https://openrouter.ai/api/v1
AI_API_KEY=your_openrouter_api_key
AI_MODEL=anthropic/claude-3-sonnet
AI_PROMPT=You are a helpful AI assistant. Please respond to the user's message in a helpful and informative way.
```

### Customizing AI Behavior

You can customize the AI behavior in two ways:

#### Option 1: Direct prompt in .env file
Modify the `AI_PROMPT` environment variable to change how the AI responds. You can use multi-line prompts in several ways:

**Single-line prompt:**
```env
AI_PROMPT=You are a friendly customer support assistant. Always be polite and helpful.
```

**Multi-line prompt with escaped newlines:**
```env
AI_PROMPT="You are a helpful AI assistant. You should:\n- Always be polite and respectful\n- Provide accurate information\n- Ask clarifying questions when needed\n- Keep responses concise but informative\n\nPlease respond to the user's message following these guidelines."
```

**Multi-line prompt with single quotes (works in some systems):**
```env
AI_PROMPT='You are a helpful AI assistant. You should:
- Always be polite and respectful
- Provide accurate information
- Ask clarifying questions when needed
- Keep responses concise but informative

Please respond to the user's message following these guidelines.'
```

**Using heredoc in shell scripts:**
```bash
export AI_PROMPT=$(cat <<EOF
You are a helpful AI assistant with the following characteristics:

1. Personality:
   - Friendly and approachable
   - Professional but not formal
   - Patient and understanding

2. Capabilities:
   - Answer questions accurately
   - Provide step-by-step explanations
   - Offer helpful suggestions

3. Guidelines:
   - Always be respectful
   - Admit when you don't know something
   - Ask clarifying questions when needed

Please respond to the user's message following these guidelines.
EOF
)
```

#### Option 2: Load prompt from file (Recommended for long prompts)
Instead of putting the prompt directly in the `.env` file, you can create a separate text file and reference it:

```env
AI_PROMPT_FILE=prompts/english-teacher.txt
```

**Available prompt files:**
- `prompts/simple-assistant.txt` - Basic helpful assistant
- `prompts/customer-support.txt` - Customer support specialist  
- `prompts/english-teacher.txt` - English teacher and translator

**Creating your own prompt file:**
1. Create a new file in the `prompts/` directory
2. Write your prompt in plain text (no need for escaping)
3. Reference it in `.env` with `AI_PROMPT_FILE=prompts/your-file.txt`

Example prompt file (`prompts/my-custom-prompt.txt`):
```
You are a specialized AI assistant for [your domain].

Your main responsibilities:
- Provide accurate information about [topic]
- Help users with [specific tasks]
- Maintain a [tone/style] communication style

Guidelines:
- Always be [specific behavior]
- Never [specific restrictions]
- When in doubt, [fallback behavior]

Remember: [key principles]
```

## Markdown Formatting Support

The bot automatically converts AI responses from Markdown to Telegram format. Supported formatting includes:

### ✅ **Supported Elements:**
- **Headers** (`#`, `##`, `###`) → **Bold text**
- **Bold text** (`**text**`, `__text__`) → **Bold text**
- **Italic text** (`*text*`) → _Italic text_
- **Code blocks** (```language → ```)
- **Inline code** (`` `code` ``)
- **Unordered lists** (`-`, `*`) → • Bullet points
- **Ordered lists** (`1.`, `2.`) → Numbered lists
- **Links** (`[text](url)`) → [text](url)

### 📝 **Example Conversion:**

**Input (AI Response):**
```markdown
# Welcome to AI Bot

This is **bold** and *italic* text.

## Features
- Feature 1
- Feature 2

Use `code` for examples.
```

**Output (Telegram):**
```
**Welcome to AI Bot**

This is **bold** and _italic_ text.

**Features**
• Feature 1
• Feature 2

Use `code` for examples.
```

## Project Structure

```
english-bot/
├── cmd/bot/                 # Application entry point
├── internal/
│   ├── ai/                  # AI service for provider integration
│   ├── bot/                 # Bot logic and handlers
│   ├── config/              # Configuration management
│   ├── logger/              # Logging configuration
│   └── utils/               # Utility functions (Markdown conversion)
├── prompts/                 # AI prompt files
│   ├── simple-assistant.txt
│   ├── customer-support.txt
│   └── english-teacher.txt
├── .github/workflows/       # GitHub Actions
├── Dockerfile               # Docker image
├── docker-compose.yml       # Docker Compose
├── Makefile                 # Development commands
└── README.md               # Documentation
```

## Commands

### Main Commands

```bash
make setup          # Initial project setup
make build          # Build application
make run            # Run application
make dev            # Run in development mode
make test           # Run tests
make lint           # Code linting
make clean          # Clean build artifacts
```

### Docker Commands

```bash
make docker-build   # Build Docker image
make docker-run     # Run with docker-compose
make docker-stop    # Stop containers
make docker-logs    # View logs
make docker-clean   # Clean Docker resources
make deploy         # Deploy with Docker
```

## Docker

### Building Image

```bash
docker build -t ai-telegram-bot .
```

### Running Container

```bash
docker run -d \
  --name ai-telegram-bot \
  -p 8080:8080 \
  -e TELEGRAM_TOKEN=your_token_here \
  -e AI_URL=https://openrouter.ai/api/v1 \
  -e AI_API_KEY=your_api_key \
  -e AI_MODEL=gpt-3.5-turbo \
  -e AI_PROMPT="You are a helpful assistant" \
  ai-telegram-bot
```

### Docker Compose

```bash
# Start
docker-compose up -d

# View logs
docker-compose logs -f

# Stop
docker-compose down
```

## API Endpoints

When running, the bot exposes the following HTTP endpoints:

- `GET /health` - Health check endpoint

## Development

### Adding New Commands

1. Edit `internal/bot/handler.go`
2. Add a new case to the `handleCommand` function
3. Implement the command logic

### Adding New Features

1. Create new packages in `internal/`
2. Add configuration if needed
3. Update the main application

## Testing

```bash
# Run all tests
make test

# Run tests with verbose output
make test-verbose

# Check coverage
go test -cover ./...
```

## Monitoring

The bot provides a health check endpoint:

```bash
curl http://localhost:8080/health
```

## Security

- Uses non-root user in Docker
- Security scanning with Gosec in CI
- Configuration validation on startup
- API keys are loaded from environment variables

## Contributing

1. Fork the project
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is distributed under the most permissive license - the MIT License. This allows maximum freedom for use, modification, and distribution. See the `LICENSE` file for details.

## Support

If you have questions or issues:

1. Check [Issues](https://github.com/your-repo/issues)
2. Create a new Issue with detailed description
3. Make sure to provide logs and configuration