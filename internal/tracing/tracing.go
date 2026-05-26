package tracing

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"

	"github.com/lookingcamel/system-framework/internal/config"
	"github.com/lookingcamel/system-framework/internal/logger"
)

var tracer trace.Tracer

func Init(cfg config.TracingConfig) (func(context.Context) error, error) {
	if !cfg.Enabled {
		tracer = otel.Tracer(cfg.ServiceName)
		return func(ctx context.Context) error { return nil }, nil
	}

	var exporter sdktrace.SpanExporter
	var err error

	switch cfg.Exporter {
	case "otlp", "jaeger":
		exporter, err = initOTLPExporter(cfg)
	case "stdout":
		exporter, err = initStdoutExporter(cfg)
	default:
		exporter, err = initStdoutExporter(cfg)
	}

	if err != nil {
		return nil, err
	}

	res, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(cfg.ServiceName),
			semconv.ServiceVersion("1.0.0"),
		),
	)
	if err != nil {
		return nil, err
	}

	sampler := sdktrace.ParentBased(sdktrace.TraceIDRatioBased(cfg.SampleRate))

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
		sdktrace.WithSampler(sampler),
	)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	tracer = tp.Tracer(cfg.ServiceName)

	logger.Log.Info("Tracing initialized",
		zap.String("exporter", cfg.Exporter),
		zap.String("service", cfg.ServiceName),
		zap.String("endpoint", cfg.Endpoint),
	)

	return tp.Shutdown, nil
}

func initOTLPExporter(cfg config.TracingConfig) (sdktrace.SpanExporter, error) {
	logger.Log.Info("Initializing OTLP exporter",
		zap.String("endpoint", cfg.Endpoint),
	)

	ctx := context.Background()

	exporter, err := otlptracegrpc.New(ctx,
		otlptracegrpc.WithEndpoint(cfg.Endpoint),
		otlptracegrpc.WithInsecure(),
	)
	if err != nil {
		logger.Log.Error("Failed to initialize OTLP exporter",
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to create OTLP exporter: %w", err)
	}

	return exporter, nil
}

func initStdoutExporter(cfg config.TracingConfig) (sdktrace.SpanExporter, error) {
	logger.Log.Info("Initializing stdout exporter")

	exporter, err := stdouttrace.New(stdouttrace.WithPrettyPrint())
	if err != nil {
		logger.Log.Error("Failed to initialize stdout exporter",
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to create stdout exporter: %w", err)
	}

	return exporter, nil
}

func Tracer() trace.Tracer {
	return tracer
}

func StartSpan(ctx context.Context, name string, opts ...trace.SpanStartOption) (context.Context, trace.Span) {
	return tracer.Start(ctx, name, opts...)
}

func StartSpanWithAttributes(ctx context.Context, name string, attrs ...trace.SpanStartOption) (context.Context, trace.Span) {
	return tracer.Start(ctx, name, attrs...)
}

func Inject(ctx context.Context, carrier propagation.TextMapCarrier) {
	otel.GetTextMapPropagator().Inject(ctx, carrier)
}

func Extract(ctx context.Context, carrier propagation.TextMapCarrier) context.Context {
	return otel.GetTextMapPropagator().Extract(ctx, carrier)
}