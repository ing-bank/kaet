package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"runtime/debug"

	"github.com/ing-bank/kaet/pkg/runner"
	"github.com/ing-bank/kaet/pkg/types"
	"github.com/projectdiscovery/goflags"
	"github.com/projectdiscovery/gologger"
	"k8s.io/client-go/util/homedir"
)

var options = types.NewOptions()

func Run() {
	set := goflags.NewFlagSet()
	configureCLI(set)
	defer handlePanic(set)

	options.Validate()
	options.Configure()

	/*
		1. Parse config file location (yaml file)
		2. Verify pre/post-requisites
		3. Run steps
		4. Gather results
		5. Write results
	*/

	gologger.Silent().Msgf("%s\n\n%s\n\n", types.KAETBanner, types.Disclaimer)

	runner := runner.FromOptions(options)
	defer runner.Close()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			gologger.Info().Msgf("CTRL+C pressed: Exiting\n")
			runner.Close()
			os.Exit(0)
		}
	}()

	runner.Run()
}

func configureCLI(set *goflags.FlagSet) {
	set.SetDescription(types.KAETBanner)

	setGroup(set, "kubernetes", "kubernetes options",
		set.StringVarP(&options.Kubernetes.BaseURL, "u", "base-url", "https://kubernetes.default.svc", "kubernetes API base url"),
		set.StringVarP(&options.Kubernetes.AuthToken, "t", "auth-token", "", "kubernetes service account token"),
		set.BoolVarP(&options.Kubernetes.IgnoreTLS, "ignore-tls", "k", false, "ignore TLS"),
		set.StringVarP(&options.Kubernetes.CustomUserAgent, "user-agent", "ua", "KAET", "custom user agent"),
		set.StringVarP(&options.Kubernetes.Namespace, "namespace", "n", "", "kubernetes namespace"),
		set.BoolVar(&options.Kubernetes.NoRunSafe, "no-safe", false, "do not explore control namespaces"),
	)

	if home := homedir.HomeDir(); home != "" {
		set.
			StringVar(
				&options.Kubernetes.KubeConfigPath,
				"kubeconfig",
				filepath.Join(home, ".kube", "config"),
				"absolute path to kubeconfig file",
			).
			Group("kubernetes")
	} else {
		set.
			StringVar(
				&options.Kubernetes.KubeConfigPath,
				"kubeconfig",
				"",
				"absolute path to kubeconfig file",
			).
			Group("kubernetes")
	}

	setGroup(set, "output", "output options",
		set.BoolVarP(&options.Output.Verbose, "verbose", "v", false, "verbose output"),
		set.BoolVarP(&options.Output.Silent, "silent", "s", false, "silent output"),
		set.BoolVarP(&options.Output.JSON, "json", "j", false, "json output"),
		set.BoolVarP(&options.Output.NoColor, "no-color", "nc", false, "colorful output"),
	)

	_ = set.Parse()
}

func setGroup(set *goflags.FlagSet, groupName, description string, flags ...*goflags.FlagData) {
	set.SetGroup(groupName, description)
	for _, currentFlag := range flags {
		currentFlag.Group(groupName)
	}
}

func handlePanic(_ *goflags.FlagSet) {
	if err := recover(); err != nil {
		gologger.Error().Msgf("something bad happened... panic error: %s\n", err)

		if options.Output.Verbose {
			fmt.Printf("\ndebug stack: %s\n", debug.Stack())
		}
	}
}
