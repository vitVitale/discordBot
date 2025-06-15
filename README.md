# Discord Bot With ChatGPT

This is a minimal example of a Discord bot that forwards user messages to OpenAI and
posts the generated reply back. Update `Token` and `ApiKey` in `app/server.go` with
your credentials or load them from environment variables before running.

Build and run:

```bash
cd app
go build -o bot
./bot
```

The bot now uses the chat completion endpoint with the `gpt-3.5-turbo` model.
