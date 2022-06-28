package models

import (
	"gdialog/global"
	"gdialog/utils"
	"regexp"
	"time"
)

type (
	User struct {
		UserID     int64   `gorm:"autoIncement;primaryKey" json:"user_id"`
		Username   string  `gorm:"unique;not null" json:"username"`
		Password   string  `gorm:"not null" json:"password"`
		Phone      string  `json:"phone"`
		Age        uint16  `json:"age"`
		Gender     string  `json:"gender"`
		Height     float32 `json:"height"`
		Weight     float32 `json:"weight"`
		CreateTime int64   `gorm:"column:create_time"`
		Deleted    bool
		// for models.Disease use
		Disease string `gorm:"-" json:"disease"`
	}
)

func (u User) TableName() string {
	return "user"
}

func (u User) Exists() bool {
	// whether User exists
	if err := global.DB.Where("username = ?", u.Username).Take(&User{}).Error; err != nil {
		return false
	}
	return true
}

func (u *User) GetUser() error {
	// get user info
	if err := global.DB.Where("username = ?", u.Username).Take(&u).Error; err != nil {
		return err
	}
	return nil
}

func (u User) ValidLogin() bool {
	// whether User exists
	if err := global.DB.Where("username = ? and password = ?", u.Username, u.Password).Take(&User{}).Error; err != nil {
		return false
	}
	return true
}

func (u User) Validate() (b bool, err string) {
	// username
	if ok, _ := regexp.MatchString("^[a-zA-Z0-9_]{4,16}$", u.Username); !ok {
		return false, "Username not valid, only contains 4-16 char A-Za-z0-9_"
	}
	// password
	if ok, _ := regexp.MatchString("^[a-zA-Z0-9.!]{4,16}$", u.Password); !ok {
		return false, "Password not valid, only contains 4-16 char A-Za-z0-9.!"
	}
	// phone
	if ok, _ := regexp.MatchString("^1[0-9]{10}$", u.Phone); u.Phone != "" && !ok {
		return false, "Input 11-digit mobile phone number"
	}
	// gender
	if u.Gender != "" && !utils.In([]any{"M", "F", "O"}, u.Gender) {
		return false, "Check your gender, values are M, F, O"
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

func (u User) Save() error {
	// add create time
	u.CreateTime = time.Now().Unix()
	// create User
	if err := global.DB.Create(&u).Error; err != nil {
		return err
	}
	return nil
}
