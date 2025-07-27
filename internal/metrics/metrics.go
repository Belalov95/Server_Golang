package metrics

import (
	"example/web-service-gin/internal/cache"
	"log/slog"
	"net/http"
	"runtime"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	HttpStatusMetric = prometheus.NewCounterVec( // создание метрики счетчик, для подсчета количества http запросов
		prometheus.CounterOpts{ // создается для описания метрики
			Name: "http_request_total",
			Help: "Number of http response labeled by status code and method",
		},
		[]string{"status", "method", "path"}, // список лейблов по которым будут группирроваться данные
	)
	CacheMemoryUsage = prometheus.NewGaugeFunc( // метрика, которая отслеживает память
		prometheus.GaugeOpts{
			Name: "cache_memory_usage_bytes",
			Help: "Amount of memory occupied by the cache in bytes",
		},
		GetCacheMemoryMetrics, // ф-ция которую NewGaugeFunc будет вызывать каждый раз, когда Prometheus запрашивает метрики
	)
)

// ф-ция возвращает примерное количество выделенной памяти программой
func GetCacheMemoryMetrics() float64 {
	var m runtime.MemStats   // содержит статистику о выделенной памяти в программе
	runtime.ReadMemStats(&m) // заполняем m актуальной информацией о памяти
	return float64(m.Alloc)  // возвращаем общее кол-во байт выделенных прямо сейчас
}

func InitMetrics(port string, cache *cache.CacheDecorator) {
	prometheus.MustRegister(HttpStatusMetric) // регистрация метрики в prometheus
	prometheus.MustRegister(CacheMemoryUsage)
	http.Handle("/metrics", promhttp.Handler()) // регистрация маршрута /metrics

	go func() { // создаем чтобы параллельно запустить с основными API, который работает на другом порту
		server := &http.Server{ // создание нового http сервера
			Addr:         ":" + port,        // адрес и порт
			ReadTimeout:  5 * time.Second,   // максимальное время чтения запроса
			WriteTimeout: 10 * time.Second,  // максимальное время отправки ответа
			IdleTimeout:  120 * time.Second, // максимальное время простоя между запросами
		}
		slog.Info("Starting metrics server on port:", slog.String("port:", port))
		if err := server.ListenAndServe(); err != nil {
			slog.Error("Failed to start metrics:", slog.Any("error:", err))
		}
	}()
}

// Используется для увеличения счетчика http - запросов в prometheus
func HttpStatusMetricInc(statusCode int, method string, path string) {
	// метод .WithLabelValues() возвращает конкретный экземпляр счетчика с лейблами и .Inc() добавляет +1 как counter++
	HttpStatusMetric.WithLabelValues(http.StatusText(statusCode), method, path).Inc()
}
