package medical_history

type MedicalHistoryHandler struct {
	repo *MedicalHistoryRepository
}

func NewMedicalHistoryHandler(repo *MedicalHistoryRepository) *MedicalHistoryHandler {
	return &MedicalHistoryHandler{repo: repo}
}
