package services

import (
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	ID        string    `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()" json:"id"`
	Email     string    `gorm:"unique" json:"email"`
	Password  string    `json:"-"`
	IsAdmin   bool      `gorm:"default:false" json:"isAdmin"`
	CreatedAt time.Time `gorm:"datetime:timestamp;default:CURRENT_TIMESTAMP" json:"createdAt"`
	UpdatedAt time.Time `gorm:"datetime:timestamp" json:"updatedAt"`
}

type AdminServices struct {
	User       User
	AdminStore *gorm.DB
}

func NewAdminServices(u User, aStore *gorm.DB) AdminServices {

	return AdminServices{
		User:       u,
		AdminStore: aStore,
	}
}

func (as *AdminServices) CreateAdmin(u User) error {

	// Hash Password & update user
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), 8)
	if err != nil {
		return fmt.Errorf("error hashing password: %s", err)
	}
	u.Password = string(hashedPassword)

	// Create user
	if err := as.AdminStore.Create(&u).Error; err != nil {
		return fmt.Errorf("error creating user: %s", err)
	}

	return nil
}

func (as *AdminServices) LoginAsAdmin(email, password string) (User, error) {
	// Find user (as admin) in DB
	if err := as.AdminStore.
		Where("email = ? AND is_admin = ?", email, true).
		First(&as.User).Error; err != nil {

		return User{}, fmt.Errorf("user not found: %s", err)
	}

	// Compare Passwords
	if err := bcrypt.CompareHashAndPassword([]byte(as.User.Password), []byte(password)); err != nil {

		return User{}, fmt.Errorf("invalid password: %s", err)
	}

	return as.User, nil
}

/* REFERENCES:
https://gorm.io/docs/models.html
https://stackoverflow.com/questions/74965975/retrieve-timestamp-from-postgresql-with-gorm-in-golang

https://stackoverflow.com/questions/24452323/whats-the-difference-between-pointer-and-value-in-struct

https://gorm.io/docs/query.html#String-Conditions

*/
