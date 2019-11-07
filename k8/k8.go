package k8

import (
	"github.com/galderz/golanzzo/util"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"os"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// Kubernetes abstracts interaction with a Kubernetes cluster
type Kubernetes struct {
	Client     client.Client
	restClient *rest.RESTClient
	RestConfig *rest.Config
}

// MapperProvider is a function that provides RESTMapper instances
type MapperProvider func(c *rest.Config) (meta.RESTMapper, error)

// NewKubernetesFromConfig creates a new Kubernetes instance from configuration.
// The configuration is resolved locally from known locations.
func NewKubernetesFromLocalConfig() *Kubernetes {
	scheme := runtime.NewScheme()
	mapperProvider := NewDynamicRESTMapper

	config := resolveConfig()
	config = setConfigDefaults(config)
	mapper, err := mapperProvider(config)
	util.ExpectNoError(err)
	kubernetes, err := client.New(config, createOptions(scheme, mapper))
	restClient, err := rest.RESTClientFor(config)
	util.ExpectNoError(err)

	return &Kubernetes{
		Client:     kubernetes,
		restClient: restClient,
		RestConfig: config,
	}
}

func createOptions(scheme *runtime.Scheme, mapper meta.RESTMapper) client.Options {
	return client.Options{
		Scheme: scheme,
		Mapper: mapper,
	}
}

func resolveConfig() *rest.Config {
	internal, _ := rest.InClusterConfig()
	if internal == nil {
		kubeConfig := findKubeConfig()
		clientConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
			&clientcmd.ClientConfigLoadingRules{ExplicitPath: kubeConfig},
			&clientcmd.ConfigOverrides{})
		external, _ := clientConfig.ClientConfig()
		return external
	}
	return internal
}

func setConfigDefaults(config *rest.Config) *rest.Config {
	gv := v1.SchemeGroupVersion
	config.GroupVersion = &gv
	config.APIPath = "/api"
	config.NegotiatedSerializer = serializer.DirectCodecFactory{CodecFactory: scheme.Codecs}
	config.UserAgent = rest.DefaultKubernetesUserAgent()
	return config
}

// FindKubeConfig returns local Kubernetes configuration
func findKubeConfig() string {
	kubeConfig := os.Getenv("KUBECONFIG")
	if kubeConfig != "" {
		return kubeConfig
	}
	home := os.Getenv("HOME")
	return home + "/.kube/config"
}
