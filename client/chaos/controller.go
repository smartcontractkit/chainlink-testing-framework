// Package chaos handles the setup necessary to run chaos tests across chainlnk nodes and blockchains
package chaos

import (
	"bytes"
	"context"
	"fmt"
	"github.com/smartcontractkit/integrations-framework/client/chaos/experiments"
	"github.com/smartcontractkit/integrations-framework/utils"
	"io/ioutil"
	"path/filepath"
	"text/template"

	"github.com/ghodss/yaml"
	"github.com/rs/zerolog/log"
	uuid "github.com/satori/go.uuid"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

const (
	// APIBasePath in form of /apis/<spec.group>/<spec.versions.name>, see Chaosmesh CRD 2.0.0
	APIBasePath = "/apis/chaos-mesh.org/v1alpha1"
	// TemplatesPath path to the chaos templates
	TemplatesPath = "client/chaos/templates"
)

// Experimentable interface for chaos experiments
type Experimentable interface {
	SetBase(base experiments.Base)
	Filename() string
	Resource() string
}

// Controller is controller that manages Chaosmesh CRD instances to run experiments
type Controller struct {
	Client   *kubernetes.Clientset
	Requests map[string]*rest.Request
	Cfg      *Config
}

// Config Chaosmesh controller config
type Config struct {
	Client    *kubernetes.Clientset
	Namespace string
}

// NewController creates controller to run and stop chaos experiments
func NewController(cfg *Config) (*Controller, error) {
	return &Controller{
		Client:   cfg.Client,
		Requests: make(map[string]*rest.Request),
		Cfg:      cfg,
	}, nil
}

// Run runs experiment and saves it's ID
func (c *Controller) Run(exp Experimentable) (string, error) {
	name := fmt.Sprintf("%s-%s", exp.Resource(), uuid.NewV4().String())
	exp.SetBase(experiments.Base{
		Name:      name,
		Namespace: c.Cfg.Namespace,
	})
	fileBytes, err := ioutil.ReadFile(filepath.Join(utils.ProjectRoot, TemplatesPath, exp.Filename()))
	if err != nil {
		return "", err
	}
	d, err := marshallTemplate(exp, "Chaos template", string(fileBytes))
	if err != nil {
		return "", err
	}
	data, err := yaml.YAMLToJSON([]byte(d))
	if err != nil {
		return "", err
	}
	log.Info().Str("Name", name).Str("Resource", exp.Resource()).Msg("Starting chaos experiment")
	req := c.Client.RESTClient().
		Post().
		AbsPath(APIBasePath).
		Name(name).
		Namespace(c.Cfg.Namespace).
		Resource(exp.Resource()).
		Body(data)
	resp := req.Do(context.Background())
	if resp.Error() != nil {
		return "", err
	}
	c.Requests[name] = req
	return name, nil
}

// Stop removes experiment's entity
func (c *Controller) Stop(name string) error {
	log.Info().Str("ID", name).Msg("Deleting chaos experiment")
	exp, ok := c.Requests[name]
	if !ok {
		return fmt.Errorf("experiment %s not found", name)
	}
	res := exp.Verb("DELETE").Do(context.Background())
	if res.Error() != nil {
		return res.Error()
	}
	delete(c.Requests, name)
	return nil
}

// StopAll removes all experiments entities
func (c *Controller) StopAll() error {
	for id := range c.Requests {
		err := c.Stop(id)
		if err != nil {
			return err
		}
	}
	return nil
}

// marshallTemplate Helper to marshall templates
func marshallTemplate(jobSpec interface{}, name, templateString string) (string, error) {
	var buf bytes.Buffer
	tmpl, err := template.New(name).Parse(templateString)
	if err != nil {
		return "", err
	}
	err = tmpl.Execute(&buf, jobSpec)
	if err != nil {
		return "", err
	}
	return buf.String(), err
}
