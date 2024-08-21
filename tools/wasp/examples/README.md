# Tests layout and local/cluster mode

## How it works
- You run your tests in plain `go test` locally
- You define cluster entrypoint in `zcluster` and use the same `go test` but with a cluster entrypoint, it builds and pushes images to ECR and run tests with `Helm`
- When each `job` got allocated your pod code will start to wait until all `pods` spun by client ( with unique `sync` id generated on start ) will have status `Running` if that won't happen until timeout client will remove all the pods
- When `N pods with unique test label are in status Running` the test will start, lag between pods start is 1sec max
- In case of any fatal error client will remove all the `jobs`
- When all `jobs` are complete test will end

## Prepare k8s namespace
- Create namespace `wasp` and apply default permissions
```
cd charts/wasp
kubectl create ns wasp
kubectl -n wasp apply -f setup.yaml
```

## Layout
In order to build the tests properly your test directory layout should follow simple rules:
- It should have `go.mod`
- It should be buildable with `go test -c ./...` from a dir where `go.mod` exists

Overall layout looks like:
```
tests
├── group1
│   ├── test_case_1_test.go
│   ├── gun.go
│   └── vu.go
├── group2
│   ├── subgroup21
│   │   ├── test_case_2_test.go
│   │   ├── gun.go
│   │   └── vu.go
└── go.mod
```

### Build

Default [Dockerfile](../DockerfileWasp) and [build_script](../build_test_image.sh) are working with AWS ECR private repos and building only for `amd64` platform

These vars are required to rebuild default layout and update an image
```
		UpdateImage:  true,
		BuildCtxPath: "..",
		
		// Helm values
		"image": ...
```
Turn off `UpdateImage` if you want to skip the build

By default, we are going one dir up from the cluster entrypoint script and building all:
```
		DockerCmdExecPath: "..",
		BuildCtxPath:      ".",
```
`BuildCtxPath` is relative to `DockerCmdExecPath`

If for some reason you don't like this layout or can't build like `go test -c ./...`, or you would like to customize your builds then you need to customize default [Dockerfile](../DockerfileWasp) and [build_script](../build_test_image.sh) and reference them in [cluster_entrypoint](zcluster/cluster_test.go)
```
		DockerfilePath: "",
		DockerIgnoreFilePath: "",
		BuildScriptPath: "",
```


### Deployment

Default helm chart is [here](../charts/wasp)

If no chart override is provided we are using this [tar](../charts/wasp/wasp-0.1.8.tgz)

You can set your custom chart as a local path or as an OCI registry URI in [cluster_entrypoint](zcluster/cluster_test.go)
- Set [wasp chart](../charts/wasp) in test params
```
ChartPath: "../../charts/wasp"
or 
ChartPath: "oci://public.ecr.aws/chainlink/wasp"
```

### Deployment values
- This Helm values are usually static or secret, use any configuration to set them, simplest is just as env vars
```
                        // secrets
			"env.loki.url":        os.Getenv("LOKI_URL"),
			"env.loki.token":      os.Getenv("LOKI_TOKEN"),
			"env.loki.basic_auth": os.Getenv("LOKI_BASIC_AUTH"),
			"env.loki.tenant_id":  os.Getenv("LOKI_TENANT_ID"),
			
			// WASP vars
			"env.wasp.log_level": "debug",

			
			// test vars
			"image":               os.Getenv("WASP_TEST_IMAGE"),
			"test.binaryName":     os.Getenv("WASP_TEST_BIN"),
			"test.name":           os.Getenv("WASP_TEST_NAME"),
			"test.timeout":       "24h",
			"jobs":               "40",
			
			// k8s pods resources
			"resources.requests.cpu":    "2000m",
			"resources.requests.memory": "512Mi",
			"resources.limits.cpu":      "2000m",
			"resources.limits.memory":   "512Mi",
			
			// your custom vars that will pass through to k8s jobs
			"test.MY_CUSTOM_VAR": "abc",
```
[Code](zcluster) examples

## Running tests

Use `go test` for local runs
```
go test -v -run TestNodeRPS
```

Use `go test` for remote runs (consume vars from any config method, example is env vars)
```
export WASP_TEST_TIMEOUT="12h"
export WASP_TEST_IMAGE="..."
export WASP_TEST_BIN="profiles.test"
export WASP_TEST_NAME="TestNodeRPS"
export WASP_UPDATE_IMAGE="true"
go test -v -run TestClusterEntrypoint
```