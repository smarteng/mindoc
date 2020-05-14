package main

import (
	"fmt"
	"log"
	"os"

	_ "github.com/astaxie/beego/session/memcache"
	_ "github.com/astaxie/beego/session/mysql"
	_ "github.com/astaxie/beego/session/redis"
	"github.com/kardianos/service"
	_ "github.com/mattn/go-sqlite3"
	"github.com/smarteng/mindoc/commands"
	"github.com/smarteng/mindoc/commands/daemon"
	_ "github.com/smarteng/mindoc/routers"
)

func main() {
	fmt.Println(os.Args)
	if len(os.Args) >= 3 && os.Args[1] == "service" {
		switch os.Args[2] {
		case commands.CmdInstall:
			daemon.Install()
		case commands.CmdRemove:
			daemon.Uninstall()
		case commands.CmdRestart:
			daemon.Restart()
		}
	}

	commands.RegisterCommand()
	d := daemon.NewDaemon()
	s, err := service.New(d, d.Config())

	if err != nil {
		fmt.Println("Create service error => ", err)
		os.Exit(1)
	}

	if err := s.Run(); err != nil {
		log.Fatal("启动程序失败 ->", err)
	}
}
