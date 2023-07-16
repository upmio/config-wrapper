package conf

import (
	"fmt"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

var clientset kubernetes.Interface

func (k *Kube) GetClientSet() (kubernetes.Interface, error) {
	k.lock.Lock()
	defer k.lock.Unlock()
	if clientset == nil {
		conn, err := k.getClientSet()
		if err != nil {
			return nil, err
		}
		clientset = conn
	}
	return clientset, nil
}

func (k *Kube) getClientSet() (kubernetes.Interface, error) {
	var (
		err    error
		config *rest.Config
	)

	// creates the config
	if len(k.KubeConfig) == 0 {
		if config, err = rest.InClusterConfig(); err != nil {
			return nil, fmt.Errorf("create in-cluster config fail, error: %v", err)
		}
	} else {
		if config, err = clientcmd.BuildConfigFromFlags("", k.KubeConfig); err != nil {
			return nil, fmt.Errorf("create out-of-cluster config fail, error: %v", err)
		}
	}

	// creates the clientset
	clientset, err = kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("create clientset fail, error: %v", err)
	}

	return clientset, nil
}
