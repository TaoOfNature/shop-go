package middleware

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

func Tracing() gin.HandlerFunc {
	tracer := otel.Tracer("shop-go/http")

	return func(c *gin.Context) {
		path := c.FullPath()
		if path == "" {
			path = c.Request.URL.Path
		}

		ctx, span := tracer.Start(c.Request.Context(), fmt.Sprintf("%s %s", c.Request.Method, path))
		span.SetAttributes(
			attribute.String("http.method", c.Request.Method),
			attribute.String("http.route", path),
			attribute.String("http.target", c.Request.URL.Path),
			attribute.String("client.ip", c.ClientIP()),
		)
		c.Request = c.Request.WithContext(ctx)

		c.Next()

		span.SetAttributes(attribute.Int("http.status_code", c.Writer.Status()))
		if len(c.Errors) > 0 {
			span.SetStatus(codes.Error, c.Errors.String())
		} else {
			span.SetStatus(codes.Ok, "ok")
		}
		span.End()
	}
}
