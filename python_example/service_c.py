import grpc
from concurrent import futures
import time
from opentelemetry.sdk.trace import TracerProvider
from opentelemetry.trace import (
    SpanKind,
    get_tracer_provider,
    set_tracer_provider,
)
from opentelemetry.propagate import extract
# Import the generated classes
import service_c_pb2
import service_c_pb2_grpc
from opentelemetry import trace 

set_tracer_provider(TracerProvider())

tracer = get_tracer_provider().get_tracer("service_c.tracer")


class ServiceC(service_c_pb2_grpc.ServiceCServicer):
    def GetResponse(self, request, context):
        with tracer.start_as_current_span(
            "endpoint_c",
            context=extract(dict(context.invocation_metadata())),
            kind=SpanKind.SERVER,
        ):
            ctx = trace.get_current_span().get_span_context()
            trace_id = '{trace:032x}'.format(trace=ctx.trace_id)
            span_id = '{span:016x}'.format(span=ctx.span_id)
            print(trace_id+" "+span_id)
            return service_c_pb2.ResponseC(message='Response from Service C. TraceID: '+ trace_id + ' SpanID: '+span_id)

def serve():
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    service_c_pb2_grpc.add_ServiceCServicer_to_server(ServiceC(), server)
    server.add_insecure_port('[::]:50051')
    server.start()
    server.wait_for_termination()

if __name__ == '__main__':
    serve()