package main

import (
	"github.com/reechou/robot-tangbao/config"
	"github.com/reechou/robot-tangbao/controller"
)

func main() {
	controller.NewLogic(config.NewConfig()).Run()
}
