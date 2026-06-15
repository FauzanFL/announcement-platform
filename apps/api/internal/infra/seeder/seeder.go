package seeder

import (
	"log"

	"gorm.io/gorm"
)

func Run(db *gorm.DB) {
	log.Println("[Seeder] Start seeding data...")
	SeedUser(db)
	log.Println("[Seeder] Finish seeding data...")
}
