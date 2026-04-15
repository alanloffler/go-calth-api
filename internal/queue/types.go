package queue

type BusinessCreatedPayload struct {
	Email        string `json:"email"`
	BusinessName string `json:"businessName"`
	BusinessLink string `json:"businessLink"`
}

type EventCreatedPayload struct {
	Email       string `json:"email"`
	CompanyName string `json:"companyName"`
	FullName    string `json:"fullName"`
	Title       string `json:"title"`
	StartDate   string `json:"startDate"`
}
