from flask import Flask, request
import requests
import grpc

# Import OpenTelemetry and related packages
from opentelemetry import trace
from opentelemetry.propagate import extract
from opentelemetry.instrumentation.flask import FlaskInstrumentor
from opentelemetry.instrumentation.grpc import GrpcInstrumentorClient
from opentelemetry.sdk.trace import TracerProvider
from opentelemetry.trace import (
    SpanKind,
    get_tracer_provider,
    set_tracer_provider,
)
from opentelemetry.trace.propagation.tracecontext import TraceContextTextMapPropagator
from opentelemetry.instrumentation.wsgi import collect_request_attributes
# Import generated gRPC classes
import service_c_pb2
import service_c_pb2_grpc

app = Flask(__name__)

set_tracer_provider(TracerProvider())

tracer = get_tracer_provider().get_tracer("service_a.tracer")
GrpcInstrumentorClient().instrument()

@app.route('/', methods=['GET'])
def endpoint_a():
     with tracer.start_as_current_span(
        "endpoint_a",
        context=extract(request.headers),
        kind=SpanKind.SERVER,
        attributes=collect_request_attributes(request.environ),
    ):  
        # Write the current context into the carrier.
        carrier = {}
        TraceContextTextMapPropagator().inject(carrier)
        ctx = trace.get_current_span().get_span_context()
        trace_id = '{trace:032x}'.format(trace=ctx.trace_id)
        span_id = '{span:016x}'.format(span=ctx.span_id)
        print(trace_id + " "+ span_id)

        
        # HTTP call to Service B
        response_b = requests.get('http://localhost:5001/endpoint_b',headers=carrier)
        
        # gRPC call to Service C
        with grpc.insecure_channel('localhost:50051') as channel:
            stub = service_c_pb2_grpc.ServiceCStub(channel)
            response_c = stub.GetResponse(service_c_pb2.RequestC(message='Request from Service A.'),metadata=carrier)
        
        return f"Response B: {response_b.text}, Response C: {response_c.message}"

if __name__ == "__main__":
    app.run(port=8080)





