package modules

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/Ak-Army/xlog"
	"gopkg.in/yaml.v3"
)

func init() {
	addModule("kubectl", func() Module {
		return &kubectl{
			Base: Get(),
		}
	})
}

type kubectl struct {
	*Base

	// Output
	Context string
	kubeContext
}

type kubeConfig struct {
	CurrentContext string `yaml:"current-context"`
	Contexts       []struct {
		Context *kubeContext `yaml:"context"`
		Name    string       `yaml:"name"`
	} `yaml:"contexts"`
}

type kubeContext struct {
	Cluster   string `yaml:"cluster"`
	User      string `yaml:"user"`
	Namespace string `yaml:"namespace"`
}

func (k *kubectl) Init() error {
	paths := append(strings.Split(os.Getenv("KUBECONFIG"), ":"), filepath.Join(k.Base.homeDir(), ".kube", "config"))
	contexts := make(map[string]*kubeContext)
	for _, kubeconfig := range paths {
		if len(kubeconfig) == 0 {
			continue
		}
		content, err := ioutil.ReadFile(kubeconfig)
		if err != nil {
			xlog.Warnf("Unable to read kube prompt-line: %s", kubeconfig, err)
			continue
		}
		var config kubeConfig
		err = yaml.Unmarshal(content, &config)
		if err != nil {
			xlog.Warnf("Unable to parse kube prompt-line: %s", kubeconfig, err)
			continue
		}
		for _, context := range config.Contexts {
			if _, exists := contexts[context.Name]; !exists {
				contexts[context.Name] = context.Context
			}
		}
		if len(k.Context) == 0 {
			k.Context = config.CurrentContext
		}
		context, exists := contexts[k.Context]
		if !exists {
			continue
		}
		if context != nil {
			k.kubeContext = *context
			return nil
		}
	}
	xlog.Info("Kube prompt-line not found use command")
	ctx, err := k.Base.runCommand("kubectl", "prompt-line", "current-context")
	if err != nil {
		return err
	}
	k.Context = strings.TrimSuffix(ctx, "\n")
	namespace, _ := k.Base.runCommand(
		"kubectl",
		"prompt-line",
		"view",
		"-o",
		fmt.Sprintf("jsonpath='{.contexts[?(@.name == \"%s\")].context.namespace}'", k.Context),
	)
	k.Namespace = strings.Trim(namespace, "'")

	return nil
}
