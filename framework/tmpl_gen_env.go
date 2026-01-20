package framework

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"text/template"
	"unicode"

	"github.com/rs/zerolog/log"
)

/* Templates */

const (
	// ProductDashboardUUID is a default product dashboard uuid, can be static since it's our local environment
	ProductDashboardUUID = "f8a04cef-653f-46d3-86df-87c532300672"

	ReadmeTmpl = `## Chainlink Developer Environment

This template provides a complete Chainlink development environment with pre-configured infrastructure and observability tools, enabling rapid development while maintaining high quality standards.

🔧 Address all **TODO** comments and implement "product_configuration.go"

💻 Enter the shell:
` + "```" + `bash
just cli && {{ .CLIName }} sh
` + "```" + `

🚀 Spin up the environment
` + "```" + `bash
up ↵
` + "```" + `

🔍 Implement system-level smoke tests (tests/smoke_test.go) and run them:
` + "```" + `bash
test smoke ↵
` + "```" + `

📈 Implement load/chaos tests (tests/load_test.go) and run them:
` + "```" + `bash
test load ↵
` + "```" + `

🔄 **Enforce** quality standards in CI: copy .github/workflows to your CI folder, commit and make them pass
`
	// ProductsInterfaceTmpl common interface for arbitrary products deployed in devenv
	ProductsInterfaceTmpl = `package {{ .PackageName }}

import (
	"context"

	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/fake"

	nodeset "github.com/smartcontractkit/chainlink-testing-framework/framework/components/simple_node_set"
)

// Product describes a minimal set of methods that each legacy product must implement
type Product interface {
	Load() error

	Store(path string, instanceIdx int) error

	GenerateNodesSecrets(
		ctx context.Context,
		fs *fake.Input,
		bc *blockchain.Input,
		ns *nodeset.Input,
	) (string, error)

	GenerateNodesConfig(
		ctx context.Context,
		fs *fake.Input,
		bc *blockchain.Input,
		ns *nodeset.Input,
	) (string, error)

	ConfigureJobsAndContracts(
		ctx context.Context,
		fs *fake.Input,
		bc *blockchain.Input,
		ns *nodeset.Input,
	) error
}
`

	// GoModTemplate go module template
	GoModTemplate = `module {{.ModuleName}}

go {{.RuntimeVersion}}

require (
	github.com/smartcontractkit/chainlink-evm v0.0.0-20250709215002-07f34ab867df
	github.com/smartcontractkit/chainlink-deployments-framework v0.17.0
	github.com/smartcontractkit/chainlink-testing-framework/framework v0.11.10
)

replace github.com/fbsobreira/gotron-sdk => github.com/smartcontractkit/chainlink-tron/relayer/gotron-sdk v0.0.5-0.20250528121202-292529af39df

replace (
	github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.3-alpha.regen.1
	github.com/prometheus/common => github.com/prometheus/common v0.62.0
	github.com/smartcontractkit/chainlink-testing-framework/lib => github.com/smartcontractkit/chainlink-testing-framework/lib v1.54.4
)
`
	// GitIgnoreTmpl default gitignore template
	GitIgnoreTmpl = `compose/
blockscout/
env-out.toml`

	// GrafanaDashboardTmpl is a Grafana dashboard template for your product
	GrafanaDashboardTmpl = `{
  "annotations": {
    "list": [
      {
        "builtIn": 1,
        "datasource": {
          "type": "grafana",
          "uid": "-- Grafana --"
        },
        "enable": true,
        "hide": true,
        "iconColor": "rgba(0, 211, 255, 1)",
        "name": "Annotations & Alerts",
        "type": "dashboard"
      }
    ]
  },
  "editable": true,
  "fiscalYearStartMonth": 0,
  "graphTooltip": 0,
  "links": [],
  "liveNow": false,
  "panels": [
    {
      "collapsed": false,
      "gridPos": {
        "h": 1,
        "w": 24,
        "x": 0,
        "y": 0
      },
      "id": 34,
      "panels": [],
      "title": "Container Resources",
      "type": "row"
    },
    {
      "datasource": {
        "type": "prometheus",
        "uid": "PBFA97CFB590B2093"
      },
      "fieldConfig": {
        "defaults": {
          "color": {
            "mode": "palette-classic"
          },
          "custom": {
            "axisCenteredZero": false,
            "axisColorMode": "text",
            "axisLabel": "",
            "axisPlacement": "auto",
            "barAlignment": 0,
            "drawStyle": "line",
            "fillOpacity": 10,
            "gradientMode": "none",
            "hideFrom": {
              "legend": false,
              "tooltip": false,
              "viz": false
            },
            "insertNulls": false,
            "lineInterpolation": "linear",
            "lineWidth": 1,
            "pointSize": 5,
            "scaleDistribution": {
              "type": "linear"
            },
            "showPoints": "never",
            "spanNulls": false,
            "stacking": {
              "group": "A",
              "mode": "normal"
            },
            "thresholdsStyle": {
              "mode": "off"
            }
          },
          "mappings": [],
          "thresholds": {
            "mode": "absolute",
            "steps": [
              {
                "color": "green",
                "value": null
              },
              {
                "color": "red",
                "value": 80
              }
            ]
          },
          "unit": "percent"
        },
        "overrides": []
      },
      "gridPos": {
        "h": 7,
        "w": 12,
        "x": 0,
        "y": 1
      },
      "id": 2,
      "options": {
        "legend": {
          "calcs": [
            "mean",
            "max"
          ],
          "displayMode": "table",
          "placement": "right",
          "showLegend": true
        },
        "tooltip": {
          "mode": "multi",
          "sort": "none"
        }
      },
      "pluginVersion": "10.1.0",
      "targets": [
        {
          "datasource": {
            "type": "prometheus",
            "uid": "PBFA97CFB590B2093"
          },
          "editorMode": "code",
          "expr": "sum(rate(container_cpu_usage_seconds_total{name=~\".*don.*|.*fake.*\"}[5m])) by (name) *100",
          "hide": false,
          "interval": "",
          "legendFormat": "{{"{{"}}name{{"}}"}}",
          "range": true,
          "refId": "A"
        }
      ],
      "title": "CPU Usage",
      "transparent": true,
      "type": "timeseries"
    },
    {
      "datasource": {
        "type": "prometheus",
        "uid": "PBFA97CFB590B2093"
      },
      "fieldConfig": {
        "defaults": {
          "color": {
            "mode": "palette-classic"
          },
          "custom": {
            "axisCenteredZero": false,
            "axisColorMode": "text",
            "axisLabel": "",
            "axisPlacement": "auto",
            "barAlignment": 0,
            "drawStyle": "line",
            "fillOpacity": 10,
            "gradientMode": "none",
            "hideFrom": {
              "legend": false,
              "tooltip": false,
              "viz": false
            },
            "insertNulls": false,
            "lineInterpolation": "linear",
            "lineWidth": 1,
            "pointSize": 5,
            "scaleDistribution": {
              "type": "linear"
            },
            "showPoints": "never",
            "spanNulls": false,
            "stacking": {
              "group": "A",
              "mode": "none"
            },
            "thresholdsStyle": {
              "mode": "off"
            }
          },
          "mappings": [],
          "thresholds": {
            "mode": "absolute",
            "steps": [
              {
                "color": "green",
                "value": null
              },
              {
                "color": "red",
                "value": 80
              }
            ]
          },
          "unit": "Bps"
        },
        "overrides": []
      },
      "gridPos": {
        "h": 7,
        "w": 12,
        "x": 12,
        "y": 1
      },
      "id": 4,
      "options": {
        "legend": {
          "calcs": [
            "mean",
            "max"
          ],
          "displayMode": "table",
          "placement": "right",
          "showLegend": true
        },
        "tooltip": {
          "mode": "multi",
          "sort": "none"
        }
      },
      "pluginVersion": "10.1.0",
      "targets": [
        {
          "datasource": {
            "type": "prometheus",
            "uid": "PBFA97CFB590B2093"
          },
          "editorMode": "code",
          "expr": "sum(rate(container_network_receive_bytes_total{name=~\".*don.*|.*fake.*\"}[5m])) by (name)",
          "hide": false,
          "interval": "",
          "legendFormat": "{{"{{"}}name{{"}}"}}",
          "range": true,
          "refId": "A"
        }
      ],
      "title": "Received Network Traffic",
      "transparent": true,
      "type": "timeseries"
    },
    {
      "datasource": {
        "type": "prometheus",
        "uid": "PBFA97CFB590B2093"
      },
      "fieldConfig": {
        "defaults": {
          "color": {
            "mode": "palette-classic"
          },
          "custom": {
            "axisCenteredZero": false,
            "axisColorMode": "text",
            "axisLabel": "",
            "axisPlacement": "auto",
            "barAlignment": 0,
            "drawStyle": "line",
            "fillOpacity": 10,
            "gradientMode": "none",
            "hideFrom": {
              "legend": false,
              "tooltip": false,
              "viz": false
            },
            "insertNulls": false,
            "lineInterpolation": "linear",
            "lineWidth": 1,
            "pointSize": 5,
            "scaleDistribution": {
              "type": "linear"
            },
            "showPoints": "never",
            "spanNulls": false,
            "stacking": {
              "group": "A",
              "mode": "normal"
            },
            "thresholdsStyle": {
              "mode": "off"
            }
          },
          "mappings": [],
          "thresholds": {
            "mode": "absolute",
            "steps": [
              {
                "color": "green",
                "value": null
              },
              {
                "color": "red",
                "value": 80
              }
            ]
          },
          "unit": "bytes"
        },
        "overrides": []
      },
      "gridPos": {
        "h": 7,
        "w": 12,
        "x": 0,
        "y": 8
      },
      "id": 3,
      "options": {
        "legend": {
          "calcs": [
            "mean",
            "max"
          ],
          "displayMode": "table",
          "placement": "right",
          "showLegend": true
        },
        "tooltip": {
          "mode": "multi",
          "sort": "none"
        }
      },
      "pluginVersion": "10.1.0",
      "targets": [
        {
          "datasource": {
            "type": "prometheus",
            "uid": "PBFA97CFB590B2093"
          },
          "editorMode": "code",
          "expr": "sum(container_memory_rss{name=~\".*don.*|.*fake.*\"}) by (name)",
          "hide": false,
          "interval": "",
          "legendFormat": "{{"{{"}}name{{"}}"}}",
          "range": true,
          "refId": "A"
        }
      ],
      "title": "Memory Usage",
      "transparent": true,
      "type": "timeseries"
    },
    {
      "datasource": {
        "type": "prometheus",
        "uid": "PBFA97CFB590B2093"
      },
      "fieldConfig": {
        "defaults": {
          "color": {
            "mode": "palette-classic"
          },
          "custom": {
            "axisCenteredZero": false,
            "axisColorMode": "text",
            "axisLabel": "",
            "axisPlacement": "auto",
            "barAlignment": 0,
            "drawStyle": "line",
            "fillOpacity": 10,
            "gradientMode": "none",
            "hideFrom": {
              "legend": false,
              "tooltip": false,
              "viz": false
            },
            "insertNulls": false,
            "lineInterpolation": "linear",
            "lineWidth": 1,
            "pointSize": 5,
            "scaleDistribution": {
              "type": "linear"
            },
            "showPoints": "never",
            "spanNulls": false,
            "stacking": {
              "group": "A",
              "mode": "none"
            },
            "thresholdsStyle": {
              "mode": "off"
            }
          },
          "mappings": [],
          "thresholds": {
            "mode": "absolute",
            "steps": [
              {
                "color": "green",
                "value": null
              },
              {
                "color": "red",
                "value": 80
              }
            ]
          },
          "unit": "Bps"
        },
        "overrides": []
      },
      "gridPos": {
        "h": 7,
        "w": 12,
        "x": 12,
        "y": 8
      },
      "id": 5,
      "options": {
        "legend": {
          "calcs": [
            "mean",
            "max"
          ],
          "displayMode": "table",
          "placement": "right",
          "showLegend": true
        },
        "tooltip": {
          "mode": "multi",
          "sort": "none"
        }
      },
      "pluginVersion": "10.1.0",
      "targets": [
        {
          "datasource": {
            "type": "prometheus",
            "uid": "PBFA97CFB590B2093"
          },
          "editorMode": "code",
          "expr": "sum(rate(container_network_transmit_bytes_total{name=~\".*don.*|.*fake.*\"}[5m])) by (name)",
          "interval": "",
          "legendFormat": "{{"{{"}}name{{"}}"}}",
          "range": true,
          "refId": "A"
        }
      ],
      "title": "Sent Network Traffic",
      "transparent": true,
      "type": "timeseries"
    },
    {
      "collapsed": true,
      "gridPos": {
        "h": 1,
        "w": 24,
        "x": 0,
        "y": 15
      },
      "id": 11,
      "panels": [
        {
          "datasource": {
            "type": "loki",
            "uid": "P8E80F9AEF21F6940"
          },
          "description": "CL Node Logs",
          "gridPos": {
            "h": 5,
            "w": 24,
            "x": 0,
            "y": 16
          },
          "id": 6,
          "options": {
            "dedupStrategy": "none",
            "enableLogDetails": true,
            "prettifyLogMessage": false,
            "showCommonLabels": false,
            "showLabels": false,
            "showTime": false,
            "sortOrder": "Descending",
            "wrapLogMessage": false
          },
          "targets": [
            {
              "datasource": {
                "type": "loki",
                "uid": "P8E80F9AEF21F6940"
              },
              "editorMode": "code",
              "expr": "{job=\"ctf\", container=~\".*don-node.*\"}",
              "queryType": "range",
              "refId": "A"
            }
          ],
          "title": "CL Node Logs",
          "transparent": true,
          "type": "logs"
        },
        {
          "datasource": {
            "type": "loki",
            "uid": "P8E80F9AEF21F6940"
          },
          "description": "P2P Discovery Report",
          "gridPos": {
            "h": 5,
            "w": 24,
            "x": 0,
            "y": 21
          },
          "id": 8,
          "options": {
            "dedupStrategy": "none",
            "enableLogDetails": true,
            "prettifyLogMessage": false,
            "showCommonLabels": false,
            "showLabels": false,
            "showTime": false,
            "sortOrder": "Descending",
            "wrapLogMessage": false
          },
          "targets": [
            {
              "datasource": {
                "type": "loki",
                "uid": "P8E80F9AEF21F6940"
              },
              "editorMode": "code",
              "expr": "{job=\"ctf\"} |= \"DiscoveryProtocol: Status report\"",
              "queryType": "range",
              "refId": "A"
            }
          ],
          "title": "P2P Discovery Report",
          "transparent": true,
          "type": "logs"
        }
      ],
      "title": "CL Node Stats",
      "type": "row"
    },
    {
      "collapsed": true,
      "gridPos": {
        "h": 1,
        "w": 24,
        "x": 0,
        "y": 16
      },
      "id": 35,
      "panels": [
        {
          "datasource": {
            "type": "loki",
            "uid": "P8E80F9AEF21F6940"
          },
          "gridPos": {
            "h": 4,
            "w": 24,
            "x": 0,
            "y": 17
          },
          "id": 40,
          "options": {
            "code": {
              "language": "plaintext",
              "showLineNumbers": false,
              "showMiniMap": false
            },
            "content": "# Pprof profiling with Pyroscope\nSelect \"service_name\" variable to start.\n\nFor more info you can also use native [UI](http://localhost:4040)",
            "mode": "markdown"
          },
          "pluginVersion": "10.1.0",
          "transparent": true,
          "type": "text"
        },
        {
          "datasource": {
            "type": "grafana-pyroscope-datasource",
            "uid": "P02E4190217B50628"
          },
          "description": "CPU",
          "gridPos": {
            "h": 9,
            "w": 24,
            "x": 0,
            "y": 21
          },
          "id": 36,
          "targets": [
            {
              "datasource": {
                "type": "grafana-pyroscope-datasource",
                "uid": "P02E4190217B50628"
              },
              "groupBy": [],
              "labelSelector": "{service_name=\"${service_name}\"}",
              "profileTypeId": "process_cpu:cpu:nanoseconds:cpu:nanoseconds",
              "queryType": "profile",
              "refId": "A"
            }
          ],
          "title": "CPU",
          "transparent": true,
          "type": "flamegraph"
        },
        {
          "datasource": {
            "type": "grafana-pyroscope-datasource",
            "uid": "P02E4190217B50628"
          },
          "description": "Memory (alloc_objects)",
          "gridPos": {
            "h": 9,
            "w": 24,
            "x": 0,
            "y": 30
          },
          "id": 37,
          "targets": [
            {
              "datasource": {
                "type": "grafana-pyroscope-datasource",
                "uid": "P02E4190217B50628"
              },
              "groupBy": [],
              "labelSelector": "{service_name=\"${service_name}\"}",
              "profileTypeId": "memory:alloc_objects:count:space:bytes",
              "queryType": "profile",
              "refId": "A"
            }
          ],
          "title": "Memory (alloc_objects)",
          "transparent": true,
          "type": "flamegraph"
        },
        {
          "datasource": {
            "type": "grafana-pyroscope-datasource",
            "uid": "P02E4190217B50628"
          },
          "description": "Goroutines",
          "gridPos": {
            "h": 9,
            "w": 24,
            "x": 0,
            "y": 39
          },
          "id": 39,
          "targets": [
            {
              "datasource": {
                "type": "grafana-pyroscope-datasource",
                "uid": "P02E4190217B50628"
              },
              "groupBy": [],
              "labelSelector": "{service_name=\"${service_name}\"}",
              "profileTypeId": "goroutines:goroutine:count:goroutine:count",
              "queryType": "profile",
              "refId": "A"
            }
          ],
          "title": "Goroutines",
          "transparent": true,
          "type": "flamegraph"
        },
        {
          "datasource": {
            "type": "grafana-pyroscope-datasource",
            "uid": "P02E4190217B50628"
          },
          "description": "Mutex",
          "gridPos": {
            "h": 9,
            "w": 24,
            "x": 0,
            "y": 48
          },
          "id": 38,
          "targets": [
            {
              "datasource": {
                "type": "grafana-pyroscope-datasource",
                "uid": "P02E4190217B50628"
              },
              "groupBy": [],
              "labelSelector": "{service_name=\"${service_name}\"}",
              "profileTypeId": "mutex:contentions:count:contentions:count",
              "queryType": "profile",
              "refId": "A"
            }
          ],
          "title": "Mutex",
          "transparent": true,
          "type": "flamegraph"
        }
      ],
      "title": "Pyroscope",
      "type": "row"
    }
  ],
  "refresh": "5s",
  "schemaVersion": 38,
  "style": "dark",
  "tags": [],
  "templating": {
    "list": [
      {
        "current": {
          "selected": false,
          "text": "chainlink-node",
          "value": "chainlink-node"
        },
        "description": "service_name",
        "hide": 0,
        "includeAll": false,
        "label": "service_name",
        "multi": false,
        "name": "service_name",
        "options": [
          {
            "selected": true,
            "text": "chainlink-node",
            "value": "chainlink-node"
          }
        ],
        "query": "chainlink-node",
        "queryValue": "",
        "skipUrlSync": false,
        "type": "custom"
      }
    ]
  },
  "time": {
    "from": "now-1h",
    "to": "now"
  },
  "timepicker": {},
  "timezone": "",
  "title": "{{ .ProductName }} Services",
  "uid": "{{ .UUID }}",
  "version": 1,
  "weekStart": ""
}`
	// ConfigTOMLTmpl is a default env.toml template for devenv describind components configuration
	ConfigTOMLTmpl = `
[[products]]
name = "{{ .ProductName }}"
instances = 1

[fake_server]
  image = "{{ .ProductName }}-fakes:latest"
  port = 9111

[[blockchains]]
  chain_id = "1337"
  docker_cmd_params = ["-b", "1", "--mixed-mining", "--slots-in-an-epoch", "1"]
  image = "ghcr.io/foundry-rs/foundry:stable"
  port = "8545"
  type = "anvil"

[[nodesets]]
  name = "don"
  nodes = {{ .Nodes }}
  override_mode = "each"

  [nodesets.db]
    image = "postgres:15.0"

	{{- range .NodeIndices }}
	[[nodesets.node_specs]]
	    [nodesets.node_specs.node]
	    image = "public.ecr.aws/chainlink/chainlink:2.26.0"
	{{- end }}
`

	// CILoadChaosTemplate is a continuous integration template for end-to-end load/chaos tests
	CILoadChaosTemplate = `name: End-to-end {{ .ProductName }} Load and Chaos Tests

on:
 pull_request:

defaults:
 run:
   working-directory: {{ .DevEnvRelPath }}

concurrency:
 group: {{"${{"}} github.workflow {{"}}"}}-{{"${{"}} github.ref {{"}}"}}
 cancel-in-progress: true

jobs:
 e2e-tests:
   permissions:
     id-token: write
     contents: read
     pull-requests: write
   runs-on: ubuntu-latest
   strategy:
     matrix:
       include:
         - name: LoadChaos
           config: env.toml
         # TODO: Add more test and environment configurations as needed
   steps:
     - name: Checkout code
       uses: actions/checkout@v5
       with:
         fetch-depth: 0

     - name: Set up Docker Buildx
       uses: docker/setup-buildx-action@e468171a9de216ec08956ac3ada2f0791b6bd435 # v3.11.1

     - name: Install Just
       uses: extractions/setup-just@e33e0265a09d6d736e2ee1e0eb685ef1de4669ff # v3
       with:
         just-version: '1.40.0'

     - name: Authenticate to AWS ECR
       uses: ./.github/actions/aws-ecr-auth
       with:
         role-to-assume: {{"${{"}} secrets.CCV_IAM_ROLE {{"}}"}}
         aws-region: us-east-1
         registry-type: public

     - name: Authenticate to AWS ECR (JD)
       uses: ./.github/actions/aws-ecr-auth
       with:
         role-to-assume: {{"${{"}} secrets.CCV_IAM_ROLE {{"}}"}}
         aws-region: us-west-2
         registry-type: private
         registries: {{"${{"}} secrets.JD_REGISTRY {{"}}"}}

     - name: Set up Go
       uses: actions/setup-go@v6 # v6
       with:
         cache: true
         go-version-file: {{ .DevEnvRelPath }}/go.mod
         cache-dependency-path: {{ .DevEnvRelPath }}/go.sum

     - name: Download Go dependencies
       run: |
         go mod download

     - name: Run CCV environment
       env:
         JD_IMAGE: {{"${{"}} secrets.JD_IMAGE {{"$}}"}}
       run: |
         cd cmd/{{ .CLIName }} && go install . && cd -
         {{ .CLIName }} u {{"${{"}} matrix.config {{"}}"}}
         {{ .CLIName }} obs u -f

     - name: Run tests
       working-directory: build/devenv/tests
       run: |
         set -o pipefail
         go test -v -count=1 -run 'Test{{"${{"}} matrix.name {{"}}"}}'

     - name: Upload Logs
       if: always()
       uses: actions/upload-artifact@v4
       with:
         name: container-logs-{{"${{"}} matrix.name {{"}}"}}
         path: {{ .DevEnvRelPath }}/tests/logs
         retention-days: 1
`

	// CISmokeTmpl is a continuous integration template for end-to-end smoke tests
	CISmokeTmpl = `name: End-to-end {{ .ProductName }} Tests

on:
  pull_request:

defaults:
  run:
    working-directory: {{ .DevEnvRelPath }}

concurrency:
  group: {{"${{"}} github.workflow {{"}}"}}-{{"${{"}} github.ref {{"}}"}}
  cancel-in-progress: true

jobs:
  e2e-tests:
    permissions:
      id-token: write
      contents: read
      pull-requests: write
    runs-on: ubuntu-latest
    strategy:
      matrix:
        include:
          - name: Smoke
            config: env.toml
          # TODO: Add more test and environment configurations as needed
    steps:
      - name: Checkout code
        uses: actions/checkout@v5
        with:
          fetch-depth: 0

      - name: Set up Docker Buildx
        uses: docker/setup-buildx-action@e468171a9de216ec08956ac3ada2f0791b6bd435 # v3.11.1

      - name: Install Just
        uses: extractions/setup-just@e33e0265a09d6d736e2ee1e0eb685ef1de4669ff # v3
        with:
          just-version: '1.40.0'

      - name: Authenticate to AWS ECR
        uses: ./.github/actions/aws-ecr-auth
        with:
          role-to-assume: {{"${{"}} secrets.CCV_IAM_ROLE {{"}}"}}
          aws-region: us-east-1
          registry-type: public

      - name: Authenticate to AWS ECR (JD)
        uses: ./.github/actions/aws-ecr-auth
        with:
          role-to-assume: {{"${{"}} secrets.CCV_IAM_ROLE {{"}}"}}
          aws-region: us-west-2
          registry-type: private
          registries: {{"${{"}} secrets.JD_REGISTRY {{"}}"}}

      - name: Set up Go
        uses: actions/setup-go@v6 # v6
        with:
          cache: true
          go-version-file: {{ .DevEnvRelPath }}/go.mod
          cache-dependency-path: {{ .DevEnvRelPath }}/go.sum

      - name: Download Go dependencies
        run: |
          go mod download

      - name: Run CCV environment
        env:
          JD_IMAGE: {{"${{"}} secrets.JD_IMAGE {{"$}}"}}
        run: |
          cd cmd/{{ .CLIName }} && go install . && cd -
          {{ .CLIName }} u {{"${{"}} matrix.config {{"}}"}}

      - name: Run tests
        working-directory: build/devenv/tests
        run: |
          set -o pipefail
          go test -v -count=1 -run 'Test{{"${{"}} matrix.name {{"}}"}}'

      - name: Upload Logs
        if: always()
        uses: actions/upload-artifact@v4
        with:
          name: container-logs-{{"${{"}} matrix.name {{"}}"}}
          path: {{ .DevEnvRelPath }}/tests/logs
          retention-days: 1
`

	// CompletionTmpl is a go-prompt library completion template providing interactive prompt
	CompletionTmpl = `package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/c-bata/go-prompt"
)

func getCommands() []prompt.Suggest {
	return []prompt.Suggest{
		{Text: "", Description: "Choose command, press <space> for more options after selecting command"},
		{Text: "up", Description: "Spin up the development environment"},
		{Text: "down", Description: "Tear down the development environment"},
		{Text: "restart", Description: "Restart the development environment"},
		{Text: "test", Description: "Perform smoke or load/chaos testing"},
		{Text: "bs", Description: "Manage the Blockscout EVM block explorer"},
		{Text: "obs", Description: "Manage the observability stack"},
		{Text: "exit", Description: "Exit the interactive shell"},
	}
}

func getSubCommands(parent string) []prompt.Suggest {
	switch parent {
	case "test":
		return []prompt.Suggest{
			{Text: "smoke", Description: "Run {{ .CLIName }} smoke test"},
			{Text: "load", Description: "Run {{ .CLIName }} load test"},
		}
	case "bs":
		return []prompt.Suggest{
			{Text: "up", Description: "Spin up Blockscout and listen to dst chain (8555)"},
			{Text: "up -u http://host.docker.internal:8545 -c 1337", Description: "Spin up Blockscout and listen to src chain (8545)"},
			{Text: "down", Description: "Remove Blockscout stack"},
			{Text: "restart", Description: "Restart Blockscout and listen to dst chain (8555)"},
			{Text: "restart -u http://host.docker.internal:8545 -c 1337", Description: "Restart Blockscout and listen to src chain (8545)"},
		}
	case "obs":
		return []prompt.Suggest{
			{Text: "up", Description: "Spin up observability stack (Loki/Prometheus/Grafana)"},
			{Text: "up -f", Description: "Spin up full observability stack (Pyroscope, cadvisor, postgres exporter)"},
			{Text: "down", Description: "Spin down observability stack"},
			{Text: "restart", Description: "Restart observability stack"},
			{Text: "restart -f", Description: "Restart full observability stack"},
		}
	case "u":
		fallthrough
	case "up":
		fallthrough
	case "r":
		fallthrough
	case "restart":
		return []prompt.Suggest{
			{Text: "env.toml,products/{{ .ProductName }}/basic.toml", Description: "Spin up Anvil <> Anvil local chains, all services, 4 CL nodes"},
			{Text: "env.toml,products/{{ .ProductName }}/basic.toml,products/{{ .ProductName }}/soak.toml", Description: "Spin up Anvil <> Anvil local chains, all services, 4 CL nodes"},
			{Text: "env.toml,products/{{ .ProductName }}/basic.toml,env-cl-rebuild.toml", Description: "Spin up Anvil <> Anvil local chains, all services, 4 CL nodes (custom build)"},
			{Text: "env.toml,products/{{ .ProductName }}/basic.toml,env-geth.toml", Description: "Spin up Geth <> Geth local chains (clique), all services, 4 CL nodes"},
			{Text: "env.toml,products/{{ .ProductName }}/basic.toml,env-fuji-fantom.toml", Description: "Spin up testnets: Fuji <> Fantom, all services, 4 CL nodes"},
		}
	default:
		return []prompt.Suggest{}
	}
}

func executor(in string) {
	checkDockerIsRunning()
	in = strings.TrimSpace(in)
	if in == "" {
		return
	}
	if in == "exit" {
		fmt.Println("Goodbye!")
		os.Exit(0)
	}

	args := strings.Fields(in)
	os.Args = append([]string{"{{ .CLIName }}"}, args...)
	if err := rootCmd.Execute(); err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}

// completer provides autocomplete suggestions for multi-word commands.
func completer(in prompt.Document) []prompt.Suggest {
	text := in.TextBeforeCursor()
	words := strings.Fields(text)
	lastCharIsSpace := len(text) > 0 && text[len(text)-1] == ' '

	switch {
	case len(words) == 0:
		return getCommands()
	case len(words) == 1:
		if lastCharIsSpace {
			return getSubCommands(words[0])
		} else {
			return prompt.FilterHasPrefix(getCommands(), words[0], true)
		}

	case len(words) >= 2:
		if lastCharIsSpace {
			return []prompt.Suggest{}
		} else {
			parent := words[0]
			currentWord := words[len(words)-1]
			return prompt.FilterHasPrefix(getSubCommands(parent), currentWord, true)
		}
	default:
		return []prompt.Suggest{}
	}
}

// resetTerm resets terminal settings to Unix defaults.
func resetTerm() {
	cmd := exec.Command("stty", "sane")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}

func StartShell() {
	defer resetTerm()
	p := prompt.New(
		executor,
		completer,
		prompt.OptionPrefix("{{ .CLIName }}> "),
		prompt.OptionPrefixTextColor(prompt.Yellow),
		prompt.OptionTitle("CCIP Interactive Shell"),
		prompt.OptionMaxSuggestion(15),
		prompt.OptionShowCompletionAtStart(),
		prompt.OptionCompletionWordSeparator(" "),
		prompt.OptionSelectedSuggestionBGColor(prompt.LightGray),
		prompt.OptionSelectedSuggestionTextColor(prompt.Black),
		prompt.OptionDescriptionBGColor(prompt.DarkGray),
		prompt.OptionDescriptionTextColor(prompt.White),
		prompt.OptionSuggestionBGColor(prompt.Black),
		prompt.OptionSuggestionTextColor(prompt.Green),
		prompt.OptionScrollbarThumbColor(prompt.DarkGray),
		prompt.OptionScrollbarBGColor(prompt.Black),
		prompt.OptionAddKeyBind(prompt.KeyBind{
			Key: prompt.ControlC,
			Fn: func(buf *prompt.Buffer) {
				fmt.Println("Interrupted, exiting...")
				resetTerm()
				os.Exit(0)
			},
		}),
	)
	p.Run()
}
`
	// CLITmpl is a Cobra library CLI template with basic devenv commands
	CLITmpl = `package main

import (
		"context"
		"fmt"
		"os"
		"os/exec"
		"syscall"

		"github.com/docker/docker/client"
		"github.com/spf13/cobra"

		"github.com/smartcontractkit/chainlink-testing-framework/framework"
		"{{ .DevEnvPkgImport }}"
)

const (
		LocalWASPLoadDashboard = "http://localhost:3000/d/WASPLoadTests/wasp-load-test?orgId=1&from=now-5m&to=now&refresh=5s"
		Local{{ .ProductName }}Dashboard      = "http://localhost:3000/d/{{ .DashboardUUID }}/{{ .ProductName}}?orgId=1&refresh=5s"
)

var rootCmd = &cobra.Command{
	Use:   "{{ .CLIName }}",
	Short: "A {{ .ProductName }} local environment tool",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		debug, err := cmd.Flags().GetBool("debug")
		if err != nil {
			return err
		}
		if debug {
			framework.L.Info().Msg("Debug mode enabled, setting CTF_CLNODE_DLV=true")
			os.Setenv("CTF_CLNODE_DLV", "true")
		}
		return nil
	},
}

var restartCmd = &cobra.Command{
	Use:     "restart",
	Aliases: []string{"r"},
	Args:    cobra.RangeArgs(0, 1),
	Short:   "Restart development environment, remove apps and apply default configuration again",
	RunE: func(cmd *cobra.Command, args []string) error {
		var configFile string
		if len(args) > 0 {
			configFile = args[0]
		} else {
			configFile = "env.toml,products/{{ .ProductName }}/basic.toml"
		}
		framework.L.Info().Str("Config", configFile).Msg("Reconfiguring development environment")
		_ = os.Setenv("CTF_CONFIGS", configFile)
		_ = os.Setenv("TESTCONTAINERS_RYUK_DISABLED", "true")
		framework.L.Info().Msg("Tearing down the development environment")
		err := framework.RemoveTestContainers()
		if err != nil {
			return fmt.Errorf("failed to clean Docker resources: %w", err)
		}
		return devenv.NewEnvironment(context.Background())
	},
}

var upCmd = &cobra.Command{
	Use:     "up",
	Aliases: []string{"u"},
	Short:   "Spin up the development environment",
	Args:    cobra.RangeArgs(0, 1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var configFile string
		if len(args) > 0 {
			configFile = args[0]
		} else {
			configFile = "env.toml,products/{{ .ProductName }}/basic.toml"
		}
		framework.L.Info().Str("Config", configFile).Msg("Creating development environment")
		_ = os.Setenv("CTF_CONFIGS", configFile)
		_ = os.Setenv("TESTCONTAINERS_RYUK_DISABLED", "true")
		return devenv.NewEnvironment(context.Background())
	},
}

var downCmd = &cobra.Command{
	Use:     "down",
	Aliases: []string{"d"},
	Short:   "Tear down the development environment",
	RunE: func(cmd *cobra.Command, args []string) error {
		framework.L.Info().Msg("Tearing down the development environment")
		err := framework.RemoveTestContainers()
		if err != nil {
			return fmt.Errorf("failed to clean Docker resources: %w", err)
		}
		return nil
	},
}

var bsCmd = &cobra.Command{
	Use:   "bs",
	Short: "Manage the Blockscout EVM block explorer",
	Long:  "Spin up or down the Blockscout EVM block explorer",
}

var bsUpCmd = &cobra.Command{
	Use:     "up",
	Aliases: []string{"u"},
	Short:   "Spin up Blockscout EVM block explorer",
	RunE: func(cmd *cobra.Command, args []string) error {
		url, _ := bsCmd.Flags().GetString("url")
		chainID, _ := bsCmd.Flags().GetString("chain-id")
		return framework.BlockScoutUp(url, chainID)
	},
}

var bsDownCmd = &cobra.Command{
	Use:     "down",
	Aliases: []string{"d"},
	Short:   "Spin down Blockscout EVM block explorer",
	RunE: func(cmd *cobra.Command, args []string) error {
		url, _ := bsCmd.Flags().GetString("url")
		return framework.BlockScoutDown(url)
	},
}

var bsRestartCmd = &cobra.Command{
	Use:     "restart",
	Aliases: []string{"r"},
	Short:   "Restart the Blockscout EVM block explorer",
	RunE: func(cmd *cobra.Command, args []string) error {
		url, _ := bsCmd.Flags().GetString("url")
		chainID, _ := bsCmd.Flags().GetString("chain-id")
		if err := framework.BlockScoutDown(url); err != nil {
			return err
		}
		return framework.BlockScoutUp(url, chainID)
	},
}

var obsCmd = &cobra.Command{
	Use:   "obs",
	Short: "Manage the observability stack",
	Long:  "Spin up or down the observability stack with subcommands 'up' and 'down'",
}

var obsUpCmd = &cobra.Command{
	Use:     "up",
	Aliases: []string{"u"},
	Short:   "Spin up the observability stack",
	RunE: func(cmd *cobra.Command, args []string) error {
		full, _ := cmd.Flags().GetBool("full")
		var err error
		if full {
			err = framework.ObservabilityUpFull()
		} else {
			err = framework.ObservabilityUp()
		}
		if err != nil {
			return fmt.Errorf("observability up failed: %w", err)
		}
		devenv.L.Info().Msgf("{{ .ProductName }} Dashboard: %s", Local{{ .ProductName }}Dashboard)
		devenv.L.Info().Msgf("{{ .ProductName }} Load Test Dashboard: %s", LocalWASPLoadDashboard)
		return nil
	},
}

var obsDownCmd = &cobra.Command{
	Use:     "down",
	Aliases: []string{"d"},
	Short:   "Spin down the observability stack",
	RunE: func(cmd *cobra.Command, args []string) error {
		return framework.ObservabilityDown()
	},
}

var obsRestartCmd = &cobra.Command{
	Use:     "restart",
	Aliases: []string{"r"},
	Short:   "Restart the observability stack (data wipe)",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := framework.ObservabilityDown(); err != nil {
			return fmt.Errorf("observability down failed: %w", err)
		}
		full, _ := cmd.Flags().GetBool("full")
		var err error
		if full {
			err = framework.ObservabilityUpFull()
		} else {
			err = framework.ObservabilityUp()
		}
		if err != nil {
			return fmt.Errorf("observability up failed: %w", err)
		}
		devenv.L.Info().Msgf("{{ .ProductName }} Dashboard: %s", Local{{ .ProductName }}Dashboard)
		devenv.L.Info().Msgf("{{ .ProductName }} Load Test Dashboard: %s", LocalWASPLoadDashboard)
		return nil
	},
}

var testCmd = &cobra.Command{
	Use:     "test",
	Aliases: []string{"t"},
	Short:   "Run the tests",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return fmt.Errorf("specify the test suite: smoke or load")
		}
		var testPattern string
		switch args[0] {
		case "smoke":
			testPattern = "TestSmoke"
		case "load":
			testPattern = "TestLoadChaos"
		default:
			return fmt.Errorf("test suite %s is unknown, choose between smoke or load", args[0])
		}

		testCmd := exec.Command("go", "test", "-v", "-run", testPattern)
		testCmd.Dir = "./tests"
		testCmd.Stdout = os.Stdout
		testCmd.Stderr = os.Stderr
		testCmd.Stdin = os.Stdin

		if err := testCmd.Run(); err != nil {
			if exitError, ok := err.(*exec.ExitError); ok {
				if status, ok := exitError.Sys().(syscall.WaitStatus); ok {
					os.Exit(status.ExitStatus())
				}
				os.Exit(1)
			}
			return fmt.Errorf("failed to run test command: %w", err)
		}
		return nil
	},
}

func init() {
	rootCmd.PersistentFlags().BoolP("debug", "d", false, "Enable running services with dlv to allow remote debugging.")

	rootCmd.AddCommand(testCmd)

	// Blockscout, on-chain debug
	bsCmd.PersistentFlags().StringP("url", "u", "http://host.docker.internal:8555", "EVM RPC node URL (default to dst chain on 8555")
	bsCmd.PersistentFlags().StringP("chain-id", "c", "2337", "RPC's Chain ID")
	bsCmd.AddCommand(bsUpCmd)
	bsCmd.AddCommand(bsDownCmd)
	bsCmd.AddCommand(bsRestartCmd)
	rootCmd.AddCommand(bsCmd)

	// observability
	obsCmd.PersistentFlags().BoolP("full", "f", false, "Enable full observability stack with additional components")
	obsCmd.AddCommand(obsRestartCmd)
	obsCmd.AddCommand(obsUpCmd)
	obsCmd.AddCommand(obsDownCmd)
	rootCmd.AddCommand(obsCmd)

	// main env commands
	rootCmd.AddCommand(upCmd)
	rootCmd.AddCommand(restartCmd)
	rootCmd.AddCommand(downCmd)
}

func checkDockerIsRunning() {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		fmt.Println("Can't create Docker client, please check if Docker daemon is running!")
		os.Exit(1)
	}
	defer cli.Close()
	_, err = cli.Ping(context.Background())
	if err != nil {
		fmt.Println("Docker is not running, please start Docker daemon first!")
		os.Exit(1)
	}
}

func main() {
	checkDockerIsRunning()
	if len(os.Args) == 2 && (os.Args[1] == "shell" || os.Args[1] == "sh") {
		_ = os.Setenv("CTF_CONFIGS", "env.toml") // Set default config for shell

		StartShell()
		return
	}
	if err := rootCmd.Execute(); err != nil {
		devenv.L.Err(err).Send()
		os.Exit(1)
	}
}`
	// LoadTestTmpl is a load/chaos test template
	LoadTestTmpl = `package devenv_test

import (
	"testing"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/require"

	de "{{ .GoModName }}"

	"github.com/smartcontractkit/{{ .ProductName }}/devenv/products"
	"github.com/smartcontractkit/{{ .ProductName }}/devenv/products/{{ .ProductName }}"

	"github.com/smartcontractkit/chainlink-testing-framework/framework/chaos"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/clclient"
	"github.com/smartcontractkit/chainlink-testing-framework/wasp"
)

type ExampleGun struct {
	target string
	client *resty.Client
	Data   []string
}

func NewExampleHTTPGun(target string) *ExampleGun {
	return &ExampleGun{
		client: resty.New(),
		target: target,
		Data:   make([]string, 0),
	}
}

// Call implements example gun call, assertions on response bodies should be done here
func (m *ExampleGun) Call(l *wasp.Generator) *wasp.Response {
	var result map[string]any
	r, err := m.client.R().
		SetResult(&result).
		Get(m.target)
	if err != nil {
		return &wasp.Response{Data: result, Error: err.Error()}
	}
	if r.Status() != "200 OK" {
		return &wasp.Response{Data: result, Error: "not 200", Failed: true}
	}
	return &wasp.Response{Data: result}
}

func TestLoadChaos(t *testing.T) {
	in, err := de.LoadOutput[de.Cfg]("../env-out.toml")
	require.NoError(t, err)
	inProduct, err := products.LoadOutput[productone.Configurator]("../env-out.toml")
	require.NoError(t, err)

	_ = inProduct

	clNodes, err := clclient.New(in.NodeSets[0].Out.CLNodes)
	require.NoError(t, err)

	// use local observability stack for Docker load/chaos tests
	t.Setenv("LOKI_URL", "http://localhost:3030/loki/api/v1/push")

	// define labels for differentiate one run from another
	labels := map[string]string{
		"go_test_name": "generator_healthcheck",
		"gen_name":     "generator_healthcheck",
		"branch":       "test",
		"commit":       "test",
	}

	// create a WASP generator
	gen, err := wasp.NewGenerator(&wasp.Config{
		LoadType: wasp.RPS,
		T:        t,
		// just use plain line profile - 1 RPS for 60s
		Schedule:   wasp.Plain(1, 60*time.Second),
		Gun:        NewExampleHTTPGun("https://example.com"),
		Labels:     labels,
		LokiConfig: wasp.NewEnvLokiConfig(),
	})
	require.NoError(t, err)
	gen.Run(false)

	// define chaos test cases
	testCases := []struct {
		name     string
		command  string
		wait     time.Duration
		validate func(c []*clclient.ChainlinkClient) error
	}{
		{
			name:    "Reboot the pods",
			wait:    20 * time.Second,
			command: "stop --duration=20s --restart re2:don-node0",
			validate: func(c []*clclient.ChainlinkClient) error {
				return nil
			},
		},
		{
			name:    "Introduce network delay",
			wait:    20 * time.Second,
			command: "netem --tc-image=gaiadocker/iproute2 --duration=20s delay --time=1000 re2:don-node.*",
			validate: func(c []*clclient.ChainlinkClient) error {
				return nil
			},
		},
	}

	// Run chaos test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Log(tc.name)
			_, err = chaos.ExecPumba(tc.command, tc.wait)
			require.NoError(t, err)
			err = tc.validate(clNodes)
			require.NoError(t, err)
		})
	}
	// wait for the workload to finish
	_, failed := gen.Wait()
	require.False(t, failed)
}
`

	SmokeTestImplTmpl = `package devenv_test

import (
	"testing"

	de "{{ .GoModName }}"

	"github.com/smartcontractkit/{{ .ProductName }}/devenv/products"
	"github.com/smartcontractkit/{{ .ProductName }}/devenv/products/{{ .ProductName }}"

	"github.com/smartcontractkit/chainlink-testing-framework/framework/clclient"
	"github.com/stretchr/testify/require"
)

var L = de.L

func TestSmoke(t *testing.T) {
	in, err := de.LoadOutput[de.Cfg]("../env-out.toml")
	require.NoError(t, err)
	inProduct, err := products.LoadOutput[productone.Configurator]("../env-out.toml")
	require.NoError(t, err)
	clNodes, err := clclient.New(in.NodeSets[0].Out.CLNodes)
	require.NoError(t, err)

	_ = in
	_ = inProduct

	tests := []struct {
		name    string
		clNodes []*clclient.ChainlinkClient
	}{
		{
			name: "feature_1",
		},
		{
			name: "feature_2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_ = clNodes
		})
	}
}
`
	// JDTmpl is a JobDistributor client wrappers
	JDTmpl = `package {{ .PackageName }}

import (
	"context"
	"errors"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	csav1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/csa"
	jobv1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/job"
	nodev1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/node"
)

type JobDistributor struct {
	nodev1.NodeServiceClient
	jobv1.JobServiceClient
	csav1.CSAServiceClient
	WSRPC string
}

type JDConfig struct {
	GRPC  string
	WSRPC string
}

func (jd JobDistributor) GetCSAPublicKey(ctx context.Context) (string, error) {
	keypairs, err := jd.ListKeypairs(ctx, &csav1.ListKeypairsRequest{})
	if err != nil {
		return "", err
	}
	if keypairs == nil || len(keypairs.Keypairs) == 0 {
		return "", errors.New("no keypairs found")
	}
	csakey := keypairs.Keypairs[0].PublicKey
	return csakey, nil
}

// ProposeJob proposes jobs through the jobService and accepts the proposed job on selected node based on ProposeJobRequest.NodeId.
func (jd JobDistributor) ProposeJob(ctx context.Context, in *jobv1.ProposeJobRequest, opts ...grpc.CallOption) (*jobv1.ProposeJobResponse, error) {
	res, err := jd.JobServiceClient.ProposeJob(ctx, in, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to propose job. err: %w", err)
	}
	if res.Proposal == nil {
		return nil, errors.New("failed to propose job. err: proposal is nil")
	}

	return res, nil
}

// NewJDConnection creates new gRPC connection with JobDistributor.
func NewJDConnection(cfg JDConfig) (*grpc.ClientConn, error) {
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
	interceptors := []grpc.UnaryClientInterceptor{}

	if len(interceptors) > 0 {
		opts = append(opts, grpc.WithChainUnaryInterceptor(interceptors...))
	}

	conn, err := grpc.NewClient(cfg.GRPC, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to connect Job Distributor service. Err: %w", err)
	}

	return conn, nil
}
`
	// DebugToolsTmpl is a template for various debug tools, tracing, tx debug, etc
	DebugToolsTmpl = `package {{ .PackageName }}

import (
	"os"
	"runtime/trace"
)


func tracing() func() {
	f, err := os.Create("trace.out")
	if err != nil {
		panic("can't create trace.out file")
	}
	if err := trace.Start(f); err != nil {
		panic("can't start tracing")
	}
	return func() {
		trace.Stop()
		f.Close()
	}
}
`
	// ConfigTmpl is a template for reading and writing devenv configuration (env.toml, env-out.toml)
	ConfigTmpl = `package {{ .PackageName }}

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/davecgh/go-spew/spew"
	"github.com/pelletier/go-toml/v2"
	"github.com/rs/zerolog"
)

const (
	// DefaultConfigDir is the default directory we are expecting TOML config to be.
	DefaultConfigDir = "."
	// EnvVarTestConfigs is the environment variable name to read config paths from, ex.: CTF_CONFIGS=env.toml,overrides.toml.
	EnvVarTestConfigs = "CTF_CONFIGS"
	// DefaultOverridesFilePath is the default overrides.toml file path.
	DefaultOverridesFilePath = "overrides.toml"
	// DefaultAnvilKey is a default, well-known Anvil first key
	DefaultAnvilKey = "ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"
)

// Load loads TOML configurations from environment variable, ex.: CTF_CONFIGS=env.toml,overrides.toml
// and unmarshalls the files from left to right overriding keys.
func Load[T any]() (*T, error) {
	var config T
	paths := strings.Split(os.Getenv(EnvVarTestConfigs), ",")
	for _, path := range paths {
		L.Info().Str("Path", path).Msg("Loading configuration input")
		data, err := os.ReadFile(filepath.Join(DefaultConfigDir, path)) //nolint:gosec
		if err != nil {
			if path == DefaultOverridesFilePath {
				L.Info().Str("Path", path).Msg("Overrides file not found or empty")
				continue
			}
			return nil, fmt.Errorf("error reading config file %s: %w", path, err)
		}
		if L.GetLevel() == zerolog.TraceLevel {
			fmt.Println(string(data)) //nolint:forbidigo
		}

		decoder := toml.NewDecoder(strings.NewReader(string(data)))

		if err := decoder.Decode(&config); err != nil {
			var details *toml.StrictMissingError
			if errors.As(err, &details) {
				fmt.Println(details.String()) //nolint:forbidigo
			}
			return nil, fmt.Errorf("failed to decode TOML config, strict mode: %s", err)
		}
	}
	if L.GetLevel() == zerolog.TraceLevel {
		L.Trace().Msg("Merged inputs")
		spew.Dump(config) //nolint:forbidigo
	}
	return &config, nil
}

// Store writes config to a file, adds -cache.toml suffix if it's an initial configuration.
func Store[T any](cfg *T) error {
	baseConfigPath, err := BaseConfigPath()
	if err != nil {
		return err
	}
	newCacheName := strings.ReplaceAll(baseConfigPath, ".toml", "")
	var outCacheName string
	if strings.Contains(newCacheName, "cache") {
		L.Info().Str("Cache", baseConfigPath).Msg("Cache file already exists, overriding")
		outCacheName = baseConfigPath
	} else {
		outCacheName = fmt.Sprintf("%s-out.toml", strings.ReplaceAll(baseConfigPath, ".toml", ""))
	}
	L.Info().Str("OutputFile", outCacheName).Msg("Storing configuration output")
	d, err := toml.Marshal(cfg)
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(DefaultConfigDir, outCacheName), d, 0o644)
}

// LoadOutput loads config output file from path.
func LoadOutput[T any](path string) (*T, error) {
	_ = os.Setenv(EnvVarTestConfigs, path)
	return Load[T]()
}

// BaseConfigPath returns base config path, ex. env.toml,overrides.toml -> env.toml.
func BaseConfigPath() (string, error) {
	configs := os.Getenv(EnvVarTestConfigs)
	if configs == "" {
		return "", fmt.Errorf("no %s env var is provided, you should provide at least one test config in TOML", EnvVarTestConfigs)
	}
	L.Debug().Str("Configs", configs).Msg("Getting base config path")
	return strings.Split(configs, ",")[0], nil
}

// GetNetworkPrivateKey gets network private key or fallback to default simulator key (Anvil's first key)
func GetNetworkPrivateKey() string {
	pk := os.Getenv("PRIVATE_KEY")
	if pk == "" {
		// that's the first Anvil and Geth private key, serves as a fallback for local testing if not overridden
		return DefaultAnvilKey
	}
	return pk
}
`

	// EnvironmentTmpl is an environment.go template - main file for environment composition
	EnvironmentTmpl = `package {{ .PackageName }}

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/fake"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/jd"
	"github.com/smartcontractkit/{{ .ProductName }}/devenv/products/{{ .ProductName }}"

	ns "github.com/smartcontractkit/chainlink-testing-framework/framework/components/simple_node_set"
)

var L = log.Output(zerolog.ConsoleWriter{Out: os.Stderr}).Level(zerolog.DebugLevel).With().Fields(map[string]any{"component": "{{ .ProductName }}"}).Logger()

type ProductInfo struct {
	Name      string ` + "`" + `toml:"name"` + "`" + `
	Instances int    ` + "`" + `toml:"instances"` + "`" + `
}

type Cfg struct {
	Products    []*ProductInfo      ` + "`" + `toml:"products"` + "`" + `
	Blockchains []*blockchain.Input ` + "`" + `toml:"blockchains" validate:"required"` + "`" + `
	FakeServer  *fake.Input         ` + "`" + `toml:"fake_server" validate:"required"` + "`" + `
	NodeSets    []*ns.Input         ` + "`" + `toml:"nodesets"    validate:"required"` + "`" + `
	JD          *jd.Input           ` + "`" + `toml:"jd"` + "`" + `
}

func newProduct(name string) (Product, error) {
	switch name {
	case "{{ .ProductName }}":
		return {{ .ProductName }}.NewConfigurator(), nil
	default:
		return nil, fmt.Errorf("unknown product type: %s", name)
	}
}

func NewEnvironment(ctx context.Context) error {
	if err := framework.DefaultNetwork(nil); err != nil {
		return err
	}
	in, err := Load[Cfg]()
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}
	_, err = blockchain.NewBlockchainNetwork(in.Blockchains[0])
	if err != nil {
		return fmt.Errorf("failed to create blockchain network 1337: %w", err)
	}
	if os.Getenv("FAKE_SERVER_IMAGE") != "" {
		in.FakeServer.Image = os.Getenv("FAKE_SERVER_IMAGE")
	}
	_, err = fake.NewDockerFakeDataProvider(in.FakeServer)
	if err != nil {
		return fmt.Errorf("failed to create fake data provider: %w", err)
	}

	// get all the product orchestrations, generate product specific overrides
	productConfigurators := make([]Product, 0)
	nodeConfigs := make([]string, 0)
	nodeSecrets := make([]string, 0)
	for _, product := range in.Products {
		p, err := newProduct(product.Name)
		if err != nil {
			return err
		}
		if err = p.Load(); err != nil {
			return fmt.Errorf("failed to load product config: %w", err)
		}

		cfg, err := p.GenerateNodesConfig(ctx, in.FakeServer, in.Blockchains[0], in.NodeSets[0])
		if err != nil {
			return fmt.Errorf("failed to generate CL nodes config: %w", err)
		}
		nodeConfigs = append(nodeConfigs, cfg)

		secrets, err := p.GenerateNodesSecrets(ctx, in.FakeServer, in.Blockchains[0], in.NodeSets[0])
		if err != nil {
			return fmt.Errorf("failed to generate CL nodes config: %w", err)
		}
		nodeSecrets = append(nodeSecrets, secrets)

		productConfigurators = append(productConfigurators, p)
	}

	// merge overrides, spin up node sets and write infrastructure outputs
	// infra is always common for all the products, if it can't be we should fail
	// user should use different infra layout in env.toml then
	for _, ns := range in.NodeSets[0].NodeSpecs {
		ns.Node.TestConfigOverrides = strings.Join(nodeConfigs, "\n")
		ns.Node.TestSecretsOverrides = strings.Join(nodeSecrets, "\n")
		if os.Getenv("CHAINLINK_IMAGE") != "" {
			ns.Node.Image = os.Getenv("CHAINLINK_IMAGE")
		}
	}
	_, err = ns.NewSharedDBNodeSet(in.NodeSets[0], nil)
	if err != nil {
		return fmt.Errorf("failed to create new shared db node set: %w", err)
	}
	if err := Store[Cfg](in); err != nil {
		return err
	}

	// deploy all products and all instances,
	// product config function controls what to read and how to orchestrate each instance
	// via their own TOML part, we only deploy N instances of product M
	for productIdx, productInfo := range in.Products {
		for productInstance := range productInfo.Instances {
			err = productConfigurators[productIdx].ConfigureJobsAndContracts(
				ctx,
				in.FakeServer,
				in.Blockchains[0],
				in.NodeSets[0],
			)
			if err != nil {
				return fmt.Errorf("failed to setup default product deployment: %w", err)
			}
			if err := productConfigurators[productIdx].Store("env-out.toml", productInstance); err != nil {
				return errors.New("failed to store product config")
			}
		}
	}
	L.Info().Str("BootstrapNode", in.NodeSets[0].Out.CLNodes[0].Node.ExternalURL).Send()
	for _, n := range in.NodeSets[0].Out.CLNodes[1:] {
		L.Info().Str("Node", n.Node.ExternalURL).Send()
	}
	return nil
}
`
	// JustFileTmpl is a Justfile template used for building and publishing Docker images
	JustFileTmpl = `set fallback

# Default: show available recipes
default:
    just --list

clean:
    rm -rf compose/ blockscout/

build-fakes:
    just fakes/build

push-fakes:
    just fakes/push

# Rebuild CLI
@cli:
    pushd cmd/{{ .CLIName }} > /dev/null && go install -ldflags="-X main.Version=1.0.0" . && popd > /dev/null`
)

