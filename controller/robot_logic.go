package controller

import (
	"fmt"
	"time"
	"math/rand"
	"strings"
	
	"github.com/reechou/holmes"
	"github.com/reechou/robot-tangbao/robot_proto"
	"github.com/reechou/robot-tangbao/models"
)

func (self *Logic) HandleReceiveMsg(msg *robot_proto.ReceiveMsgInfo) {
	holmes.Debug("receive robot msg: %v", msg)
	switch msg.BaseInfo.ReceiveEvent {
	case robot_proto.RECEIVE_EVENT_MSG:
		self.handleMsg(msg)
	}
}

func (self *Logic) handleMsg(msg *robot_proto.ReceiveMsgInfo) {
	switch msg.BaseInfo.FromType {
	case robot_proto.FROM_TYPE_PEOPLE:
		self.doPeopleMsg(msg)
	}
}

func (self *Logic) doPeopleMsg(msg *robot_proto.ReceiveMsgInfo) {
	switch msg.MsgType {
	case robot_proto.RECEIVE_MSG_TYPE_TEXT:
		if strings.Contains(msg.Msg, "领衣服") || strings.Contains(msg.Msg, "糖宝") {
			self.doTangbao(msg)
		}
	case robot_proto.RECEIVE_MSG_TYPE_IMG:
		self.doShareImg(msg)
	}
}

func (self *Logic) doShareImg(msg *robot_proto.ReceiveMsgInfo) {
	select {
	case self.msgChan <- msg:
	case <-self.stop:
		return
	}
}

func (self *Logic) userKey(msg *robot_proto.ReceiveMsgInfo) string {
	return fmt.Sprintf("%s_%s", msg.BaseInfo.WechatNick, msg.BaseInfo.FromNickName)
}

func (self *Logic) doLuck(robot, user string) (int, bool) {
	randNum := rand.Intn(1000)
	if randNum < 10 {
		// get the luck
		lc := &models.LuckCode{
			LuckCode: int64(self.cfg.LuckCode1),
		}
		has, err := models.GetLuckCode(lc)
		if err != nil {
			holmes.Error("get luck code error: %v", err)
			return self.getRandLuckCode(), false
		}
		if !has {
			holmes.Error("get luck code[%d] has node", self.cfg.LuckCode1)
			return self.getRandLuckCode(), false
		}
		if lc.NowNum < int64(self.cfg.Luck1Num) {
			lc.NowNum = lc.NowNum + 1
			err, ok := models.UpdateLuckCode(lc, int64(self.cfg.Luck1Num))
			if err != nil {
				holmes.Error("update luck code error: %v", err)
				return self.getRandLuckCode(), false
			}
			if ok {
				holmes.Info("robot[%s] user[%s] randnum[%d] get the luck", robot, user, randNum)
				return self.cfg.LuckCode1, true
			}
		}
	}
	
	return self.getRandLuckCode(), false
}

func (self *Logic) getRandLuckCode() int {
	luckCode := rand.Intn(1701)
	luckCode += 188
	
	if luckCode == self.cfg.LuckCode1 {
		luckCode++
	}
	
	return luckCode
}

func (self *Logic) doTangbao(msg *robot_proto.ReceiveMsgInfo) {
	holmes.Debug("in tangbao")
	var posterImg string
	if len(self.cfg.PosterImgList) > 0 {
		offset := rand.Intn(len(self.cfg.PosterImgList))
		posterImg = self.cfg.PosterImgList[offset]
	}
	msgs := []MsgInfo{
		MsgInfo{Msg: "您好，「baby sugar潮童」联合《昕薇》杂志，为您的孩子准备了【一整年】的童装大礼包，价值2888元，免费领取中！\n" +
			"新加的家长请第一时间完成报名！报名成功即可免费给宝宝领取一整年的衣服啦！\n\n" +
			"【报名方式】\n1、复制下方文字，长按保存图片\n" +
			"2、将文字和图片转发到【朋友圈】和【2个微信群】（最好是家长群），后截图\n" +
			"3、将3张截图发给我，审核通过后报名成功",
			MsgType: robot_proto.RECEIVE_MSG_TYPE_TEXT},
		MsgInfo{Msg: "太爽啦！免费给孩子拿了【一整年】的衣服，今年春夏秋冬都不用买啦！欢乐等包裹中……点开图片长按关注这个号并发送“领衣服”就行，仅限前100名哟！",
			MsgType: robot_proto.RECEIVE_MSG_TYPE_TEXT},
		MsgInfo{Msg: posterImg,
			MsgType: robot_proto.RECEIVE_MSG_TYPE_IMG},
	}
	self.sendMsgBack(msgs, msg)
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

type MsgInfo struct {
	Msg     string
	MsgType string
}

func (self *Logic) sendMsgBack(msgs []MsgInfo, msg *robot_proto.ReceiveMsgInfo) error {
	robotHost := self.cfg.DefaultRobotHost
	var sendReq robot_proto.SendMsgInfo
	for _, v := range msgs {
		sendReq.SendMsgs = append(sendReq.SendMsgs, robot_proto.SendBaseInfo{
			WechatNick: msg.BaseInfo.WechatNick,
			ChatType:   msg.BaseInfo.FromType,
			UserName:   msg.BaseInfo.FromUserName,
			NickName:   msg.BaseInfo.FromNickName,
			MsgType:    v.MsgType,
			Msg:        v.Msg,
		})
	}
	err := self.robotExt.SendMsgs(robotHost, &sendReq)
	if err != nil {
		holmes.Error("send msg[%v] back error: %v", sendReq, err)
		return err
	}
	return nil
}


