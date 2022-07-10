package main

import (
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/contrib/propagators/b3"
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
	log.SetFormatter(customLogger{
		formatter: log.JSONFormatter{FieldMap: log.FieldMap{
			"msg": "message",
		}},
	})
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
	b3 := b3.New()
	otel.SetTextMapPropagator(b3)
	return tp, err
}

func main() {
	ctx := context.Background()
	tp, err := initTracer()
	if err != nil {
		log.Error(err)
	}
	defer func() {
		if err := tp.Shutdown(ctx); err != nil {
			log.Error("Error shutting down tracer provider: %v", err)
		}
	}()

	otelHandler := otelhttp.NewHandler(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		//log.Info("About End Point Printning...")
		// inject req.Context
		log.WithContext(req.Context()).Info("About End Point Printning...")
		fmt.Fprintln(w, "about page")
	}), "About")
	http.Handle("/about", otelHandler)
	log.Info("Listening...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

type customLogger struct {
	formatter log.JSONFormatter
}

func (l customLogger) Format(entry *log.Entry) ([]byte, error) {
	span := trace.SpanFromContext(entry.Context)
	entry.Data["trace_id"] = span.SpanContext().TraceID().String()
	entry.Data["span_id"] = span.SpanContext().SpanID().String()
	//Below injection is Just to understand what Context has
	entry.Data["Context"] = span.SpanContext()
	return l.formatter.Format(entry)
}
