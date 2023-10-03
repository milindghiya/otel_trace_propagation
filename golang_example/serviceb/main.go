package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/milindghiya/otel_trace_propagation/golang_example/otel_utils"
)

func commonMiddleware(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	})
}

var serviceName string = "serviceB"
var serviceVersion string = "0.1.0"

func main() {

	om := otel_utils.InitOtelManager(serviceName, serviceVersion)
	// Initialize TracerProvider
	otelShutdown, err := om.SetupOTelSDK(context.Background())
	defer func() {
		err = errors.Join(err, otelShutdown(context.Background()))
	}()
	r := mux.NewRouter()
	r.Use(commonMiddleware)
	r.HandleFunc("/new", func(w http.ResponseWriter, r *http.Request) {
		tracer := om.GetTracer("serviceB")
		_, span := tracer.Start(r.Context(), "main")
		defer span.End()

		// Log TraceID and SpanID
		traceID := span.SpanContext().TraceID()
		spanID := span.SpanContext().SpanID()
		log.Printf("Service B - TraceID: %s, SpanID: %s\n", traceID, spanID)

		// Respond to the request
		w.Write([]byte(fmt.Sprintf("\nService B responding - TraceID: %s, SpanID: %s", traceID, spanID)))
	})

	handler := om.AddOtelMiddlewareforMuxRouter(r, "/")
	log.Fatal(http.ListenAndServe(":8081", handler))
}
