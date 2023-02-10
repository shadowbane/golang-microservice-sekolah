package models

import (
	"github.com/google/uuid"
	"time"
)

type School struct {
	ID            uint      `json:"id" gorm:"auto_increment;unique"`
	UUID          uuid.UUID `json:"uuid" gorm:"primary_key;type:char(36);not null;unique"`
	Name          string    `json:"name" gorm:"not null"`
	KodeProvinsi  string    `json:"kode_provinsi"`
	KodeKabKota   string    `json:"kode_kab_kota"`
	KodeKecamatan string    `json:"kode_kecamatan"`
	NPSN          string    `json:"npsn"`
	Bentuk        string    `json:"bentuk" gorm:"not null"`
	Status        string    `json:"status" gorm:"not null"`
	AlamatJalan   string    `json:"alamat_jalan"`
	Lintang       float64   `json:"lintang"`
	Bujur         float64   `json:"bujur"`
	CreatedAt     time.Time `json:"created_at" gorm:"type:timestamp"`
	UpdatedAt     time.Time `json:"updated_at" gorm:"type:timestamp"`
}
