package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	v1 "k8s.io/api/core/v1"
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

const namespace = "kube-system"

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
	
	checkCoreDNSAvailability(clientset)

	fmt.Println("=====================")

	dnsDaemonsetRes, dnsDaemonset := createDaemonSet(client)

	fmt.Println("=====================")

	fmt.Println("Waiting 5 seconds for information to propagate...")
	time.Sleep(5 * time.Second)

	checkPodLogs(clientset)
	fmt.Println("=====================")


	deleteDaemonSet(client, dnsDaemonsetRes, dnsDaemonset)



}

func checkCoreDNSAvailability(clientset *kubernetes.Clientset) {
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
}

func createDaemonSet(client dynamic.Interface) (schema.GroupVersionResource, *unstructured.Unstructured){
	dnsDaemonsetRes := schema.GroupVersionResource{Group: "apps", Version: "v1", Resource: "daemonsets"}
	dnsDaemonset := &unstructured.Unstructured{}

	dnsDaemonYaml, err := ioutil.ReadFile("yml/daemon-dns.yaml")
	check(err)

	dec := yaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme)
	_, _, err = dec.Decode(dnsDaemonYaml, nil, dnsDaemonset)
	check(err)

	//fmt.Println(dnsDaemonset.GetName(), gvk.String())

	fmt.Println("Creating daemonset...")
	result, err := client.Resource(dnsDaemonsetRes).Namespace(namespace).Create(context.TODO(), dnsDaemonset, metav1.CreateOptions{})
	check(err)

	fmt.Printf("Created daemonset %q.\n", result.GetName())

	return dnsDaemonsetRes, result
}

func deleteDaemonSet(client dynamic.Interface, dnsDaemonsetRes schema.GroupVersionResource, dnsDaemonset *unstructured.Unstructured) {
	deletePolicy := metav1.DeletePropagationForeground
	deleteOptions := metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	}
	fmt.Println("Deleting daemonset...")
	err := client.Resource(dnsDaemonsetRes).Namespace(dnsDaemonset.GetNamespace()).Delete(context.TODO(), dnsDaemonset.GetName(), deleteOptions)
	check(err)
	fmt.Printf("Deleted daemonset %q.\n", dnsDaemonset.GetName())
}

func checkPodLogs(clientset *kubernetes.Clientset) {
	//dnsDaemonSet, err := clientset.AppsV1().DaemonSets("kube-system").List(context.TODO(), )
	//check(err)

	fmt.Println("Checking pod logs...")

	dnsPods, err := clientset.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{LabelSelector: "name=dnsresolver-spread"})

	check(err)

	fmt.Printf("Number of DNS Pods: %d\n", len(dnsPods.Items))
	//fmt.Println(dnsPods.String())

	for i, pod := range dnsPods.Items {
		fmt.Printf("Pod #%d Name: %s\n", i, pod.Name)
		fmt.Printf("Logs:\n\t%s", getPodLogs(clientset, pod))
		fmt.Println("----------------")
	}
}

func getPodLogs(clientset *kubernetes.Clientset, pod v1.Pod) string {
	req := clientset.CoreV1().Pods(pod.Namespace).GetLogs(pod.Name, &v1.PodLogOptions{})

	podLogs, err := req.Stream(context.TODO())
	check(err)
	defer podLogs.Close()

	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, podLogs)
	check(err)

	return buf.String()
}