package models

import (
	"time"
	
	"github.com/reechou/holmes"
)

type LuckCode struct {
	ID        int64 `xorm:"pk autoincr" json:"id"`
	LuckCode  int64 `xorm:"not null default 0 int unique" json:"luckCode"`
	NowNum    int64 `xorm:"not null default 0 int" json:"luckCode"`
	CreatedAt int64 `xorm:"not null default 0 int" json:"createAt"`
	UpdatedAt int64 `xorm:"not null default 0 int" json:"-"`
}

func GetLuckCode(info *LuckCode) (bool, error) {
	has, err := x.Where("luck_code = ?", info.LuckCode).Get(info)
	if err != nil {
		return false, err
	}
	if !has {
		holmes.Debug("cannot find robot luck_code from luck_code[%d]", info.LuckCode)
		return false, nil
	}
	return true, nil
}

func UpdateLuckCode(info *LuckCode, maxNum int64) (error, bool) {
	info.UpdatedAt = time.Now().Unix()
	affected, err := x.ID(info.ID).Cols("now_num", "updated_at").Where("luck_code = ?", info.LuckCode).And("now_num < ?", maxNum).Update(info)
	if err != nil {
		holmes.Error("update luck code error: %v", err)
		return err, false
	}
	if affected == 0 {
		holmes.Debug("luck code[%s] has max", info.LuckCode)
		return nil, false
	}
	return nil, true
}
