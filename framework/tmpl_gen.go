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
	ConfigTOMLTmpl = `[on_chain]
  link_contract_address = "0xDc64a140Aa3E981100a9becA4E685f962f0cF6C9"
  cl_nodes_funding_eth = 50
  cl_nodes_funding_link = 50
  verification_timeout_sec = 400
  contracts_configuration_timeout_sec = 60
  verify = false

  [on_chain.gas_settings]
  fee_cap_multiplier = 2
  tip_cap_multiplier = 2


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
			{Text: "env.toml", Description: "Spin up Anvil <> Anvil local chains, all services, 4 CL nodes"},
			{Text: "env.toml,env-cl-rebuild.toml", Description: "Spin up Anvil <> Anvil local chains, all services, 4 CL nodes (custom build)"},
			{Text: "env.toml,env-geth.toml", Description: "Spin up Geth <> Geth local chains (clique), all services, 4 CL nodes"},
			{Text: "env.toml,env-fuji-fantom.toml", Description: "Spin up testnets: Fuji <> Fantom, all services, 4 CL nodes"},
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
			configFile = "env.toml"
		}
		framework.L.Info().Str("Config", configFile).Msg("Reconfiguring development environment")
		_ = os.Setenv("CTF_CONFIGS", configFile)
		_ = os.Setenv("TESTCONTAINERS_RYUK_DISABLED", "true")
		framework.L.Info().Msg("Tearing down the development environment")
		err := framework.RemoveTestContainers()
		if err != nil {
			return fmt.Errorf("failed to clean Docker resources: %w", err)
		}
		_, err = devenv.NewEnvironment()
		return err
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
			configFile = "env.toml"
		}
		framework.L.Info().Str("Config", configFile).Msg("Creating development environment")
		_ = os.Setenv("CTF_CONFIGS", configFile)
		_ = os.Setenv("TESTCONTAINERS_RYUK_DISABLED", "true")
		_, err := devenv.NewEnvironment()
		if err != nil {
			return err
		}
		return nil
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
		devenv.Plog.Info().Msgf("{{ .ProductName }} Dashboard: %s", Local{{ .ProductName }}Dashboard)
		devenv.Plog.Info().Msgf("{{ .ProductName }} Load Test Dashboard: %s", LocalWASPLoadDashboard)
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
		devenv.Plog.Info().Msgf("{{ .ProductName }} Dashboard: %s", Local{{ .ProductName }}Dashboard)
		devenv.Plog.Info().Msgf("{{ .ProductName }} Load Test Dashboard: %s", LocalWASPLoadDashboard)
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
		devenv.Plog.Err(err).Send()
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
	// SmokeTestTmpl is a smoke test template
	SmokeTestTmpl = `package devenv_test

import (
	"fmt"
	"strconv"
	"testing"
	"time"

	de "{{ .GoModName }}"

	"github.com/smartcontractkit/chainlink-testing-framework/framework/clclient"
	"github.com/stretchr/testify/require"

	f "github.com/smartcontractkit/chainlink-testing-framework/framework"
)

var L = de.Plog

func TestSmoke(t *testing.T) {
	in, err := de.LoadOutput[de.Cfg]("../env-out.toml")
	require.NoError(t, err)
	c, _, _, err := de.ETHClient(in)
	require.NoError(t, err)
	clNodes, err := clclient.New(in.NodeSets[0].Out.CLNodes)
	require.NoError(t, err)

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
			_, _ = c, clNodes
		})
	}
}

