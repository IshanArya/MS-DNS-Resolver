package main

import (
	"context"
	"fmt"
	"io/ioutil"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer/yaml"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"log"
	"time"
)

func check(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func main() {
	config, err := rest.InClusterConfig()
	check(err)

	time.Sleep(1 * time.Second)

	clientset, err := kubernetes.NewForConfig(config)
	check(err)

	client, err := dynamic.NewForConfig(config)
	check(err)

	time.Sleep(1 * time.Second)

	//for {

	dnsDeployments, err := clientset.AppsV1().Deployments("").List(context.TODO(), metav1.ListOptions{LabelSelector: `k8s-app=kube-dns`})
	check(err)
	fmt.Printf("There are %d coredns deployments in the cluster\n", len(dnsDeployments.Items))

	for i, deployment := range dnsDeployments.Items {
		fmt.Printf("DNS Deployment %d: \n\tReplicas - %d/%d\n\tCondition - %v\n\n", i,
			deployment.Status.ReadyReplicas,
			deployment.Status.Replicas,
			deployment.Status.Conditions)
		fmt.Printf("Complete Status: %s\n", deployment.Status.String())
		fmt.Println("----------------")
	}
	fmt.Println("=====================")
	//}

	testDeploymentRes := schema.GroupVersionResource{Group: "apps", Version: "v1", Resource: "deployments"}
	testDeployment := &unstructured.Unstructured{}

	emptyLinuxYaml, err := ioutil.ReadFile("yml/empty-linux.yaml")
	check(err)

	//fmt.Printf("File Contents: %s\n", emptyLinuxYaml)

	dec := yaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme)
	_, gvk, err := dec.Decode(emptyLinuxYaml, nil, testDeployment)

	fmt.Println(testDeployment.GetName(), gvk.String())

	fmt.Println("=====================")

	fmt.Println("Creating deployment...")
	result, err := client.Resource(testDeploymentRes).Namespace("kube-system").Create(context.TODO(), testDeployment, metav1.CreateOptions{})
	check(err)

	fmt.Printf("Created deployment %q.\n", result.GetName())


	//for{
	//	time.Sleep(10 * time.Second)
	//}
}