package entity

// BoardBackgroundImage is model of board_background_images table.
type BoardBackgroundImage struct {
	BoardID           uint `json:"board_id" gorm:"primary_key;auto_increment:false"`
	BackgroundImageID uint `json:"background_image_id" gorm:"not null"`
}
