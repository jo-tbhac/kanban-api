package entity

// Cover is model of covers table.
type Cover struct {
	CardID uint `json:"card_id" gorm:"primary_key;auto_increment:false"`
	FileID uint `json:"file_id" gorm:"not null"`
}
