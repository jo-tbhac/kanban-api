package entity

// CardLabel is model of card_labels table.
type CardLabel struct {
	CardID  uint `gorm:"primary_key;auto_increment:false"`
	LabelID uint `gorm:"primary_key;auto_increment:false"`
}
