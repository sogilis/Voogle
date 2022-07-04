// Package docs GENERATED BY SWAG; DO NOT EDIT
// This file was generated by swaggo/swag
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/api/v1/videos/list/{attribute}/{order}/{page}/{limit}": {
            "get": {
                "description": "Get list of all videos",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "video"
                ],
                "summary": "Get list of all videos",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Sort attribute",
                        "name": "attribute",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Sort order",
                        "name": "order",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Page number",
                        "name": "page",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Video per page",
                        "name": "limit",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Video list and Hateoas links",
                        "schema": {
                            "$ref": "#/definitions/controllers.VideoListResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/api/v1/videos/upload": {
            "post": {
                "description": "Upload video file",
                "consumes": [
                    "multipart/form-data"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "video"
                ],
                "summary": "Upload video file",
                "parameters": [
                    {
                        "type": "file",
                        "description": "video",
                        "name": "file",
                        "in": "formData",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Video and Links (HATEOAS)",
                        "schema": {
                            "$ref": "#/definitions/controllers.Response"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "409": {
                        "description": "This title already exists",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "415": {
                        "description": "Unsupported Media Type",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/api/v1/videos/{id}/delete": {
            "delete": {
                "description": "Delete video",
                "produces": [
                    "text/plain"
                ],
                "tags": [
                    "video"
                ],
                "summary": "Delete video",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Video ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/api/v1/videos/{id}/info": {
            "get": {
                "description": "Get video informations",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "video"
                ],
                "summary": "Get video informations",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Video ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Video Informations",
                        "schema": {
                            "$ref": "#/definitions/json.VideoInfo"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/api/v1/videos/{id}/status": {
            "get": {
                "description": "Get video status",
                "produces": [
                    "text/plain"
                ],
                "tags": [
                    "video"
                ],
                "summary": "Get video status",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Video ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Describe video status",
                        "schema": {
                            "$ref": "#/definitions/json.VideoStatus"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/api/v1/videos/{id}/streams/master.m3u8": {
            "get": {
                "description": "Get video master",
                "produces": [
                    "text/plain"
                ],
                "tags": [
                    "video"
                ],
                "summary": "Get video master",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Video ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "HLS video master",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/api/v1/videos/{id}/streams/{quality}/{filename}": {
            "get": {
                "description": "Get sub part stream video",
                "produces": [
                    "text/plain"
                ],
                "tags": [
                    "video"
                ],
                "summary": "Get sub part stream video",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Video ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Video quality",
                        "name": "quality",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Video sub part name",
                        "name": "filename",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "array",
                        "items": {
                            "type": "string"
                        },
                        "description": "List of required filters",
                        "name": "filter",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Video sub part (.ts)",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/ws": {
            "get": {
                "description": "Send Update to Front",
                "consumes": [
                    "text/plain"
                ],
                "produces": [
                    "text/plain"
                ],
                "tags": [
                    "websocket"
                ],
                "summary": "Send Update to Front",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Authentication cookie",
                        "name": "Cookie",
                        "in": "header",
                        "required": true
                    }
                ],
                "responses": {
                    "101": {
                        "description": "Switching Protocols",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "controllers.Response": {
            "type": "object",
            "properties": {
                "_links": {
                    "type": "object",
                    "additionalProperties": {
                        "$ref": "#/definitions/json.LinkJson"
                    }
                },
                "video": {
                    "$ref": "#/definitions/json.VideoJson"
                }
            }
        },
        "controllers.VideoInfo": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "string",
                    "example": "1"
                },
                "title": {
                    "type": "string",
                    "example": "my title"
                }
            }
        },
        "controllers.VideoListResponse": {
            "type": "object",
            "properties": {
                "_lastpage": {
                    "type": "integer"
                },
                "_links": {
                    "type": "object",
                    "additionalProperties": {
                        "$ref": "#/definitions/json.LinkJson"
                    }
                },
                "videos": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/controllers.VideoInfo"
                    }
                }
            }
        },
        "json.LinkJson": {
            "type": "object",
            "properties": {
                "href": {
                    "type": "string"
                },
                "method": {
                    "type": "string"
                }
            }
        },
        "json.VideoInfo": {
            "type": "object",
            "properties": {
                "title": {
                    "type": "string",
                    "example": "amazingtitle"
                },
                "uploadDateUnix": {
                    "type": "integer",
                    "example": 1652173257
                }
            }
        },
        "json.VideoJson": {
            "type": "object",
            "properties": {
                "createdAt": {
                    "type": "string",
                    "example": "2022-04-15T12:59:52Z"
                },
                "id": {
                    "type": "string",
                    "example": "aaaa-b56b-..."
                },
                "status": {
                    "type": "string",
                    "example": "VIDEO_STATUS_ENCODING"
                },
                "title": {
                    "type": "string",
                    "example": "A Title"
                },
                "updatedAt": {
                    "type": "string",
                    "example": "2022-04-15T12:59:52Z"
                },
                "uploadedAt": {
                    "type": "string",
                    "example": "2022-04-15T12:59:52Z"
                }
            }
        },
        "json.VideoStatus": {
            "type": "object",
            "properties": {
                "status": {
                    "type": "string",
                    "example": "UPLOADED"
                },
                "title": {
                    "type": "string",
                    "example": "AmazingTitle"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "",
	Host:             "",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "",
	Description:      "",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