/* Template params in heirarchical order, module -> file(table test) -> test */

// SmokeTestParams params for generating end-to-end test template
type SmokeTestParams struct {
	GoModName   string
	ProductName string
}

// LoadTestParams params for generating end-to-end test template
type LoadTestParams struct {
	GoModName   string
	ProductName string
}

// CISmokeParams params for generating CI smoke tests file
type CISmokeParams struct {
	ProductName   string
	DevEnvRelPath string
	CLIName       string
}

// CILoadChaosParams params for generating CI load&chaos tests file
type CILoadChaosParams struct {
	ProductName   string
	DevEnvRelPath string
	CLIName       string
}

// GoModParams params for generating go.mod file
type GoModParams struct {
	ModuleName     string
	RuntimeVersion string
}

// ReadmeParams params for generating README.md file
type ReadmeParams struct {
	CLIName string
}

// GitIgnoreParams default .gitignore params
type GitIgnoreParams struct{}

// GrafanaDashboardParams default Grafana dashboard params
type GrafanaDashboardParams struct {
	ProductName string
	UUID        string
}

// ConfigTOMLParams default env.toml params
type ConfigTOMLParams struct {
	PackageName string
	ProductName string
	Nodes       int
	NodeIndices []int
}

// JustfileParams Justfile params
type JustfileParams struct {
	PackageName string
	CLIName     string
}

