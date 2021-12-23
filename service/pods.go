package service

import (
	"context"
	"fmt"
	"time"

	"github.com/ez-deploy/protobuf/model"
	pb "github.com/ez-deploy/protobuf/project"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (s *Service) ListPods(ctx context.Context, req *pb.ListPodsReq) (*pb.ListPodsResp, error) {
	listOption := metav1.ListOptions{
		LabelSelector: fmt.Sprintf("%s:%s", "name", req.ServiceName),
	}

	listRes, err := s.k8sClientset.CoreV1().Pods(req.ProjectName).List(ctx, listOption)
	if err != nil {
		return &pb.ListPodsResp{Error: model.NewError(err)}, nil
	}

	res := &pb.ListPodsResp{}

	for _, pod := range listRes.Items {
		res.Pods = append(res.Pods, &model.Pod{
			Name:        pod.Name,
			Status:      pod.Status.String(),
			Age:         time.Now().Unix() - pod.CreationTimestamp.Unix(),
			MachineName: pod.Spec.NodeName,
		})
	}

	return res, nil
}
