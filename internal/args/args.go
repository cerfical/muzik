package args

import (
	"flag"

	"github.com/cerfical/muzik/internal/log"
)

type Args struct {
	LogLevel   log.Level
	ServerAddr string
}

func Parse(args []string) *Args {
	a := Args{}

	flags := flag.NewFlagSet(args[0], flag.ExitOnError)
	flags.TextVar(&a.LogLevel, "log", log.LevelInfo, "severity `level` of log messages")
	flags.Parse(args[1:])

	// Handle positional arguments
	if flags.NArg() >= 1 {
		a.ServerAddr = flags.Arg(0)
	} else {
		a.ServerAddr = "localhost:8080"
	}

	return &a
}