// CLICompletionParams cli.go file params
type CLICompletionParams struct {
	PackageName string
	ProductName string
	CLIName     string
}

// CLIParams cli.go file params
type CLIParams struct {
	PackageName     string
	CLIName         string
	DevEnvPkgImport string
	ProductName     string
	DashboardUUID   string
}

// CLDFParams cldf.go file params
type CLDFParams struct {
	PackageName string
}

// ToolsParams tools.go file params
type ToolsParams struct {
	PackageName string
}

// ConfigParams config.go file params
type ConfigParams struct {
	PackageName string
}

// DevEnvInterfaceParams interface.go file params
type DevEnvInterfaceParams struct {
	PackageName string
}

// EnvParams environment.go file params
type EnvParams struct {
	PackageName string
	ProductName string
}

// ProductConfigurationSimple product_configuration.go file params
type ProductConfigurationSimple struct {
	PackageName string
}

// TableTestParams params for generating a table test
type TableTestParams struct {
	Package       string
	TableTestName string
	TestCases     []TestCaseParams
	WorkloadCode  string
	GunCode       string
}

// TestCaseParams params for generating a test case
type TestCaseParams struct {
	Name    string
	RunFunc string
}

/* Codegen logic */

// EnvBuilder builder for load test codegen
type EnvBuilder struct {
	productName string
	nodes       int
	outputDir   string
	packageName string
	cliName     string
	moduleName  string
}

