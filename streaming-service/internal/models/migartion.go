package models

import (
	"fmt"
	"stream-mesh/streaming/internal/config"

	"gorm.io/gorm"
)

func InitModels(db *gorm.DB, cfg *config.AppConfig) error {
	err := db.AutoMigrate(
		&Video{},
	)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}
