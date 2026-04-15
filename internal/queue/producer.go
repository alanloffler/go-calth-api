package queue

import (
	"encoding/json"
	"fmt"

	"github.com/hibiken/asynq"
)

func EnqueueBusinessCreated(client *asynq.Client, payload BusinessCreatedPayload) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal business_created payload: %w", err)
	}

	task := asynq.NewTask("email:business_created", data)

	if _, err = client.Enqueue(task, asynq.MaxRetry(2), asynq.Queue("default")); err != nil {
		return fmt.Errorf("enqueue business_created: %w", err)
	}

	return nil
}

func EnqueueEventCreated(client *asynq.Client, payload EventCreatedPayload) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal event_created payload: %w", err)
	}

	task := asynq.NewTask("email:event_created", data)

	if _, err := client.Enqueue(task, asynq.MaxRetry(2), asynq.Queue("default")); err != nil {
		return fmt.Errorf("enqueue event_created: %w", err)
	}

	return nil
}