// assertResources is a simple assertion on resources if you run with the observability stack (obs up)
func assertResources(t *testing.T, in *de.Cfg, start, end time.Time) {
	pc := f.NewPrometheusQueryClient(f.LocalPrometheusBaseURL)
	// no more than 10% CPU for this test
	maxCPU := 10.0
	cpuResp, err := pc.Query("sum(rate(container_cpu_usage_seconds_total{name=~\".*don.*\"}[5m])) by (name) *100", end)
	require.NoError(t, err)
	cpu := f.ToLabelsMap(cpuResp)
	for i := 0; i < in.NodeSets[0].Nodes; i++ {
		nodeLabel := fmt.Sprintf("name:don-node%d", i)
		nodeCpu, err := strconv.ParseFloat(cpu[nodeLabel][0].(string), 64)
		L.Info().Int("Node", i).Float64("CPU", nodeCpu).Msg("CPU usage percentage")
		require.NoError(t, err)
		require.LessOrEqual(t, nodeCpu, maxCPU)
	}
	// no more than 200mb for this test
	maxMem := int(200e6) // 200mb
	memoryResp, err := pc.Query("sum(container_memory_rss{name=~\".*don.*\"}) by (name)", end)
	require.NoError(t, err)
	mem := f.ToLabelsMap(memoryResp)
	for i := 0; i < in.NodeSets[0].Nodes; i++ {
		nodeLabel := fmt.Sprintf("name:don-node%d", i)
		nodeMem, err := strconv.Atoi(mem[nodeLabel][0].(string))
		L.Info().Int("Node", i).Int("Memory", nodeMem).Msg("Total memory")
		require.NoError(t, err)
		require.LessOrEqual(t, nodeMem, maxMem)
	}
}
`
	// CLDFTmpl is a Chainlink Deployments Framework template
	CLDFTmpl = `package {{ .PackageName }}

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/shared/generated/link_token"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	chainsel "github.com/smartcontractkit/chain-selectors"
	cldfchain "github.com/smartcontractkit/chainlink-deployments-framework/chain"
	cldfevm "github.com/smartcontractkit/chainlink-deployments-framework/chain/evm"
	cldfevmprovider "github.com/smartcontractkit/chainlink-deployments-framework/chain/evm/provider"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	csav1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/csa"
	jobv1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/job"
	nodev1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/node"
)

const (
	AnvilKey0                     = "ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"
	DefaultNativeTransferGasPrice = 21000
)

const LinkToken cldf.ContractType = "LinkToken"

var _ cldf.ChangeSet[[]uint64] = DeployLinkToken

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

// DeployLinkToken deploys a link token contract to the chain identified by the ChainSelector.
func DeployLinkToken(e cldf.Environment, chains []uint64) (cldf.ChangesetOutput, error) { //nolint:gocritic
	newAddresses := cldf.NewMemoryAddressBook()
	deployGrp := errgroup.Group{}
	for _, chain := range chains {
		family, err := chainsel.GetSelectorFamily(chain)
		if err != nil {
			return cldf.ChangesetOutput{AddressBook: newAddresses}, err
		}
		var deployFn func() error
		switch family {
		case chainsel.FamilyEVM:
			// Deploy EVM LINK token
			deployFn = func() error {
				_, err := deployLinkTokenContractEVM(
					e.Logger, e.BlockChains.EVMChains()[chain], newAddresses,
				)
				return err
			}
		default:
			return cldf.ChangesetOutput{}, fmt.Errorf("unsupported chain family %s", family)
		}
		deployGrp.Go(func() error {
			err := deployFn()
			if err != nil {
				e.Logger.Errorw("Failed to deploy link token", "chain", chain, "err", err)
				return fmt.Errorf("failed to deploy link token for chain %d: %w", chain, err)
			}
			return nil
		})
	}
	return cldf.ChangesetOutput{AddressBook: newAddresses}, deployGrp.Wait()
}

func deployLinkTokenContractEVM(
		lggr logger.Logger,
		chain cldfevm.Chain, //nolint:gocritic
		ab cldf.AddressBook,
) (*cldf.ContractDeploy[*link_token.LinkToken], error) {
	linkToken, err := cldf.DeployContract[*link_token.LinkToken](lggr, chain, ab,
		func(chain cldfevm.Chain) cldf.ContractDeploy[*link_token.LinkToken] {
			var (
				linkTokenAddr common.Address
				tx            *types.Transaction
				linkToken     *link_token.LinkToken
				err2          error
			)
			if !chain.IsZkSyncVM {
				linkTokenAddr, tx, linkToken, err2 = link_token.DeployLinkToken(
					chain.DeployerKey,
					chain.Client,
				)
			} else {
				linkTokenAddr, _, linkToken, err2 = link_token.DeployLinkTokenZk(
					nil,
					chain.ClientZkSyncVM,
					chain.DeployerKeyZkSyncVM,
					chain.Client,
				)
			}
			return cldf.ContractDeploy[*link_token.LinkToken]{
				Address:  linkTokenAddr,
				Contract: linkToken,
				Tx:       tx,
				Tv:       cldf.NewTypeAndVersion(LinkToken, *semver.MustParse("1.0.0")),
				Err:      err2,
			}
		})
	if err != nil {
		lggr.Errorw("Failed to deploy link token", "chain", chain.String(), "err", err)
		return linkToken, err
	}
	return linkToken, nil
}

