package main

import (
	"fmt"
	"net/url"
	"os"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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

	fmt.Println(PublicIp(config))

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	fmt.Printf("There are %d pods in the cluster\n", NumberOfPods(clientset))

	fmt.Printf("Named service: %s\n", NamedService(clientset))

	fmt.Printf("Selectors for service: %s\n", SelectoForService(clientset))

	fmt.Printf("Pod ip: %s\n", PodIp("example-infinispan-0", clientset))
}

func SelectoForService(clientset *kubernetes.Clientset) map[string]string {
	svc, err := clientset.
		CoreV1().
		Services("local-operators").
		Get("example-infinispan", metav1.GetOptions{})
	if err != nil {
		panic(err.Error())
	}
	return svc.Spec.Selector
}

func NamedService(clientset *kubernetes.Clientset) string {
	svc, err := clientset.
		CoreV1().
		Services("local-operators").
		Get("example-infinispan", metav1.GetOptions{})
	if err != nil {
		panic(err.Error())
	}
	return svc.Name
}

func NumberOfPods(clientset *kubernetes.Clientset) int {
	pods, err := clientset.CoreV1().Pods("").List(metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	return len(pods.Items)
}

func PublicIp(c *rest.Config) string {
	u, err := url.Parse(c.Host)
	if err != nil {
		panic(err.Error())
	}
	return u.Hostname()
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
