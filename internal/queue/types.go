package queue

type BusinessCreatedPayload struct {
	Email        string `json:"email"`
	BusinessName string `json:"businessName"`
	BusinessLink string `json:"businessLink"`
}
