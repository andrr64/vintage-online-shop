package model

// Province merepresentasikan tabel 'provinces'
type Province struct {
	ID   string `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
}

// Regency merepresentasikan tabel 'regencies'
type Regency struct {
	ID         string `json:"id" db:"id"`
	ProvinceID string `json:"province_id" db:"province_id"`
	Name       string `json:"name" db:"name"`
}

// District merepresentasikan tabel 'districts'
type District struct {
	ID        string `json:"id" db:"id"`
	RegencyID string `json:"regency_id" db:"regency_id"`
	Name      string `json:"name" db:"name"`
}
