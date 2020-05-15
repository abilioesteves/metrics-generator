package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
)

var instance *Collector

// Collector defines the available metric collectors for prometheus
type Collector struct {
	RequestSecondsHistogram *prometheus.HistogramVec

	ResponseBytesCounter *prometheus.CounterVec

	DependencyUp *prometheus.GaugeVec

	ApplicationInfo *prometheus.GaugeVec
}

// Init properly initializes system metrics and registers them to the prometheus registry
func Init() *Collector {
	logrus.Infof("Registering metrics collectors...")
	if instance == nil {
		instance = &Collector{
			RequestSecondsHistogram: getRequestSecondsHistogram(),
			ResponseBytesCounter:    getResponseBytesCounter(),
			DependencyUp:            getDependencyUp(),
			ApplicationInfo:         getApplicationInfo(),
		}

		prometheus.MustRegister(instance.RequestSecondsHistogram, instance.ResponseBytesCounter, instance.DependencyUp,
			instance.ApplicationInfo)
	}

	logrus.Infof("Now collecting HTTP Requestes metrics!")
	return instance
}

func getRequestSecondsHistogram() *prometheus.HistogramVec {
	return prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "request_seconds",
		Help:    "HTTP requests count and latency histogram",
		Buckets: []float64{0.1, 0.3, 2},
	}, []string{
		"type",         // request type (http, grpc, etc)
		"status",       // response status
		"method",       // method used to reach the endpoint
		"addr",         // endpoint address
		"isError",      // flag indicating if the status means an error or not
		"errorMessage", // error message if status is an error
	})
}

func getResponseBytesCounter() *prometheus.CounterVec {
	return prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "response_size_bytes",
		Help: "Response size bytes gauge",
	}, []string{
		"type",         // request type (http, grpc, etc)
		"status",       // response status
		"method",       // method used to reach the endpoint
		"addr",         // endpoint address
		"isError",      // flag indicating if the status means an error or not
		"errorMessage", // error message if status is an error
	})
}

func getDependencyUp() *prometheus.GaugeVec {
	return prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "dependency_up",
		Help: "dependencies status",
	}, []string{
		"name",
	})
}

func getApplicationInfo() *prometheus.GaugeVec {
	return prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "application_info",
		Help: "static information about the application",
	}, []string{
		"version",
	})
}