// LoadCLDFEnvironment loads CLDF environment with a memory data store and JD client.
func LoadCLDFEnvironment(in *Cfg) (cldf.Environment, error) {
	ctx := context.Background()

	getCtx := func() context.Context {
		return ctx
	}

	// This only generates a brand new datastore and does not load any existing data.
	// We will need to figure out how data will be persisted and loaded in the future.
	ds := datastore.NewMemoryDataStore().Seal()

	lggr, err := logger.NewWith(func(config *zap.Config) {
		config.Development = true
		config.Encoding = "console"
	})
	if err != nil {
		return cldf.Environment{}, fmt.Errorf("failed to create logger: %w", err)
	}

	blockchains, err := loadCLDFChains(in.Blockchains)
	if err != nil {
		return cldf.Environment{}, fmt.Errorf("failed to load CLDF chains: %w", err)
	}

	jd, err := NewJDClient(ctx, JDConfig{
		GRPC:  in.JD.Out.ExternalGRPCUrl,
		WSRPC: in.JD.Out.ExternalWSRPCUrl,
	})
	if err != nil {
		return cldf.Environment{},
			fmt.Errorf("failed to load offchain client: %w", err)
	}

	opBundle := operations.NewBundle(
		getCtx,
		lggr,
		operations.NewMemoryReporter(),
		operations.WithOperationRegistry(operations.NewOperationRegistry()),
	)

	return cldf.Environment{
		Name:              "local",
		Logger:            lggr,
		ExistingAddresses: cldf.NewMemoryAddressBook(),
		DataStore:         ds,
		Offchain:          jd,
		GetContext:        getCtx,
		OperationsBundle:  opBundle,
		BlockChains:       cldfchain.NewBlockChainsFromSlice(blockchains),
	}, nil
}

func loadCLDFChains(bcis []*blockchain.Input) ([]cldfchain.BlockChain, error) {
	blockchains := make([]cldfchain.BlockChain, 0)
	for _, bci := range bcis {
		switch bci.Type {
		case "anvil":
			bc, err := loadEVMChain(bci)
			if err != nil {
				return blockchains, fmt.Errorf("failed to load EVM chain %s: %w", bci.ChainID, err)
			}

			blockchains = append(blockchains, bc)
		default:
			return blockchains, fmt.Errorf("unsupported chain type %s", bci.Type)
		}
	}

	return blockchains, nil
}

func loadEVMChain(bci *blockchain.Input) (cldfchain.BlockChain, error) {
	if bci.Out == nil {
		return nil, fmt.Errorf("output configuration for %s blockchain %s is not set", bci.Type, bci.ChainID)
	}

	chainDetails, err := chainsel.GetChainDetailsByChainIDAndFamily(bci.ChainID, chainsel.FamilyEVM)
	if err != nil {
		return nil, fmt.Errorf("failed to get chain details for %s: %w", bci.ChainID, err)
	}

	chain, err := cldfevmprovider.NewRPCChainProvider(
		chainDetails.ChainSelector,
		cldfevmprovider.RPCChainProviderConfig{
			DeployerTransactorGen: cldfevmprovider.TransactorFromRaw(
				// TODO: we need to figure out a reliable way to get secrets here that is
				// TODO: - easy for developers
				// TODO: - works the same way locally, in K8s and in CI
				// TODO: - do not require specific AWS access like AWSSecretsManager
				// TODO: for now it's just an Anvil 0 key
				AnvilKey0,
			),
			RPCs: []cldf.RPC{
				{
					Name:               "default",
					WSURL:              bci.Out.Nodes[0].ExternalWSUrl,
					HTTPURL:            bci.Out.Nodes[0].ExternalHTTPUrl,
					PreferredURLScheme: cldf.URLSchemePreferenceHTTP,
				},
			},
			ConfirmFunctor: cldfevmprovider.ConfirmFuncGeth(1 * time.Minute),
		},
	).Initialize(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to initialize EVM chain %s: %w", bci.ChainID, err)
	}

	return chain, nil
}

