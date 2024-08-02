package cliflags

import (
	"fmt"
	"github.com/jessevdk/go-flags"
	"os"
)

type FlagOptions struct {
	LogFile   string `long:"logfile" description:"Log file to write to" default:"titbits.log" required:"false"`
	LogToFile bool   `long:"logtofile" description:"Log to file" required:"false"`
}

func ParseFlags() *FlagOptions {
	opts := FlagOptions{}
	_, err := flags.ParseArgs(&opts, os.Args)

	if err != nil {
		fmt.Println("failed to parse flags", err)
		os.Exit(2)
	}
	return &opts
}
