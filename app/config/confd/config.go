package confd

import (
	"github.com/upmio/config-wrapper/app/config/confd/backends"
	"github.com/upmio/config-wrapper/app/config/confd/template"
)

type TemplateConfig = template.Config
type BackendsConfig = backends.Config

// A Config structure is used to configure confd.
type Config struct {
	TemplateConfig
	BackendsConfig
}