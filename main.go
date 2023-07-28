package main

import (
	"context"
	"fmt"
	"log"
	"os"

	examplecrdv1 "github.com/jasonhancock/examplecrd2/api/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func main() {
	scheme := runtime.NewScheme()
	examplecrdv1.AddToScheme(scheme)

	var config *rest.Config
	var err error
	if kubeconfig := os.Getenv("KUBECONFIG"); kubeconfig != "" {
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
	} else {
		// we're running in-cluster. Try using the service account
		config, err = rest.InClusterConfig()
	}

	if err != nil {
		log.Fatal(fmt.Errorf("getting kubeconfig: %w", err))
	}

	controllerClient, err := client.New(config, client.Options{Scheme: scheme})
	if err != nil {
		log.Fatal(err)
	}

	cl := NewClient(controllerClient)
	ctx := context.Background()

	resp, err := cl.List(ctx, "")
	if err != nil {
		log.Fatal(err)
	}

	for _, v := range resp.Items {
		fmt.Printf("%s %s %s\n", v.ObjectMeta.Name, v.Spec.Color, v.Spec.Size)
	}

	getResp, err := cl.Get(ctx, "", "testing")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s %s %s\n", getResp.ObjectMeta.Name, getResp.Spec.Color, getResp.Spec.Size)
}

type Client struct {
	client client.Client
}

func NewClient(client client.Client) *Client {
	return &Client{client: client}
}

func (c *Client) List(ctx context.Context, namespace string) (*examplecrdv1.ExampleResourceList, error) {
	var list examplecrdv1.ExampleResourceList
	err := c.client.List(ctx, &list, &client.ListOptions{Namespace: namespace})
	return &list, err
}

func (c *Client) Get(ctx context.Context, namespace, name string) (*examplecrdv1.ExampleResource, error) {
	var obj examplecrdv1.ExampleResource
	err := c.client.Get(ctx, types.NamespacedName{Namespace: namespace, Name: name}, &obj)
	return &obj, err
}
