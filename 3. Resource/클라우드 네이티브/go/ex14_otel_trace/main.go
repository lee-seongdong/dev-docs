package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"go.opentelemetry.io/otel/trace"
)

const (
	jaegerEndpoint = "http://localhost:4317"
	serviceName    = "Fibonacci"
)

func Fibonacci(ctx context.Context, n int) chan int {
	ch := make(chan int)

	go func() {
		tr := otel.GetTracerProvider().Tracer(serviceName)

		cctx, sp := tr.Start(ctx, fmt.Sprintf("Fibonacci(%d)", n), trace.WithAttributes(attribute.Int("fibonacci.n", n)))
		defer sp.End()

		result := 1
		if n > 1 {
			a := Fibonacci(cctx, n-1)
			b := Fibonacci(cctx, n-2)
			result = <-a + <-b
		}
		sp.SetAttributes(attribute.Int("result", result))
		ch <- result
	}()

	return ch
}

func fibonacciHandler(w http.ResponseWriter, req *http.Request) {
	var err error
	var n int

	if len(req.URL.Query()["n"]) != 1 {
		err = fmt.Errorf("wrong number of arguments")
	} else {
		n, err = strconv.Atoi(req.URL.Query()["n"][0])
	}

	if err != nil {
		http.Error(w, "couldn't parse index n", 400)
		return
	}

	ctx := req.Context()

	result := <-Fibonacci(ctx, n)

	if sp := trace.SpanFromContext(ctx); sp != nil {
		sp.SetAttributes(
			attribute.Int("parameter", n),
			attribute.Int("result", result),
		)
	}

	fmt.Fprintln(w, result)
}

func newTracerProvider(ctx context.Context) (*sdktrace.TracerProvider, error) {
	res, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(serviceName),
		),
	)

	stdExporter, err := stdouttrace.New(
		stdouttrace.WithPrettyPrint(),
	)

	if err != nil {
		return nil, fmt.Errorf("failed to build StdoutExporter: %w", err)
	}

	otlpExporter, err := otlptracegrpc.New(
		ctx,
		otlptracegrpc.WithEndpoint(jaegerEndpoint),
		otlptracegrpc.WithInsecure(),
	)

	if err != nil {
		return nil, fmt.Errorf("failed to build OtlpExporter: %w", err)
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithResource(res),
		sdktrace.WithBatcher(stdExporter),
		sdktrace.WithBatcher(otlpExporter),
	)

	return tp, nil
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	tp, err := newTracerProvider(ctx)
	if err != nil {
		slog.ErrorContext(ctx, err.Error())
		return
	}

	// Handle shutdown properly so nothing leaks.
	defer func() { tp.Shutdown(ctx) }()

	otel.SetTracerProvider(tp)
	fmt.Println("Browse to localhost:3000?n=6")

	http.Handle("/", otelhttp.NewHandler(http.HandlerFunc(fibonacciHandler), "root"))
	if err := http.ListenAndServe(":3000", nil); err != nil {
		slog.ErrorContext(ctx, err.Error())
		return
	}
}
