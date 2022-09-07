package main

import (
	"fmt"
	"net/http"
	"net/http/pprof"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/zipkin"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

func handle(c *gin.Context) {
	c.Status(200)
	_, _ = c.Writer.Write([]byte("OK\n"))
}

func main() {
	err := setupZipkin()
	if err != nil {
		panic(err)
	}

	router := gin.New()
	router.Use(func(c *gin.Context) {
		fmt.Println("middleware")
	})
	router.GET("/", handle)
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))
	router.GET("/debug/pprof/*profile", gin.WrapF(pprof.Index))
	err = http.ListenAndServe(":8080", otelhttp.NewHandler(router.Handler(), "root"))
	if err != nil {
		panic(err)
	}
}
func setupZipkin() error {
	zipkin, err := zipkin.New("http://zipkin:9411/api/v2/spans")
	if err != nil {
		return err
	}
	tp := sdktrace.NewTracerProvider(sdktrace.WithBatcher(zipkin), sdktrace.WithResource(resource.NewSchemaless(attribute.String("service.name", "leaker"))))

	otel.SetTracerProvider(tp)
	return nil
}
