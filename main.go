package main

import (
	"github.com/gin-gonic/gin"
)

func setupRouter() *gin.Engine {
	r := gin.Default()
	r.GET("/health", HealthCheck)
	v1 := r.Group("/v1")
	{
		v1.Use(authRequiredMiddleware())
		v1.GET("/payments", GetPayments)
		v1.GET("/payments/:id", GetSinglePayment)
		v1.POST("/payments", CreatePayment)
		v1.PATCH("/payments", UpdatePayment)
		v1.DELETE("/payments/:id", DeletePayment)
	}
	return r
}

func main() {
	r := setupRouter()
	r.Run(":6543")
}
