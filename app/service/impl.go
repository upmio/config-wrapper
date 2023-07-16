package service

import (
	"config-wrapper/app"
	"config-wrapper/conf"
	"context"
	"fmt"
	"github.com/abrander/go-supervisord"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

var (
	// service 服务实例
	svr = &service{}
)

type service struct {
	service ServiceLifecycleServer
	UnimplementedServiceLifecycleServer
	logger *zap.SugaredLogger
	client *supervisord.Client
}

func (s *service) Config() error {
	client, err := conf.GetConf().Supervisor.GetSupervisorClient()
	if err != nil {
		return err
	}

	s.client = client
	s.service = app.GetGrpcApp(appName).(ServiceLifecycleServer)
	s.logger = zap.L().Named("[SERVICE]").Sugar()

	return nil
}

func (s *service) Name() string {
	return appName
}

func (s *service) Registry(server *grpc.Server) {
	RegisterServiceLifecycleServer(server, svr)
}

func (s *service) StartService(_ context.Context, _ *ServiceRequest) (*ServiceResponse, error) {
	err := s.client.StartProcess("unit_app", false)
	if err != nil {
		errMsg := fmt.Sprintf("Start service failed, error: %v", err)
		s.logger.Errorf(errMsg)
		return &ServiceResponse{
			Message: errMsg,
		}, fmt.Errorf(errMsg)

	}

	successMsg := "Start service success."
	s.logger.Info(successMsg)

	return &ServiceResponse{Message: successMsg}, nil
}

func (s *service) StopService(_ context.Context, _ *ServiceRequest) (*ServiceResponse, error) {
	err := s.client.StopProcess("unit_app", false)
	if err != nil {
		errMsg := fmt.Sprintf("Stop service failed, error: %v", err)
		s.logger.Errorf(errMsg)
		return &ServiceResponse{
			Message: errMsg,
		}, fmt.Errorf(errMsg)

	}

	successMsg := "Stop service success."
	s.logger.Info(successMsg)

	return &ServiceResponse{Message: successMsg}, nil
}

func init() {
	app.RegistryGrpcApp(svr)
}