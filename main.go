package main

import (
	"flag"
	"fmt"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	l   *zap.Logger
	log *zap.SugaredLogger
)

func init() {
	logConfig := zap.NewDevelopmentConfig()
	logConfig.EncoderConfig.TimeKey = ""
	logConfig.EncoderConfig.NameKey = ""
	logConfig.EncoderConfig.CallerKey = ""
	logConfig.EncoderConfig.StacktraceKey = ""
	l, _ = logConfig.Build()
}

var Usage = `
mygvm 
	install VERSION
	enable VERSION
	show
	help
`

func printHelp() {
	fmt.Println(Usage)
	os.Exit(0)
}

type ScriptParams struct {
	Command   string
	GoVersion string
}

var Debug = flag.Bool("debug", false, "run script in debug mode")

func getScriptParams() ScriptParams {
	flag.Parse()
	if !*Debug {
		l = l.WithOptions(zap.IncreaseLevel(zapcore.InfoLevel))
	}
	log = l.Sugar()
	c := ScriptParams{
		Command:   flag.Arg(0),
		GoVersion: flag.Arg(1),
	}
	valid := true
	if c.Command == "" {
		log.Errorf("command required as a first script param")
		valid = false
	}
	if !valid {
		printHelp()
		os.Exit(0)
	}
	return c
}

func main() {
	defer l.Sync()

	c := getScriptParams()

	switch c.Command {
	case "show":
		commandShow()
	case "install":
	case "help":
		printHelp()
	default:
		log.Errorf("unknown command: '%s'\n", c.Command)
		printHelp()
	}
	log.Debug("Process finished in DEBUG mode")
}
