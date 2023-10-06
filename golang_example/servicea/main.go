package main

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/milindghiya/otel_trace_propagation/golang_example/hello"
	"google.golang.org/grpc"

	"github.com/go-resty/resty/v2"
	"github.com/milindghiya/otel_trace_propagation/golang_example/otel_utils"
	"go.opentelemetry.io/otel/propagation"
)

var serviceName string = "serviceA"
var serviceVersion string = "0.1.0"

func main() {
	om := otel_utils.InitOtelManager(serviceName, serviceVersion)
	// Initialize TracerProvider
	otelShutdown, err := om.SetupOTelSDK(context.Background())
	defer func() {
		err = errors.Join(err, otelShutdown(context.Background()))
	}()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Start a new span
		tracer := om.GetTracer("serviceA")
		ctx, span := tracer.Start(r.Context(), "main")
		defer span.End()

		// Log TraceID and SpanID
		traceID := span.SpanContext().TraceID()
		spanID := span.SpanContext().SpanID()
		log.Printf("Service A - TraceID: %s, SpanID: %s\n", traceID, spanID)
		output := ""
		// Call Service B using http client
		output += CallServiceByUsingHttpClient(ctx)
		// Call Service B using resty client
		output += CallServiceByUsingResty(ctx)
		// Call Service C using grpc
		output += CallServiceUsingGrpc(ctx)
		// Respond to the original request
		w.Write([]byte(fmt.Sprintf("Service A calling Service B and Service C - TraceID: %s, SpanID: %s\n", traceID, spanID) + output))
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func CallServiceByUsingResty(ctx context.Context) string {
	// Call Service B using resty client
	om, _ := otel_utils.GetOtelManager()
	cli := resty.New()
	restyReq := cli.R()
	om.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(restyReq.Header))
	restyReq.SetContext(ctx)
	res, err := restyReq.Get("http://localhost:8081/new")
	if err != nil {
		// handle the error
		return ""
	}
	return string(res.Body()) + " Using Resty client"

}

func CallServiceByUsingHttpClient(ctx context.Context) string {
	om, _ := otel_utils.GetOtelManager()
	req, _ := http.NewRequestWithContext(ctx, "GET", "http://localhost:8081/new", nil)
	client := http.Client{Transport: om.GetOtelTransportForHttp()}
	resp, err := client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		// handle the error
		return ""
	}
	text, _ := ioutil.ReadAll(resp.Body)
	return string(text) + " Using Http client"
}

func CallServiceUsingGrpc(ctx context.Context) string {
	target := "localhost:50051"
	om, _ := otel_utils.GetOtelManager()
	ctx, span := om.GetTracer("serviceA").Start(ctx, "CallServiceUsingGrpc")
	defer span.End()
	conn, err := grpc.Dial(target, grpc.WithInsecure(), grpc.WithStatsHandler(om.GetOtelGrpcHandler()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := hello.NewHelloServiceClient(conn)
	name := "world"
	if len(os.Args) > 1 {
		name = os.Args[1]
	}
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	out, err := c.SayHello(ctx, &hello.HelloRequest{Name: name})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	return out.Message
}
