package cmd

import (
	"errors"
	"github.com/kardianos/service"
	"github.com/urfave/cli/v2"
	"github.com/veypi/utils"
	"github.com/veypi/utils/log"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"math"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var (
	UnSupportWin = errors.New("un support on windows")
)

type Service interface {
	// 设置最多重试次数
	SetExecMax(uint)
	// 设置服务停止触发函数
	SetStopFunc(func())

	Run()
}

func NewCli(appName string, cfgArgs ...interface{}) *cli.App {
	if appName == "" {
		panic("invalid app")
	}
	app := cli.NewApp()
	var cfgFilePath = ""
	if len(cfgArgs) > 1 {
		cfgFilePath = cfgArgs[1].(string)
		if len(app.Flags) == 0 {
			app.Flags = make([]cli.Flag, 0, 10)
		}
	}
	app.Flags = []cli.Flag{&cli.StringFlag{
		Name:    "cfg",
		Aliases: []string{"c"},
		Value:   cfgFilePath,
	}}
	install := &cli.Command{
		Name:  "install",
		Usage: "install cli and write cli args to cfg file",
		Action: func(_ *cli.Context) error {
			p, err := installCli(app.Name)
			if err != nil {
				return err
			}
			log.Info().Msgf("install %s to %s", app.Name, p)
			if len(cfgArgs) == 0 {
				return nil
			}
			path := GetCfgPath(app.Name, app.Name)
			if len(cfgArgs) > 1 && cfgArgs[1].(string) != "" {
				path = cfgArgs[1].(string)
			}
			err = DumpCfg(path, cfgArgs[0])
			if err != nil {
				return err
			}
			log.Info().Msgf("write %s.yml to %s", app.Name, path)
			return nil
		},
	}
	unInstall := &cli.Command{
		Name:  "uninstall",
		Usage: "remove cli",
		Action: func(c *cli.Context) error {
			p, err := unInstallCli(app.Name)
			if err != nil {
				return err
			}
			log.Info().Msg("uninstall " + p)
			return nil
		},
	}
	app.Commands = []*cli.Command{install, unInstall}
	app.ExitErrHandler = func(c *cli.Context, err error) {
		HandleExitCoder(err)
	}
	app.CommandNotFound = func(c *cli.Context, s string) {
		log.Warn().Msgf("%s command not found", s)
	}
	return app
}

// 与 urfave/cli 配合使用
// runnerFunc 是 服务的阻塞启动程序，函数如果执行完会自动重新执行, 间隔时间1ms开始指数增长
// cfgArgs 是参数的相关配置， 第一项是cfg，是个可以序列化的对象， 第二项是cfgFilePath, 字符串
// cfg 是 参数对象， cfgFilePath 是参数文件地址， install 时 创建该文件，并填入cfg中数据, cfg 为空则不创建， cfgFilePath 为空则在默认位置创建
// 默认位置 `C:\Program Files\name\name.yml` 或者 /etc/name/name.yaml
func NewSrv(app *cli.App, runnerFunc cli.ActionFunc, cfgArgs ...interface{}) (Service, error) {
	if app == nil || app.Name == "" {
		panic("invalid app")
	}
	sc := &srvCommand{
		name:       app.Name,
		runnerFunc: runnerFunc,
		app:        app,
	}
	if len(cfgArgs) > 0 {
		sc.cfg = cfgArgs[0]
	}
	if len(cfgArgs) > 1 {
		sc.cfgFilePath = cfgArgs[1].(string)
		if len(app.Flags) == 0 {
			app.Flags = make([]cli.Flag, 0, 10)
		}
		app.Flags = append(app.Flags, &cli.StringFlag{
			Name:    "cfg",
			Aliases: []string{"c"},
			Value:   sc.cfgFilePath,
		})
	}
	svcConfig := &service.Config{Name: sc.name}
	srv, err := service.New(sc, svcConfig)
	if err != nil {
		return nil, err
	}
	sc.srv = srv
	sc.init()

	if len(app.Commands) == 0 {
		app.Commands = make([]*cli.Command, 0, 10)
	}
	app.Commands = append(app.Commands, sc.install, sc.unInstall, sc.start, sc.stop, sc.restart, sc.run)
	app.Action = sc.run.Action
	app.ExitErrHandler = func(c *cli.Context, err error) {
		HandleExitCoder(err)
	}
	app.CommandNotFound = func(c *cli.Context, s string) {
		log.Warn().Msgf("%s command not found", s)
	}
	return sc, nil
}

// TODO: status command
type srvCommand struct {
	app         *cli.App
	srv         service.Service
	name        string
	install     *cli.Command
	unInstall   *cli.Command
	start       *cli.Command
	stop        *cli.Command
	restart     *cli.Command
	run         *cli.Command
	cliCtx      *cli.Context
	runnerFunc  cli.ActionFunc
	stopFunc    func()
	exeCount    uint
	execMax     uint
	exit        chan uint8
	cfg         interface{}
	cfgFilePath string
}

func (sc *srvCommand) Run() {
	err := sc.app.Run(os.Args)
	if err != nil {
		log.Warn().Msg(err.Error())
	}
}

func (sc *srvCommand) SetExecMax(c uint) {
	sc.execMax = c
}

func (sc *srvCommand) SetStopFunc(fc func()) {
	sc.stopFunc = fc
}

// 这个被service 包的Run调用, 是服务开始执行的地方
func (sc *srvCommand) Start(service.Service) error {
	go sc.running()
	return nil
}

func (sc *srvCommand) running() {
	exit := make(chan uint8, 1)
	sc.exit = exit
	exit <- 1
	for {
		select {
		case n := <-exit:
			if n > 0 {
				go func() {
					defer func() {
						exit <- 1
						if e := recover(); e != nil {
							log.Error().Err(nil).Msgf("%v", e)
						}
					}()
					if sc.execMax > 0 && sc.execMax == sc.exeCount {
						err := sc.Stop(sc.srv)
						if err != nil {
							log.Warn().Msg(err.Error())
						}
						return
					}
					delta := time.Millisecond * time.Duration(math.Pow(2, float64(sc.exeCount)))
					if delta > time.Minute {
						delta = time.Minute
					}
					sc.exeCount++
					time.Sleep(delta)
					err := sc.runnerFunc(sc.cliCtx)
					if err != nil {
						log.Warn().Msg(err.Error())
					}
				}()
			} else {
				return
			}
		}
	}
}

func (sc *srvCommand) Stop(service.Service) error {
	close(sc.exit)
	if sc.stopFunc != nil {
		sc.stopFunc()
	}
	if service.Interactive() {
		os.Exit(0)
	}
	return nil
}

func (sc *srvCommand) init() {
	sc.install = &cli.Command{
		Name:  "install",
		Usage: "install service and write cli args to config file",
		Action: func(c *cli.Context) error {
			p, err := installCli(sc.name)
			if err != nil {
				return err
			}
			log.Info().Msgf("install %s to %s", sc.name, p)
			if sc.cfg != nil {
				path := utils.PathJoin(GetSrvPath(sc.name), sc.name+".yml")
				if sc.cfgFilePath != "" {
					path = sc.cfgFilePath
				}
				err = DumpCfg(path, sc.cfg)
				if err != nil {
					return err
				}
				log.Info().Msgf("write %s.yml to %s", sc.name, path)
			}
			svcConfig := &service.Config{Name: sc.name, Executable: p}
			if sc.cfgFilePath != "" {
				svcConfig.Arguments = []string{"-c", sc.cfgFilePath}
			}
			srv, err := service.New(sc, svcConfig)
			if err != nil {
				return err
			}
			_ = srv.Stop()
			_ = srv.Uninstall()
			err = srv.Install()
			if err != nil {
				return err
			}
			return nil
		},
	}
	sc.unInstall = &cli.Command{
		Name:  "uninstall",
		Usage: "remove service",
		Action: func(c *cli.Context) error {
			err := sc.Srv().Uninstall()
			if err != nil {
				if e, ok := err.(ExitCoder); ok && e.ExitCode() == 1 {
					return errors.New("service has been uninstalled")
				}
				return err
			}
			p, err := unInstallCli(sc.name)
			if err != nil {
				return err
			}
			log.Info().Msg("uninstall " + p)
			return nil
		},
	}

	sc.start = &cli.Command{
		Name:  "start",
		Usage: "start service",
		Action: func(c *cli.Context) error {
			return sc.Srv().Start()
		},
	}
	sc.run = &cli.Command{
		Name:  "run",
		Usage: "(default) program entry, blocked until running 'Stop' ",
		Action: func(c *cli.Context) error {
			sc.cliCtx = c
			return sc.Srv().Run()
		},
	}
	sc.stop = &cli.Command{
		Name:  "stop",
		Usage: "stop service",
		Action: func(c *cli.Context) error {
			return sc.Srv().Stop()
		},
	}
	sc.restart = &cli.Command{
		Name:  "restart",
		Usage: "restart service",
		Action: func(c *cli.Context) error {
			return sc.Srv().Restart()
		},
	}
}

func (sc *srvCommand) Srv() service.Service {
	return sc.srv
}

func GetSrvPath(name string) string {
	var path string
	if utils.IsWindows() {
		path = `.\` + name
	} else {
		path = GetLocalCfg(name)
	}
	return path
}

func GetCfgPath(item1, item2 string) string {
	name := ""
	for i, v := range os.Args {
		if v == "-c" || v == "--cfg" {
			if i < len(os.Args)-1 {
				name = os.Args[i+1]
				break
			}
		} else if strings.HasPrefix(v, "-c=") || strings.HasPrefix(v, "--cfg=") {
			name = strings.SplitN(v, "=", 2)[1]
			break
		}
	}
	if name != "" {
		s, err := utils.Abs(name)
		if err != nil {
			log.Warn().Msgf("parse cfg file path error: %s", err.Error())
		} else {
			return s
		}

	}
	return utils.PathJoin(GetSrvPath(item1), item2+".yml")
}

func GetLocalCfg(name string) string {
	home, err := utils.Home()
	if err != nil {
		log.Warn().Msg(err.Error())
		home = utils.GetRunnerPath()
	}
	return utils.PathJoin(home, ".config", name)
}

func installCli(bin string) (string, error) {
	binPath, err := os.Executable()
	if err != nil {
		binPath, _ = filepath.Abs(os.Args[0])
	}
	if strings.HasPrefix(binPath, "/usr/bin") || strings.HasPrefix(binPath, "/usr/local/bin") {
		return binPath, nil
	}
	if !utils.IsWindows() {
		if _, err := utils.CopyFile(binPath, "/usr/bin/"+bin); err != nil {
			if _, err := utils.CopyFile(binPath, "/usr/local/bin/"+bin); err != nil {
				return "", err
			} else {
				binPath = "/usr/local/bin/" + bin
			}
		} else {
			binPath = "/usr/bin/" + bin
		}
	} else {
		// TODO add bin to path
		//utils.CopyFile(filepath.Join(srcPath, bin+".exe"), filepath.Join(utils.GetRunnerPath(), bin+".exe"))
	}
	utils.ChMod(binPath, 0755)
	return binPath, nil
}

func unInstallCli(bin string) (rp string, re error) {
	path := "/usr/bin/" + bin
	if !utils.IsWindows() {
		if ok, err := utils.PathExists(path); ok {
			rp = path
			re = os.Remove(path)
		} else if err != nil {
			rp = path
			re = os.Remove(path)
		}
		path = "/usr/local/bin/" + bin
		if ok, err := utils.PathExists(path); ok {
			rp = path
			re = os.Remove(path)
		} else if err != nil {
			rp = path
			re = os.Remove(path)
		}
		return
	}
	return "", UnSupportWin
}

type ExitCoder interface {
	error
	ExitCode() int
}

func HandleExitCoder(err error) {
	if err == nil {
		return
	}
	if err != nil {
		log.Warn().Err(err).Msg("exit")
	}
	if exitErr, ok := err.(ExitCoder); ok {
		os.Exit(exitErr.ExitCode())
		return
	}
	os.Exit(1)
}

func LoadCfg(path string, cfg interface{}) {
	yamlFile, err := ioutil.ReadFile(path)
	if err != nil && !os.IsNotExist(err) {
		log.Warn().Msg(err.Error())
		return
	}
	err = yaml.Unmarshal(yamlFile, cfg)
	if err != nil {
		log.Warn().Msg(err.Error())
	}
}

// 会覆盖写入
func DumpCfg(path string, cfg interface{}) error {
	body, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}
	f, err := utils.MkFile(path)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.Write(body)
	return err
}
