package main

import (
	"fmt"
	"github.com/galderz/golanzzo/k8"
	"github.com/galderz/golanzzo/util"
	"strconv"
	"strings"
)

func main() {
	podName := "example-infinispan-0"
	namespace := "local-operators"

	kubernetes := k8.NewKubernetesFromLocalConfig()

	printMaxUnboundedMemoryBytes(podName, namespace, kubernetes)
}

func printMaxUnboundedMemoryBytes(podName string, namespace string, kubernetes *k8.Kubernetes) {
	command := []string{"cat", "/proc/meminfo"}

	execOptions := k8.ExecOptions{Command: command, PodName: podName, Namespace: namespace}
	execOut, execErr, err := kubernetes.Execute(execOptions)
	if err == nil {
		result := execOut.String()
		lines := strings.Split(result, "\n")
		for _, line := range lines {
			if strings.Contains(line, "MemTotal:") {
				tokens := strings.Fields(line)
				maxUnbound, err := strconv.ParseUint(tokens[1], 10, 64)
				util.ExpectNoError(err)
				fmt.Printf("Max unbounded memory: %d", maxUnbound * 1024)
				return
			}
		}
	}
	err = fmt.Errorf("unexpected error getting max unbounded memory, stderr: %v, err: %v", execErr, err)
	util.ExpectNoError(err)
}
