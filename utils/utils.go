package utils

import (
	"GinChat/db"
	"GinChat/entity"
	"github.com/gin-gonic/gin"
	"github.com/mssola/user_agent"
	"gorm.io/gorm"
	"math/rand/v2"
	"time"
)

func GetClientIP(request *gin.Context) {

}

func GetExpiryTime() time.Time {
	const expTime int8 = 60
	return time.Now().Add(time.Second * time.Duration(expTime))

}

func SmsTokenGenerate() int {
	return rand.IntN(8999) + 1000
}

func UserLoggedIn(request *gin.Context, user entity.User) error {
	var postDb *gorm.DB = db.ConnectPostgres()
	var userLogins entity.UserLogins
	if res := postDb.Where("user_id = ?", user.ID).Take(&userLogins); res.Error != nil {
		return res.Error
	}
	userLogins.NoLogins += 1
	if res := postDb.Save(&userLogins); res.Error != nil {
		return res.Error
	}

	var userIP entity.UserIP
	userIP.UserLoginsID = userLogins.ID
	userIP.IP = request.ClientIP()
	if res := postDb.Save(&userIP); res.Error != nil {
		return res.Error
	}

	userAgent := request.GetHeader("User-Agent")
	ua := user_agent.New(userAgent)
	var userDevice entity.UserDevice
	userDevice.UserLoginsID = userLogins.ID
	userDevice.Os = ua.OS()
	userDevice.Browser, _ = ua.Browser()
	userDevice.DeviceName = ua.Model()
	userDevice.IsPhone = ua.Mobile()

	if res := postDb.Save(&userDevice); res.Error != nil {
		return res.Error
	}

	return nil
}
