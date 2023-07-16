package config

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/upmio/config-wrapper/app"
	"github.com/upmio/config-wrapper/protocol"
)

var (
	// handler 服务实例
	h = &handler{}
)

type handler struct {
	service SyncConfigServiceServer
}

func (h *handler) Config() error {
	h.service = app.GetGrpcApp(appName).(SyncConfigServiceServer)
	return nil
}

func (h *handler) Name() string {
	return appName
}

func (h *handler) Registry(r *gin.Engine, subPath string) {
	userSubRouter := r.Group(protocol.ApiV1).Group(subPath)

	userSubRouter.POST("/sync", h.SyncConfigRouter)

}

// @Summary 同步配置文件接口
// @Description 在kubernetes集群中获取ConfigMap的内容,使用confd落地生成为对应软件格式的配置文件
// @Tags 同步配置文件接口
// @Accept application/json
// @Produce application/json
// @Param SyncConfigRequest body SyncConfigRequest true "ConfigMap信息"
// @Success 200 {object} SyncConfigResponse
// @Failure 400 {object} SyncConfigResponse
// @Failure 500 {object} SyncConfigResponse
// @Router /config/sync [post]
func (h *handler) SyncConfigRouter(c *gin.Context) {
	req := &SyncConfigRequest{}

	if err := c.BindJSON(req); err != nil {
		c.JSON(http.StatusBadRequest, SyncConfigResponse{
			Message: fmt.Sprintf("Request body binding failed, error: %v", err),
		})
	} else {
		resp, err := h.service.SyncConfig(c, req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, resp)
		} else {
			c.JSON(http.StatusCreated, resp)
		}
	}
}

func init() {
	app.RegistryHttpApp(h)
}