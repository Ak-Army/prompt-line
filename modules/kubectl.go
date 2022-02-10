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
	contexts map[string]*kubeContext
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
	k.contexts = make(map[string]*kubeContext)
	for _, kc := range paths {
		if k.readKubeConfig(kc) {
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

func (k *kubectl) readKubeConfig(kc string) bool {
	if len(kc) == 0 {
		return false
	}
	content, err := ioutil.ReadFile(kc)
	if err != nil {
		xlog.Warnf("Unable to read kube prompt-line: %s", kc, err)
		return false
	}
	var config kubeConfig
	err = yaml.Unmarshal(content, &config)
	if err != nil {
		xlog.Warnf("Unable to parse kube prompt-line: %s", kc, err)
		return false
	}
	for _, context := range config.Contexts {
		if _, exists := k.contexts[context.Name]; !exists {
			k.contexts[context.Name] = context.Context
		}
	}
	if len(k.Context) == 0 {
		k.Context = config.CurrentContext
	}
	context, exists := k.contexts[k.Context]
	if !exists {
		return false
	}
	if context != nil {
		k.kubeContext = *context
		return true
	}
	return false
}
