package service

import (
	"context"

	"github.com/ez-deploy/protobuf/model"
	pb "github.com/ez-deploy/protobuf/project"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const tryExistLimit = 3

// CreateProject via k8s client-go (exactly, create a namespace).
func (s *Service) CreateProject(ctx context.Context, req *pb.CreateProjectReq) (*model.CommonResp, error) {
	namespace := &corev1.Namespace{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Namespace",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:        req.Project.Name,
			Annotations: map[string]string{},
		},
	}

	if len(req.Project.Describe) != 0 {
		namespace.ObjectMeta.Annotations["describe"] = req.Project.Describe
	}

	options := metav1.CreateOptions{}

	if _, err := s.k8sClientset.CoreV1().Namespaces().Create(ctx, namespace, options); err != nil {
		return model.NewCommonRespWithError(err), nil
	}

	return &model.CommonResp{Error: nil}, nil
}

// DeleteProject via k8s client-go (exactly, delete a namespace).
// ATTENTION: Can delete only when no service under project.
func (s *Service) DeleteProject(ctx context.Context, req *pb.DeleteProjectReq) (*model.CommonResp, error) {
	listOptions := metav1.ListOptions{Limit: tryExistLimit}

	listResp, err := s.k8sClientset.CoreV1().Services(req.ProjectName).List(ctx, listOptions)
	if err != nil {
		return model.NewCommonRespWithError(err), nil
	}

	if len(listResp.Items) != 0 {
		return model.NewCommonRespWithErrorMessage("Can delete only when no service under project."), nil
	}

	deleteOptions := metav1.DeleteOptions{}
	if err := s.k8sClientset.CoreV1().Namespaces().Delete(ctx, req.ProjectName, deleteOptions); err != nil {
		return nil, err
	}

	return &model.CommonResp{}, nil
}