// EnvCodegen is a load test code generator that creates workload and chaos experiments
type EnvCodegen struct {
	cfg *EnvBuilder
}

// NewEnvBuilder creates a new Chainlink Cluster developer environment
func NewEnvBuilder(cliName string, nodes int, productName string) *EnvBuilder {
	return &EnvBuilder{
		productName: productName,
		cliName:     cliName,
		nodes:       nodes,
		packageName: "devenv",
		outputDir:   "devenv",
	}
}

// OutputDir sets the output directory for generated files
func (g *EnvBuilder) OutputDir(dir string) *EnvBuilder {
	g.outputDir = dir
	return g
}

// validate verifier that we can build codegen with provided params
// fixing input params if possible
func (g *EnvBuilder) validate() error {
	// to avoid issues in different contexts (go.mod, CI, env variables) we only allow lower-case letters, no special symbols or numerics
	g.productName = onlyLetters(g.productName)
	g.cliName = onlyLetters(g.cliName)
	g.moduleName = fmt.Sprintf("github.com/smartcontractkit/%s/devenv", g.productName)
	return nil
}

// onlyLetters strip any symbol except letters
func onlyLetters(name string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsUpper(r) {
			r = []rune(strings.ToLower(string(r)))[0]
			return r
		}
		if unicode.IsLetter(r) {
			return r
		}
		return -1
	}, name)
}

