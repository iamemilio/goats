package input

import (
	"flag"
	"fmt"
	"os"
)

const (
	appName = "goats"
)

type RunInputs struct {
	ConfigFile []byte
}

func (inputs *RunInputs) Get() error {
	var configFile = flag.String("c", "", "the path to your runConfig.yaml file")
	flag.Parse()

	if configFile == nil || *configFile == "" {
		return fmt.Errorf("a config file must be passed: %s -c configFile.yaml", appName)
	}

	file, err := os.ReadFile(*configFile)
	if err != nil {
		return err
	}

	inputs.ConfigFile = file
	return nil
}

func New() RunInputs {
	return RunInputs{}
}
