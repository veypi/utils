package cmd

import (
	"errors"
	"github.com/kardianos/service"
	"github.com/urfave/cli/v2"
	"github.com/veypi/utils"
	"github.com/veypi/utils/log"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

var (
	UnSupportWin = errors.New("un support on windows")
)

func NewCli(app *cli.App, cfgArgs ...interface{}) error {
	if app == nil || app.Name == "" {
		panic("invalid app")
	}
	c := new(cliCommand)
	c.Name = app.Name
	c.Install = &cli.Command{
		Name:  "install",
		Usage: "install cli and write cli args to cfg file",
		Action: func(_ *cli.Context) error {
			p, err := installCli(utils.GetRunnerPath(), app.Name)
			if err != nil {
				return err
			}
			log.Info().Msgf("install %s to %s", app.Name, p)
			if len(cfgArgs) == 0 {
				return nil
			}
			path := utils.PathJoin(GetSrvPath(c.Name), c.Name+".yml")
			if len(cfgArgs) > 1 && cfgArgs[1].(string) != "" {
				path = cfgArgs[1].(string)
			}
			err = DumpCfg(path, cfgArgs[0])
			if err != nil {
				return err
			}
			log.Info().Msgf("write %s.yml to %s", c.Name, path)
			return nil
		},
	}
	c.UnInstall = &cli.Command{
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
	if len(app.Commands) == 0 {
		app.Commands = make([]*cli.Command, 0, 10)
	}
	app.Commands = append(app.Commands, c.Install, c.UnInstall)
	app.ExitErrHandler = func(c *cli.Context, err error) {
		HandleExitCoder(err)
	}
	app.CommandNotFound = func(c *cli.Context, s string) {
		log.Warn().Msgf("%s command not found", s)
	}
	return nil
}

// 与 urfave/cli 配合使用
// startFunc 是 服务的启动程序，不要阻塞住，会注册到app.Action上
// stopFunc 是 服务的停止接口
// cfgArgs 是参数的相关配置， 第一项是cfg，是个可以序列化的对象， 第二项是cfgFilePath, 字符串
// cfg 是 参数对象， cfgFilePath 是参数文件地址， install 时 创建该文件，并填入cfg中数据, cfg 为空则不创建， cfgFilePath 为空则在默认位置创建
// 默认位置 `C:\Program Files\name\name.yml` 或者 /etc/name/name.yaml
func NewSrv(app *cli.App, startFunc cli.ActionFunc, stopFunc cli.ActionFunc, cfgArgs ...interface{}) error {
	if app == nil || app.Name == "" {
		panic("invalid app")
	}
	sc := &srvCommand{
		name:      app.Name,
		startFunc: startFunc,
		stopFunc:  stopFunc,
	}
	if len(cfgArgs) > 0 {
		sc.cfg = cfgArgs[0]
	}
	if len(cfgArgs) > 1 {
		sc.cfgFilePath = cfgArgs[1].(string)
	}
	svcConfig := &service.Config{Name: sc.name}
	srv, err := service.New(sc, svcConfig)
	if err != nil {
		return err
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
	return nil
}

// 这个被service 包的Run调用, 是服务开始执行的地方
func (sc *srvCommand) Start(service.Service) error {
	return sc.startFunc(sc.cliCtx)
}
func (sc *srvCommand) Stop(service.Service) error {
	return sc.stopFunc(sc.cliCtx)
}

// TODO: status command
type srvCommand struct {
	srv         service.Service
	name        string
	install     *cli.Command
	unInstall   *cli.Command
	start       *cli.Command
	stop        *cli.Command
	restart     *cli.Command
	run         *cli.Command
	cliCtx      *cli.Context
	startFunc   cli.ActionFunc
	stopFunc    cli.ActionFunc
	cfg         interface{}
	cfgFilePath string
}

func (sc *srvCommand) init() {
	sc.install = &cli.Command{
		Name:  "install",
		Usage: "install service and write cli args to config file",
		Action: func(c *cli.Context) error {
			p, err := installCli(utils.GetRunnerPath(), sc.name)
			if err != nil {
				return err
			}
			log.Info().Msgf("install %s to %s", sc.name, p)
			svcConfig := &service.Config{Name: sc.name, Executable: p}
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
			if sc.cfg == nil {
				return nil
			}
			path := utils.PathJoin(GetSrvPath(sc.name), sc.name+".yml")
			if sc.cfgFilePath != "" {
				path = sc.cfgFilePath
			}
			err = DumpCfg(path, sc.cfg)
			if err != nil {
				return err
			}
			log.Info().Msgf("write %s.yml to %s", sc.name, path)
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

type cliCommand struct {
	Name      string
	Install   *cli.Command
	UnInstall *cli.Command
}

func GetSrvPath(name string) string {
	var path string
	if utils.IsWindows() {
		path = `C:\Program Files\` + name
	} else {
		path = "/etc/" + name
	}
	return path
}

func GetLocalCfg(name string) string {
	home, err := utils.Home()
	if err != nil {
		log.Warn().Msg(err.Error())
		home = utils.GetRunnerPath()
	}
	return utils.PathJoin(home, ".config", name)
}

func installCli(srcPath, bin string) (string, error) {
	binPath, err := os.Executable()
	if err != nil {
		binPath, _ = filepath.Abs(os.Args[0])
	}
	if strings.HasPrefix(binPath, "/usr/bin") || strings.HasPrefix(binPath, "/usr/local/bin") {
		return binPath, nil
	}
	if !utils.IsWindows() {
		if _, err := utils.CopyFile(filepath.Join(srcPath, bin), "/usr/bin/"+bin); err != nil {
			if _, err := utils.CopyFile(filepath.Join(srcPath, bin), "/usr/local/bin/"+bin); err != nil {
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

func LoadCfg(path string, cfg interface{}) error {
	yamlFile, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(yamlFile, cfg)
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
