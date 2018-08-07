package common

import (
	"errors"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/micro/go-os/config"
)

func combinePath(basePath, relPath string) string {
	if relPath == "" || relPath[0] == '/' {
		return relPath
	}
	return path.Join(path.Dir(basePath), relPath)
}

// Pops a parameter from os.Args and removes it form there, to not spoil the rest argsparsing
func PopParameter(name string) (string, error) {
	var param_idx int
	for idx, arg := range os.Args {
		if strings.HasSuffix(arg, name) {
			param_idx = idx
			break
		}
	}

	if param_idx == 0 {
		return "", errors.New(fmt.Sprintf("%v argument not found", name))
	}

	var newArgs []string
	oldArgs := make([]string, len(os.Args))

	if param_idx+2 > len(os.Args) {
		return "", errors.New(fmt.Sprintf("%v path must be in format '-%v val'", name, name))
	}

	copy(oldArgs, os.Args)
	os.Args = os.Args[:param_idx+2]
	newArgs = append(oldArgs[:param_idx], oldArgs[param_idx+2:]...)

	paramVal := os.Args[param_idx+1]
	if !strings.HasPrefix(paramVal, "./") {
		paramVal = combinePath(os.Args[0], paramVal)
	}
	os.Args = newArgs
	return paramVal, nil
}

// Find a parameter from os.Args and returns its value or empty string
func FindParameter(name string) (string, error) {
	var param_idx int
	for idx, arg := range os.Args {
		if strings.HasPrefix(arg, name) {
			param_idx = idx
			break
		}
	}
	if param_idx == 0 {
		return "", errors.New(fmt.Sprintf("%v argument not found", name))
	}

	paramVal := os.Args[param_idx][len(name):]
	return paramVal, nil
}

// Return address to use for metrics service
func MetricAddress() string {
	collectorAddress := "localhost"
	metricsAddress, err := PopParameter("metrics_address")
	if err != nil && len(metricsAddress) > 0 {
		collectorAddress = metricsAddress
	}
	collectorAddress = fmt.Sprintf("%v:8125", collectorAddress)
	return collectorAddress
}

// Set dynamic value function
type SetValue func(v config.Value)

// Function that watches for changes in dynamic values and sets them
func ConfigValueWatcher(c config.Config, s SetValue, key ...string) {
	w, err := c.Watch(key...)
	if err != nil {
		return
	}

	go func() {
		for {
			v, err := w.Next()
			if err != nil {
				return
			}
			s(v)
		}
	}()
}
