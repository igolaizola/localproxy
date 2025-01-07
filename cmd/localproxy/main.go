package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime/debug"
	"strings"

	"github.com/igolaizola/localproxy"
	"github.com/peterbourgon/ff/v3"
	"github.com/peterbourgon/ff/v3/ffcli"
	"github.com/peterbourgon/ff/v3/ffyaml"
)

// Build flags
var version = ""
var commit = ""
var date = ""

func main() {
	// Create signal based context
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	// Launch command
	cmd := newCommand()
	if err := cmd.ParseAndRun(ctx, os.Args[1:]); err != nil {
		log.Fatal(err)
	}
}

func newCommand() *ffcli.Command {
	fs := flag.NewFlagSet("localproxy", flag.ExitOnError)

	var debug bool
	var upstream, addr string
	fs.StringVar(&upstream, "upstream", "", "upstream URL")
	fs.StringVar(&addr, "addr", "0.0.0.0:0", "listen address")
	fs.BoolVar(&debug, "debug", false, "enable debug logging")

	return &ffcli.Command{
		ShortUsage: "localproxy [flags]",
		FlagSet:    fs,
		Options: []ff.Option{
			ff.WithConfigFileFlag("config"),
			ff.WithConfigFileParser(ffyaml.Parser),
			ff.WithEnvVarPrefix("LOCALPROXY"),
		},
		Exec: func(ctx context.Context, args []string) error {
			log.Println("starting local proxy")
			localProxyURL, err := localproxy.Run(ctx, debug, upstream, addr)
			if err != nil {
				return err
			}
			log.Println("local proxy running at", localProxyURL)
			<-ctx.Done()
			return nil
		},
		Subcommands: []*ffcli.Command{
			newVersionCommand(),
		},
	}
}

func newVersionCommand() *ffcli.Command {
	return &ffcli.Command{
		Name:       "version",
		ShortUsage: "localproxy version",
		ShortHelp:  "print version",
		Exec: func(ctx context.Context, args []string) error {
			v := version
			if v == "" {
				if buildInfo, ok := debug.ReadBuildInfo(); ok {
					v = buildInfo.Main.Version
				}
			}
			if v == "" {
				v = "dev"
			}
			versionFields := []string{v}
			if commit != "" {
				versionFields = append(versionFields, commit)
			}
			if date != "" {
				versionFields = append(versionFields, date)
			}
			fmt.Println(strings.Join(versionFields, " "))
			return nil
		},
	}
}
