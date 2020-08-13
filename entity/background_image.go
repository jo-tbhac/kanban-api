package entity

// BackgroundImage is model of background_images table.
type BackgroundImage struct {
	ID    uint   `json:"id"`
	URL   string `json:"url" gorm:"not null"`
	Theme string `json:"theme" gorm:"type:enum('light','dark')"`
}
