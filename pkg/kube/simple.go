package kube

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (k *Kube) ListPods(ctx context.Context, namespace string) ([]corev1.Pod, error) {
	pods := k.c.CoreV1().Pods(namespace)
	list, err := pods.List(ctx, v1.ListOptions{})
	if err != nil {
		return nil, err
	}

	return list.Items, nil
}
