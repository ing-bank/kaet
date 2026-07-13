package kubernetes

import (
	"context"
	"time"

	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	watch2 "k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/watch"
)

func isPodReady(event watch2.Event) (bool, error) {
	pod, ok := event.Object.(*corev1.Pod)
	if !ok {
		return false, nil
	}

	for _, cond := range pod.Status.Conditions {
		if cond.Type == corev1.PodReady {
			return cond.Status == corev1.ConditionTrue, nil
		}
	}
	return false, nil
}

func (c *Client) WaitForPodReady(parentCtx context.Context, ns, name string,
	timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(parentCtx, timeout)
	defer cancel()

	lw := &cache.ListWatch{
		ListWithContextFunc: func(ctx context.Context, options v1.ListOptions) (runtime.Object, error) {
			options.FieldSelector = "metadata.name=" + name
			return c.ClientSet.CoreV1().
				Pods(ns).
				List(ctx, options)
		},
		WatchFuncWithContext: func(ctx context.Context, options v1.ListOptions) (watch2.Interface, error) {
			options.FieldSelector = "metadata.name=" + name
			return c.ClientSet.CoreV1().Pods(ns).
				Watch(ctx, options)
		},
	}

	_, err := watch.UntilWithSync(
		ctx,
		lw,
		&corev1.Pod{},
		nil,
		isPodReady,
	)

	return err
}
