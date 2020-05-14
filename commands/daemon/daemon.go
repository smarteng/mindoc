package daemon

import (
	"fmt"
	"os"

	"path/filepath"

	"github.com/astaxie/beego"
	"github.com/kardianos/service"
	"github.com/smarteng/mindoc/commands"
	"github.com/smarteng/mindoc/conf"
	"github.com/smarteng/mindoc/controllers"
)

type Daemon struct {
	config *service.Config
	errs   chan error
}

func NewDaemon() *Daemon {
	config := &service.Config{
		Name:             "mindocd",                               //服务显示名称
		DisplayName:      "MinDoc service",                        //服务名称
		Description:      "A document online management program.", //服务描述
		WorkingDirectory: conf.WorkingDirectory,
		Arguments:        os.Args[1:],
	}

	return &Daemon{
		config: config,
		errs:   make(chan error, 100),
	}
}

func (d *Daemon) Config() *service.Config {
	return d.config
}
func (d *Daemon) Start(s service.Service) error {
	go d.Run()
	return nil
}

func (d *Daemon) Run() {
	commands.ResolveCommand(d.config.Arguments)
	commands.RegisterFunction()
	commands.RegisterAutoLoadConfig()
	commands.RegisterError()
	beego.ErrorController(&controllers.ErrorController{})

	f, err := filepath.Abs(os.Args[0])
	if err != nil {
		f = os.Args[0]
	}
	msg := "MinDoc version:%s\nbuild time:%s\nstart directory:%s\ngo version:%s\n"
	fmt.Printf(msg, conf.VERSION, conf.BUILD_TIME, f, conf.GO_VERSION)
	beego.Run()
}

func (d *Daemon) Stop(s service.Service) error {
	if service.Interactive() {
		os.Exit(0)
	}
	return nil
}

func Install() {
	d := NewDaemon()
	d.config.Arguments = os.Args[3:]
	s, err := service.New(d, d.config)

	if err != nil {
		beego.Error("Create service error => ", err)
		os.Exit(1)
	}
	err = s.Install()
	if err != nil {
		beego.Error("Install service error:", err)
		os.Exit(1)
	} else {
		beego.Info("Service installed!")
	}

	os.Exit(0)
}

func Uninstall() {
	d := NewDaemon()
	s, err := service.New(d, d.config)

	if err != nil {
		beego.Error("Create service error => ", err)
		os.Exit(1)
	}
	err = s.Uninstall()
	if err != nil {
		beego.Error("Install service error:", err)
		os.Exit(1)
	} else {
		beego.Info("Service uninstalled!")
	}
	os.Exit(0)
}

func Restart() {
	d := NewDaemon()
	s, err := service.New(d, d.config)

	if err != nil {
		beego.Error("Create service error => ", err)
		os.Exit(1)
	}
	err = s.Restart()
	if err != nil {
		beego.Error("Install service error:", err)
		os.Exit(1)
	} else {
		beego.Info("Service Restart!")
	}
	os.Exit(0)
}
