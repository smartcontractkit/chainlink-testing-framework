package havoc

import (
	"fmt"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/pkg/errors"
	"github.com/samber/lo"
	"regexp"
	"strings"
)

const (
	ErrParsingOpenAPISpec = "failed to parse OpenAPISpec"
)

var (
	OpenAPIPathParam = regexp.MustCompile(`({.*})`)
)

type OAPISpecData struct {
	Port     int64
	RawPaths []string
	SpecData map[string]*openapi3.PathItem
}

// ParseOpenAPISpecs parses OpenAPI spec methods
func (m *Controller) ParseOpenAPISpecs() ([]*OAPISpecData, error) {
	data := make([]*OAPISpecData, 0)
	for _, oapiData := range m.cfg.Havoc.OpenAPI.Mapping {
		for _, p := range oapiData.SpecToPortMappings {
			loader := openapi3.NewLoader()
			doc, err := loader.LoadFromFile(p.Path)
			if err != nil {
				return nil, errors.Wrap(err, ErrParsingOpenAPISpec)
			}
			oa := &OAPISpecData{
				Port:     p.Port,
				RawPaths: make([]string, 0),
				SpecData: doc.Paths.Map(),
			}
			for rawPath := range doc.Paths.Map() {
				L.Info().Str("Path", rawPath).Msg("Found API path")
				oa.RawPaths = append(oa.RawPaths, rawPath)
			}
			data = append(data, oa)
		}
	}
	return data, nil
}

// generateOAPIExperiments generates HTTP experiments for a component group (entry), for each method type
func (m *Controller) generateOAPIExperiments(experiments map[string]string, namespace string, entry lo.Entry[string, int], oapiSpecs []*OAPISpecData) error {
	for _, apiSpec := range oapiSpecs {
		for _, rawPath := range apiSpec.RawPaths {
			pathData := apiSpec.SpecData[rawPath]
			if pathData.Connect != nil {
				if err := m.generateHTTPExperiment(experiments, namespace, entry, rawPath, "CONNECT", apiSpec.Port); err != nil {
					return err
				}
			}
			if pathData.Delete != nil {
				if err := m.generateHTTPExperiment(experiments, namespace, entry, rawPath, "DELETE", apiSpec.Port); err != nil {
					return err
				}
			}
			if pathData.Get != nil {
				if err := m.generateHTTPExperiment(experiments, namespace, entry, rawPath, "GET", apiSpec.Port); err != nil {
					return err
				}
			}
			if pathData.Head != nil {
				if err := m.generateHTTPExperiment(experiments, namespace, entry, rawPath, "HEAD", apiSpec.Port); err != nil {
					return err
				}
			}
			if pathData.Options != nil {
				if err := m.generateHTTPExperiment(experiments, namespace, entry, rawPath, "OPTIONS", apiSpec.Port); err != nil {
					return err
				}
			}
			if pathData.Patch != nil {
				if err := m.generateHTTPExperiment(experiments, namespace, entry, rawPath, "PATCH", apiSpec.Port); err != nil {
					return err
				}
			}
			if pathData.Post != nil {
				if err := m.generateHTTPExperiment(experiments, namespace, entry, rawPath, "POST", apiSpec.Port); err != nil {
					return err
				}
			}
			if pathData.Put != nil {
				if err := m.generateHTTPExperiment(experiments, namespace, entry, rawPath, "PUT", apiSpec.Port); err != nil {
					return err
				}
			}
			if pathData.Trace != nil {
				if err := m.generateHTTPExperiment(experiments, namespace, entry, rawPath, "TRACE", apiSpec.Port); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (m *Controller) generateHTTPExperiment(
	experiments map[string]string,
	namespace string,
	entry lo.Entry[string, int],
	rawPath string,
	method string,
	port int64,
) error {
	sanitizedLabel := sanitizeLabel(entry.Key)
	sanitizedRawPath := sanitizeLabel(rawPath)
	sanitizedLabel = fmt.Sprintf("%s-%s-%s", sanitizedLabel, sanitizedRawPath, method)
	experiment, err := HTTPExperiment{
		Namespace:      namespace,
		ExperimentName: strings.ToLower(fmt.Sprintf("%s-%s", ChaosTypeHTTP, sanitizedLabel)),
		Duration:       m.cfg.Havoc.StressCPU.Duration,
		Mode:           "all",
		Selector:       entry.Key,
		Target:         "Response",
		Abort:          true,
		Path:           pathToWildcardExpr(rawPath),
		Method:         method,
		Port:           port,
	}.String()
	if err != nil {
		return err
	}
	experiments[sanitizedLabel] = experiment
	return nil
}

// pathToWildcardExpr transforms path params into wildcard expressions
// TODO: this need thorough testing though, since it can be much more complex
func pathToWildcardExpr(path string) string {
	return string(OpenAPIPathParam.ReplaceAll([]byte(path), []byte("*")))
}
