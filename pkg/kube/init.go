package kube

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type Kube struct {
	c *kubernetes.Clientset
}

func New() (*Kube, error) {
	cfg, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}

	cli, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		return nil, err
	}

	return &Kube{
		c: cli,
	}, nil
}
