package models

import (
	"gorm.io/gorm"
)

type Alamat struct {
	gorm.Model
	AlamatLengkap string
	Kel           int
	Kec           int
	Kab           int
	Prov          int
	KodePos       string
}