package utils

import (
	"GinChat/db"
	"GinChat/entity"
	"github.com/gin-gonic/gin"
	"github.com/mssola/user_agent"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
	"math/rand/v2"
	"mime/multipart"
	"net"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func GetClientIP(request *gin.Context) string {
	ip := request.GetHeader("X-Forwarded-For")
	if ip != "" {
		// first X-Forwarded-For
		ips := strings.Split(ip, ",")
		if len(ips) > 0 {
			return strings.TrimSpace(ips[0])
		}
	}
	// Check the X-Real-IP header
	ip = request.GetHeader("X-Real-IP")
	if ip != "" {
		return ip
	}
	// Fallback to RemoteAddr
	ip, _, err := net.SplitHostPort(request.Request.RemoteAddr)
	if err != nil {
		return request.Request.RemoteAddr
	}
	return ip
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
	userIP.IP = GetClientIP(request)
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
func ImageValidate(image *multipart.FileHeader) bool {
	// format Check
	ext := filepath.Ext(image.Filename)
	var ImageFormats = []string{".png", ".jpg", ".jpeg", ".webp"}
	var formatCheck bool
	for _, item := range ImageFormats {
		if ext == item {
			formatCheck = true
		}
	}
	if !formatCheck {
		return false
	}
	//todo: make image quality less
	if image.Size > 600000 {
		return false
	}
	return true
}

func FilePathController(imagePath string, imageName string) string {
	if _, err := os.Open(imagePath + imageName); err == nil {
		index := strings.Index(imageName, ".")
		noFormatName := imageName[:index]
		FormatName := filepath.Ext(imageName)
		imagePath += noFormatName + uuid.NewV4().String() + FormatName
	} else {
		imagePath += imageName
	}
	return imagePath
}
