package tracing

import (
	"context"
	"log/slog"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.23.1"
)

func InitTracer(ctx context.Context) (*sdktrace.TracerProvider, error) {
	//создаем экспортер, который отправляет данные трасировки(спаны) в jaeger по http
	exp, err := jaeger.New(
		jaeger.WithCollectorEndpoint( //обозначаем что у нас http запрос (collector)
			jaeger.WithEndpoint("http://jaeger:14268/api/traces"), //указываем url адрес куда будут отправлятся спаны
		),
	)
	if err != nil {
		return nil, err
	}

	//создаем информацию о нашем сервисе
	res, err := resource.New( //создаем ресурс и лобавляем в него атрибуты
		ctx,
		resource.WithSchemaURL(semconv.SchemaURL),
		resource.WithAttributes(semconv.ServiceName("api-gateway")),
	)
	if err != nil {
		slog.Error("failed to create resource", slog.Any("error", err))
		return nil, err
	}

	tp := sdktrace.NewTracerProvider( //создаем менеджер для всех трейсеров
		sdktrace.WithBatcher(exp),  //говорим чтобы спаны отправлялись не по одному, а много
		sdktrace.WithResource(res), //описываем ресурс (имя, версия и тд) откуда пришли трейсы
	)

	otel.SetTracerProvider(tp) //передаем что через этот объект все трессируем (активируем трассировку)
	return tp, nil
}
