package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strconv"

	"k8s.io/client-go/kubernetes/scheme"

	v1 "k8s.io/api/core/v1"
	coreV1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/remotecommand"

	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func main() {
	helper := NewIspnCliHelper()
	size, err := helper.GetClusterSize("local-operators", "example-infinispan-0")
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("Cluster size: %v\n", size)

	ls, err := helper.Ls("local-operators", "example-infinispan-0")
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("Ls: %v\n", ls)
}

// IspnCliHelper represent an helper for running CLI commands
type IspnCliHelper struct {
	coreClient *coreV1.CoreV1Client
	restConfig *rest.Config
	cliCmd     string
}

func getConfigLocation() string {
	var Home = os.Getenv("HOME")
	var ConfigLocation = Home + "/.kube/config"
	return ConfigLocation
}

var configLocation = getConfigLocation()

// NewIspnCliHelper create an IspnCliHelper
func NewIspnCliHelper() *IspnCliHelper {
	help := new(IspnCliHelper)
	help.restConfig, _ = rest.InClusterConfig()
	help.cliCmd = GetEnvWithDefault("CLI_CMD", "/opt/jboss/infinispan-server/bin/ispn-cli.sh")
	if help.restConfig == nil {
		clientConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
			&clientcmd.ClientConfigLoadingRules{ExplicitPath: configLocation},
			&clientcmd.ConfigOverrides{})
		help.restConfig, _ = clientConfig.ClientConfig()

	}
	help.coreClient, _ = coreV1.NewForConfig(help.restConfig)
	return help
}

type ClusterHealth struct {
	Nodes []string `json:"node_names"`
}

type Health struct {
	ClusterHealth  ClusterHealth `json:"cluster_health"`
}

// OldGetClusterSize get the cluster size via the ISPN cli
func (help *IspnCliHelper) GetClusterSize(namespace, namePod string) (int, error) {
	podIp := help.PodIp(namePod)
	httpUrl := fmt.Sprintf("http://%v:11222/rest/v2/cache-managers/DefaultCacheManager/health", podIp)
	commands := []string{"curl", "-u", "operator:LkzJ3OiMGKaldofl", httpUrl}
	var execIn, execOut, execErr bytes.Buffer
	//execIn.WriteString(cliCommand)
	err := help.ExecuteCmdOnPod(namespace, namePod, commands,
		&execIn, &execOut, &execErr)
	if err == nil {
		result := execOut.Bytes()

		var health Health
		err = json.Unmarshal(result, &health)
		if err != nil {
			panic(fmt.Errorf("unable to decode"))
		}

		return len(health.ClusterHealth.Nodes), nil
	}
	return -1, err
}

func (help *IspnCliHelper) PodIp(podName string) string {
	pod, err := help.coreClient.
		Pods("local-operators").
		Get(podName, metaV1.GetOptions{})
	if err != nil {
		panic(err.Error())
	}
	return pod.Status.PodIP
}

// OldGetClusterSize get the cluster size via the ISPN cli
func (help *IspnCliHelper) OldGetClusterSize(namespace, namePod string) (int, error) {
	cliCommand := "/subsystem=datagrid-infinispan/cache-container=clustered/:read-attribute(name=cluster-size)\n"
	commands := []string{help.cliCmd, "--connect"}
	var execIn, execOut, execErr bytes.Buffer
	execIn.WriteString(cliCommand)
	err := help.ExecuteCmdOnPod(namespace, namePod, commands,
		&execIn, &execOut, &execErr)
	if err == nil {
		result := execOut.String()
		// Match the correct line in the output
		resultRegExp := regexp.MustCompile("\"result\" => \"\\d+\"")
		// Match the result value
		valueRegExp := regexp.MustCompile("\\d+")
		resultLine := resultRegExp.FindString(result)
		resultValueStr := valueRegExp.FindString(resultLine)
		return strconv.Atoi(resultValueStr)
	}
	return -1, err
}

func (help *IspnCliHelper) Ls(namespace, namePod string) (string, error) {
	//cliCommand := "/subsystem=datagrid-infinispan/cache-container=clustered/:read-attribute(name=cluster-size)\n"
	//commands := []string{help.cliCmd, "--connect"}
	commands := []string{"ls"}
	var execIn, execOut, execErr bytes.Buffer
	//execIn.WriteString(cliCommand)
	err := help.ExecuteCmdOnPod(namespace, namePod, commands,
		&execIn, &execOut, &execErr)
	if err == nil {
		result := execOut.String()
		return result, nil
	}
	return "", err
}

// ExecuteCmdOnPod Excecutes command on pod
// commands array example { "/usr/bin/ls", "folderName" }
// execIn, execOut, execErr stdin, stdout, stderr stream for the command
func (help *IspnCliHelper) ExecuteCmdOnPod(namespace, podName string, commands []string,
	execIn, execOut, execErr *bytes.Buffer) error {
	// Create a POST request
	execRequest := help.coreClient.RESTClient().Post().
		Resource("pods").
		Name(podName).
		Namespace(namespace).
		SubResource("exec").
		VersionedParams(&v1.PodExecOptions{
			Container: "infinispan",
			Command:   commands,
			Stdin:     true,
			Stdout:    true,
			Stderr:    true,
			TTY:       false,
		}, scheme.ParameterCodec)
	// Create an executor
	exec, err := remotecommand.NewSPDYExecutor(help.restConfig, "POST", execRequest.URL())
	if err != nil {
		return err
	}
	// Run the command
	err = exec.Stream(remotecommand.StreamOptions{
		Stdin:  execIn,
		Stdout: execOut,
		Stderr: execErr,
		Tty:    false,
	})
	return err
}

// GetEnvWithDefault return GetEnv(name) if exists else
// return defVal
func GetEnvWithDefault(name, defVal string) string {
	str := os.Getenv(name)
	if str != "" {
		return str
	}
	return defVal
}
