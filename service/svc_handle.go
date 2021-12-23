package service

import (
	"context"

	"github.com/ez-deploy/protobuf/model"
	pb "github.com/ez-deploy/protobuf/project"
)

// CreateService under a project.
// will create a svc & deployment under the given namespace.
func (s *Service) SetService(ctx context.Context, req *pb.SetServiceReq) (*model.CommonResp, error) {
	if err := s.applyServiceFromCfg(ctx, req); err != nil {
		return model.NewCommonRespWithError(err), nil
	}

	if err := s.applyDeploymentFromCfg(ctx, req); err != nil {
		return model.NewCommonRespWithError(err), nil
	}

	return &model.CommonResp{}, nil
}

func (s *Service) GetService(ctx context.Context, req *pb.GetServiceReq) (*pb.GetServiceResp, error) {
	pbsvc, err := s.getPBSVCFromDeployment(ctx, req)
	if err != nil {
		return &pb.GetServiceResp{Error: model.NewError(err)}, nil
	}

	return &pb.GetServiceResp{Service: pbsvc}, nil
}

func (s *Service) ListService(ctx context.Context, req *pb.ListServiceReq) (*pb.ListServiceResp, error) {
	pbsvcList, err := s.listPBSVC(ctx, req)
	if err != nil {
		return &pb.ListServiceResp{Error: model.NewError(err)}, nil
	}

	return &pb.ListServiceResp{Service: pbsvcList}, nil
}

func (s *Service) DeleteService(ctx context.Context, req *pb.DeleteServiceReq) (*model.CommonResp, error) {
	if err := s.deleteServiceFromCfg(ctx, req); err != nil {
		return model.NewCommonRespWithError(err), nil
	}

	if err := s.deleteDeploymentFromCfg(ctx, req); err != nil {
		return model.NewCommonRespWithError(err), nil
	}

	return &model.CommonResp{}, nil
}
