package types

import (
	"github.com/logrusorgru/aurora/v4"
	"github.com/projectdiscovery/gologger"
	"github.com/projectdiscovery/gologger/formatter"
	"github.com/projectdiscovery/gologger/levels"
)

var (
	AU            *aurora.Aurora
	jsonOutput    bool
	silentOutput  bool
	noColorOutput bool
)

func NewOptions() *Options {
	return &Options{
		Kubernetes: &OptionsKubernetes{},
		Output:     &OptionsOutput{},
	}
}

func IsJSONOutput() bool    { return jsonOutput }
func IsSilentOutput() bool  { return silentOutput }
func IsNoColorOutput() bool { return noColorOutput }

type Options struct {
	Kubernetes *OptionsKubernetes
	Output     *OptionsOutput
}

func (o *Options) Validate() {
	if o.Output.Verbose && o.Output.Silent {
		gologger.Fatal().Msg("verbose and silent output selected")
	}
}

func (o *Options) Configure() {
	jsonOutput = o.Output.JSON
	silentOutput = o.Output.Silent
	noColorOutput = o.Output.NoColor

	// singleton
	if AU == nil {
		AU = aurora.New(aurora.WithColors(!o.Output.NoColor))
	}

	gologger.DefaultLogger.SetTimestamp(true, levels.LevelSilent)
	gologger.DefaultLogger.SetFormatter(formatter.NewCLI(o.Output.NoColor))

	if o.Output.Verbose {
		gologger.DefaultLogger.SetMaxLevel(levels.LevelDebug)
	} else if o.Output.Silent {
		gologger.DefaultLogger.SetMaxLevel(levels.LevelSilent)
	} else {
		gologger.DefaultLogger.SetMaxLevel(levels.LevelWarning)
	}

	if o.Output.JSON {
		gologger.DefaultLogger.SetFormatter(&formatter.JSON{})
		AU = aurora.New(aurora.WithColors(false)) // when json, always no color
	}
}

type OptionsExecution struct {
	Batch       bool
	Interactive bool
}

type OptionsKubernetes struct {
	AuthToken       string
	BaseURL         string
	CustomUserAgent string
	IgnoreTLS       bool
	KubeConfigPath  string
	Namespace       string
	NoRunSafe       bool
}

type OptionsOutput struct {
	Verbose bool
	Silent  bool
	JSON    bool
	NoColor bool
}