// Validate validate generation params
// for now it's empty but for more complex mutually exclusive cases we should
// add validation here
func (g *EnvBuilder) Build() (*EnvCodegen, error) {
	if err := g.validate(); err != nil {
		return nil, err
	}
	return &EnvCodegen{g}, nil
}

// Read read K8s namespace and find all the pods
// some pods may be crashing but it doesn't matter for code generation
func (g *EnvCodegen) Read() error {
	return nil
}

// Write generates a complete devenv boilerplate, can be multiple files
func (g *EnvCodegen) Write() error {
	// Create output directory
	if err := os.MkdirAll( //nolint:gosec
		g.cfg.outputDir,
		os.ModePerm,
	); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Generate go.mod file
	goModContent, err := g.GenerateGoMod()
	if err != nil {
		return err
	}
	if err := os.WriteFile( //nolint:gosec
		filepath.Join(g.cfg.outputDir, "go.mod"),
		[]byte(goModContent),
		os.ModePerm,
	); err != nil {
		return fmt.Errorf("failed to write go.mod: %w", err)
	}

	// Generate Justfile
	justContents, err := g.GenerateJustfile()
	if err != nil {
		return err
	}
	if err := os.WriteFile( //nolint:gosec
		filepath.Join(g.cfg.outputDir, "Justfile"),
		[]byte(justContents),
		os.ModePerm,
	); err != nil {
		return fmt.Errorf("failed to write CLI completion file: %w", err)
	}

	// Generate default env.toml file
	tomlContents, err := g.GenerateDefaultTOMLConfig()
	if err != nil {
		return err
	}
	if err := os.WriteFile( //nolint:gosec
		filepath.Join(g.cfg.outputDir, "env.toml"),
		[]byte(tomlContents),
		os.ModePerm,
	); err != nil {
		return fmt.Errorf("failed to write default env.toml file: %w", err)
	}

	cliDir := filepath.Join(g.cfg.outputDir, "cmd", g.cfg.cliName)
	if err := os.MkdirAll( //nolint:gosec
		cliDir,
		os.ModePerm,
	); err != nil {
		return fmt.Errorf("failed to create CLI directory: %w", err)
	}

	// Generate CLI file
	cliContents, err := g.GenerateCLI(ProductDashboardUUID)
	if err != nil {
		return err
	}
	if err := os.WriteFile( //nolint:gosec
		filepath.Join(cliDir, fmt.Sprintf("%s.go", g.cfg.cliName)),
		[]byte(cliContents),
		os.ModePerm,
	); err != nil {
		return fmt.Errorf("failed to write CLI file: %w", err)
	}

	// Generate completion file
	completionContents, err := g.GenerateCLICompletion()
	if err != nil {
		return err
	}
	if err := os.WriteFile( //nolint:gosec
		filepath.Join(cliDir, "completion.go"),
		[]byte(completionContents),
		os.ModePerm,
	); err != nil {
		return fmt.Errorf("failed to write CLI completion file: %w", err)
	}

	// Generate README.md file
	readmeContents, err := g.GenerateReadme()
	if err != nil {
		return err
	}
	if err := os.WriteFile( //nolint:gosec
		filepath.Join(g.cfg.outputDir, "README.md"),
		[]byte(readmeContents),
		os.ModePerm,
	); err != nil {
		return fmt.Errorf("failed to write README.md file: %w", err)
	}

	// Generate tools.go
	toolsContents, err := g.GenerateDebugTools()
	if err != nil {
		return err
	}
	if err := os.WriteFile( //nolint:gosec
		filepath.Join(g.cfg.outputDir, "tools.go"),
		[]byte(toolsContents),
		os.ModePerm,
	); err != nil {
		return fmt.Errorf("failed to write tools file: %w", err)
	}

	// Generate config.go
	configFileContents, err := g.GenerateConfig()
	if err != nil {
		return err
	}
	if err := os.WriteFile( //nolint:gosec
		filepath.Join(g.cfg.outputDir, "config.go"),
		[]byte(configFileContents),
		os.ModePerm,
	); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	// Generate jd.go
	cldfContents, err := g.GenerateJD()
	if err != nil {
		return err
	}
	if err := os.WriteFile( //nolint:gosec
		filepath.Join(g.cfg.outputDir, "jd.go"),
		[]byte(cldfContents),
		os.ModePerm,
	); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	// Generate interface.go
	interfaceContents, err := g.GenerateProductsInterface()
	if err != nil {
		return err
	}
	if err := os.WriteFile( //nolint:gosec
		filepath.Join(g.cfg.outputDir, "interface.go"),
		[]byte(interfaceContents),
		os.ModePerm,
	); err != nil {
		return fmt.Errorf("failed to write products interface file: %w", err)
	}

	// Generate environment.go
	envFileContents, err := g.GenerateEnvironment()
	if err != nil {
		return err
	}
	if err := os.WriteFile( //nolint:gosec
		filepath.Join(g.cfg.outputDir, "environment.go"),
		[]byte(envFileContents),
		os.ModePerm,
	); err != nil {
		return fmt.Errorf("failed to write environment file: %w", err)
	}

	// create CI directory
	ciDir := filepath.Join(g.cfg.outputDir, ".github", "workflows")
	if err := os.MkdirAll( //nolint:gosec
		ciDir,
		os.ModePerm,
	); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Generate GitHub CI smoke test workflow
	ciSmokeContents, err := g.GenerateCISmoke()
	if err != nil {
		return err
	}
	if err := os.WriteFile( //nolint:gosec
		filepath.Join(ciDir, "devenv-smoke-test.yml"),
		[]byte(ciSmokeContents),
		os.ModePerm,
	); err != nil {
		return fmt.Errorf("failed to write CI smoke workflow file: %w", err)
	}

	// Generate GitHub CI load&chaos test workflow
	ciLoadChaosContents, err := g.GenerateCILoadChaos()
	if err != nil {
		return err
	}
	if err := os.WriteFile( //nolint:gosec
		filepath.Join(ciDir, "devenv-load-chaos-test.yml"),
		[]byte(ciLoadChaosContents),
		os.ModePerm,
	); err != nil {
		return fmt.Errorf("failed to write CI load&chaos workflow file: %w", err)
	}

	// create e2e tests directory
	e2eDir := filepath.Join(g.cfg.outputDir, "tests")
	if err := os.MkdirAll( //nolint:gosec
		e2eDir,
		os.ModePerm,
	); err != nil {
		return fmt.Errorf("failed to create tests directory: %w", err)
	}

	// generate smoke tests
	smokeTestsContent, err := g.GenerateSmokeTests()
	if err != nil {
		return err
	}
	if err := os.WriteFile( //nolint:gosec
		filepath.Join(e2eDir, "smoke_test.go"),
		[]byte(smokeTestsContent),
		os.ModePerm,
	); err != nil {
		return fmt.Errorf("failed to write smoke tests file: %w", err)
	}

	// generate load/chaos tests
	loadTestsContent, err := g.GenerateLoadTests()
	if err != nil {
		return err
	}
	if err := os.WriteFile( //nolint:gosec
		filepath.Join(e2eDir, "load_test.go"),
		[]byte(loadTestsContent),
		os.ModePerm,
	); err != nil {
		return fmt.Errorf("failed to write load tests file: %w", err)
	}

	// create Grafana dashboards directory
	dashboardsDir := filepath.Join(g.cfg.outputDir, "dashboards")
	if err := os.MkdirAll( //nolint:gosec
		dashboardsDir,
		os.ModePerm,
	); err != nil {
		return fmt.Errorf("failed to create dashboards directory: %w", err)
	}

	// generate Grafana dashboard
	grafanaDashboardContents, err := g.GenerateGrafanaDashboard(ProductDashboardUUID)
	if err != nil {
		return err
	}
	if err := os.WriteFile( //nolint:gosec
		filepath.Join(dashboardsDir, "environment.json"),
		[]byte(grafanaDashboardContents),
		os.ModePerm,
	); err != nil {
		return fmt.Errorf("failed to write Grafana dashboard file: %w", err)
	}

	// generate .gitignore file
	gitIgnoreContents, err := g.GenerateGitIgnore()
	if err != nil {
		return err
	}
	if err := os.WriteFile( //nolint:gosec
		filepath.Join(g.cfg.outputDir, ".gitignore"),
		[]byte(gitIgnoreContents),
		os.ModePerm,
	); err != nil {
		return fmt.Errorf("failed to write gitignore file: %w", err)
	}

	return nil
}

