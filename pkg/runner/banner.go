package runner

import (
	"github.com/zan8in/gologger"
)

var Version = "0.0.1"

func ShowBanner() {
	gologger.Print().Msgf("\n\tM A S S C A N\t%s\n\n", Version)
}
