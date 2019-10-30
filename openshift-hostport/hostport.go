package main

import (
	"context"
	"fmt"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func main() {
	openshiftViaSessionToken()
	openshiftViaServiceAccountToken()

	// Via username/password does not work
	// openshiftViaUsernamePassword()
}

func openshiftViaUsernamePassword() {
	url := "https://api.<...>"
	username := "<...>"
	password := "<...>"
	restConfig := GetOpenShiftUsernamePasswordRESTConfig(url, username, password)
	fmt.Printf("Rest config %+v", restConfig)
	openshift := NewOpenShiftFromConfig(restConfig)
	printPodList(openshift)
	printService(openshift)
}

func openshiftViaServiceAccountToken() {
	url := "https://api.<...>"
	token := "<...>"
	restConfig := GetOpenShiftTokenRESTConfig(url, token)
	openshift := NewOpenShiftFromConfig(restConfig)
	//printPodList(openshift)
	printService(openshift)
}

func openshiftViaSessionToken() {
	url := "https://api.<...>"
	token := "<...>"
	restConfig := GetOpenShiftTokenRESTConfig(url, token)
	openshift := NewOpenShiftFromConfig(restConfig)
	printPodList(openshift)
	printService(openshift)
}

func printService(openShift *OpenShift) {
	fmt.Printf("Service list:\n")
	namespacedName := types.NamespacedName{Name: "example-infinispan-site", Namespace: "local-operators"}
	service := &corev1.Service{}
	err := openShift.Client.Get(context.TODO(), namespacedName, service)
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("Service: %v\n", service)
}

func printPodList(openShift *OpenShift) {
	fmt.Printf("Pod list:\n")
	podList := &corev1.PodList{}
	listOps := &client.ListOptions{Namespace: "kube-system"}
	err := openShift.Client.List(context.TODO(), listOps, podList)
	if err != nil {
		panic(err.Error())
	}
	for _, pod := range podList.Items {
		fmt.Printf("Pod: %s\n", pod.Name)
	}
	fmt.Printf("End pod list.\n")
}

func NewOpenShiftFromConfig(config *rest.Config) *OpenShift {
	kubeClient, err := client.New(config, client.Options{})
	if err != nil {
		panic(err.Error())
	}
	config = setConfigDefaults(config)
	restClient, err := rest.RESTClientFor(config)
	if err != nil {
		panic(err.Error())
	}
	openShift := &OpenShift{
		Client:     kubeClient,
		restClient: restClient,
		RestConfig: config,
	}
	return openShift
}

func GetOpenShiftUsernamePasswordRESTConfig(masterURL, username, password string) *rest.Config {
	restConfig, err := clientcmd.BuildConfigFromFlags(masterURL, "")
	if err != nil {
		panic(err.Error())
	}

	restConfig.Username = username
	restConfig.Password = password
	//restConfig.Impersonate.UserName = username
	restConfig.Insecure = true
	return restConfig
}

func GetOpenShiftTokenRESTConfig(masterURL, token string) *rest.Config {
	restConfig, err := clientcmd.BuildConfigFromFlags(masterURL, "")
	if err != nil {
		panic(err.Error())
	}

	restConfig.Insecure = true
	restConfig.BearerToken = token
	return restConfig
}

func setConfigDefaults(config *rest.Config) *rest.Config {
	gv := v1.SchemeGroupVersion
	config.GroupVersion = &gv
	config.APIPath = "/api"
	config.NegotiatedSerializer = serializer.DirectCodecFactory{CodecFactory: scheme.Codecs}
	config.UserAgent = rest.DefaultKubernetesUserAgent()
	return config
}

type OpenShift struct{
	Client     client.Client
	restClient *rest.RESTClient
	RestConfig *rest.Config
}
