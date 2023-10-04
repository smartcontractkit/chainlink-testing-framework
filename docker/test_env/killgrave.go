package test_env

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"path"
	"runtime"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	tc "github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/smartcontractkit/chainlink-testing-framework/logging"
)

type Killgrave struct {
	EnvComponent
	ExternalEndpoint  string
	InternalPort      string
	InternalEndpoint  string
	InternalImposters []*url.URL
	impostersPath     string
	t                 *testing.T
	l                 zerolog.Logger
}

// Imposter define an imposter structure
type KillgraveImposter struct {
	Request  KillgraveRequest  `json:"request"`
	Response KillgraveResponse `json:"response"`
}

type KillgraveRequest struct {
	Method     string             `json:"method"`
	Endpoint   string             `json:"endpoint,omitempty"`
	SchemaFile *string            `json:"schemaFile,omitempty"`
	Params     *map[string]string `json:"params,omitempty"`
	Headers    *map[string]string `json:"headers"`
}

// Response represent the structure of real response
type KillgraveResponse struct {
	Status   int                     `json:"status"`
	Body     string                  `json:"body,omitempty"`
	BodyFile *string                 `json:"bodyFile,omitempty"`
	Headers  *map[string]string      `json:"headers,omitempty"`
	Delay    *KillgraveResponseDelay `json:"delay,omitempty"`
}

// ResponseDelay represent time delay before server responds.
type KillgraveResponseDelay struct {
	Delay  int64 `json:"delay,omitempty"`
	Offset int64 `json:"offset,omitempty"`
}

// AdapterResponse represents a response from an adapter
type KillgraveAdapterResponse struct {
	Id    string                 `json:"id"`
	Data  KillgraveAdapterResult `json:"data"`
	Error interface{}            `json:"error"`
}

// AdapterResult represents an int result for an adapter
type KillgraveAdapterResult struct {
	Result interface{} `json:"result"`
}

func NewKillgrave(networks []string, impostersDirectoryPath string, opts ...EnvComponentOption) *Killgrave {
	k := &Killgrave{
		EnvComponent: EnvComponent{
			ContainerName: fmt.Sprintf("%s-%s", "killgrave", uuid.NewString()[0:3]),
			Networks:      networks,
		},
		InternalPort:  "3000",
		impostersPath: impostersDirectoryPath,
		l:             log.Logger,
	}
	for _, opt := range opts {
		opt(&k.EnvComponent)
	}
	return k
}

func (k *Killgrave) WithTestLogger(t *testing.T) *Killgrave {
	k.l = logging.GetTestLogger(t)
	k.t = t
	return k
}

func (k *Killgrave) StartContainer() error {
	l := tc.Logger
	if k.t != nil {
		l = logging.CustomT{
			T: k.t,
			L: k.l,
		}
	}
	c, err := tc.GenericContainer(context.Background(), tc.GenericContainerRequest{
		ContainerRequest: k.getContainerRequest(),
		Started:          true,
		Reuse:            true,
		Logger:           l,
	})
	if err != nil {
		return errors.Wrapf(err, "cannot start Killgrave container")
	}
	endpoint, err := c.Endpoint(context.Background(), "http")
	if err != nil {
		return err
	}
	k.Container = c
	k.ExternalEndpoint = endpoint
	k.InternalEndpoint = fmt.Sprintf("http://%s:%s", k.ContainerName, k.InternalPort)

	log.Info().Str("External Endpoint", k.ExternalEndpoint).
		Str("Internal Endpoint", k.InternalEndpoint).
		Str("Container Name", k.ContainerName).
		Msgf("Started Killgrave Container")
	return nil
}

func (k *Killgrave) getContainerRequest() tc.ContainerRequest {
	if len(k.impostersPath) == 0 {
		_, f, _, _ := runtime.Caller(0)
		k.impostersPath = path.Join(path.Dir(f), "/killgrave_imposters")
	}
	return tc.ContainerRequest{
		Name:         k.ContainerName,
		Networks:     k.Networks,
		Image:        "friendsofgo/killgrave",
		ExposedPorts: []string{NatPortFormat(k.InternalPort)},
		Cmd:          []string{"-host=0.0.0.0", "-imposters=/imposters", "-watcher"},
		Mounts: tc.ContainerMounts{
			tc.ContainerMount{
				Source: tc.GenericBindMountSource{
					HostPath: k.impostersPath,
				},
				Target: "/imposters",
			},
		},
		WaitingFor: wait.ForLog("The fake server is on tap now"),
	}
}

// AddImposter adds an imposter to the killgrave container
func (k *Killgrave) AddImposter(imposter KillgraveImposter) error {
	req := imposter.Request

	imposters := []KillgraveImposter{imposter}
	data, err := json.Marshal(imposters)
	if err != nil {
		return err
	}

	f, err := os.CreateTemp("", "imposter.imp.json")
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.WriteString(string(data))
	if err != nil {
		return err
	}

	localFile := f.Name()
	containerFile := fmt.Sprintf("/imposters%s.imp.json", req.Endpoint)

	err = k.Container.CopyFileToContainer(context.Background(), localFile, containerFile, 0644)
	if err != nil {
		return err
	}

	// wait for the log saying the imposter was loaded
	logWaitStrategy := wait.ForLog(fmt.Sprintf("imposter %s loaded", containerFile)).WithStartupTimeout(15 * time.Second)
	err = logWaitStrategy.WaitUntilReady(context.Background(), k.Container)
	return err
}

// SetStringValuePath sets a path to return a string value
func (k *Killgrave) SetStringValuePath(path string, method string, headers map[string]string, v string) error {
	imp := KillgraveImposter{
		Request: KillgraveRequest{
			Method:   method,
			Endpoint: path,
		},
		Response: KillgraveResponse{
			Status:  200,
			Body:    v,
			Headers: &headers,
		},
	}

	return k.AddImposter(imp)
}

// SetAdapterBasedAnyValuePath sets a path to return a value as though it was from an adapter
func (k *Killgrave) SetAdapterBasedAnyValuePath(path string, method string, v interface{}) error {
	ar := KillgraveAdapterResponse{
		Id: "",
		Data: KillgraveAdapterResult{
			Result: v,
		},
		Error: nil,
	}
	data, err := json.Marshal(ar)
	if err != nil {
		return err
	}

	return k.SetStringValuePath(path, method, map[string]string{
		"Content-Type": "application/json",
	}, string(data))
}

// SetAdapterBasedAnyValuePathObject sets a path to return a value as though it was from an adapter
func (k *Killgrave) SetAdapterBasedIntValuePath(path string, method string, v int) error {
	return k.SetAdapterBasedAnyValuePath(path, method, v)
}
