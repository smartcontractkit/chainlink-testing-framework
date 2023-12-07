package test_env

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/otiai10/copy"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	tc "github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/smartcontractkit/chainlink-testing-framework/logging"
	"github.com/smartcontractkit/chainlink-testing-framework/mirror"
	"github.com/smartcontractkit/chainlink-testing-framework/utils/testcontext"
)

type Killgrave struct {
	EnvComponent
	ExternalEndpoint    string
	InternalPort        string
	InternalEndpoint    string
	impostersPath       string
	impostersDirBinding string
	t                   *testing.T
	l                   zerolog.Logger
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

func (k *Killgrave) WithTestInstance(t *testing.T) *Killgrave {
	k.l = logging.GetTestLogger(t)
	k.t = t
	return k
}

func (k *Killgrave) StartContainer() error {
	err := k.setupImposters()
	if err != nil {
		return err
	}
	if k.t != nil {
		k.t.Cleanup(func() {
			os.RemoveAll(k.impostersDirBinding)
		})
	}
	l := logging.GetTestContainersGoTestLogger(k.t)
	cr, err := k.getContainerRequest()
	if err != nil {
		return err
	}
	c, err := tc.GenericContainer(testcontext.Get(k.t), tc.GenericContainerRequest{
		ContainerRequest: cr,
		Started:          true,
		Reuse:            true,
		Logger:           l,
	})
	if err != nil {
		return fmt.Errorf("cannot start Killgrave container: %w", err)
	}
	endpoint, err := GetEndpoint(testcontext.Get(k.t), c, "http")
	if err != nil {
		return err
	}
	k.Container = c
	k.ExternalEndpoint = endpoint
	k.InternalEndpoint = fmt.Sprintf("http://%s:%s", k.ContainerName, k.InternalPort)

	k.l.Info().Str("External Endpoint", k.ExternalEndpoint).
		Str("Internal Endpoint", k.InternalEndpoint).
		Str("Container Name", k.ContainerName).
		Msgf("Started Killgrave Container")
	return nil
}

func (k *Killgrave) getContainerRequest() (tc.ContainerRequest, error) {
	killgraveImage, err := mirror.GetImage("friendsofgo/killgrave")
	if err != nil {
		return tc.ContainerRequest{}, err
	}
	return tc.ContainerRequest{
		Name:         k.ContainerName,
		Networks:     k.Networks,
		Image:        killgraveImage,
		ExposedPorts: []string{NatPortFormat(k.InternalPort)},
		Cmd:          []string{"-host=0.0.0.0", "-imposters=/imposters", "-watcher"},
		Mounts: tc.ContainerMounts{
			tc.ContainerMount{
				Source: tc.GenericBindMountSource{
					HostPath: k.impostersDirBinding,
				},
				Target: "/imposters",
			},
		},
		WaitingFor: wait.ForLog("The fake server is on tap now"),
	}, nil
}

func (k *Killgrave) setupImposters() error {
	// create temporary directory for imposters
	var err error
	k.impostersDirBinding, err = os.MkdirTemp(k.impostersDirBinding, "imposters*")
	if err != nil {
		return err
	}
	k.l.Info().Str("Path", k.impostersDirBinding).Msg("Imposters directory created at")

	// copy user imposters
	if len(k.impostersPath) != 0 {
		err = copy.Copy(k.impostersPath, k.impostersDirBinding)
		if err != nil {
			return err
		}
	}

	// add default five imposter
	return k.SetAdapterBasedIntValuePath("/five", []string{http.MethodGet, http.MethodPost}, 5)
}

// AddImposter adds an imposter to the killgrave container
func (k *Killgrave) AddImposter(imposters []KillgraveImposter) error {
	// if the endpoint paths do not start with '/' then add it
	for i, imposter := range imposters {
		if !strings.HasPrefix(imposter.Request.Endpoint, "/") {
			imposter.Request.Endpoint = fmt.Sprintf("/%s", imposter.Request.Endpoint)
			imposters[i] = imposter
		}
	}

	req := imposters[0].Request
	data, err := json.Marshal(imposters)
	if err != nil {
		return err
	}

	// build the file name from the req.Endpoint
	unsafeFileName := strings.TrimPrefix(req.Endpoint, "/")
	safeFileName := strings.ReplaceAll(unsafeFileName, "/", ".")
	f, err := os.Create(filepath.Join(k.impostersDirBinding, fmt.Sprintf("%s.imp.json", safeFileName)))
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString(string(data))
	if err != nil {
		return err
	}

	// when adding default imposters, the container is not yet started and the container will be nil
	// this allows us to add them without having to wait for the imposter to load later
	if k.Container != nil {
		// wait for the log saying the imposter was loaded
		containerFile := filepath.Join("/imposters", fmt.Sprintf("%s.imp.json", safeFileName))
		logWaitStrategy := wait.ForLog(fmt.Sprintf("imposter %s loaded", containerFile)).WithStartupTimeout(15 * time.Second)
		err = logWaitStrategy.WaitUntilReady(testcontext.Get(k.t), k.Container)
	}
	return err
}

// SetStringValuePath sets a path to return a string value
func (k *Killgrave) SetStringValuePath(path string, methods []string, headers map[string]string, v string) error {
	imposters := []KillgraveImposter{}
	for _, method := range methods {
		imposters = append(imposters, KillgraveImposter{
			Request: KillgraveRequest{
				Method:   method,
				Endpoint: path,
			},
			Response: KillgraveResponse{
				Status:  200,
				Body:    v,
				Headers: &headers,
			},
		})
	}

	return k.AddImposter(imposters)
}

// SetAdapterBasedAnyValuePath sets a path to return a value as though it was from an adapter
func (k *Killgrave) SetAdapterBasedAnyValuePath(path string, methods []string, v interface{}) error {
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

	return k.SetStringValuePath(path, methods, map[string]string{
		"Content-Type": "application/json",
	}, string(data))
}

// SetAdapterBasedAnyValuePathObject sets a path to return a value as though it was from an adapter
func (k *Killgrave) SetAdapterBasedIntValuePath(path string, methods []string, v int) error {
	return k.SetAdapterBasedAnyValuePath(path, methods, v)
}
