package models

import (
	"github.com/google/uuid"
	"time"
)

type School struct {
	ID            uint      `json:"id" gorm:"auto_increment"`
	UUID          uuid.UUID `json:"uuid" gorm:"primary_key;type:uuid;default:uuid_generate_v4()"`
	Name          string    `json:"name"`
	KodeProvinsi  string    `json:"kode_provinsi"`
	KodeKabKota   string    `json:"kode_kab_kota"`
	KodeKecamatan string    `json:"kode_kecamatan"`
	NPSN          string    `json:"npsn"`
	Bentuk        string    `json:"bentuk"`
	Status        string    `json:"status"`
	AlamatJalan   string    `json:"alamat_jalan"`
	Lintang       float64   `json:"lintang"`
	Bujur         float64   `json:"bujur"`
	//DeletedAt     time.Time `json:"deleted_at"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