// GenerateSmokeTests generates a smoke test template
func (g *EnvCodegen) GenerateLoadTests() (string, error) {
	log.Info().Msg("Generating load test template")
	data := LoadTestParams{
		GoModName:   g.cfg.moduleName,
		ProductName: g.cfg.productName,
	}
	return render(LoadTestTmpl, data)
}

// GenerateSmokeTests generates a smoke test template
func (g *EnvCodegen) GenerateSmokeTests() (string, error) {
	log.Info().Msg("Generating smoke test template")
	data := SmokeTestParams{
		GoModName:   g.cfg.moduleName,
		ProductName: g.cfg.productName,
	}
	return render(SmokeTestImplTmpl, data)
}

// GenerateCILoadChaos generates a load&chaos test CI workflow
func (g *EnvCodegen) GenerateCILoadChaos() (string, error) {
	log.Info().Msg("Generating GitHub CI load&chaos test")
	p, err := dirPathRelFromGitRoot(g.cfg.outputDir)
	if err != nil {
		return "", fmt.Errorf("failed to find relative devenv path from Git root: %w", err)
	}
	data := CILoadChaosParams{
		DevEnvRelPath: p,
		ProductName:   g.cfg.productName,
		CLIName:       g.cfg.cliName,
	}
	return render(CILoadChaosTemplate, data)
}

