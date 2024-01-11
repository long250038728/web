package auth

import "github.com/gin-gonic/gin"

type Auth interface {
	Allow(c *gin.Context) error
}
