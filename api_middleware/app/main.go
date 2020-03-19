package main

import (
	"fmt"
	"github.com/Semior001/mdcd-travelhack/app/cmd"
	"log"
	"os"

	"github.com/hashicorp/logutils"
	"github.com/jessevdk/go-flags"
)

const revision = "unknown"

type Opts struct {
	ServerCmd cmd.ServerCmd `command:"server"`
	Dbg       bool          `long:"dbg" env:"DEBUG" description:"debug mode"`
}

var logFlags int = log.Ldate | log.Ltime

func main() {
	fmt.Printf("mdcd_api_middleware revision %s", revision)

	var opts Opts
	p := flags.NewParser(&opts, flags.Default)

	p.CommandHandler = func(command flags.Commander, args []string) error {
		setupLog(opts.Dbg)
		err := command.Execute(args)
		if err != nil {
			log.Printf("[ERROR] failed to execute command %+v", err)
		}
		return nil
	}

	if _, err := p.Parse(); err != nil {
		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			os.Exit(0)
		} else {
			os.Exit(1)
		}
	}

}

func setupLog(dbg bool) {
	filter := &logutils.LevelFilter{
		Levels:   []logutils.LogLevel{"DEBUG", "INFO", "WARN", "ERROR"},
		MinLevel: "INFO",
		Writer:   os.Stdout,
	}

	if dbg {
		logFlags = log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile
		filter.MinLevel = "DEBUG"
	}

	log.SetFlags(logFlags)
	log.SetOutput(filter)
}
