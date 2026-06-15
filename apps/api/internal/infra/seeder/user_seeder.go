package seeder

import (
	"announcement-api/internal/domain/entity"
	"log"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func SeedUser(db *gorm.DB) {
	users := []entity.User{
		{
			Username: "admin",
			Password: "Admin123!",
			Role:     entity.RoleAdmin,
		},
		{
			Username: "user",
			Password: "User123!",
			Role:     entity.RoleUser,
		},
	}

	for _, user := range users {

		hashed, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			log.Printf("[Seeder] Failed to hash password for user %s: %v", user.Username, err)
			return
		}

		user.Password = string(hashed)

		if err := db.Where(entity.User{Username: user.Username}).FirstOrCreate(&user).Error; err != nil {
			log.Printf("[Seeder] Failed to create user %s: %v", user.Username, err)
		} else {
			log.Printf("[Seeder] User synchronized: %s", user.Username)
		}
	}
}
