package models

import (
	"fmt"
	"gdialog/global"
	"gdialog/utils"
	"sort"
	"strings"
)

var DiseaseList = []interface{}{"糖尿病", "高血压"}

type Disease struct {
	ID      int64
	UserID  int64 `gorm:"unique"`
	Disease string
	Deleted bool
}

func (d Disease) TableName() string {
	return "disease"
}

func (d *Disease) Valid() (b bool, msg string) {
	// if no disease
	if d.Disease == "" {
		return true, ""
	}
	// split by ;
	ds := strings.Split(d.Disease, ";")
	new_ds := []string{}
	for _, v := range ds {
		if v == "" {
			continue
		}
		// not in given disease list
		if !utils.In(DiseaseList, v) {
			return false, fmt.Sprintf("Check %s disease", v)
		}
		new_ds = append(new_ds, v)
	}
	// sort disease
	sort.Strings(new_ds)
	d.Disease = strings.Join(new_ds, ";")
	return true, ""
}

func (d Disease) Save() {
	global.DB.Create(&d)
}
