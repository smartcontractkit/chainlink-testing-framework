package pods

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var (
	DefaultDstDir = filepath.Join("/", "go", "local")
)

type TestRunner struct {
	namespace string
	srcDir    string
	podLabel  string
	dstDir    string
	container string
	podName   string
}

// NewTestRunner creates a new TestRunner instance
func NewTestRunner(namespace, srcDir string) (*TestRunner, error) {
	t := &TestRunner{
		namespace: namespace,
		srcDir:    srcDir,
		podLabel:  "app=ubuntu",
		dstDir:    DefaultDstDir,
		container: "ubuntu-container",
	}

	// Get pod name
	getPodCmd := exec.Command("kubectl", "get", "pod", "-n", t.namespace, "-l", t.podLabel, "-o", "jsonpath={.items[0].metadata.name}") //nolint
	podNameBytes, err := getPodCmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get pod name: %v", err)
	}
	t.podName = string(podNameBytes)
	if t.podName == "" {
		L.Error().Str("namespace", t.namespace).Str("label", t.podLabel).Msg("No pod found with specified label")
		return nil, fmt.Errorf("no pod found with label %s in namespace %s", t.podLabel, t.namespace)
	}

	if err = t.Run(srcDir, "go", "install", "github.com/go-delve/delve/cmd/dlv@latest"); err != nil {
		return nil, err
	}

	L.Info().Str("pod", t.podName).Str("namespace", t.namespace).Msg("Found target pod")
	return t, t.copyTests()
}

// CopyTests copies local directory to the pod
func (t *TestRunner) copyTests() error {
	if _, err := os.Stat(t.srcDir); os.IsNotExist(err) {
		L.Error().Str("directory", t.srcDir).Msg("Local directory does not exist")
		return fmt.Errorf("local directory %s does not exist", t.srcDir)
	}

	dest := fmt.Sprintf("%s:%s", t.podName, t.dstDir)
	args := []string{"cp", t.srcDir, dest, "-n", t.namespace, "--container", t.container}
	L.Debug().Str("command", "kubectl").Strs("args", args).Msg("Copy arguments")

	cmd := exec.Command("kubectl", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		L.Error().Err(err).Str("source", t.srcDir).Str("destination", dest).Msg("Copy failed")
		return fmt.Errorf("copy failed: %v", err)
	}

	L.Info().Str("source", t.srcDir).Str("destination", dest).Msg("Copy completed successfully")
	return nil
}

func TestPodConfig() *PodConfig {
	config := &PodConfig{
		Name:        S("ubuntu"),
		Image:       S("golang:1.24.3-alpine"),
		Limits:      ResourcesMedium(),
		Ports:       []string{"40000:40000"},
		Command:     S("sleep 999999999"),
		StatefulSet: true,
	}
	L.Debug().Interface("config", config).Msg("Created pod configuration")
	return config
}

// Run executes a command in the container at specified working directory
func (t *TestRunner) Run(fromDir string, command string, args ...string) error {
	if t.podName == "" {
		L.Error().Msg("Pod name not initialized - call Initialize() first")
		return fmt.Errorf("pod name not initialized - call Initialize() first")
	}

	targetDir := filepath.Join(DefaultDstDir, fromDir)
	execArgs := []string{
		"exec", t.podName,
		"-n", t.namespace,
		"--container", t.container,
		"--",
		"sh", "-ClientSet",
		fmt.Sprintf("cd %s && %s %s", targetDir, command, strings.Join(args, " ")),
	}

	L.Info().
		Str("pod", t.podName).
		Str("command", command).
		Strs("args", args).
		Str("directory", targetDir).
		Msg("Executing command in pod")

	cmd := exec.Command("kubectl", execArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		L.Error().
			Err(err).
			Str("pod", t.podName).
			Str("command", command).
			Msg("Command execution failed")
		return fmt.Errorf("command execution failed: %v", err)
	}

	return nil
}
