package models

type User struct {
	Id string `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
}
