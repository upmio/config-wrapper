package protocol

import (
	"github.com/upmio/config-wrapper/app"
	"github.com/upmio/config-wrapper/conf"
	"go.uber.org/zap"
	"net"

	"google.golang.org/grpc"
)

type GrpcService struct {
	l *zap.SugaredLogger
	s *grpc.Server
}

func NewGrpcService() *GrpcService {
	server := grpc.NewServer()
	return &GrpcService{
		s: server,
		l: zap.L().Named("[GRPC SERVICE]").Sugar(),
	}
}

func (g *GrpcService) Start() {
	app.LoadGrpcApp(g.s)

	addr := conf.GetConf().GrpcAddr()
	lsr, err := net.Listen("tcp", addr)
	if err != nil {
		g.l.Errorf("listen grpc tcp conn error, %s", err)
		return
	}

	g.l.Infof("GRPC服务启动成功, 监听地址: %s", addr)

	if err := g.s.Serve(lsr); err != nil {
		if err == grpc.ErrServerStopped {
			g.l.Info("service is stopped")
		}

		g.l.Error("start grpc service error, %s", err.Error())
		return
	}
}

func (g *GrpcService) Stop() {
	g.l.Info("start graceful shutdown")
	g.s.GracefulStop()
	g.l.Info("service is stopped")
}