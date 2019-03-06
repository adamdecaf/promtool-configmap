package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"
)

const Version = "0.3.0-dev"

var (
	flagVerbose = flag.Bool("verbose", false, "verbose output (show promtool output)")
	flagVersion = flag.Bool("version", false, fmt.Sprintf("Show the version (%s)", Version))
)

func main() {
	flag.Parse()

	// early exits
	if *flagVersion {
		fmt.Println(Version)
		os.Exit(0)
	}
	if flag.NArg() == 0 || strings.EqualFold(flag.Arg(0), "help") {
		showHelp()
		os.Exit(1)
	}

	args := flag.Args()
	var foundErrors bool
	for i := range args {
		rawPath := flag.Arg(i)

		// read from stdin if we see '--'
		if rawPath == "--" {
			if err := check(os.Stdin); err != nil {
				foundErrors = true
				fmt.Printf("ERROR validating rules: \n%v\n", err)
			}
			continue
		}

		// read arg as a filepath
		path, err := filepath.Abs(rawPath)
		if err != nil {
			foundErrors = true
			fmt.Printf("ERROR attempting to read %s, err=%v\n", rawPath, err)
			continue
		}

		f, err := os.Open(path)
		if err != nil {
			foundErrors = true
			fmt.Printf("ERROR opening %s, err=%v\n", path, err)
			continue
		}

		if *flagVerbose {
			fmt.Printf(" Checking %s...\n", rawPath)
		}

		if err := check(f); err != nil {
			foundErrors = true
			fmt.Printf("ERROR validating rules for %s \n%v\n", path, err)
			continue
		}

		fmt.Printf("%s passed checks\n", rawPath)
	}

	if foundErrors {
		os.Exit(1)
	}
}

func showHelp() {
	fmt.Println(`Usage of promtool-configmap

This tool is a utility to run the same promtool config and rule checks against a kubernetes ConfigMap object.

USAGE

  promtool-configmap [-verbose] [file ...]

  cat rules.yaml | promtool-configmap --`)
}

func check(r io.Reader) error {
	// be naive and read whole file
	bs, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}

	var cfg ConfigMap
	parts := bytes.Split(bs, []byte("---"))
	if len(parts) == 0 {
		return errors.New("no objecs found")
	}

	checked := 0
	for i := range parts {
		if bytes.Equal(parts[i], []byte("")) {
			continue
		}
		checked++

		// read as json
		if err := json.Unmarshal(parts[i], &cfg); err == nil {
			if cfg.Kind != "ConfigMap" {
				continue
			}
			if err := cfg.validate(); err == nil {
				return nil // return early if successful
			}
		}

		// read as yaml
		if err := yaml.Unmarshal(parts[i], &cfg); err == nil {
			if cfg.Kind != "ConfigMap" {
				continue
			}
			if err := cfg.validate(); err != nil {
				return fmt.Errorf("ERROR: %v", err)
			}
		} else {
			return errors.New("unknown format for ConfigMap, tried json and yaml")
		}
	}
	if checked > 0 {
		return nil
	}
	return errors.New("no objects checked")
}
