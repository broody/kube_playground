package main

import (
	"flag"
	"fmt"
	"path/filepath"

	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/restmapper"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"k8s.io/metrics/pkg/client/custom_metrics"
)

func main() {
	// config, err := rest.InClusterConfig()
	// if err != nil {
	// 	panic(err.Error())
	// }
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	discoveryClient := discovery.NewDiscoveryClientForConfigOrDie(config)
	apiGroupResources, err := restmapper.GetAPIGroupResources(discoveryClient)
	if err != nil {
        panic(err.Error())
    }
	mapper := restmapper.NewDiscoveryRESTMapper(apiGroupResources)

	cmClient := custom_metrics.NewForConfig(config, mapper, custom_metrics.NewAvailableAPIsGetter(discoveryClient))

	metricName := "katana_jemalloc_active"
	namespace := "external-nobody"
	podName := "katana-74f5fff546-gqgfp"
	gk := schema.GroupKind{Group: "", Kind: "Pod"}
	
	metricValue, err := cmClient.NamespacedMetrics(namespace).GetForObject(gk, podName, metricName, labels.Everything())
	if err != nil {
		fmt.Println("shit")
	}

	// // Process the retrieved metrics
	// for _, metricValue := range metrics.Items {
	// 	// Here you can access each metric's value, timestamp, etc.
	// 	fmt.Println(metricValue)
	// }

	fmt.Println(metricValue)

}