# Telegram Bot Skeleton

Minimal skeleton for creating a Telegram bot in Go with full infrastructure for development, testing, and deployment.

## Features

- Simple architecture with minimal dependencies
- Support for webhook and long polling modes
- Configuration via environment variables
- Structured logging with zap
- Docker and docker-compose for containerization
- CI/CD with GitHub Actions
- Linting and testing
- Makefile for convenient development
- Health check endpoint
- Graceful shutdown

## Requirements

- Go 1.23+
- Docker and Docker Compose (optional)
- Telegram Bot Token

## Installation and Setup

### 1. Cloning and Setup

```bash
git clone git@github.com:positron48/tgbot-skeleton.git
cd tgbot-skeleton
make setup
```

### 2. Bot Configuration

1. Create a bot via [@BotFather](https://t.me/botfather)
2. Get the bot token
3. Edit the `.env` file:

```bash
# Required settings
TELEGRAM_TOKEN=your_bot_token_here

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

## Project Structure

```
tgbot-skeleton/
├── cmd/bot/                 # Application entry point
├── internal/
│   ├── bot/                 # Bot logic
│   ├── config/              # Configuration
│   └── logger/              # Logging
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

## Configuration

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `TELEGRAM_TOKEN` | Telegram bot token | **Required** |
| `TELEGRAM_DEBUG` | Debug mode | `false` |
| `TELEGRAM_UPDATES_TIMEOUT` | Updates timeout | `30` |
| `TELEGRAM_WEBHOOK_ENABLE` | Enable webhook | `false` |
| `TELEGRAM_WEBHOOK_DOMAIN` | Webhook domain | - |
| `TELEGRAM_WEBHOOK_PATH` | Webhook path | `/webhook` |
| `SERVER_ADDRESS` | Server address | `:8080` |
| `LOG_LEVEL` | Logging level | `info` |

### Webhook Mode

To use webhook mode:

1. Set `TELEGRAM_WEBHOOK_ENABLE=true`
2. Specify `TELEGRAM_WEBHOOK_DOMAIN` (e.g., `https://yourdomain.com`)
3. Ensure your server is accessible via HTTPS
4. Restart the bot

## Docker

### Building Image

```bash
docker build -t tgbot-skeleton .
```

### Running Container

```bash
docker run -d \
  --name tgbot-skeleton \
  -p 8080:8080 \
  -e TELEGRAM_TOKEN=your_token_here \
  tgbot-skeleton
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

## CI/CD

The project uses GitHub Actions for continuous integration and deployment:

- **CI**: Runs on every push and pull request
  - Tests
  - Linting with golangci-lint
  - Build verification
- **Release**: Automatic creation of releases with binary files when you push a tag
- **Deploy**: Automatic deployment to server when you push a tag (requires secrets configuration)

### Deployment Commands

```bash
# Deploy from GitHub releases
make deploy

# Update deployed binary
make update

# Check service status
make status

# View service logs
make logs
```

### Docker Deployment

```bash
# Build and run with Docker
make docker-deploy
```

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

## Development

### Adding New Commands

1. Edit `internal/bot/handler.go`
2. Add a new case to the `handleCommand` function
3. Implement the command logic

### Adding Middleware

1. Create a middleware function
2. Add it to the processing chain in `HandleUpdate`

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
