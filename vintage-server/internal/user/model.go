package user

import "time"

type User struct {
	ID             int64  `gorm:"primaryKey;autoIncrement"`
	Username       string `gorm:"size:255;unique;not null"`
	Email          string `gorm:"size:255;unique;not null"`
	PasswordHash   string `gorm:"size:255;not null"`
	ProfilePicture string `gorm:"size:255"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      *time.Time `gorm:"index"`
}

type Address struct {
	ID          int64  `gorm:"primaryKey;autoIncrement"`
	UserID      int64  `gorm:"not null;index"`
	KecamatanID int64  `gorm:"not null;index"`
	Label       string `gorm:"size:50"`
	Recipient   string `gorm:"size:100"`
	PhoneNumber string `gorm:"size:20"`
	FullAddress string `gorm:"type:text"`
	PostalCode  string `gorm:"size:10"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time `gorm:"index"`
}

type Cart struct {
	ID        int64 `gorm:"primaryKey;autoIncrement"`
	UserID    int64 `gorm:"unique;not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `gorm:"index"`
}

type CartDetail struct {
	CartID    int64 `gorm:"primaryKey;autoIncrement:false"`
	VariantID int64 `gorm:"primaryKey;autoIncrement:false"`
	Quantity  int
	AddedAt   time.Time
}

type Wishlist struct {
	UserID    int64 `gorm:"primaryKey;autoIncrement:false"`
	VariantID int64 `gorm:"primaryKey;autoIncrement:false"`
	CreatedAt time.Time
}

type Review struct {
	ID        int64 `gorm:"primaryKey;autoIncrement"`
	UserID    int64 `gorm:"not null;index"`
	ProductID int64 `gorm:"not null;index"`
	Rating    int
	Comment   string `gorm:"type:text"`
	CreatedAt time.Time
}
