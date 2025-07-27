package tracing

import (
	"context"
	"log/slog"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.23.1"
)

func Init(ctx context.Context) (*sdktrace.TracerProvider, error) {
	// создаём OTLP-экспортёр (через HTTP)
	exp, err := otlptracehttp.New(ctx,
		otlptracehttp.WithEndpoint("jaeger:4318"), // OTLP HTTP порт по умолчанию
		otlptracehttp.WithInsecure(),              // если без TLS
	)
	if err != nil {
		slog.Error("failed to create OTLP exporter", slog.Any("error", err))
		return nil, err
	}

	// создаем информацию о нашем сервисе
	res, err := resource.New( // создаем ресурс и лобавляем в него атрибуты
		ctx,
		resource.WithSchemaURL(semconv.SchemaURL),
		resource.WithAttributes(semconv.ServiceName("api-gateway")),
	)
	if err != nil {
		slog.Error("failed to create resource", slog.Any("error", err))
		return nil, err
	}

	tp := sdktrace.NewTracerProvider( // создаем менеджер для всех трейсеров
		sdktrace.WithBatcher(exp),  // говорим чтобы спаны отправлялись не по одному, а много
		sdktrace.WithResource(res), // описываем ресурс (имя, версия и тд) откуда пришли трейсы
	)

	otel.SetTracerProvider(tp) // передаем что через этот объект все трессируем (активируем трассировку)
	return tp, nil
}
