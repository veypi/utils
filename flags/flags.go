//
// flags.go
// Copyright (C) 2024 veypi <i@veypi.com>
// 2024-08-12 14:52
// Distributed under terms of the MIT license.
//

package flags

import (
	"flag"
	"fmt"
	"os"

	"github.com/veypi/utils"
	"github.com/veypi/utils/logx"
	"gopkg.in/yaml.v3"
)

func New(name string, des string) *Flags {
	return &Flags{
		FlagSet: *flag.NewFlagSet(name, flag.ExitOnError),
		des:     des,
		depth:   0,
	}
}

type Flags struct {
	flag.FlagSet
	des     string
	depth   int
	subs    []*Flags
	parent  *Flags
	help    string
	Command func() error
	Before  func() error
	After   func(error) error
}

func (f *Flags) Parse() error {
	return f.parse(os.Args[1:])
}

func (f *Flags) str() string {
	if f.parent == nil {
		return f.Name()
	}
	return f.parent.str() + " " + f.Name()
}

func (f *Flags) fire_before_func() error {
	var err error
	if f.Before != nil {
		err = f.Before()
		if err != nil {
			return err
		}
	}
	if f.parent != nil {
		return f.parent.fire_before_func()
	}
	return nil
}

func (f *Flags) fire_after_func(e error) error {
	var err error
	if f.After != nil {
		err = f.After(e)
		if err != nil {
			return err
		}
	}
	if f.parent != nil {
		return f.parent.fire_after_func(e)
	}
	return e
}

func (f *Flags) parse(arguments []string) (err error) {
	f.FlagSet.Usage = f.Usage
	err = f.FlagSet.Parse(arguments)
	if err != nil {
		return err
	}
	arg0 := f.Arg(0)
	for _, c := range f.subs {
		if c.Name() == arg0 {
			err = c.parse(f.Args()[1:])
			return
		}
	}
	if f.Command != nil {
		if f.NArg() != 0 {
			return fmt.Errorf("unexpected argument: %s", f.Arg(0))
		}
		defer func() {
			if e := recover(); e != nil {
				err = fmt.Errorf("%v", e)
			}
		}()
		err = f.fire_before_func()
		if err != nil {
			return
		}
		err = f.Command()
		return f.fire_after_func(err)
	} else {
		f.Usage()
	}
	return
}

func (f *Flags) Usage() {
	f.usage()
	fmt.Fprintf(os.Stderr, "%s\n", f.help)
}

func (f *Flags) usage() {
	f.help = ""
	f.SetOutput(f)
	f.FlagSet.PrintDefaults()
	if f.help != "" {
		f.help = fmt.Sprintf("%s:\n  %s\n\nOptions:\n", f.str(), f.des) + f.help
	} else {
		f.help = fmt.Sprintf("%s:\n  %s", f.str(), f.des) + f.help
	}
	if len(f.subs) > 0 {
		fmt.Fprint(f, "\nCommands:\n")
		for _, c := range f.subs {
			fmt.Fprintf(f, "  %s:  %s\n", c.Name(), c.des)
		}
	}
}

func (f *Flags) Write(p []byte) (n int, err error) {
	f.help += string(p)
	return len(p), nil
}

func (f *Flags) SubCommand(name, des string) *Flags {
	s := New(name, des)
	s.depth = f.depth + 1
	f.subs = append(f.subs, s)
	s.parent = f
	return s
}

func LoadCfg(path string, cfg interface{}) {
	yamlFile, err := os.ReadFile(path)
	if err != nil && !os.IsNotExist(err) {
		logx.Warn().Msg(err.Error())
		return
	}
	err = yaml.Unmarshal(yamlFile, cfg)
	if err != nil {
		logx.Warn().Msg(err.Error())
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