package runner

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

var (
	tempfile = "appbase-temp-targets-*"
)

func (runner *Runner) InitTargetTempFile() error {
	var err error

	tempTargets, err := os.CreateTemp("", tempfile)
	if err != nil {
		return err
	}
	defer tempTargets.Close()

	if len(runner.Opts.Target) > 0 {
		for _, target := range runner.Opts.Target {
			fmt.Fprintf(tempTargets, "%s\n", target)
		}
	}

	if len(runner.Opts.TargetFile) > 0 {
		f, err := os.Open(runner.Opts.TargetFile)
		if err != nil {
			return err
		}
		defer f.Close()

		if _, err := io.Copy(tempTargets, f); err != nil {
			return err
		}
	}

	runner.TargetTempName = tempTargets.Name()

	return nil
}

func (runner *Runner) GetTargets() (chan string, error) {
	if len(runner.TargetTempName) == 0 {
		return nil, fmt.Errorf("no target")
	}

	out := make(chan string)
	go func() {
		defer close(out)

		f, err := os.Open(runner.TargetTempName)
		if err != nil {
			return
		}
		defer f.Close()

		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			out <- scanner.Text()
		}
	}()

	return out, nil
}