// NewJDClient creates a new JobDistributor client.
func NewJDClient(ctx context.Context, cfg JDConfig) (cldf.OffchainClient, error) {
	conn, err := NewJDConnection(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to connect Job Distributor service. Err: %w", err)
	}
	jd := &JobDistributor{
		WSRPC:             cfg.WSRPC,
		NodeServiceClient: nodev1.NewNodeServiceClient(conn),
		JobServiceClient:  jobv1.NewJobServiceClient(conn),
		CSAServiceClient:  csav1.NewCSAServiceClient(conn),
	}

	return jd, err
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

// FundNodeEIP1559 funds CL node using RPC URL, recipient address and amount of funds to send (ETH).
// Uses EIP-1559 transaction type.
func FundNodeEIP1559(c *ethclient.Client, pkey, recipientAddress string, amountOfFundsInETH float64) error {
	amount := new(big.Float).Mul(big.NewFloat(amountOfFundsInETH), big.NewFloat(1e18))
	amountWei, _ := amount.Int(nil)
	Plog.Info().Str("Addr", recipientAddress).Str("Wei", amountWei.String()).Msg("Funding Node")

	chainID, err := c.NetworkID(context.Background())
	if err != nil {
		return err
	}
	privateKeyStr := strings.TrimPrefix(pkey, "0x")
	privateKey, err := crypto.HexToECDSA(privateKeyStr)
	if err != nil {
		return err
	}
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return errors.New("error casting public key to ECDSA")
	}
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	nonce, err := c.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		return err
	}
	feeCap, err := c.SuggestGasPrice(context.Background())
	if err != nil {
		return err
	}
	tipCap, err := c.SuggestGasTipCap(context.Background())
	if err != nil {
		return err
	}
	recipient := common.HexToAddress(recipientAddress)
	tx := types.NewTx(&types.DynamicFeeTx{
		ChainID:   chainID,
		Nonce:     nonce,
		To:        &recipient,
		Value:     amountWei,
		Gas:       DefaultNativeTransferGasPrice,
		GasFeeCap: feeCap,
		GasTipCap: tipCap,
	})
	signedTx, err := types.SignTx(tx, types.NewLondonSigner(chainID), privateKey)
	if err != nil {
		return err
	}
	err = c.SendTransaction(context.Background(), signedTx)
	if err != nil {
		return err
	}
	if _, err := bind.WaitMined(context.Background(), c, signedTx); err != nil {
		return err
	}
	Plog.Info().Str("Wei", amountWei.String()).Msg("Funded with ETH")
	return nil
}

/*
This is just a basic ETH client, CLDF should provide something like this
*/

// ETHClient creates a basic Ethereum client using PRIVATE_KEY env var and tip/cap gas settings
func ETHClient(in *Cfg) (*ethclient.Client, *bind.TransactOpts, string, error) {
	rpcURL := in.Blockchains[0].Out.Nodes[0].ExternalWSUrl
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		return nil, nil, "", fmt.Errorf("could not connect to eth client: %w", err)
	}
	privateKey, err := crypto.HexToECDSA(GetNetworkPrivateKey())
	if err != nil {
		return nil, nil, "", fmt.Errorf("could not parse private key: %w", err)
	}
	publicKey := privateKey.PublicKey
	address := crypto.PubkeyToAddress(publicKey).String()
	chainID, err := client.ChainID(context.Background())
	if err != nil {
		return nil, nil, "", fmt.Errorf("could not get chain ID: %w", err)
	}
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		return nil, nil, "", fmt.Errorf("could not create transactor: %w", err)
	}
	gasSettings := in.OnChain.GasSettings
	fc, tc, err := MultiplyEIP1559GasPrices(client, gasSettings.FeeCapMultiplier, gasSettings.TipCapMultiplier)
	if err != nil {
		return nil, nil, "", fmt.Errorf("could not get bumped gas price: %w", err)
	}
	auth.GasFeeCap = fc
	auth.GasTipCap = tc
	Plog.Info().
		Str("GasFeeCap", fc.String()).
		Str("GasTipCap", tc.String()).
		Msg("Default gas prices set")
	return client, auth, address, nil
}

