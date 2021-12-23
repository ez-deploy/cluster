package service

import (
	"context"

	"github.com/ez-deploy/protobuf/model"
	pb "github.com/ez-deploy/protobuf/project"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	corev1 "k8s.io/client-go/applyconfigurations/core/v1"
)

const tryExistLimit = 3

// CreateProject via k8s client-go (exactly, create a namespace).
func (s *Service) CreateProject(ctx context.Context, req *pb.CreateProjectReq) (*model.CommonResp, error) {
	namespace := corev1.Namespace(req.Project.Name)
	options := metav1.ApplyOptions{}

	if _, err := s.k8sClientset.CoreV1().Namespaces().Apply(ctx, namespace, options); err != nil {
		return model.NewCommonRespWithError(err), nil
	}

	return &model.CommonResp{}, nil
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
