package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

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

const (
	TarInstallation = "tar"
)

var AllInstallationSources = []string{TarInstallation}

type Installation struct {
	Installer goInstaller
}

type ScriptParams struct {
	Command   string
	GoVersion string // FIXME make goData object here and validate format during parse
	nextArg   int

	InstallParams *Installation
}

func (s *ScriptParams) NextScriptArg() string {
	r := flag.Arg(s.nextArg)
	s.nextArg++
	return r
}

func (s *ScriptParams) FullfillInstallParams() {
	if *InstallationSource == "" {
		log.Fatalf("installation source doesn't chosen, use '--'src' flag")
	}
	switch *InstallationSource {
	case TarInstallation:
		tarArchivePath := s.NextScriptArg()
		s.InstallParams = &Installation{}
		var err error
		s.InstallParams.Installer, err = NewTarInstaller(tarArchivePath)
		if err != nil {
			log.Fatalf("Failed to setup tar installation %s", err)
		}
	default:
		log.Fatalf(
			"Unknown installation source: %v; only %s supported",
			*InstallationSource,
			strings.Join(AllInstallationSources, ", "),
		)
	}
}

var (
	Debug = flag.Bool("debug", false, "run script in debug mode")

	// Installation
	InstallationSource = flag.String("src", "", "source of go code")
)

func getScriptParams() ScriptParams {
	flag.Parse()
	if !*Debug {
		l = l.WithOptions(zap.IncreaseLevel(zapcore.InfoLevel))
	}
	log = l.Sugar()
	c := ScriptParams{
		Command:   flag.Arg(0),
		GoVersion: flag.Arg(1),
		nextArg:   2,
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
		c.FullfillInstallParams()
		cmdInstall(c.GoVersion, c.InstallParams.Installer)
	case "enable":
		cmdEnable(c.GoVersion)
	case "help":
		printHelp()
	default:
		log.Errorf("unknown command: '%s'\n", c.Command)
		printHelp()
	}
	log.Debug("Process finished in DEBUG mode")
}