// MultiplyEIP1559GasPrices returns bumped EIP1159 gas prices increased by multiplier
func MultiplyEIP1559GasPrices(client *ethclient.Client, fcMult, tcMult int64) (*big.Int, *big.Int, error) {
	feeCap, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		return nil, nil, err
	}
	tipCap, err := client.SuggestGasTipCap(context.Background())
	if err != nil {
		return nil, nil, err
	}

	return new(big.Int).Mul(feeCap, big.NewInt(fcMult)), new(big.Int).Mul(tipCap, big.NewInt(tcMult)), nil
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

/*
This file provides a simple boilerplate for TOML configuration with overrides
It has 4 functions: Load[T], Store[T], LoadCache[T] and GetNetworkPrivateKey

To configure the environment we use a set of files we read from the env var CTF_CONFIGS=env.toml,overrides.toml (can be more than 2) in Load[T]
To store infra or product component outputs we use Store[T] that creates env-cache.toml file.
This file can be used in tests or in any other code that integrated with dev environment.
LoadCache[T] is used if you need to write outputs the second time.

GetNetworkPrivateKey is used to get your network private key from the env var we are using across all our environments, or fallback to default Anvil's key.
*/

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/davecgh/go-spew/spew"
	"github.com/pelletier/go-toml/v2"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
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

var L = log.Output(zerolog.ConsoleWriter{Out: os.Stderr}).Level(zerolog.InfoLevel)

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
		decoder.DisallowUnknownFields()

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
	"fmt"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/jd"

	ns "github.com/smartcontractkit/chainlink-testing-framework/framework/components/simple_node_set"
)

type Cfg struct {
    OnChain         *OnChain                ` + "`" + `toml:"on_chain"` + "`" + `
    Blockchains []*blockchain.Input ` + "`" + `toml:"blockchains" validate:"required"` + "`" + `
    NodeSets    []*ns.Input         ` + "`" + `toml:"nodesets"    validate:"required"` + "`" + `
    JD          *jd.Input           ` + "`" + `toml:"jd"` + "`" + `
}

func NewEnvironment() (*Cfg, error) {
	endTracing := tracing()
	defer endTracing()

	if err := framework.DefaultNetwork(nil); err != nil {
		return nil, err
	}
	in, err := Load[Cfg]()
	if err != nil {
		return nil, fmt.Errorf("failed to load configuration: %w", err)
	}
	_, err = blockchain.NewBlockchainNetwork(in.Blockchains[0])
	if err != nil {
		return nil, fmt.Errorf("failed to create blockchain network 1337: %w", err)
	}
	if err := DefaultProductConfiguration(in, ConfigureNodesNetwork); err != nil {
		return nil, fmt.Errorf("failed to setup default CLDF orchestration: %w", err)
	}
	_, err = ns.NewSharedDBNodeSet(in.NodeSets[0], nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create new shared db node set: %w", err)
	}
	if err := DefaultProductConfiguration(in, ConfigureProductContractsJobs); err != nil {
		return nil, fmt.Errorf("failed to setup default CLDF orchestration: %w", err)
	}
	return in, Store[Cfg](in)
}
`
	// SingleNetworkProductConfigurationTmpl is an single-network EVM product configuration template
	SingleNetworkEVMProductConfigurationTmpl = `package {{ .PackageName }}

import (
	"context"
	"fmt"
	"math/big"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/shared/generated/link_token"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/clclient"
)

const (
		ConfigureNodesNetwork ConfigPhase = iota
		ConfigureProductContractsJobs
)

var Plog = log.Output(zerolog.ConsoleWriter{Out: os.Stderr}).Level(zerolog.DebugLevel).With().Fields(map[string]any{"component": "on_chain"}).Logger()

type OnChain struct {
	LinkContractAddress              string                 ` + "`" + `toml:"link_contract_address"` + "`" + `
	CLNodesFundingETH                float64                ` + "`" + `toml:"cl_nodes_funding_eth"` + "`" + `
	CLNodesFundingLink               float64                ` + "`" + `toml:"cl_nodes_funding_link"` + "`" + `
	VerificationTimeoutSec           time.Duration          ` + "`" + `toml:"verification_timeout_sec"` + "`" + `
	ContractsConfigurationTimeoutSec time.Duration          ` + "`" + `toml:"contracts_configuration_timeout_sec"` + "`" + `
	GasSettings                      *GasSettings           ` + "`" + `toml:"gas_settings"` + "`" + `
	Verify                           bool                   ` + "`" + `toml:"verify"` + "`" + `
	DeployedContracts                *DeployedContracts     ` + "`" + `toml:"deployed_contracts"` + "`" + `
}

