// Code generated by swaggo/swag. DO NOT EDIT.

package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {
            "name": "俞桢浩",
            "url": "http://bsgchind.io",
            "email": "yuzhenhao0906@gmail.com"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/config/sync": {
            "post": {
                "description": "在kubernetes集群中获取ConfigMap的内容,使用confd落地生成为对应软件格式的配置文件",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "同步配置文件接口"
                ],
                "summary": "同步配置文件接口",
                "parameters": [
                    {
                        "description": "ConfigMap信息",
                        "name": "SyncConfigRequest",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/config.SyncConfigRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/config.SyncConfigResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/config.SyncConfigResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/config.SyncConfigResponse"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "config.SyncConfigRequest": {
            "type": "object",
            "required": [
                "configmap_name",
                "namespace"
            ],
            "properties": {
                "configmap_name": {
                    "type": "string"
                },
                "namespace": {
                    "type": "string"
                }
            }
        },
        "config.SyncConfigResponse": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "127.0.0.1:22104",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "github.com/upmio/config-wrapper api demo",
	Description:      "swagger测试",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}