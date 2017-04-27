package controller

import (
	"fmt"
	
	"github.com/reechou/holmes"
	"github.com/reechou/robot-tangbao/robot_proto"
	"github.com/reechou/robot-tangbao/models"
)

func (self *Logic) runShareImgMsg() {
	for {
		select {
		case msg := <-self.msgChan:
			self.handleShareImgMsg(msg)
		case <-self.stop:
			close(self.done)
			return 
		}
	}
}

func (self *Logic) handleShareImgMsg(msg *robot_proto.ReceiveMsgInfo) {
	rul := &models.RobotUserLuck{
		Robot: msg.BaseInfo.WechatNick,
		User: msg.BaseInfo.FromNickName,
	}
	has, err := models.GetRobotUserLuck(rul)
	if err != nil {
		holmes.Error("get robot user luck error: %v", err)
		return
	}
	if !has {
		err = models.CreateRobotUserLuck(rul)
		if err != nil {
			holmes.Error("create robot user luck error: %v", err)
		}
	}
	if rul.LuckCode != 0 {
		msgs := []MsgInfo{
			MsgInfo{Msg: fmt.Sprintf("领奖号：%d\n恭喜您审核通过！\n我们将于5月10日开奖，请勿删除朋友圈！", rul.LuckCode),
				MsgType: robot_proto.RECEIVE_MSG_TYPE_TEXT},
		}
		self.sendMsgBack(msgs, msg)
		return
	}
	err = models.UpdateRobotUserLuckShareNum(rul)
	if err != nil {
		holmes.Error("update robot user luck share num error: %v", err)
		return
	}
	_, err = models.GetRobotUserLuck(rul)
	if err != nil {
		holmes.Error("get robot user luck error: %v", err)
		return
	}
	holmes.Debug("get rul: %v", rul)
	if rul.ShareNum >= 3 {
		luckCode, _ := self.doLuck(msg.BaseInfo.WechatNick, msg.BaseInfo.FromNickName)
		rul.LuckCode = int64(luckCode)
		err, ok := models.UpdateRobotUserLuckCode(rul)
		if err != nil {
			holmes.Error("update robot user luck code error: %v", err)
			return
		}
		if ok {
			msgs := []MsgInfo{
				MsgInfo{Msg: fmt.Sprintf("领奖号：%d\n恭喜您审核通过！\n我们将于5月10日开奖，请勿删除朋友圈！", luckCode),
					MsgType: robot_proto.RECEIVE_MSG_TYPE_TEXT},
				MsgInfo{Msg: "重要提示：开奖规则请看下图，或者访问 http://t.cn/RXECF8r \n" +
					"请家长放心分享，我们承诺免费赠送，不收任何费用。",
					MsgType: robot_proto.RECEIVE_MSG_TYPE_TEXT},
				MsgInfo{Msg: "http://oe3slowqt.bkt.clouddn.com/FpHmwnKW7peN_1IKS3Jwq6aafIzU",
					MsgType: robot_proto.RECEIVE_MSG_TYPE_IMG},
			}
			self.sendMsgBack(msgs, msg)
		}
	}
}
