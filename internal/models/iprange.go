package models

import (
	"time"
)

type IPRange struct {
	Start     uint32    `gorm:"primaryKey;autoIncrement:false"`
	End       uint32    `gorm:"primaryKey;autoIncrement:false"`
	List      string    `gorm:"type:varchar(100)"`
	Timestamp time.Time `gorm:"not null"`
}
