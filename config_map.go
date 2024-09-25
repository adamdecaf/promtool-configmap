package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/prometheus/prometheus/pkg/rulefmt"
)

// ConfigMap is the minimum representation of a kubernetes ConfigMap object
// for us to properly parse and validate the nested rules (or prometheus config)
//
// https://v1-8.docs.kubernetes.io/docs/api-reference/v1.8/#configmap-v1-core
type ConfigMap struct {
	ApiVersion string            `json:"apiVersion" yaml:"apiVersion"`
	Kind       string            `json:"kind" yaml:"kind"`
	Data       map[string]string `json:"data" yaml:"data"`
}

func (c ConfigMap) validate() error {
	if c.Kind != "ConfigMap" {
		return fmt.Errorf("Unknown kind: %s", c.Kind)
	}
	if len(c.Data) == 0 {
		return errors.New("empty ConfigMap")
	}

	// Check each yaml blob
	for k, v := range c.Data {
		v = strings.TrimSpace(v)
		if v == "" {
			return fmt.Errorf("%s contained nothing", k)
		}

		e1 := checkAsPromConfig(v)
		e2 := checkAsPromRules(v)

		// TODO(adam): clean this up, virtual fs to drop
		// extra files into?
		if e1 != nil && strings.Contains(e1.Error(), "does not point to an existing file") {
			e1 = nil
		}

		// See if we've failed
		if e1 == nil || e2 == nil {
			// return nil if one check passed
			return nil
		}
		if e1 != nil {
			return fmt.Errorf("when validating %s prom config: %v", k, e1)
		}
		return fmt.Errorf("when validating %s rule file: %v", k, e2)
	}
	return nil
}

func checkAsPromConfig(raw string) error {
	// For now, just find promtool on PATH // TODO(adam)
	_, err := exec.LookPath("promtool")
	if err == nil {
		fd, err := os.CreateTemp("", "promtool-configmap")
		if err != nil {
			return err
		}
		defer os.Remove(fd.Name())

		n, err := io.Copy(fd, strings.NewReader(raw))
		if err != nil || n == 0 {
			return fmt.Errorf("problem copying 'prom config', n=%d, err=%v", n, err)
		}
		out, err := exec.Command("promtool", "check", "config", fd.Name()).CombinedOutput() //nolint:gosec
		if err != nil {
			return errors.New(string(out))
		}

		// verbose output
		if *flagVerbose {
			fmt.Printf("\n\n promtool check config:\n")
			fmt.Println(string(out))
		}
	}
	if err != nil {
		return err
	}
	return nil
}

func checkAsPromRules(raw string) error {
	groups, errs := rulefmt.Parse([]byte(raw))
	if len(errs) != 0 {
		var buf strings.Builder
		for i := range errs {
			buf.WriteString(fmt.Sprintf(" %v\n", errs[i].Error()))
		}
		return errors.New(buf.String())
	}
	for i := range groups.Groups {
		rules := groups.Groups[i].Rules
		for j := range rules {
			if err := rules[j].Validate(); err != nil {
				return fmt.Errorf("when validating %s: %v", groups.Groups[i].Name, err)
			}
		}
	}
	return nil
}