type DeployedContracts struct {
	SomeContractAddr string ` + "`" + `toml:"some_contract_addr"` + "`" + `
}


type GasSettings struct {
	FeeCapMultiplier int64 ` + "`" + `toml:"fee_cap_multiplier"` + "`" + `
	TipCapMultiplier int64 ` + "`" + `toml:"tip_cap_multiplier"` + "`" + `
}

type Jobs struct {
	ConfigPollIntervalSeconds time.Duration ` + "`" + `toml:"config_poll_interval_sec"` + "`" + `
	MaxTaskDurationSec        time.Duration ` + "`" + `toml:"max_task_duration_sec"` + "`" + `
}

type ConfigPhase int

// deployLinkAndMint is a universal action that deploys link token and mints required amount of LINK token for all the nodes.
func deployLinkAndMint(ctx context.Context, in *Cfg, c *ethclient.Client, auth *bind.TransactOpts, rootAddr string, transmitters []common.Address) (*link_token.LinkToken, error) {
	addr, tx, lt, err := link_token.DeployLinkToken(auth, c)
	if err != nil {
		return nil, fmt.Errorf("could not create link token contract: %w", err)
	}
	_, err = bind.WaitDeployed(ctx, c, tx)
	if err != nil {
		return nil, err
	}
	Plog.Info().Str("Address", addr.Hex()).Msg("Deployed link token contract")
	tx, err = lt.GrantMintRole(auth, common.HexToAddress(rootAddr))
	if err != nil {
		return nil, fmt.Errorf("could not grant mint role: %w", err)
	}
	_, err = bind.WaitMined(ctx, c, tx)
	if err != nil {
		return nil, err
	}
	// mint for public keys of nodes directly instead of transferring
	for _, transmitter := range transmitters {
		amount := new(big.Float).Mul(big.NewFloat(in.OnChain.CLNodesFundingLink), big.NewFloat(1e18))
		amountWei, _ := amount.Int(nil)
		Plog.Info().Msgf("Minting LINK for transmitter address: %s", transmitter.Hex())
		tx, err = lt.Mint(auth, transmitter, amountWei)
		if err != nil {
			return nil, fmt.Errorf("could not transfer link token contract: %w", err)
		}
		_, err = bind.WaitMined(ctx, c, tx)
		if err != nil {
			return nil, err
		}
	}
	return lt, nil
}


func configureContracts(in *Cfg, c *ethclient.Client, auth *bind.TransactOpts, cl []*clclient.ChainlinkClient, rootAddr string, transmitters []common.Address) (*DeployedContracts, error) {
	ctx, cancel := context.WithTimeout(context.Background(), in.OnChain.ContractsConfigurationTimeoutSec*time.Second)
	defer cancel()
	Plog.Info().Msg("Deploying LINK token contract")
	_, err := deployLinkAndMint(ctx, in, c, auth, rootAddr, transmitters)
	if err != nil {
		return nil, fmt.Errorf("could not create link token contract and mint: %w", err)
	}
	// TODO: use client and deploy your contracts
	return &DeployedContracts{
		SomeContractAddr: "",
	}, nil
}

func configureJobs(in *Cfg, clNodes []*clclient.ChainlinkClient, contracts *DeployedContracts) error {
	bootstrapNode := clNodes[0]
	workerNodes := clNodes[1:]
	// TODO: define your jobs
	job := ""
	_, _, err := bootstrapNode.CreateJobRaw(job)
	if err != nil {
		return fmt.Errorf("creating bootstrap job have failed: %w", err)
	}

	for _, chainlinkNode := range workerNodes {
		// TODO: define your job for nodes here
		job := ""
		_, _, err = chainlinkNode.CreateJobRaw(job)
		if err != nil {
			return fmt.Errorf("creating job on node have failed: %w", err)
		}
	}
	return nil
}

