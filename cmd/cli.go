package main

import (
	"context"
	"regexp"

	"github.com/Water-W/PVP/pkg/biz"
	"github.com/c-bata/go-prompt"
)

func cli(ctrl *biz.MasterController) {
	mctrl := (*mctrl)(ctrl)
	c := CLI{
		cmds: []Cmd{
			{
				Pattern: `echo (\S+)`,
				Action:  mctrl.echo,
			},
			{
				Pattern: "workers",
				Action:  mctrl.workers,
			},
		},
	}
	c.Run()
}

type CLI struct {
	cmds []Cmd
}

type Cmd struct {
	Pattern string
	Action  func([]string)
}

func (c *CLI) Run() {
	noopCompleter := func(d prompt.Document) []prompt.Suggest {
		return []prompt.Suggest{}
	}
	for {
		input := prompt.Input("> ", noopCompleter)
		for _, cmd := range c.cmds {
			re := regexp.MustCompile(cmd.Pattern)
			matches := re.FindStringSubmatch(input)
			if matches == nil {
				continue
			}
			cmd.Action(matches)
		}
	}
}

type mctrl biz.MasterController

func (c *mctrl) echo(s []string) {
	results, err := (*biz.MasterController)(c).Echo(context.Background(), s[1])
	if err != nil {
		log.Error(err)
		return
	}
	log.Infof("echo:%+v", results)
}

func (c *mctrl) workers(s []string) {
	log.Infof("workers:%+v", (*biz.MasterController)(c).WorkerAddrs())
}
