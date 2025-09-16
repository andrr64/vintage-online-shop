package models;

type CartDetail struct {
	CartID   uint `gorm:"primaryKey"`
	ProdukID uint `gorm:"primaryKey"`
	Jumlah   int
	Produk   Produk
}