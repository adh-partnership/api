package models

type TrainingSessionRequest struct {
	UUIDBase
	UserID uint  `json:"user_id"`
	User   *User `json:"user"`
	// Should be one of: Simulation, Live, OTS, Other -- not enforced
	TrainingType string `json:"training_type"`
	// Must be one of the fields specified in config facility.training.positions array
	TrainingFor string `json:"training_for"`
	Notes       string `json:"notes"`
	// Must be one of: none, open, accepted, completed, cancelled
	Status       string `json:"status"`
	InstructorID *uint  `json:"instructor_id"`
	Instructor   *User  `json:"instructor"`
}
