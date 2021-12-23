package service

import (
	pb "github.com/ez-deploy/protobuf/project"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

// Service controll service.
type Service struct {
	k8sClientset *kubernetes.Clientset

	pb.UnimplementedOpsServer
}

// New in-cluster service.
func NewInClusterService() (*Service, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	res := &Service{
		k8sClientset: clientset,
	}

	return res, nil
}
