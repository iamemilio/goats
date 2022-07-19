package main

import (
	"context"
	"goats/pkg/generator"
	"goats/pkg/input"
	"log"

	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	inputs := input.New()
	err := inputs.Get()
	if err != nil {
		log.Fatal(err)
	}

	var runner generator.RunConfig
	err = runner.Parse(inputs.ConfigFile)
	if err != nil {
		log.Fatal(err)
	}

	for _, run := range runner.Runs {
		config, err := clientcmd.BuildConfigFromFlags("", run.Kubeconfig)
		if err != nil {
			log.Fatal(err)
		}

		client, err := kubernetes.NewForConfig(config)
		if err != nil {
			log.Fatal(err)
		}

		job := run.CreateJob()
		jobClient := client.BatchV1().Jobs(apiv1.NamespaceDefault)

		log.Println("Creating Kubernetes Job...")
		result, err := jobClient.Create(context.TODO(), &job, metav1.CreateOptions{})
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("Created Job %qs\n", result.GetObjectMeta().GetName())
	}
}
