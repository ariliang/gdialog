package models

import (
	"gdialog/global"
	"gdialog/utils"
	"regexp"
	"time"
)

type (
	UserWX struct {
		UserID     int64  `gorm:"autoIncement;primaryKey"`
		OpenID     string `gorm:"unique;not null"`
		Username   string `gorm:"not null"`
		Phone      string `gorm:"default:null;unique"`
		Age        uint16
		Gender     uint8 `gorm:"default:null"`
		Height     float32
		Weight     float32
		CreateTime int64 `gorm:"column:create_time"`
		Deleted    bool
	}
)

func (u UserWX) TableName() string {
	return "userwx"
}

func (u UserWX) Exists() bool {
	// whether User exists
	if err := global.DB.Where("open_id = ?", u.OpenID).Take(&UserWX{}).Error; err != nil {
		return false
	}
	return true
}

func (u *UserWX) GetUser() error {
	// get user info
	if err := global.DB.Where("open_id = ?", u.OpenID).Take(&u).Error; err != nil {
		return err
	}
	return nil
}

func (u UserWX) Validate() (b bool, err string) {
	// gender
	if !utils.In([]any{0, 1}, u.Gender) {
		return false, "Check your gender"
	}
	// phone
	if ok, _ := regexp.MatchString("^1[0-9]{10}$", u.Phone); u.Phone != "" && !ok {
		return false, "Input 11-digit mobile phone number"
	}
	// height
	if u.Height < 0 {
		return false, "Check your height"
	}
	// weight
	if u.Weight < 0 {
		return false, "Check your weight"
	}
	return true, ""
}

func (u UserWX) Update() {
}

func (u UserWX) Save() error {
	// add create time
	u.CreateTime = time.Now().Unix()
	// create User
	if err := global.DB.Create(&u).Error; err != nil {
		return err
	}
	return nil
}