// GenerateCISmoke generates a smoke test CI workflow
func (g *EnvCodegen) GenerateCISmoke() (string, error) {
	log.Info().Msg("Generating GitHub CI smoke test")
	p, err := dirPathRelFromGitRoot(g.cfg.outputDir)
	if err != nil {
		return "", fmt.Errorf("failed to find relative devenv path from Git root: %w", err)
	}
	data := CISmokeParams{
		DevEnvRelPath: p,
		ProductName:   g.cfg.productName,
		CLIName:       g.cfg.cliName,
	}
	return render(CISmokeTmpl, data)
}

// GenerateGoMod generates a go.mod file
func (g *EnvCodegen) GenerateGoMod() (string, error) {
	log.Info().Msg("Generating Go module")
	data := GoModParams{
		ModuleName:     g.cfg.moduleName,
		RuntimeVersion: strings.ReplaceAll(runtime.Version(), "go", ""),
	}
	return render(GoModTemplate, data)
}

// GenerateReadme generates a readme file
func (g *EnvCodegen) GenerateReadme() (string, error) {
	log.Info().Msg("Generating README file")
	data := ReadmeParams{
		CLIName: g.cfg.cliName,
	}
	return render(ReadmeTmpl, data)
}

// GenerateGitIgnore generate .gitignore file
func (g *EnvCodegen) GenerateGitIgnore() (string, error) {
	log.Info().Msg("Generating .gitignore file")
	p := GitIgnoreParams{}
	return render(GitIgnoreTmpl, p)
}

