package models

type Rating struct {
	ID    int    `json:"id" example:"1"`
	Long  string `json:"long" gorm:"type:varchar(25)" example:"Observer"`
	Short string `json:"short" gorm:"type:varchar(3)" example:"OBS"`
}
