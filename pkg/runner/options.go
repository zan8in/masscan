package runner

import (
	"fmt"

	"github.com/zan8in/goflags"
	"github.com/zan8in/gologger"
	"github.com/zan8in/gologger/levels"
)

type Options struct {
	Target     goflags.StringSlice
	TargetFile string

	Concurrency int
	RateLimit   int

	Debug bool

	Timeout int
	Retries int

	Proxy string
}

func NewOptions() (*Options, error) {

	ShowBanner()

	opts := &Options{}

	flagSet := goflags.NewFlagSet()
	flagSet.SetDescription(`AppBase`)

	flagSet.CreateGroup("input", "Input",
		flagSet.StringSliceVar(&opts.Target, "t", nil, "target (comma-separated)", goflags.NormalizedStringSliceOptions),
		flagSet.StringVar(&opts.TargetFile, "T", "", "list of target (one per line)"),
	)

	flagSet.CreateGroup("rate-limit", "Rate-limit",
		flagSet.IntVar(&opts.Concurrency, "c", DefaultConcurrency, "general internal worker threads"),
		flagSet.IntVar(&opts.RateLimit, "rate", DefaultRateLimit, "packets to send per second"),
	)

	flagSet.CreateGroup("optimization", "Optimization",
		flagSet.IntVar(&opts.Timeout, "timeout", DefaultTimeout, "time to wait in seconds before timeout (default 10)"),
		flagSet.IntVar(&opts.Retries, "retries", DefaultRetries, "number of times to retry a failed request (default 1)"),
		flagSet.StringVar(&opts.Proxy, "proxy", "", "list of http/socks5 proxy to use (comma separated or file input)"),
	)

	flagSet.CreateGroup("debug", "Debug",
		flagSet.BoolVar(&opts.Debug, "debug", false, "display debugging information"),
	)

	_ = flagSet.Parse()

	if err := opts.validateOptions(); err != nil {
		return nil, err
	}

	return opts, nil
}

func (opts *Options) validateOptions() error {

	if len(opts.Target) == 0 && len(opts.TargetFile) == 0 {
		return fmt.Errorf("validate error: must specify a target file or a target")
	}

	if opts.Debug {
		gologger.DefaultLogger.SetMaxLevel(levels.LevelDebug)
	}

	return nil
}
