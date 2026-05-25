package framework

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/rs/zerolog/log"
)

/* Templates */

const (
	ProductFakesJustfile = `IMAGE_NAME := "{{ .ProductName }}-fakes"

run:
    docker run --rm -it -v $(pwd):/app -p 9111:9111 {{ "{{" }}IMAGE_NAME{{ "}}" }}:latest

build:
    docker build -f Dockerfile -t {{ "{{" }}IMAGE_NAME{{ "}}" }}:latest .

push registry:
    docker build --platform linux/amd64 -f Dockerfile -t {{ "{{" }}IMAGE_NAME{{ "}}" }}:latest .
    aws ecr get-login-password --region us-west-2 | docker login --username AWS --password-stdin {{ "{{" }}registry{{ "}}" }}
    docker tag {{ "{{" }}IMAGE_NAME{{ "}}" }}:latest {{ "{{" }}registry{{ "}}" }}/{{ "{{" }}IMAGE_NAME{{ "}}" }}
    docker push {{ "{{" }}registry{{ "}}" }}/{{ "{{" }}IMAGE_NAME{{ "}}" }}

clean:
    docker rmi {{ "{{" }}IMAGE_NAME{{ "}}" }}:latest
`
	ProductFakesImplTmpl = `package main

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/gin-gonic/gin"

	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/fake"
)

var L = log.Output(zerolog.ConsoleWriter{Out: os.Stderr}).Level(zerolog.DebugLevel).With().Fields(map[string]any{"component": "ocr2"}).Logger()

// example mock service
func main() {
	_, err := fake.NewFakeDataProvider(&fake.Input{Port: fake.DefaultFakeServicePort})
	if err != nil {
		panic(err)
	}
	err = fake.Func("POST", "/example_fake", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"data": map[string]any{
				"result": "ok",
			},
		})
	})
	if err != nil {
		panic(err)
	}
	select {}
}
`
	ProductFakesGoModuleTmpl = `module github.com/smartcontractkit/{{ .ProductName}}/devenv/fakes

go {{.RuntimeVersion}}

require (
	github.com/gin-gonic/gin v1.10.1
	github.com/rs/zerolog v1.34.0
	github.com/smartcontractkit/chainlink-testing-framework/framework v0.10.1
	github.com/smartcontractkit/chainlink-testing-framework/framework/components/fake v0.10.1-0.20250711120409-5078050f9db4
)`

	ProductFakesDockerfileTmpl = `FROM golang:1.25 AS builder

ENV GOPRIVATE=github.com/smartcontractkit/*

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY ../.. .
RUN CGO_ENABLED=0 GOOS=linux go build -o /fake main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=builder /fake /fake
EXPOSE 9111
CMD ["/fake"]
`
)

type ProductFakesImplParams struct {
	ProductName string
}

func (g *EnvCodegen) GenerateFakesImpl() (string, error) {
	log.Info().Msg("Generating fakes implementation")
	p := ProductFakesImplParams{}
	return render(ProductFakesImplTmpl, p)
}

type ProductFakesJustfileParams struct {
	ProductName string
}

func (g *EnvCodegen) GenerateFakesJustfile() (string, error) {
	log.Info().Msg("Generating fakes Justfile")
	p := ProductFakesJustfileParams{
		ProductName: g.cfg.productName,
	}
	return render(ProductFakesJustfile, p)
}

type ProductFakesDockerfileParams struct{}

func (g *EnvCodegen) GenerateFakesDockerfile() (string, error) {
	log.Info().Msg("Generating fakes Dockerfile")
	p := ProductFakesDockerfileParams{}
	return render(ProductFakesDockerfileTmpl, p)
}

type ProductFakesGoModuleParams struct {
	ProductName    string
	RuntimeVersion string
}

func (g *EnvCodegen) GenerateFakesGoModule() (string, error) {
	log.Info().Msg("Generating fakes go.mod")
	p := ProductFakesGoModuleParams{
		ProductName:    g.cfg.productName,
		RuntimeVersion: strings.ReplaceAll(runtime.Version(), "go", ""),
	}
	return render(ProductFakesGoModuleTmpl, p)
}

// WriteFakes writes all files related to fake services used in testing
func (g *EnvCodegen) WriteFakes() error {
	fakesRoot := filepath.Join(g.cfg.outputDir, "fakes")
	if err := os.MkdirAll( //nolint:gosec
		fakesRoot,
		os.ModePerm,
	); err != nil {
		return fmt.Errorf("failed to create fakes directory: %w", err)
	}

	// generate Dockerfile
	dockerfileContents, err := g.GenerateFakesDockerfile()
	if err != nil {
		return err
	}
	if err := os.WriteFile( //nolint:gosec
		filepath.Join(fakesRoot, "Dockerfile"),
		[]byte(dockerfileContents),
		os.ModePerm,
	); err != nil {
		return fmt.Errorf("failed to write fakes Dockerfile file: %w", err)
	}

	// generate fakes go.mod
	fakesGoModContents, err := g.GenerateFakesGoModule()
	if err != nil {
		return err
	}
	if err := os.WriteFile( //nolint:gosec
		filepath.Join(fakesRoot, "go.mod"),
		[]byte(fakesGoModContents),
		os.ModePerm,
	); err != nil {
		return fmt.Errorf("failed to write fakes Go module file: %w", err)
	}

	// generate fakes implementation
	implContents, err := g.GenerateFakesImpl()
	if err != nil {
		return err
	}
	if err := os.WriteFile( //nolint:gosec
		filepath.Join(fakesRoot, "main.go"),
		[]byte(implContents),
		os.ModePerm,
	); err != nil {
		return fmt.Errorf("failed to write fakes implementation file: %w", err)
	}

	// generate fakes Justfile
	justfileContents, err := g.GenerateFakesJustfile()
	if err != nil {
		return err
	}
	if err := os.WriteFile( //nolint:gosec
		filepath.Join(fakesRoot, "Justfile"),
		[]byte(justfileContents),
		os.ModePerm,
	); err != nil {
		return fmt.Errorf("failed to write fakes Just file: %w", err)
	}

	// tidy and finalize
	currentDir, err := os.Getwd()
	if err != nil {
		return err
	}

	// nolint
	defer os.Chdir(currentDir)
	if err := os.Chdir(fakesRoot); err != nil {
		return err
	}
	log.Info().Msg("Downloading dependencies and running 'go mod tidy' (fakes) ..")
	_, err = exec.Command("go", "mod", "tidy").CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to tidy generated module for fakes: %w", err)
	}

	log.Info().Msg("Building fakes image")
	_, err = exec.Command("just", "build").CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to build fakes Docker image: %w", err)
	}
	return nil
}
