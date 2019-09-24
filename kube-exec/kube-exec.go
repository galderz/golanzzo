package main

import (
	"bytes"
	"fmt"
	v1 "k8s.io/api/core/v1"
	"os"
	"path/filepath"
	"regexp"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/remotecommand"
)

func main() {
	fmt.Println("Hello, 世界")

	var Home = os.Getenv("HOME")
	var ConfigLocation = Home + "/.kube/config"

	clientConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		&clientcmd.ClientConfigLoadingRules{ExplicitPath: ConfigLocation},
		&clientcmd.ConfigOverrides{})

	config, err := clientConfig.ClientConfig()
	if err != nil {
		panic(err.Error())
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	//fmt.Printf("Cluster members: %s\n", GetClusterMembers("example-infinispan-0", "operator", "Gk2qKjx8ekoxDHCk", clientset))
	members, err := GetClusterMembers("local-operators", "example-infinispan-0", config, clientset)
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("Cluster members: %v\n", members)
}

func GetClientConfig() (*rest.Config, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		//fmt.Printf("Unable to create config. Error: %+v\n", err)

		err1 := err
		kubeconfig := filepath.Join(os.Getenv("HOME"), ".kube", "config")
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			err = fmt.Errorf("InClusterConfig as well as BuildConfigFromFlags Failed. Error in InClusterConfig: %+v\nError in BuildConfigFromFlags: %+v", err1, err)
			return nil, err
		}
	}

	return config, nil
}

func PodIp(podName string, clientset *kubernetes.Clientset) string {
	pod, err := clientset.
		CoreV1().
		Pods("local-operators").
		Get(podName, metav1.GetOptions{})
	if err != nil {
		panic(err.Error())
	}
	return pod.Status.PodIP
}

type ClusterHealth struct {
	Nodes []string `json:"node_names"`
}

type Health struct {
	ClusterHealth ClusterHealth `json:"cluster_health"`
}

// GetClusterMembers get the cluster members via the ISPN cli
func GetClusterMembers(namespace, namePod string, clientcfg *rest.Config, clientset *kubernetes.Clientset) (string, error) {
	//cliCommand := "/subsystem=datagrid-infinispan/cache-container=clustered/:read-attribute(name=members)\n"
	commands := []string{"ls"}
	var execIn, execOut, execErr bytes.Buffer
	//execIn.WriteString(cliCommand)
	err := ExecuteCmdOnPod(namespace, namePod, commands,
		&execIn, &execOut, &execErr, clientcfg, clientset)
	if err == nil {
		result := execOut.String()
		// Match the correct line in the output
		resultRegExp := regexp.MustCompile("\"result\" => \"\\[.*\\]\"")
		// Match the result value
		valueRegExp := regexp.MustCompile("\\[.*\\]")
		resultLine := resultRegExp.FindString(result)
		resultValue := valueRegExp.FindString(resultLine)
		return resultValue, nil
	}
	return "-error-", err
}

// ExecuteCmdOnPod Excecutes command on pod
// commands array example { "/usr/bin/ls", "folderName" }
// execIn, execOut, execErr stdin, stdout, stderr stream for the command
func ExecuteCmdOnPod(namespace, podName string, commands []string,
	execIn, execOut, execErr *bytes.Buffer, clientcfg *rest.Config, clientset *kubernetes.Clientset) error {
	// Create a POST request
	execRequest := clientset.RESTClient().Post().
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
	fmt.Println("Request URL:", execRequest.URL())
	exec, err := remotecommand.NewSPDYExecutor(clientcfg, "POST", execRequest.URL())
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
	msg := fmt.Errorf("error in Stream: %v", err)
	fmt.Printf("Error running command: %v\n", msg)
	fmt.Printf("Standard in: %v\n", execIn.String())
	fmt.Printf("Standard error: %v\n", execErr.String())
	fmt.Printf("Standard out: %v\n", execOut.String())
	return err
}

//func GetClusterMembers2(podName string, usr string, pass string, clientset *kubernetes.Clientset) []string {
//	podIp := PodIp(podName, clientset)
//
//	httpUrl := "http://" + podIp + ":11222/rest/v2/cache-managers/DefaultCacheManager/health"
//	commands := []string{"curl", "-u", usr + ":" + pass, httpUrl}
//	var execIn, execOut, execErr bytes.Buffer
//	fmt.Printf("Commands: %v", commands)
//	err := ExecuteCmdOnPod(podName, commands, &execIn, &execOut, &execErr, clientset)
//
//	if err == nil {
//		result := execOut.Bytes()
//
//		var health Health
//		err = json.Unmarshal(result, &health)
//		if err != nil {
//			panic(fmt.Errorf("unable to decode"))
//		}
//
//		return health.ClusterHealth.Nodes
//	}
//
//	panic(err.Error())
//}
//
//// ExecuteCmdOnPod Excecutes command on pod
//// commands array example { "/usr/bin/ls", "folderName" }
//// execIn, execOut, execErr stdin, stdout, stderr stream for the command




















//func ExecuteCmdOnPod(podName string, commands []string, execIn, execOut, execErr *bytes.Buffer, clientset *kubernetes.Clientset) error {
//	// Create a POST request
//	execRequest := clientset.RESTClient().Post().
//		Resource("pods").
//		Name(podName).
//		Namespace("local-operators").
//		SubResource("exec").
//		VersionedParams(&v1.PodExecOptions{
//			Container: "infinispan",
//			Command:   commands,
//			Stdin:     true,
//			Stdout:    true,
//			Stderr:    true,
//			TTY:       false,
//		}, scheme.ParameterCodec)
//	// Create an executor
//	restConfig, err := GetClientConfig()
//	if err != nil {
//		return err
//	}
//
//	exec, err := remotecommand.NewSPDYExecutor(restConfig, "POST", execRequest.URL())
//	if err != nil {
//		return err
//	}
//	// Run the command
//	err = exec.Stream(remotecommand.StreamOptions{
//		Stdin:  nil,
//		Stdout: execOut,
//		Stderr: execErr,
//		Tty:    false,
//	})
//	fmt.Println(err)
//	return err
//}