// DefaultProductConfiguration is default product configuration that includes:
// - Deploying required prerequisites (LINK token, shared contracts)
// - Applying product-specific changesets
// - Creating cldf.Environment, connecting to components, see *Cfg fields
// - Generating CL nodes configs
// All the data can be added *Cfg struct like and is synced between local machine and remote environment
// so later both local and remote tests can use it.
func DefaultProductConfiguration(in *Cfg, phase ConfigPhase) error {
	pkey := GetNetworkPrivateKey()
	if pkey == "" {
		return fmt.Errorf("PRIVATE_KEY environment variable not set")
	}
	switch phase {
	case ConfigureNodesNetwork:
		Plog.Info().Msg("Applying default CL nodes configuration")
		node := in.Blockchains[0].Out.Nodes[0]
		chainID := in.Blockchains[0].ChainID
		// configure node set and generate CL nodes configs
		netConfig := fmt.Sprintf(` + "`" + `
    [[EVM]]
    LogPollInterval = '1s'
    BlockBackfillDepth = 100
    LinkContractAddress = '%s'
    ChainID = '%s'
    MinIncomingConfirmations = 1
    MinContractPayment = '0.0000001 link'
    FinalityDepth = %d

    [[EVM.Nodes]]
    Name = 'default'
    WsUrl = '%s'
    HttpUrl = '%s'

    [Feature]
    FeedsManager = true
    LogPoller = true
    UICSAKeys = true
    [OCR2]
    Enabled = true
    SimulateTransactions = false
    DefaultTransactionQueueDepth = 1
    [P2P.V2]
    Enabled = true
    ListenAddresses = ['0.0.0.0:6690']

	   [Log]
   JSONConsole = true
   Level = 'debug'
   [Pyroscope]
   ServerAddress = 'http://host.docker.internal:4040'
   Environment = 'local'
   [WebServer]
    SessionTimeout = '999h0m0s'
    HTTPWriteTimeout = '3m'
   SecureCookies = false
   HTTPPort = 6688
   [WebServer.TLS]
   HTTPSPort = 0
    [WebServer.RateLimit]
    Authenticated = 5000
    Unauthenticated = 5000
   [JobPipeline]
   [JobPipeline.HTTPRequest]
   DefaultTimeout = '1m'
    [Log.File]
    MaxSize = '0b'
` + "`" + `, in.OnChain.LinkContractAddress, chainID, 5, node.InternalWSUrl, node.InternalHTTPUrl)
		for _, nodeSpec := range in.NodeSets[0].NodeSpecs {
			nodeSpec.Node.TestConfigOverrides = netConfig
		}
		Plog.Info().Msg("Nodes network configuration is finished")
	case ConfigureProductContractsJobs:
		Plog.Info().Msg("Connecting to CL nodes")
		nodeClients, err := clclient.New(in.NodeSets[0].Out.CLNodes)
		if err != nil {
			return err
		}
		transmitters := make([]common.Address, 0)
		ethKeyAddresses := make([]string, 0)
		for i, nc := range nodeClients {
			addr, err := nc.ReadPrimaryETHKey(in.Blockchains[0].ChainID)
			if err != nil {
				return err
			}
			ethKeyAddresses = append(ethKeyAddresses, addr.Attributes.Address)
			transmitters = append(transmitters, common.HexToAddress(addr.Attributes.Address))
			Plog.Info().
				Int("Idx", i).
				Str("ETH", addr.Attributes.Address).
				Msg("Node info")
		}
		// ETH examples
		c, auth, rootAddr, err := ETHClient(in)
		if err != nil {
			return fmt.Errorf("could not create basic eth client: %w", err)
		}
		for _, addr := range ethKeyAddresses {
			if err := FundNodeEIP1559(c, pkey, addr, in.OnChain.CLNodesFundingETH); err != nil {
				return err
			}
		}
		contracts, err := configureContracts(in, c, auth, nodeClients, rootAddr, transmitters)
		if err != nil {
			return err
		}
		if err := configureJobs(in, nodeClients, contracts); err != nil {
			return err
		}
		Plog.Info().Str("BootstrapNode", in.NodeSets[0].Out.CLNodes[0].Node.ExternalURL).Send()
		for _, n := range in.NodeSets[0].Out.CLNodes[1:] {
			Plog.Info().Str("Node", n.Node.ExternalURL).Send()
		}
		in.OnChain.DeployedContracts = contracts
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
	GoModName string
}

// LoadTestParams params for generating end-to-end test template
type LoadTestParams struct {
	GoModName string
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

// EnvParams environment.go file params
type EnvParams struct {
	PackageName string
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
	productType string
	moduleName  string
}

// EnvCodegen is a load test code generator that creates workload and chaos experiments
type EnvCodegen struct {
	cfg *EnvBuilder
}

// NewEnvBuilder creates a new Chainlink Cluster developer environment
func NewEnvBuilder(cliName string, nodes int, productType string, productName string) *EnvBuilder {
	return &EnvBuilder{
		productName: productName,
		cliName:     cliName,
		nodes:       nodes,
		packageName: "devenv",
		outputDir:   "devenv",
		productType: productType,
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

// Write generates a complete boilerplate, can be multiple files
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

	// Generate cldf.go
	cldfContents, err := g.GenerateCLDF()
	if err != nil {
		return err
	}
	if err := os.WriteFile( //nolint:gosec
		filepath.Join(g.cfg.outputDir, "cldf.go"),
		[]byte(cldfContents),
		os.ModePerm,
	); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
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

	// Generate product_configuration.go
	switch g.cfg.productType {
	case "evm-single":
		prodConfigFileContents, err := g.GenerateSingleNetworkProductConfiguration()
		if err != nil {
			return err
		}
		if err := os.WriteFile( //nolint:gosec
			filepath.Join(g.cfg.outputDir, "product_configuration.go"),
			[]byte(prodConfigFileContents),
			os.ModePerm,
		); err != nil {
			return fmt.Errorf("failed to write product configuration file: %w", err)
		}
	case "multi-network":
		return fmt.Errorf("product configuration 'multi-network' is not supported yet")
	default:
		return fmt.Errorf("unknown product configuration type: %s, known types are 'evm-single' or 'multi-network'", g.cfg.productType)
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

	// tidy and finalize
	currentDir, err := os.Getwd()
	if err != nil {
		return err
	}

	// nolint
	defer os.Chdir(currentDir)
	if err := os.Chdir(g.cfg.outputDir); err != nil {
		return err
	}
	log.Info().Msg("Downloading dependencies and running 'go mod tidy' ..")
	_, err = exec.Command("go", "mod", "tidy").CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to tidy generated module: %w", err)
	}
	log.Info().
		Str("OutputDir", g.cfg.outputDir).
		Str("Module", g.cfg.moduleName).
		Msg("Developer environment generated")
	return nil
}

// GenerateSmokeTests generates a smoke test template
func (g *EnvCodegen) GenerateLoadTests() (string, error) {
	log.Info().Msg("Generating load test template")
	data := LoadTestParams{
		GoModName: g.cfg.moduleName,
	}
	return render(LoadTestTmpl, data)
}

// GenerateSmokeTests generates a smoke test template
func (g *EnvCodegen) GenerateSmokeTests() (string, error) {
	log.Info().Msg("Generating smoke test template")
	data := SmokeTestParams{
		GoModName: g.cfg.moduleName,
	}
	return render(SmokeTestTmpl, data)
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

// GenerateSingleNetworkProductConfiguration generate a single-network EVM product configuration
func (g *EnvCodegen) GenerateSingleNetworkProductConfiguration() (string, error) {
	log.Info().Msg("Configuring EVM network")
	p := ProductConfigurationSimple{
		PackageName: g.cfg.packageName,
	}
	return render(SingleNetworkEVMProductConfigurationTmpl, p)
}

// GenerateEnvironment generate environment.go, our environment composition function
func (g *EnvCodegen) GenerateEnvironment() (string, error) {
	log.Info().Msg("Generating environment composition (environment.go)")
	p := EnvParams{
		PackageName: g.cfg.packageName,
	}
	return render(EnvironmentTmpl, p)
}

// GenerateCLDF generate CLDF helpers
func (g *EnvCodegen) GenerateCLDF() (string, error) {
	log.Info().Msg("Generating CLDF helpers")
	p := CLDFParams{
		PackageName: g.cfg.packageName,
	}
	return render(CLDFTmpl, p)
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
