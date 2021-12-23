package service

import (
	"context"

	pb "github.com/ez-deploy/protobuf/project"
	apimetav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	corev1 "k8s.io/client-go/applyconfigurations/core/v1"
	metav1 "k8s.io/client-go/applyconfigurations/meta/v1"
)

var (
	apiVersionV1      = "v1"
	restKindService   = "Service"
	headlessClusterIP = "None"
)

var typeMetaApplyConfigureation_Service = metav1.TypeMetaApplyConfiguration{
	APIVersion: &apiVersionV1,
	Kind:       &restKindService,
}

// create a Headless k8s service from given pbSVC.
func (s *Service) applyServiceFromCfg(ctx context.Context, pbSVC *pb.SetServiceReq) error {
	svcCfg := generateSVCApplyConfigFromPBSVC(pbSVC)
	applyOptions := apimetav1.ApplyOptions{}

	_, err := s.k8sClientset.CoreV1().Services(pbSVC.ProjectName).Apply(ctx, svcCfg, applyOptions)

	return err
}

// delete a Headless k8s service from given pbSVC.
func (s *Service) deleteServiceFromCfg(ctx context.Context, pbSVC *pb.DeleteServiceReq) error {
	deleteOptions := apimetav1.DeleteOptions{}

	return s.k8sClientset.CoreV1().Services(pbSVC.ProjectName).Delete(ctx, pbSVC.ServiceName, deleteOptions)
}

// generate service kind k8s object from given pbsvc.
func generateSVCApplyConfigFromPBSVC(pbSVC *pb.SetServiceReq) *corev1.ServiceApplyConfiguration {
	metaData := &metav1.ObjectMetaApplyConfiguration{
		Name:      &pbSVC.Service.Name,
		Namespace: &pbSVC.ProjectName,
		Labels: map[string]string{
			"name": pbSVC.Service.Name,
		},
	}

	spec := &corev1.ServiceSpecApplyConfiguration{
		Selector: map[string]string{
			"name": pbSVC.Service.Name,
		},
		ClusterIP: &headlessClusterIP,
	}

	for _, exposePort := range pbSVC.Service.ExposePorts {
		spec.Ports = append(spec.Ports, corev1.ServicePortApplyConfiguration{
			Name: &exposePort.Name,
			Port: &exposePort.Port,
		})
	}

	svc := &corev1.ServiceApplyConfiguration{
		TypeMetaApplyConfiguration:   typeMetaApplyConfigureation_Service,
		ObjectMetaApplyConfiguration: metaData,
		Spec:                         spec,
	}

	return svc
}
