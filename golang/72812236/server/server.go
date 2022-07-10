package main

import (
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	stdout "go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.10.0"
	"go.opentelemetry.io/otel/trace"
	"net/http"
	"os"
)

func init() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
}

func initTracer() (*sdktrace.TracerProvider, error) {
	exporter, err := stdout.New(stdout.WithPrettyPrint())
	if err != nil {
		return nil, err
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(resource.NewWithAttributes(semconv.SchemaURL, semconv.ServiceNameKey.String("PreferencesService"))),
	)
	otel.SetTracerProvider(tp)
	//b3 := b3.New()
	//otel.SetTextMapPropagator(b3)
	return tp, err
}

func main() {
	tp, err := initTracer()
	if err != nil {
		log.Error(err)
	}
	defer func() {
		ctx := context.Background()
		if err := tp.Shutdown(ctx); err != nil {
			log.Error("Error shutting down tracer provider: %v", err)
		}
	}()

	otelHandler := otelhttp.NewHandler(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		spanCtx := trace.SpanContextFromContext(ctx)
		traceId := spanCtx.TraceID().String()
		spanId := spanCtx.SpanID().String()
		log.AddHook(NewTraceIdHook(traceId, spanId, spanCtx))
		log.Info("About End Point Printning...")
		fmt.Fprintln(w, "about page")
	}), "About")
	http.Handle("/about", otelHandler)
	log.Info("Listening...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

type TraceIdHook struct {
	TraceId string
	SpanId  string
	Context trace.SpanContext
}

func NewTraceIdHook(traceId, spanId string, ctx trace.SpanContext) log.Hook {
	hook := TraceIdHook{
		TraceId: traceId,
		SpanId:  spanId,
		Context: ctx,
	}
	return &hook
}

func (hook *TraceIdHook) Fire(entry *log.Entry) error {
	entry.Data["span_id"] = hook.SpanId
	entry.Data["traceId"] = hook.TraceId
	entry.Data["context"] = hook.Context
	return nil
}

func (hook *TraceIdHook) Levels() []log.Level {
	return log.AllLevels
}
