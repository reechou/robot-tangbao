package controller

import (
	"encoding/json"
	"net/http"
	"sync"

	"github.com/reechou/holmes"
	"github.com/reechou/robot-tangbao/config"
	"github.com/reechou/robot-tangbao/ext"
	"github.com/reechou/robot-tangbao/models"
	"github.com/reechou/robot-tangbao/robot_proto"
)

type Logic struct {
	sync.Mutex
	
	robotExt *ext.RobotExt
	msgChan chan *robot_proto.ReceiveMsgInfo

	cfg *config.Config
	
	stop chan struct{}
	done chan struct{}
}

func NewLogic(cfg *config.Config) *Logic {
	l := &Logic{
		cfg: cfg,
		msgChan: make(chan *robot_proto.ReceiveMsgInfo, 1024),
		stop: make(chan struct{}),
		done: make(chan struct{}),
	}
	l.robotExt = ext.NewRobotExt(cfg)
	models.InitDB(cfg)
	
	go l.runShareImgMsg()
	
	l.init()

	return l
}

func (self *Logic) Stop() {
	close(self.stop)
	<-self.done
}

func (self *Logic) init() {
	http.HandleFunc("/robot/receive_msg", self.RobotReceiveMsg)
}

func (self *Logic) Run() {
	defer holmes.Start(holmes.LogFilePath("./log"),
		holmes.EveryDay,
		holmes.AlsoStdout,
		holmes.DebugLevel).Stop()

	if self.cfg.Debug {
		EnableDebug()
	}

	holmes.Info("server starting on[%s]..", self.cfg.Host)
	holmes.Infoln(http.ListenAndServe(self.cfg.Host, nil))
}

func WriteJSON(w http.ResponseWriter, code int, v interface{}) error {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "x-requested-with,content-type")
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(code)
	return json.NewEncoder(w).Encode(v)
}

func WriteBytes(w http.ResponseWriter, code int, v []byte) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "x-requested-with,content-type")
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(code)
	w.Write(v)
}

func EnableDebug() {

}
