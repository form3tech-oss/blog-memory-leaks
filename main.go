package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/zipkin"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"net/http"
	"net/http/pprof"
)

func handle(c *gin.Context) {
	c.Status(200)
	_, _ = c.Writer.Write([]byte("OK\n"))
}

type handler struct {
	c *gin.Context
}

func (h *handler) ServeHTTP(w http.ResponseWriter, request *http.Request) {
	h.c.Request = request
	h.c.Next()
}

func main() {
	err := setupZipkin()
	if err != nil {
		panic(err)
	}

	router := gin.New()
	router.Use(func(c *gin.Context) {
		fmt.Println("middleware")
		h := &handler{c: c}
		handlerWithMetrics := otelhttp.NewHandler(h, "root")
		handlerWithMetrics.ServeHTTP(c.Writer, c.Request)
	})
	router.GET("/", handle)
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))
	router.GET("/debug/pprof/*profile", gin.WrapF(pprof.Index))
	err = router.Run(":8080")
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
