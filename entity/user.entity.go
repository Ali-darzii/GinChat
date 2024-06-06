package entity

import (
	"time"
)

type User struct {
	ID       uint   `gorm:"primary_key:auto_increment" json:"id"`
	Name     string `gorm:"type:varchar(50)" json:"name"`
	Username string `gorm:"type:varchar(50)" json:"username"`
	Phone    Phone  `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

type Phone struct {
	ID      uint   `gorm:"primary_key:auto_increment" json:"id"`
	UserID  uint   `gorm:"uniqueIndex" json:"user_id"`
	PhoneNo *int64 `gorm:"type:int" json:"phone_no"`
	Token   *int64 `gorm:"type:int" json:"token"`
	ExpTime *time.Time
}

func SetToken() {

}
