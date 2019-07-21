package utils

import (
	"fmt"
	"io"
	"log"
	"os"
	"time"

	opentracing "github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/zipkin"
)

//GetOsEnvironmentVariable gets an environment variable or fallsback to the default
func GetOsEnvironmentVariable(name, fallback string) string {
	if value, ok := os.LookupEnv(name); ok {
		return value
	}
	return fallback
}

// InitJaeger returns an instance of Jaeger Tracer that samples 100% of traces and logs all spans to stdout.
func InitJaeger(service string) (opentracing.Tracer, io.Closer) {
	zipkinPropagator := zipkin.NewZipkinB3HTTPHeaderPropagator()
	injector := jaeger.TracerOptions.Injector(opentracing.HTTPHeaders, zipkinPropagator)
	extractor := jaeger.TracerOptions.Extractor(opentracing.HTTPHeaders, zipkinPropagator)

	// Zipkin shares span ID between client and server spans; it must be enabled via the following option.
	zipkinSharedRPCSpan := jaeger.TracerOptions.ZipkinSharedRPCSpan(true)

	agentHost := GetOsEnvironmentVariable("jaeger_agent_host", "localhost")
	agentPort := GetOsEnvironmentVariable("jaeger_agent_port", "5775")
	sender, err := jaeger.NewUDPTransport(fmt.Sprintf("%s:%s", agentHost, agentPort), 0)
	if err != nil {
		log.Fatal(fmt.Sprintf("ERROR: cannot init Jaeger: %v\n", err))
		return nil, nil
	}

	tracer, closer := jaeger.NewTracer(
		service,
		jaeger.NewConstSampler(true),
		jaeger.NewRemoteReporter(
			sender,
			jaeger.ReporterOptions.BufferFlushInterval(1*time.Second),
			jaeger.ReporterOptions.Logger(jaeger.StdLogger)),
		injector,
		extractor,
		zipkinSharedRPCSpan,
	)

	opentracing.SetGlobalTracer(tracer)

	return tracer, closer
}
