package main

import (
	"context"
	"fmt"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func main() {
	// operatorMasterURL := "https://192.168.99.140:8443"
	siteAMasterURL := "https://192.168.99.143:8443"
	siteBMasterURL := "https://192.168.99.144:8443"
	fmt.Printf("Pod list for SiteA at %s:\n", siteAMasterURL)
	podList(siteAMasterURL)
	fmt.Printf("Pod list for SiteB at %s:\n", siteBMasterURL)
	podList(siteBMasterURL)
}

func podList(masterURL string) {
	kubernetes := NewKubernetesFromMasterURL(masterURL)
	printPodList(kubernetes)
}

func printPodList(kubernetes *Kubernetes) {
	podList := &corev1.PodList{}
	listOps := &client.ListOptions{Namespace: "kube-system"}
	err := kubernetes.Client.List(context.TODO(), listOps, podList)
	if err != nil {
		panic(err.Error())
	}
	for _, pod := range podList.Items {
		fmt.Printf("Pod: %s\n", pod.Name)
	}
}

func NewKubernetesFromMasterURL(masterURL string) *Kubernetes {
	restConfig, err := clientcmd.BuildConfigFromFlags(masterURL, "")
	if err != nil {
		panic(err.Error())
	}
	restConfig.Insecure = true
	kubeClient, err := client.New(restConfig, client.Options{})
	if err != nil {
		panic(err.Error())
	}
	gv := corev1.SchemeGroupVersion
	restConfig.GroupVersion = &gv
	// restConfig.APIPath = "/api"
	restConfig.NegotiatedSerializer = serializer.DirectCodecFactory{CodecFactory: scheme.Codecs}
	// restConfig.UserAgent = rest.DefaultKubernetesUserAgent()
	restClient, err := rest.RESTClientFor(restConfig)
	if err != nil {
		panic(err.Error())
	}
	kubernetes := &Kubernetes{
		Client:     kubeClient,
		restClient: restClient,
		RestConfig: restConfig,
	}
	return kubernetes
}

//func addToScheme(schemeBuilder *runtime.SchemeBuilder, scheme *runtime.Scheme) {
//	err := schemeBuilder.AddToScheme(scheme)
//	ExpectNoError(err)
//}
//
//func createOptions() client.Options {
//	var clientScheme = runtime.NewScheme()
//	addToScheme(&corev1.SchemeBuilder, clientScheme)
//	return client.Options{
//		Scheme: clientScheme,
//	}
//}
//
//func ExpectNoError(err error) {
//	if err != nil {
//		panic(err.Error())
//	}
//}

// Kubernetes abstracts interaction with a Kubernetes cluster
type Kubernetes struct {
	Client     client.Client
	restClient *rest.RESTClient
	RestConfig *rest.Config
}
