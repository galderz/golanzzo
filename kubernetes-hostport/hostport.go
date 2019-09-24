package main

import (
	"context"
	"fmt"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"os"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func main() {
	// operatorMasterURL := "https://192.168.99.140:8443"
	siteAMasterURL := "https://192.168.99.141:8443"
	siteBMasterURL := "https://192.168.99.142:8443"
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
	kubeconfig := FindKubeConfig()
	restConfig, err := clientcmd.BuildConfigFromFlags(masterURL, kubeconfig)
	if err != nil {
		panic(err.Error())
	}
	kubeClient, err := client.New(restConfig, client.Options{})
	if err != nil {
		panic(err.Error())
	}
	gv := v1.SchemeGroupVersion
	restConfig.GroupVersion = &gv
	restConfig.NegotiatedSerializer = serializer.DirectCodecFactory{CodecFactory: scheme.Codecs}
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

func resolveConfig() *rest.Config {
	internal, _ := rest.InClusterConfig()
	if internal == nil {
		kubeConfig := FindKubeConfig()
		clientConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
			&clientcmd.ClientConfigLoadingRules{ExplicitPath: kubeConfig},
			&clientcmd.ConfigOverrides{})
		external, _ := clientConfig.ClientConfig()
		return external
	}
	return internal
}

// FindKubeConfig returns local Kubernetes configuration
func FindKubeConfig() string {
	kubeConfig := os.Getenv("KUBECONFIG")
	if kubeConfig != "" {
		return kubeConfig
	}
	return "../../openshift.local.clusterup/kube-apiserver/admin.kubeconfig"
}

func setConfigDefaults(config *rest.Config) *rest.Config {
	//gv := v1.SchemeGroupVersion
	//config.GroupVersion = &gv
	config.APIPath = "/api"
//	config.NegotiatedSerializer = serializer.DirectCodecFactory{CodecFactory: scheme.Codecs}
	config.UserAgent = rest.DefaultKubernetesUserAgent()
	return config
}

// Kubernetes abstracts interaction with a Kubernetes cluster
type Kubernetes struct {
	Client     client.Client
	restClient *rest.RESTClient
	RestConfig *rest.Config
}
