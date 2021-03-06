package main

import (
	"context"
	"os"
	"regexp"
	"time"

	"github.com/Water-W/PVP/influxdb"
	"github.com/Water-W/PVP/pkg/biz"
	"github.com/c-bata/go-prompt"
)

var suggestions = []prompt.Suggest{
	{
		Text:        "workers",
		Description: "show worker infos",
	},
	{
		Text:        "startPeriodlyDump",
		Description: "begin to dump periodly",
	},
	{
		Text:        "storepoint",
		Description: "store the dump data to database",
	},
	{
		Text:        "Querydata",
		Description: "Try to query the dump data",
	},
	{
		Text:        "dump",
		Description: "dump node and links info",
	},
	{
		Text:        "echo",
		Description: "use echo service",
	},
	{
		Text:        "exit",
		Description: "exit pvp",
	},
}

func cli(ctrl *biz.MasterController) {
	mctrl := &mctrl{
		MasterController: ctrl,
		cancel:           nil,
	}
	c := CLI{
		cmds: []Cmd{
			{
				Pattern: `echo (\S+)`,
				Action:  mctrl.echo,
			},
			{
				Pattern: "startPeriodlyDump",
				Action:  mctrl.startPeriodlyDump,
			},
			{
				Pattern: `dump`,
				Action:  mctrl.dump,
			},
			{
				Pattern: "workers",
				Action:  mctrl.workers,
			},
			{
				Pattern: "storepoint",
				Action:  mctrl.storepoint,
			},
			{
				Pattern: "Querydata",
				Action:  mctrl.querydump,
			},
			{
				Pattern: "exit",
				Action: func([]string) {
					os.Exit(0)
				},
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

func (c *CLI) exec(input string) {
	for _, cmd := range c.cmds {
		re := regexp.MustCompile(cmd.Pattern)
		matches := re.FindStringSubmatch(input)
		if matches == nil {
			continue
		}
		cmd.Action(matches)
	}
}

func (c *CLI) Run() {
	completer := func(d prompt.Document) []prompt.Suggest {
		return prompt.FilterHasPrefix(suggestions, d.GetWordBeforeCursor(), false)
	}
	p := prompt.New(c.exec, completer)
	p.Run()
}

type mctrl struct {
	*biz.MasterController
	cancel func()
}

func (c *mctrl) echo(s []string) {
	results, err := c.MasterController.Echo(context.Background(), s[1])
	if err != nil {
		log.Error(err)
		return
	}
	log.Infof("echo:%+v", results)
}

func (c *mctrl) dump(s []string) {
	results, err := c.MasterController.Dump(context.Background())
	if err != nil {
		log.Error(err)
		return
	}
	log.Infof("dump:%+v", results)
}

func (c *mctrl) workers(s []string) {
	log.Infof("workers:%+v", c.MasterController.WorkerAddrs())
}

func (c *mctrl) storepoint(s []string) {
	results, err := c.MasterController.Dump(context.Background())
	if err != nil {
		log.Error(err)
		return
	}
	influxdb.Storedata(results)
	log.Infof("storepoint: finish")
}
func (c *mctrl) startPeriodlyDump(s []string) {
	ch, cancel := c.MasterController.StartPeriodlyDump(30 * time.Second)
	c.cancel = cancel
	go func(ch <-chan biz.DumpResults) {
		for res := range ch {
			if res.Err != nil {
				log.Error(res.Err)
				continue
			}
			log.Debug("results len:%d", len(res.Results))
			influxdb.Storedata(res.Results)
		}
	}(ch)
}

func (c *mctrl) stopPeriodlyDump(s []string) {
	c.cancel()
	c.cancel = nil
}

func (c *mctrl) querydump(s []string) {
	influxdb.Querydata()
	log.Infof("querydump: finish")
}
