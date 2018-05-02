package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/prometheus/prometheus/pkg/rulefmt"
	"gopkg.in/yaml.v2"
)

// ConfigMap is the minimum representation of a kubernetes ConfigMap object
// for us to properly parse and validate the nested rules (or prometheus config)
//
// https://v1-8.docs.kubernetes.io/docs/api-reference/v1.8/#configmap-v1-core
type ConfigMap struct {
	ApiVersion string            `json:"apiVersion", yaml:"apiVersion"`
	Kind       string            `json:"kind", yaml:"kind"`
	Data       map[string]string `json:"data", yaml:"data"`
}

func (c ConfigMap) validate() error {
	// if c.ApiVersion != "v1" { // TODO(adam): why does this fail?
	// 	return fmt.Errorf("unknown apiVersion %q", c.ApiVersion)
	// }
	if c.Kind != "ConfigMap" {
		return fmt.Errorf("got other k8s object %s", c.Kind)
	}
	if len(c.Data) == 0 {
		return errors.New("empty ConfigMap")
	}

	// Check each yaml blob
	for k,v := range c.Data {
		groups, errs := rulefmt.Parse([]byte(v))
		if len(errs) != 0 {
			buf := strings.Builder{}
			buf.WriteString("error validating rule file\n")
			for i := range errs {
				buf.WriteString(fmt.Sprintf(" %v\n", errs[i].Error()))
			}
			return errors.New(buf.String())
		}

		if err := groups.Validate(); err != nil {
			return fmt.Errorf("error validating %s, err=%v", k, err)
		}
	}
	return nil
}

func main() {
	if len(os.Args) == 1 || (len(os.Args) == 2 && strings.ToLower(os.Args[1]) == "help") {
		showHelp()
		os.Exit(1)
	}

	foundErrors := false

	for i := range os.Args[1:] {
		rawPath := os.Args[1:][i]

		// read from stdin if we see '--'
		if rawPath == "--" {
			if err := check(os.Stdin); err != nil {
				foundErrors = true
				fmt.Printf("ERROR validating rules: %v\n", err)
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

		if err := check(f); err != nil {
			foundErrors = true
			fmt.Printf("ERROR validating rules for %s, err=%v\n", path, err)
			continue
		}

		fmt.Printf("%s passed checks\n", rawPath)
	}

	if foundErrors {
		os.Exit(1)
	}
}

func showHelp() {
	fmt.Println(`Usage of promtool-rules-configmap

This tool is a utility to run the same promtool config and rule checks against a kubernetes ConfigMap object.

USAGE

  promtool-rules-configmap [file ...]

  cat rules.yaml | promtool-rules-configmap --`)
}

func check(r io.Reader) error {
	// be naive and read whole file
	bs, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}

	var cfg ConfigMap

	// read as json
	err = json.Unmarshal(bs, &cfg)
	if err == nil {
		return cfg.validate()
	}

	// read as yaml
	err = yaml.Unmarshal(bs, &cfg)
	if err == nil {
		return cfg.validate()
	}

	return errors.New("unable to parse file as json or yaml")
}
