# Insider Case - Go Backend Project

This is an Auto Message Sender Project sending 2 messages within 2 minutes long periods. These values are configurable through environment variables.

## Project Structure

```
├── cmd/
│   └── rest/
│       └── main.go                 # App entry point
├── internal/
│   ├── config/
│   │   └── config.go              # Global config management
│   ├── container/
│   │   └── container.go           # Dependency injection area
│   ├── handler/
│   │   └── message_handler.go     # Top layer interact with user through HTTP
│   ├── repositories/
│   │   ├── interface.go           # Common interface for repositories
│   │   └── message_repository.go  # Repository layer communicating with data store
│   └── services/
│       └── message_service.go     # Service layer contains all business logic
├── pkg/
│   ├── database/
│   │   └── db.go                  # DB package for connecting PostgreSQL db through Gorm
│   ├── logger/
│   │   └── logger.go              # A minimal logger package
│   ├── redis/
│   │   └── redis.go               # Redis client
│   ├── sms_client/
│   │   └── sms_client.go          # SMS service client
│   └── ticker/
│       └── ticker.go              # Ticker package for doing some works periodically
├── go.mod
└── go.sum
```

## API Endpoints

- `GET /api/v1/messages` - Lists all sent messages
- `GET /api/v1/messages/start` - Starts auto sending if it stopped
- `GET /api/v1/messages/stop` - Stops auto sending if it started

You can see them `/swagger` page as well.

## Configurations

You can use .env file to configure db, redis, server etc.

## How to execute

```bash
# Clone the project
git clone git@github.com:frkntplglu/insider.git

# Create db with default value
docker-compose up -d

# Load all dependencies
go mod tidy

# Execute the app
go run cmd/rest/main.go
```

## Extras

I tried to make the project as testable and modular as possible. Additionally, since I couldn’t consistently get different responses from webhook.site, I used the `x-request-id` parameter from the header as the messageId. 

This project fetches records from the database where the status is `pending`, ordered so that the most recently added ones come first. From each record, it takes the phone number and message, then sends an SMS. If the message is sent successfully, it updates the status to `sent` and stores it in Redis along with the `messageId`. If it fails, the status is updated to `failed`, but nothing is written to the cache.

Since the system is designed to send only 2 messages every 2 minutes, I didn’t see the need to make the sending process concurrent. If there were more messages to handle, I would build a concurrent structure. For the ones marked as `failed`, we could also set up a separate retry mechanism to process them again.