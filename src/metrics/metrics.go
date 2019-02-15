package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
)

var instance *Collector

// Collector defines the available metric collectors for prometheus
type Collector struct {
	HTTPRequestsPerServiceVersionSummary *prometheus.SummaryVec

	HTTPRequestsPerServiceVersion *prometheus.HistogramVec

	HTTPRequestsPerAppVersion *prometheus.CounterVec

	HTTPRequestsPerDevice *prometheus.CounterVec

	HTTPPendingRequests *prometheus.GaugeVec
}

// Init properly initializes system metrics and registers them to the prometheus registry
func Init() *Collector {
	logrus.Infof("Registering metrics collectors...")
	if instance == nil {
		instance = &Collector{
			HTTPRequestsPerServiceVersion: getHTTPRequestsPerServiceVersion(),
			HTTPRequestsPerAppVersion:     getHTTPRequestsPerAppVersion(),
			HTTPRequestsPerDevice:         getHTTPRequestsPerDevice(),
			HTTPPendingRequests:           getHTTPPendingRequests(),
		}

		prometheus.MustRegister(instance.HTTPRequestsPerServiceVersion, instance.HTTPRequestsPerAppVersion, instance.HTTPRequestsPerDevice, instance.HTTPPendingRequests)
	}

	logrus.Infof("Now collecting HTTP Requestes metrics!")
	return instance
}

func getHTTPRequestsPerServiceVersionSummary() *prometheus.SummaryVec {
	return prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Name:       "http_requests_seconds_summary",
		Help:       "HTTP requests count and latency summary",
		Objectives: map[float64]float64{0.75: 0.2, 0.95: 0.05},
	}, []string{
		"uri",             // requested resource
		"method",          // HTTP method
		"status",          // status of the HTTP request
		"service_version", // version of the back end system
	})
}

func getHTTPRequestsPerServiceVersion() *prometheus.HistogramVec {
	return prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "http_requests_seconds_histogram",
		Help:    "HTTP requests count and latency histogram",
		Buckets: []float64{0.3, 4, 35},
	}, []string{
		"uri",             // requested resource
		"method",          // HTTP method
		"status",          // status of the HTTP request
		"service_version", // version of the back end system
	})
}

func getHTTPRequestsPerAppVersion() *prometheus.CounterVec {
	return prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "http_requests_app_version_count",
		Help: "HTTP requests count per app version",
	}, []string{
		"uri",         // requested resource
		"method",      // HTTP method
		"status",      // status of the HTTP request
		"app_version", // version of the mobile app
	})
}

func getHTTPRequestsPerDevice() *prometheus.CounterVec {
	return prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "http_requests_device_count",
		Help: "HTTP requests count per device",
	}, []string{
		"uri",    // requested resource
		"method", // HTTP method
		"status", // status of the HTTP request
		"device", // version of the mobile app
	})
}

func getHTTPPendingRequests() *prometheus.GaugeVec {
	return prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "http_pending_requests",
		Help: "HTTP pending requests",
	}, []string{
		"service_version", // version of the back end system
	})
}
