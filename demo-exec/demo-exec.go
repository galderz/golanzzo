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

	cacheName := "default"
	defaultCacheXml := `<infinispan><cache-container>
        <distributed-cache name="default" mode="SYNC" owners="1">
            <memory>
                <off-heap 
                    size="96468992"
                    eviction="MEMORY"
                    strategy="REMOVE" 
                />
            </memory>
            <partition-handling when-split="DENY_READ_WRITES" merge-policy="REMOVE_ALL" />
        </distributed-cache>
    </cache-container></infinispan>`
	podIp := "172.17.0.2"
	password := "TtenNOI33MZDOgdq"

	err := CreateCache(cacheName, defaultCacheXml, podName, namespace, podIp, password, kubernetes)
	util.ExpectNoError(err)
	//printMaxUnboundedMemoryBytes(podName, namespace, kubernetes)
}

func CreateCache(cacheName, cacheXml, podName, namespace, podIP, password string, kubernetes *k8.Kubernetes) error {
	httpURL := fmt.Sprintf("http://%s:11222/rest/v2/caches/%s", podIP, cacheName)
	commands := []string{"curl",
		//"-s", "-o", "/dev/null", "-w", "%{http_code}", // get std out to only have code
		"-w", "\n%{http_code}", // get std out to only have code
		"-d", fmt.Sprintf("%s", cacheXml),
		"-H", "Content-Type: application/xml",
		"-u", fmt.Sprintf("operator:%s", password),
		"-X", "POST",
		httpURL,
	}

	//logger := log.WithValues("Request.Namespace", namespace, "Secret.Name", secretName, "Pod.Name", podName)
	//logger.Info("create cache", "url", httpURL, "cache name", cacheName, "cache configuration", cacheXml)

	execOptions := k8.ExecOptions{Command: commands, PodName: podName, Namespace: namespace}
	execOut, execErr, err := kubernetes.Execute(execOptions)
	if err != nil {
		return fmt.Errorf("execute error creating cache, stderr: %v, err: %v", execErr, err)
	}

	// Reverse sort the standard output so that HTTP status code is first
	execOutLines := strings.Split(execOut.String(), "\n")
	reverse(execOutLines)
	//fmt.Println(execOutLines)

	httpCode, err := strconv.ParseUint(execOutLines[0], 10, 64)
	if err != nil {
		return err
	}

	if httpCode > 299 || httpCode < 200 {
		return fmt.Errorf("server side error creating cache: %s", execOut.String())
	}

	fmt.Printf("Create cache completed successfully. Std out: %s\n", execOut.String())
	return nil
}

func reverse(ss []string) {
	last := len(ss) - 1
	for i := 0; i < len(ss)/2; i++ {
		ss[i], ss[last-i] = ss[last-i], ss[i]
	}
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
	err = fmt.Errorf("unexpected error getting max unbounded memory, stderr: %s, err: %v", execErr, err)
	util.ExpectNoError(err)
}
