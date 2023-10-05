from flask import Flask, request
from opentelemetry import trace
from opentelemetry.propagate import extract
from opentelemetry.sdk.trace import TracerProvider
from opentelemetry.trace import (
    SpanKind,
    get_tracer_provider,
    set_tracer_provider,
)
from opentelemetry.instrumentation.wsgi import collect_request_attributes
set_tracer_provider(TracerProvider())

tracer = get_tracer_provider().get_tracer("service_b.tracer")

app = Flask(__name__)

@app.route('/endpoint_b', methods=['GET'])
def endpoint_b():
    with tracer.start_as_current_span(
        "endpoint_b",
        context=extract(request.headers),
        kind=SpanKind.SERVER,
        attributes=collect_request_attributes(request.environ),
    ):
        ctx = trace.get_current_span().get_span_context()
        trace_id = '{trace:032x}'.format(trace=ctx.trace_id)
        span_id = '{span:016x}'.format(span=ctx.span_id)
        return "Response from Service B. TraceID: " + trace_id + " SpanID: " + span_id

if __name__ == "__main__":
    app.run(port=5001)