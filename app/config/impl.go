package config

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	"github.com/upmio/config-wrapper/app"
	"github.com/upmio/config-wrapper/app/config/confd"
	"github.com/upmio/config-wrapper/app/config/confd/backends"
	"github.com/upmio/config-wrapper/app/config/confd/template"
	"github.com/upmio/config-wrapper/conf"
)

var (
	// service 服务实例
	svr = &service{}
)

const (
	templateKey = "template"
	pathKey     = "path"
	valueKey    = "value"
)

type service struct {
	syncConfig SyncConfigServiceServer
	clientSet  kubernetes.Interface
	UnimplementedSyncConfigServiceServer
	logger      *zap.SugaredLogger
	confdConfig *confd.Config
}

func (s *service) Config() error {
	clientSet, err := conf.GetConf().Kube.GetClientSet()
	if err != nil {
		return err
	}

	s.clientSet = clientSet
	s.syncConfig = app.GetGrpcApp(appName).(SyncConfigServiceServer)
	s.logger = zap.L().Named("[CONFIG]").Sugar()
	s.confdConfig = &confd.Config{
		BackendsConfig: confd.BackendsConfig{
			Backend: "content",
		},
	}
	return nil
}

func (s *service) Name() string {
	return appName
}

func (s *service) Registry(server *grpc.Server) {
	RegisterSyncConfigServiceServer(server, svr)
}

func (s *service) SyncConfig(ctx context.Context, req *SyncConfigRequest) (*SyncConfigResponse, error) {
	configMapObj, err := s.clientSet.CoreV1().ConfigMaps(req.GetNamespace()).Get(ctx, req.GetConfigmapName(), metav1.GetOptions{})
	if err != nil {
		errMsg := fmt.Sprintf("Can't found ConfigMap %s in Namespace %s, error: %s", req.GetConfigmapName(), req.GetNamespace(), err)
		s.logger.Errorf(errMsg)
		return &SyncConfigResponse{
			Message: errMsg,
		}, fmt.Errorf(errMsg)
	}

	// Check configmap has value、template、path key
	if value, ok := configMapObj.Data[templateKey]; !ok {
		errMsg := fmt.Sprintf("ConfigMap %s in Namespace %s does't has key %s", req.GetConfigmapName(), req.GetNamespace(), templateKey)
		s.logger.Errorf(errMsg)
		return &SyncConfigResponse{
			Message: errMsg,
		}, fmt.Errorf(errMsg)
	} else {
		tmplFile := filepath.Join("/tmp", "template.tmpl")
		err := os.WriteFile(tmplFile, []byte(value), 0644)
		if err != nil {
			errMsg := fmt.Sprintf("Write template file %s failed, error: %v", tmplFile, err)
			s.logger.Errorf(errMsg)
			return &SyncConfigResponse{
				Message: errMsg,
			}, fmt.Errorf(errMsg)
		} else {
			defer os.Remove(tmplFile)
			s.logger.Debugf("Write template file %s Success", tmplFile)
			s.confdConfig.TemplateConfig.TemplateFile = tmplFile
		}

	}

	if value, ok := configMapObj.Data[pathKey]; !ok {
		errMsg := fmt.Sprintf("ConfigMap %v in Namespace %v does't has key %s", pathKey)
		s.logger.Errorf(errMsg)
		return &SyncConfigResponse{
			Message: errMsg,
		}, fmt.Errorf(errMsg)
	} else {
		s.confdConfig.TemplateConfig.DestFile = os.ExpandEnv(value)
	}

	if value, ok := configMapObj.Data[valueKey]; !ok {
		errMsg := fmt.Sprintf("ConfigMap %v in Namespace %v does't has key %s", valueKey)
		s.logger.Errorf(errMsg)
		return &SyncConfigResponse{
			Message: errMsg,
		}, fmt.Errorf(errMsg)
	} else {
		s.confdConfig.BackendsConfig.Content = value
	}

	// Initialize the storage client
	storeClient, err := backends.New(s.confdConfig.BackendsConfig)

	s.confdConfig.TemplateConfig.StoreClient = storeClient
	if err := template.Process(s.confdConfig.TemplateConfig); err != nil {
		errMsg := fmt.Sprintf("Generate config file with configmap %v in namespace %v failed, error: %v", req.GetConfigmapName(), req.GetNamespace(), err)
		s.logger.Errorf(errMsg)
		return &SyncConfigResponse{
			Message: errMsg,
		}, fmt.Errorf(errMsg)
	}

	successMsg := fmt.Sprintf("Generate config %v success.", req.GetConfigmapName())
	s.logger.Info(successMsg)

	return &SyncConfigResponse{Message: successMsg}, nil
}

func init() {
	app.RegistryGrpcApp(svr)
}