// GenerateGrafanaDashboard generate default Grafana dashboard
func (g *EnvCodegen) GenerateGrafanaDashboard(uuid string) (string, error) {
	log.Info().Msg("Generating default environment dashboard for Grafana")
	p := GrafanaDashboardParams{
		ProductName: g.cfg.productName,
		UUID:        uuid,
	}
	return render(GrafanaDashboardTmpl, p)
}

// GenerateDefaultTOMLConfig generate default env.toml config
func (g *EnvCodegen) GenerateDefaultTOMLConfig() (string, error) {
	log.Info().Msg("Generating default environment config (env.toml)")
	p := ConfigTOMLParams{
		PackageName: g.cfg.packageName,
		ProductName: g.cfg.productName,
		Nodes:       g.cfg.nodes,
		NodeIndices: make([]int, g.cfg.nodes),
	}
	return render(ConfigTOMLTmpl, p)
}

// GenerateJustfile generate Justfile to build and publish Docker images
func (g *EnvCodegen) GenerateJustfile() (string, error) {
	log.Info().Msg("Generating Justfile")
	p := JustfileParams{
		PackageName: g.cfg.packageName,
		CLIName:     g.cfg.cliName,
	}
	return render(JustFileTmpl, p)
}

// GenerateCLICompletion generate CLI completion for "go-prompt" library
func (g *EnvCodegen) GenerateCLICompletion() (string, error) {
	log.Info().Msg("Generating shell completion")
	p := CLICompletionParams{
		PackageName: g.cfg.packageName,
		ProductName: g.cfg.productName,
		CLIName:     g.cfg.cliName,
	}
	return render(CompletionTmpl, p)
}

// GenerateCLI generate Cobra CLI
func (g *EnvCodegen) GenerateCLI(dashboardUUID string) (string, error) {
	log.Info().Msg("Generating Cobra CLI")
	p := CLIParams{
		PackageName:     g.cfg.packageName,
		CLIName:         g.cfg.cliName,
		ProductName:     g.cfg.productName,
		DevEnvPkgImport: g.cfg.moduleName,
		DashboardUUID:   dashboardUUID,
	}
	return render(CLITmpl, p)
}

// GenerateEnvironment generate environment.go, our environment composition function
func (g *EnvCodegen) GenerateEnvironment() (string, error) {
	log.Info().Msg("Generating environment composition (environment.go)")
	p := EnvParams{
		PackageName: g.cfg.packageName,
		ProductName: g.cfg.productName,
	}
	return render(EnvironmentTmpl, p)
}

// GenerateJD generate JD helpers
func (g *EnvCodegen) GenerateJD() (string, error) {
	log.Info().Msg("Generating JD helpers")
	p := CLDFParams{
		PackageName: g.cfg.packageName,
	}
	return render(JDTmpl, p)
}

// GenerateDebugTools generate debug tools (tracing)
func (g *EnvCodegen) GenerateDebugTools() (string, error) {
	log.Info().Msg("Generating debug tools")
	p := ToolsParams{
		PackageName: g.cfg.packageName,
	}
	return render(DebugToolsTmpl, p)
}

// GenerateConfig generate read/write utilities for TOML configs
func (g *EnvCodegen) GenerateConfig() (string, error) {
	log.Info().Msg("Generating config tools")
	p := ConfigParams{
		PackageName: g.cfg.packageName,
	}
	return render(ConfigTmpl, p)
}

// GenerateProductsInterface generate devenv interface to run arbitrary products
func (g *EnvCodegen) GenerateProductsInterface() (string, error) {
	log.Info().Msg("Generating devenv interface")
	p := DevEnvInterfaceParams{
		PackageName: g.cfg.packageName,
	}
	return render(ProductsInterfaceTmpl, p)
}

// GenerateTableTest generates all possible experiments for a namespace
// first generate all small pieces then insert into a table test template
func (g *EnvCodegen) GenerateTableTest() (string, error) {
	log.Info().Msg("Generating end-to-end table tests template")
	// TODO: generate a table test when we'll have chain-specific interface solidified
	return "", nil
}

// GenerateTestCases generates table test cases
func (g *EnvCodegen) GenerateTestCases() ([]TestCaseParams, error) {
	log.Info().Msg("Generating test cases")
	// TODO: generate test cases when we'll have chain-specific interface solidified
	return []TestCaseParams{}, nil
}

/* Utility */

func gitRoot() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("not a git repository or git not installed: %w", err)
	}
	return strings.TrimSpace(string(output)), nil
}

// dirPathRelFromGitRoot gets directory path relative from Git root
// used in CI templates to set default execution directory
func dirPathRelFromGitRoot(name string) (string, error) {
	gitRoot, err := gitRoot()
	if err != nil {
		return "", err
	}
	var devenvPath string
	err = filepath.Walk(gitRoot, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() && info.Name() == name {
			devenvPath, err = filepath.Rel(gitRoot, path)
			if err != nil {
				return err
			}
			return filepath.SkipAll
		}
		return nil
	})
	if err != nil {
		return "", fmt.Errorf("devenv directory not found: %w", err)
	}
	return devenvPath, nil
}

// render is just an internal function to parse and render template
func render(tmpl string, data any) (string, error) {
	parsed, err := template.New("").Parse(tmpl)
	if err != nil {
		return "", fmt.Errorf("failed to parse text template: %w", err)
	}
	var buf bytes.Buffer
	if err := parsed.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to generate text: %w", err)
	}
	return buf.String(), err
}
