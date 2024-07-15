package entity

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	ID       uint    `gorm:"primary_key:auto_increment" json:"id"`
	Name     *string `gorm:"type:varchar(50);NULL" json:"name"`
	Username *string `gorm:"type:varchar(50);min=5;unique;NULL" json:"username"`

	Phone          Phone          `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"phone"`
	UserLogins     UserLogins     `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"user_logins"`
	PrivateMessage PrivateMessage `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"private_message"`
	IsActive       bool           `gorm:"type:bool;default:true" json:"is_active"`
	IsAdmin        bool           `gorm:"type:bool;default:false" json:"is_admin"`
}

func (u *User) AfterCreate(db *gorm.DB) error {
	u.UserLogins.UserID = u.ID
	if res := db.Save(u); res != nil {
		return res.Error
	}
	return nil

}

type Phone struct {
	ID      uint   `gorm:"primary_key:auto_increment" json:"id"`
	UserID  uint   `gorm:"uniqueIndex;NOT NULL" json:"user_id"`
	PhoneNo string `gorm:"type:varchar(11);min=11;unique;NOT NULL" json:"phone_no"`
	Token   *int   `gorm:"type:int;min=4,max=4" json:"token"`
	ExpTime *time.Time
}

// statistics

type UserLogins struct {
	ID            uint         `gorm:"primary_key:auto_increment" json:"id"`
	UserID        uint         `gorm:"uniqueIndex;NOT NULL" json:"user_id"`
	NoLogins      uint64       `gorm:"default:0;NOT NULL" json:"no_logins"`
	NOLoginFailed uint64       `gorm:"default:0;NOT NULL" json:"no_login_failed"`
	IP            []UserIP     `gorm:"foreignKey:UserLoginsID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"ip"`
	UserDevice    []UserDevice `gorm:"foreignKey:UserLoginsID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"user_device"`
}

type UserIP struct {
	ID           uint   `gorm:"primary_key:auto_increment" json:"id"`
	UserLoginsID uint   `gorm:"NOT NULL" json:"user_logins"`
	IP           string `gorm:"type:varchar(20);NOT NULL" json:"ip"`
	Date         time.Time
}

type UserDevice struct {
	ID           uint   `gorm:"primary_key:auto_increment" json:"id"`
	UserLoginsID uint   `gorm:"NOT NULL" json:"user_logins"`
	DeviceName   string `gorm:"type:varchar(100);NOT NULL" json:"device_name"`
	IsPhone      bool   `gorm:"type:bool;NOT NULL" json:"is_phone"`
	Browser      string `gorm:"type:varchar(100);NOT NULL" json:"browser"`
	Os           string `gorm:"type:varchar(100);NOT NULL" json:"os"`
}
