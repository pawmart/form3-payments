package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// AuthRequiredMiddleware function.
func authRequiredMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Header.Get("X-APIKEY") != "abc" {
			c.JSON(http.StatusForbidden, &APIError{
				Code:    string(http.StatusForbidden),
				Message: "forbidden"})
			return
		}
		c.Next()
	}
}
