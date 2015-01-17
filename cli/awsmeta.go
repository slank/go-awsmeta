package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/slank/go-awsmeta"
)

type CliOpts struct {
	ShowVersion bool
	ApiVersion  string
	Timeout     int
	ShortName   string
	ListNames   bool
}

func Usage() {
	fmt.Fprintf(os.Stderr, "usage: awsmeta [options] -n shortname\n")
	fmt.Fprintf(os.Stderr, "       awsmeta [options] /metadata/path\n")
	fmt.Fprintf(os.Stderr, "       awsmeta --list\n")
	fmt.Fprintf(os.Stderr, "       awsmeta --version\n")
}

func setup_flags(arguments []string) (*flag.FlagSet, *CliOpts) {
	var opts CliOpts
	flags := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	flags.BoolVar(&opts.ShowVersion, "version", false, "Show the awsmeta version")
	flags.StringVar(
		&opts.ApiVersion, "v", awsmeta.DEFAULT_API_VERSION,
		fmt.Sprintf("Metadata API version (default=%d)", awsmeta.DEFAULT_API_VERSION),
	)
	flags.IntVar(
		&opts.Timeout, "t", int(awsmeta.DEFAULT_TIMEOUT/time.Second),
		fmt.Sprintf("HTTP request timeout (seconds) (default=%d)", awsmeta.DEFAULT_TIMEOUT),
	)
	flags.StringVar(&opts.ShortName, "n", "", "Look up a named metadata value")
	flags.BoolVar(&opts.ListNames, "list", false, "List defined metadata names")

	flags.Parse(arguments)

	return flags, &opts
}

func show_version() {
	fmt.Printf("CLI Version: %s\n", VERSION)
	fmt.Printf("Awsmeta Version: %s\n", awsmeta.VERSION)
}

func main() {
	flags, opts := setup_flags(os.Args[1:])

	if opts.ShowVersion {
		show_version()
		os.Exit(0)
	}

	if opts.ListNames {
		for name, _ := range awsmeta.ShortNames {
			fmt.Println(name)
		}
		os.Exit(0)
	}

	var md awsmeta.MetaDataServer
	md.Timeout = opts.Timeout
	md.ApiVersion = opts.ApiVersion

	var result string
	var err error
	if opts.ShortName != "" {
		path := awsmeta.ShortNames[opts.ShortName]
		if path == "" {
			fmt.Fprintf(os.Stderr, "No such shortname: %s\n", opts.ShortName)
			os.Exit(1)
		}
		result, err = md.Get(path)
	} else {
		args := flags.Args()
		if len(args) == 0 {
			result, err = md.Get("")
		} else {
			result, err = md.Get(args[0])
		}
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Println(result)
	os.Exit(0)
}
