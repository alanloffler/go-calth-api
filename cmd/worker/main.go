package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/alanloffler/go-calth-api/internal/email"
	"github.com/alanloffler/go-calth-api/internal/queue"
	"github.com/hibiken/asynq"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found, using environment variables")
	}

	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "127.0.0.1:6379"
	}

	apiKey := os.Getenv("SENDGRID_API_KEY")
	if apiKey == "" {
		log.Fatal("SENDGRID_API_KEY is required")
	}

	fromEmail := os.Getenv("SENDGRID_FROM_EMAIL")
	fromName := os.Getenv("SENDGRID_FROM_NAME")

	emailSvc := email.NewSendGridService(apiKey, fromEmail, fromName)

	srv := asynq.NewServer(
		asynq.RedisClientOpt{Addr: redisAddr},
		asynq.Config{Concurrency: 10},
	)

	mux := asynq.NewServeMux()
	mux.HandleFunc("email:business_created", handleBusinessCreated(emailSvc))
	mux.HandleFunc("email:event_created", handleEventCreated(emailSvc))

	if err := srv.Run(mux); err != nil {
		log.Fatal(err)
	}
}

func handleBusinessCreated(emailSvc *email.SendGridService) asynq.HandlerFunc {
	return func(ctx context.Context, t *asynq.Task) error {
		// log.Printf("[worker] handler called, payload: %s", string(t.Payload()))

		var payload queue.BusinessCreatedPayload
		if err := json.Unmarshal(t.Payload(), &payload); err != nil {
			// log.Printf("[worker] unmarshal error: %v", err)
			return fmt.Errorf("unmarshal business_created payload: %w", err)
		}

		log.Printf("[worker] calling SendBusinessCreated")
		if err := emailSvc.SendBusinessCreated(payload.Email, payload.BusinessName, payload.BusinessLink); err != nil {
			// log.Printf("[worker] SendBusinessCreated error: %v", err)
			return err
		}

		// log.Printf("[worker] done")

		return nil
	}
}

func handleEventCreated(emailSvc *email.SendGridService) asynq.HandlerFunc {
	return func(ctx context.Context, t *asynq.Task) error {
		var payload queue.EventCreatedPayload
		if err := json.Unmarshal(t.Payload(), &payload); err != nil {
			return fmt.Errorf("unmarshal event_created payload: %w", err)
		}

		if err := emailSvc.SendEventCreated(payload.Email, payload.CompanyName, payload.FullName, payload.Title, payload.StartDate); err != nil {
			return err
		}

		return nil
	}
}
