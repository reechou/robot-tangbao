package models

import (
	"fmt"
	"time"
	
	"github.com/reechou/holmes"
)

type RobotUserLuck struct {
	ID        int64  `xorm:"pk autoincr" json:"id"`
	Robot     string `xorm:"not null default '' varchar(64) unique(robot_user)" json:"robot"`
	User      string `xorm:"not null default '' varchar(190) unique(robot_user)" json:"user"`
	ShareNum  int64  `xorm:"not null default 0 int" json:"shareNum"`
	LuckCode  int64  `xorm:"not null default 0 int" json:"type"`
	CreatedAt int64  `xorm:"not null default 0 int" json:"createdAt"`
	UpdatedAt   int64  `xorm:"not null default 0 int" json:"-"`
}

func CreateRobotUserLuck(info *RobotUserLuck) error {
	if info.Robot == "" {
		return fmt.Errorf("robot[%s] cannot be nil.", info.Robot)
	}
	
	now := time.Now().Unix()
	info.CreatedAt = now
	info.UpdatedAt = now
	
	_, err := x.Insert(info)
	if err != nil {
		holmes.Error("create robot user luck error: %v", err)
		return err
	}
	holmes.Info("create robot user luck[%v] success.", info)
	
	return nil
}

func UpdateRobotUserLuckCode(info *RobotUserLuck) (error, bool) {
	info.UpdatedAt = time.Now().Unix()
	affected, err := x.ID(info.ID).Cols("luck_code", "updated_at").Where("luck_code = 0").Update(info)
	if err != nil {
		return err, false
	}
	if affected == 0 {
		holmes.Debug("update luck code affected == 0 user[%d]", info.ID)
		return nil, false
	}
	return nil, true
}

func GetRobotUserLuck(info *RobotUserLuck) (bool, error) {
	tempInfo := new(RobotUserLuck)
	has, err := x.Where("robot = ?", info.Robot).And("user = ?", info.User).Get(tempInfo)
	if err != nil {
		return false, err
	}
	if !has {
		holmes.Debug("cannot find user[%s] from robot[%s]", info.User, info.Robot)
		return false, nil
	}
	info.ID = tempInfo.ID
	info.ShareNum = tempInfo.ShareNum
	info.LuckCode = tempInfo.LuckCode
	return true, nil
}

func UpdateRobotUserLuckShareNum(info *RobotUserLuck) error {
	info.UpdatedAt = time.Now().Unix()
	var err error
	result, err := x.Exec("update robot_user_luck set share_num=share_num+1, updated_at=? where robot=? and user=?",
		info.UpdatedAt, info.Robot, info.User)
	if err != nil {
		return err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		holmes.Error("update robot share num error affected == 0")
		return fmt.Errorf("update robot share num error affected == 0")
	}
	holmes.Info("robot[%s] user[%s] update share num success.", info.Robot, info.User)
	return nil
}
