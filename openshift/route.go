package main

import (
	"fmt"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func main() {
	fmt.Printf("Hello world!\n")

	url := "https://api.<...>"
	token := "<...>"
	restConfig := GetOpenShiftRESTConfig(url, token)
	openshift := NewOpenShiftFromConfig(restConfig)

	createNamespace("route-test", openshift)
	createRoute("route-test", openshift)
	printRoute("route-test", openshift)
}

func printRoute(s string, shift *OpenShift) {

}

func createRoute(s string, shift *OpenShift) {

}

func createNamespace(s string, shift *OpenShift) {

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

func GetOpenShiftRESTConfig(masterURL, token string) *rest.Config {
	restConfig, err := clientcmd.BuildConfigFromFlags(masterURL, "")
	if err != nil {
		panic(err.Error())
	}

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
