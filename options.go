package main

import (
	"code.google.com/p/gcfg"
	"github.com/Hranoprovod/api-client"
	"github.com/Hranoprovod/processor"
	"github.com/Hranoprovod/reporter"
	"github.com/codegangsta/cli"
	"os"
	"os/user"
	"time"
)

const (
	optionsFileName = "/.hranoprovod/config"
)

// Options contains the options structure
type Options struct {
	Global struct {
		DbFileName  string
		LogFileName string
		DateFormat  string
	}
	Resolver struct {
		ResolverMaxDepth int
	}
	Processor processor.Options
	Reporter  reporter.Options
	API       client.Options
}

// NewOptions returns new options structure.
func NewOptions() *Options {
	o := &Options{}
	o.Reporter = *reporter.NewDefaultOptions()
	o.Processor = *processor.NewDefaultOptions()
	o.API = *client.NewDefaultOptions()
	return o
}

// Load loads the settigns from config file/command line params/defauls from given context.
func (o *Options) Load(c *cli.Context) *Options {
	fileName := c.GlobalString("o")
	// First try to load the o file
	if exists(fileName) {
		if err := gcfg.ReadFileInto(o, fileName); err != nil {
			// o file is not valid
			panic(err)
		}
	}
	o.populateGlobals(c)
	o.populateLocals(c)
	return o
}

// GetDefaultFileName returns the default filename for the config file
func GetDefaultFileName() string {
	usr, err := user.Current()
	if err != nil {
		return ""
	}
	return usr.HomeDir + optionsFileName
}

func exists(name string) bool {
	_, err := os.Stat(name)
	return !os.IsNotExist(err)
}

func (o *Options) populateGlobals(c *cli.Context) {
	if c.GlobalIsSet("database") || o.Global.DbFileName == "" {
		o.Global.DbFileName = c.GlobalString("database")
	}

	if c.GlobalIsSet("logfile") || o.Global.LogFileName == "" {
		o.Global.LogFileName = c.GlobalString("logfile")
	}

	if c.GlobalIsSet("date-format") || o.Global.DateFormat == "" {
		o.Global.DateFormat = c.GlobalString("date-format")
	}
}

func (o *Options) populateLocals(c *cli.Context) {
	o.populateResolver(c)
	o.populateProcessor(c)
	o.populateReporter(c)
}

func (o *Options) populateResolver(c *cli.Context) {
	if c.IsSet("maxdepth") || o.Resolver.ResolverMaxDepth == 0 {
		o.Resolver.ResolverMaxDepth = c.Int("maxdepth")
	}
}

func (o *Options) populateProcessor(c *cli.Context) {
	var err error

	if c.IsSet("beginning") {
		o.Processor.BeginningTime, err = time.Parse(o.Global.DateFormat, c.String("beginning"))
		if err != nil {
			panic(err)
		}
		o.Processor.HasBeginning = true
	}

	if c.IsSet("end") {
		o.Processor.EndTime, err = time.Parse(o.Global.DateFormat, c.String("end"))
		if err != nil {
			panic(err)
		}
		o.Processor.HasEnd = true
	}

	o.Processor.Unresolved = c.Bool("unresolved")
	o.Processor.SingleFood = c.String("single-food")
	o.Processor.SingleElement = c.String("single-element")
}

func (o *Options) populateReporter(c *cli.Context) {
	if c.IsSet("csv") {
		o.Reporter.CSV = true
	}
	if c.IsSet("no-color") {
		o.Reporter.Color = false
	}

	if c.IsSet("no-totals") {
		o.Processor.Totals = false
	}
}