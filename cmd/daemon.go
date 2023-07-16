package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"sync"
	"syscall"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/upmio/config-wrapper/app"
	"github.com/upmio/config-wrapper/conf"
	"github.com/upmio/config-wrapper/protocol"
	"github.com/upmio/config-wrapper/version"

	// 新增服务需要在import
	_ "github.com/upmio/config-wrapper/app/config"
)

var (
	configPath string
	file       *os.File
)

// RootCmd represents the base command when called without any subcommands
var daemonCmd = &cobra.Command{
	Use:   "daemon",
	Short: "Run as a daemon process",
	RunE: func(cmd *cobra.Command, args []string) error {
		defer file.Close()

		// 初始化全局变量
		if err := conf.LoadConfigFromToml(configPath); err != nil {
			return err
		}

		// 初始化全局日志配置
		if err := loadGlobalLogger(); err != nil {
			return err
		}

		defer zap.L().Sync()

		// 初始化全局app
		if err := app.InitAllApp(); err != nil {
			return err
		}

		// Make sure global variable config has been initialized
		_ = conf.GetConf()

		// 启动服务
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL, syscall.SIGHUP, syscall.SIGQUIT)

		// 初始化服务
		svr, err := newService()
		if err != nil {
			return err
		}

		// 等待信号处理
		wg := &sync.WaitGroup{}
		wg.Add(1)
		go svr.waitSign(ch, wg)

		// 启动服务
		if err := svr.start(); err != nil {
			if !strings.Contains(err.Error(), "http: Server closed") {
				return err
			}
		}
		wg.Wait()
		return nil
	},
}

func init() {
	daemonCmd.PersistentFlags().StringVarP(&configPath, "file", "f", "/etc/github.com/upmio/config-wrapper/config.toml", "Specify the config file path")
	rootCmd.AddCommand(daemonCmd)
}

type service struct {
	http   *protocol.HTTPService
	grpc   *protocol.GrpcService
	logger *zap.SugaredLogger
}

func newService() (*service, error) {
	http := protocol.NewHTTPService()
	grpc := protocol.NewGrpcService()
	svr := &service{
		http:   http,
		grpc:   grpc,
		logger: zap.L().Named("[Service]").Sugar(),
	}

	return svr, nil
}

func (s *service) start() error {
	s.logger.Info(fmt.Sprintf("loaded http app: %s", app.LoadedHttpApp()))
	s.logger.Info(fmt.Sprintf("loaded grpc apps %s", app.LoadedGrpcApp()))
	go s.grpc.Start()

	return s.http.Start()
}

func (s *service) waitSign(sign chan os.Signal, wg *sync.WaitGroup) {
	for {
		select {
		case sg := <-sign:
			switch v := sg.(type) {
			default:
				zap.L().Named("[HTTP SERVICE]").Sugar().Infof("receive signal '%v', start graceful shutdown", v.String())
				if err := s.http.Stop(); err != nil {
					zap.L().Named("[HTTP SERVICE]").Sugar().Errorf("http graceful shutdown err: %s, force exit", err)
				} else {
					zap.L().Named("[HTTP SERVICE]").Info("http service stop complete")
					s.grpc.Stop()
				}
				wg.Done()
				return
			}
		}
	}
}

func loadGlobalLogger() error {

	//logger, _ := zap.NewProduction()
	cfg := zap.NewProductionEncoderConfig()
	cfg.TimeKey = "timestamp"
	cfg.MessageKey = "message"
	cfg.NameKey = "module"
	cfg.EncodeTime = zapcore.ISO8601TimeEncoder
	cfg.EncodeLevel = zapcore.CapitalLevelEncoder

	if _, err := os.Stat(conf.GetConf().Log.PathDir); os.IsNotExist(err) {
		err := os.Mkdir(conf.GetConf().Log.PathDir, 0755)
		if err != nil {
			return fmt.Errorf("Create %s directory fail, error: %v ", conf.GetConf().Log.PathDir, err)
		}
	}

	logJsonfile := filepath.Join(conf.GetConf().Log.PathDir, version.ServiceName+"-json.log")
	fileJson, err := os.Create(logJsonfile)
	if err != nil {
		return fmt.Errorf("Create log json file fail, error: %v ", err)
	}

	logfile := filepath.Join(conf.GetConf().Log.PathDir, version.ServiceName+".log")
	file, err := os.Create(logfile)
	if err != nil {
		return fmt.Errorf("Create log file fail, error: %v ", err)
	}

	core := zapcore.NewTee(
		zapcore.NewCore(zapcore.NewJSONEncoder(cfg), zapcore.AddSync(fileJson), conf.GetConf().Log.GetLogLevel()),
		zapcore.NewCore(zapcore.NewConsoleEncoder(cfg), zapcore.AddSync(file), conf.GetConf().Log.GetLogLevel()),
		zapcore.NewCore(zapcore.NewConsoleEncoder(cfg), zapcore.Lock(os.Stdout), conf.GetConf().Log.GetLogLevel()),
	)
	//logger := zap.New(core, zap.AddCaller())
	logger := zap.New(core)

	zap.ReplaceGlobals(logger)

	zap.L().Named("[INIT]").Info(conf.Banner)
	zap.L().Named("[INIT]").Info(fmt.Sprintf("log level: %s", conf.GetConf().Level))

	return nil
}