package event

type EventHandler struct {
	repo *EventRepository
}

func NewEventHandler(repo *EventRepository) *EventHandler {
	return &EventHandler{repo: repo}
}
