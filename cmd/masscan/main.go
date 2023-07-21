package main

import (
	"github.com/zan8in/gologger"
	"github.com/zan8in/masscan/pkg/runner"
)

func main() {
	opts, err := runner.NewOptions()
	if err != nil {
		gologger.Error().Msg(err.Error())
		return
	}

	runner, err := runner.NewRunner(opts)
	if err != nil {
		gologger.Error().Msg(err.Error())
		return
	}

	ts, err := runner.GetTargets()
	if err != nil {
		gologger.Error().Msg(err.Error())
		return
	}

	for t := range ts {
		gologger.Debug().Msg(t)
	}
}
