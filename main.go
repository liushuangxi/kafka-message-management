package main

import (
	_ "github.com/liushuangxi/kafka-message-management/routers"
	_ "github.com/liushuangxi/kafka-message-management/sysinit"

	"github.com/astaxie/beego"
)

func main() {
	beego.Run()
}
