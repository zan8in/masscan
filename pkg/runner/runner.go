package runner

type Runner struct {
	Opts           *Options
	TargetTempName string
}

func NewRunner(opts *Options) (*Runner, error) {

	runner := &Runner{
		Opts: opts,
	}

	if err := runner.InitTargetTempFile(); err != nil {
		return runner, err
	}

	return runner, nil
}
