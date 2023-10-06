package otel_utils

import (
	"context"
	"errors"
	"net/http"
	"sync"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/stats"
)

var initOtelOnce sync.Once
var initOtelManagerOnce sync.Once
var tracerProvider *sdktrace.TracerProvider
var otelManager *OtelManager

type OtelManager struct {
	serviceName    string
	serviceVersion string
}

func InitOtelManager(serviceName, serviceVersion string) *OtelManager {
	initOtelManagerOnce.Do(func() {
		otelManager = &OtelManager{
			serviceName:    serviceName,
			serviceVersion: serviceVersion,
		}
	})
	return otelManager
}

func GetOtelManager() (*OtelManager, error) {
	if otelManager == nil {
		return nil, errors.New("Call InitOtelManager before calling GetOtelManager")
	}
	return otelManager, nil
}

// setupOTelSDK bootstraps the OpenTelemetry pipeline.
// If it does not return an error, make sure to call shutdown for proper cleanup.
func (om *OtelManager) SetupOTelSDK(ctx context.Context) (shutdown func(context.Context) error, err error) {
	var shutdownFuncs []func(context.Context) error

	// shutdown calls cleanup functions registered via shutdownFuncs.
	// The errors from the calls are joined.
	// Each registered cleanup will be invoked once.
	shutdown = func(ctx context.Context) error {
		var err error
		for _, fn := range shutdownFuncs {
			err = errors.Join(err, fn(ctx))
		}
		shutdownFuncs = nil
		return err
	}

	// handleErr calls shutdown for cleanup and makes sure that all errors are returned.
	handleErr := func(inErr error) {
		err = errors.Join(inErr, shutdown(ctx))
	}

	// Setup trace provider.
	tracerProvider, err := om.GetTracerProvider()
	if err != nil {
		handleErr(err)
		return
	}
	shutdownFuncs = append(shutdownFuncs, tracerProvider.Shutdown)
	otel.SetTracerProvider(tracerProvider)
	otel.SetTextMapPropagator(propagation.TraceContext{})
	return
}

func (om *OtelManager) GetTextMapPropagator() propagation.TextMapPropagator {
	return otel.GetTextMapPropagator()
}

func (om *OtelManager) NewResource() (*resource.Resource, error) {
	return resource.Merge(resource.Default(),
		resource.NewWithAttributes(semconv.SchemaURL,
			semconv.ServiceName(om.serviceName),
			semconv.ServiceVersion(om.serviceVersion),
		))
}

func (om *OtelManager) GetTracerProvider() (*sdktrace.TracerProvider, error) {
	var err error
	initOtelOnce.Do(
		func() {
			var res *resource.Resource
			res, _ = om.NewResource()
			tracerProvider = sdktrace.NewTracerProvider(
				sdktrace.WithResource(res),
			)
		})
	return tracerProvider, err
}

func (om *OtelManager) GetTracer(tracerName string) trace.Tracer {
	tp, err := om.GetTracerProvider()
	if err != nil {
		// handle error log it
	}
	return tp.Tracer(tracerName)
}

func (om *OtelManager) AddOtelMiddlewareforMuxRouter(h http.Handler, route string) http.Handler {
	return otelhttp.NewHandler(h, route)
}

func (om *OtelManager) GetOtelGrpcHandler() stats.Handler {
	return otelgrpc.NewClientHandler()
}

func (om *OtelManager) GetOtelTransportForHttp() http.RoundTripper {
	return otelhttp.NewTransport(http.DefaultTransport)
}
