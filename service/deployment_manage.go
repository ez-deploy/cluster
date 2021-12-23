package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/ez-deploy/protobuf/model"
	pb "github.com/ez-deploy/protobuf/project"
	v1 "k8s.io/api/apps/v1"
	apimetav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	appsv1 "k8s.io/client-go/applyconfigurations/apps/v1"
	corev1 "k8s.io/client-go/applyconfigurations/core/v1"
	metav1 "k8s.io/client-go/applyconfigurations/meta/v1"
)

const describeAnnotationName = "describe"

var restKindDeployment = "Deployment"
var typeMetaApplyConfigureation_Deployment = metav1.TypeMetaApplyConfiguration{
	APIVersion: &apiVersionV1,
	Kind:       &restKindDeployment,
}

// create a Headless k8s deployment from given pbSVC.
func (s *Service) applyDeploymentFromCfg(ctx context.Context, pbSVC *pb.SetServiceReq) error {
	svcCfg := generateDeploymentApplyConfigFromPBSVC(pbSVC)
	applyOptions := apimetav1.ApplyOptions{}

	_, err := s.k8sClientset.AppsV1().Deployments(pbSVC.ProjectName).Apply(ctx, svcCfg, applyOptions)

	return err
}

func (s *Service) getPBSVCFromDeployment(ctx context.Context, req *pb.GetServiceReq) (*model.Service, error) {
	getOption := apimetav1.GetOptions{}

	deployment, err := s.k8sClientset.AppsV1().Deployments(req.ProjectName).Get(ctx, req.ServiceName, getOption)
	if err != nil {
		return nil, err
	}

	return generatePBSVCFromDeployment(deployment), nil
}

func (s *Service) listPBSVC(ctx context.Context, req *pb.ListServiceReq) ([]*model.Service, error) {
	listOption := apimetav1.ListOptions{}

	deployments, err := s.k8sClientset.AppsV1().Deployments(req.ProjectName).List(ctx, listOption)
	if err != nil {
		return nil, err
	}

	resPBSVCList := []*model.Service{}
	for _, deployment := range deployments.Items {
		resPBSVCList = append(resPBSVCList, generatePBSVCFromDeployment(&deployment))
	}

	return resPBSVCList, nil
}

// delete a deployment from given pbSVC.
func (s *Service) deleteDeploymentFromCfg(ctx context.Context, pbSVC *pb.DeleteServiceReq) error {
	deleteOptions := apimetav1.DeleteOptions{}

	return s.k8sClientset.AppsV1().Deployments(pbSVC.ProjectName).Delete(ctx, pbSVC.ServiceName, deleteOptions)
}

// generate deployment kind k8s object from given pbsvc.
func generateDeploymentApplyConfigFromPBSVC(pbSVC *pb.SetServiceReq) *appsv1.DeploymentApplyConfiguration {
	metaData := &metav1.ObjectMetaApplyConfiguration{
		Name:      &pbSVC.Service.Name,
		Namespace: &pbSVC.ProjectName,
		Labels: map[string]string{
			"name": pbSVC.Service.Name,
		},
		Annotations: map[string]string{
			describeAnnotationName: pbSVC.Service.Describe,
		},
	}

	spec := &appsv1.DeploymentSpecApplyConfiguration{
		Replicas: &pbSVC.Service.Replica,
		Selector: metav1.LabelSelector().WithMatchLabels(map[string]string{"name": pbSVC.Service.Name}),
		Template: generateDeplymentTemplateConfigFromPBSVC(pbSVC.Service),
	}

	deployment := &appsv1.DeploymentApplyConfiguration{
		TypeMetaApplyConfiguration:   typeMetaApplyConfigureation_Deployment,
		ObjectMetaApplyConfiguration: metaData,
		Spec:                         spec,
	}

	return deployment
}

func generateDeplymentTemplateConfigFromPBSVC(svc *model.Service) *corev1.PodTemplateSpecApplyConfiguration {
	imageFullURL := fmt.Sprintf("%s:%s", svc.Image.Url, svc.Image.Version)
	container := corev1.ContainerApplyConfiguration{
		Name:  &svc.Name,
		Image: &imageFullURL,
	}

	for _, env := range svc.Envs {
		container.Env = append(container.Env, corev1.EnvVarApplyConfiguration{
			Name:  &env.Key,
			Value: &env.Value,
		})
	}

	for _, exposePort := range svc.ExposePorts {
		container.Ports = append(container.Ports, corev1.ContainerPortApplyConfiguration{
			Name:          &exposePort.Name,
			ContainerPort: &exposePort.Port,
		})
	}

	template := &corev1.PodTemplateSpecApplyConfiguration{
		ObjectMetaApplyConfiguration: &metav1.ObjectMetaApplyConfiguration{
			Labels: map[string]string{"name": svc.Name},
		},
		Spec: &corev1.PodSpecApplyConfiguration{
			Containers: []corev1.ContainerApplyConfiguration{container},
		},
	}

	return template
}

func generatePBSVCFromDeployment(deployment *v1.Deployment) *model.Service {
	pbsvc := &model.Service{
		Name:     deployment.ObjectMeta.Name,
		Describe: deployment.ObjectMeta.Annotations[describeAnnotationName],
		Replica:  *deployment.Spec.Replicas,
		Image:    &model.Image{},
	}

	// only one container in a template.
	container := deployment.Spec.Template.Spec.Containers[0]

	splitedImageFullURL := strings.SplitN(container.Image, ":", 2)
	pbsvc.Image.Url, pbsvc.Image.Version = splitedImageFullURL[0], splitedImageFullURL[1]

	for _, exposePort := range container.Ports {
		pbsvc.ExposePorts = append(pbsvc.ExposePorts, &model.Port{
			Name: exposePort.Name,
			Port: exposePort.ContainerPort,
		})
	}

	for _, env := range container.Env {
		pbsvc.Envs = append(pbsvc.Envs, &model.EnvironmentVariable{
			Key:   env.Name,
			Value: env.Value,
		})
	}

	return pbsvc
}
