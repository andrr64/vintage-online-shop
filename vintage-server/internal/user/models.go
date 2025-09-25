package user

import (
	"time"
	"gorm.io/gorm"
	"vintage-server/internal/location"
	"vintage-server/internal/order"
	"vintage-server/internal/product"
)

type User struct {
	ID             uint   `gorm:"primaryKey"`
	Username       string `gorm:"uniqueIndex"`
	PasswordHash   string
	Email          string `gorm:"uniqueIndex"`
	ProfilePicture string
	Addresses      []location.Address
	Cart           Cart
	Wishlist       []Wishlist
	Reviews        []product.Review
	Orders         []order.Order
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      gorm.DeletedAt `gorm:"index"`
}

type Cart struct {
	ID        uint `gorm:"primaryKey"`
	UserID    uint `gorm:"uniqueIndex"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Details   []CartDetail
}

type CartDetail struct {
	CartID           uint `gorm:"primaryKey;autoIncrement:false"`
	ProductVariantID uint `gorm:"primaryKey;autoIncrement:false"`
	Quantity         int
	AddedAt          time.Time
}

type Wishlist struct {
	UserID           uint `gorm:"primaryKey;autoIncrement:false"`
	ProductVariantID uint `gorm:"primaryKey;autoIncrement:false"`
	CreatedAt        time.Time
}
