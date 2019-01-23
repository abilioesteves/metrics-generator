package generator

import (
	"context"
	"time"

	"github.com/abilioesteves/metrics-generator-tabajara/src/generator/accidenttypes"
	"github.com/abilioesteves/metrics-generator-tabajara/src/metrics"
	"github.com/sirupsen/logrus"
)

// Tabajara generates metrics
type Tabajara struct {
	Collector *metrics.Collector
	Entropy   Entropy
	Accidents map[string]Accident
}

// NewGeneratorTabajara instantiates a new
func NewGeneratorTabajara(collector *metrics.Collector, entropy Entropy) *Tabajara {
	return &Tabajara{
		Collector: collector,
		Entropy:   entropy,
	}
}

// Init initializes the generation of the dummy metrics
func (gen *Tabajara) Init(ctx context.Context) {
	logrus.Infof("Initializing metrics generator...")
	go func() {
		c := time.Tick(10 * time.Millisecond)
		for {
			select {
			case <-ctx.Done():
				logrus.Info("Generator Tabajara stopped!")
				return
			case <-c:
				gen.FillMetrics()
			}
		}
	}()
	logrus.Infof("Metrics generator initialized!")
}

// CreateAccident creates observation accidents to an specific resource
func (gen *Tabajara) CreateAccident(accident Accident) (err error) {
	gen.Accidents[GetAccidentKey(accident.Type, accident.ResourceName)] = accident
	return
}

// DeleteAccident deletes observation accidents to an specific resource
func (gen *Tabajara) DeleteAccident(accidentType, resourceName string) (err error) {
	delete(gen.Accidents, GetAccidentKey(accidentType, resourceName))
	return
}

// DeleteAccidents deletes all accidents
func (gen *Tabajara) DeleteAccidents() (err error) {
	gen.Accidents = make(map[string]Accident)
	return
}

// SetEntropy increases the number of returned time-series by n
func (gen *Tabajara) SetEntropy(e Entropy) (err error) {
	gen.Entropy = e
	return
}

// FillMetrics advances the state of the registered generator metrics with configurable random values
func (gen *Tabajara) FillMetrics() {
	methods := []string{"POST", "GET", "DELETE", "PUT"}
	oss := []string{"ios", "android"}

	uri := getRandomElemNormal(gen.getUris())
	serviceVersion := getRandomElemNormal(gen.getServiceVersions())
	calls := int(gen.getValueAccident(accidenttypes.Calls, accidenttypes.DefaultNumberOfCalls, uri))

	for i := 0; i < calls; i++ {
		appVersion := getRandomElemNormal(gen.getAppVersions())
		device := getRandomElemNormal(gen.getDevices())
		os := getRandomElemNormal(oss)
		method := methods[randomInt(int64(hash(uri)), len(methods))]
		status := getStatusWithErrorAccident(gen.getValueAccident(accidenttypes.ErrorRate, accidenttypes.DefaultErrorRate, uri))

		gen.FillHTTPRequestsPerServiceVersion(uri, method, status, serviceVersion)
		gen.FillHTTPRequestsPerAppVersion(uri, method, status, appVersion)
		gen.FillHTTPRequestsPerDevice(uri, method, status, os, device)
	}

	gen.FillHTTPPendingRequests(serviceVersion)
}

// FillHTTPRequestsPerServiceVersion fills the HTTPRequestsPerServiceVersion metric
func (gen *Tabajara) FillHTTPRequestsPerServiceVersion(uri, method, status, serviceVersion string) {
	gen.Collector.HTTPRequestsPerServiceVersion.WithLabelValues(
		uri,
		method,
		status,
		serviceVersion,
	).Observe(gen.getValueAccident(accidenttypes.Latency, getSampleRequestTime(uri), uri))
}

// FillHTTPRequestsPerAppVersion fills the HTTPRequestsPerAppVersion metric
func (gen *Tabajara) FillHTTPRequestsPerAppVersion(uri, method, status, appVersion string) {
	gen.Collector.HTTPRequestsPerAppVersion.WithLabelValues(
		uri,
		method,
		status,
		appVersion,
	).Inc()
}

// FillHTTPPendingRequests fills the HTTPPendingRequests metric
func (gen *Tabajara) FillHTTPPendingRequests(serviceVersion string) {
	gen.Collector.HTTPPendingRequests.WithLabelValues(
		serviceVersion,
	).Set(float64(randomRangeNormal(0, 400)))
}

// FillHTTPRequestsPerDevice fills the HTTPRequestsPerDevice metric
func (gen *Tabajara) FillHTTPRequestsPerDevice(uri, method, status, os, device string) {
	gen.Collector.HTTPRequestsPerDevice.WithLabelValues(
		uri,
		method,
		status,
		os+device,
	).Inc()
}

func (gen *Tabajara) getUris() []string {
	return generateItems("/resource/test-", gen.Entropy.URICount)
}

func (gen *Tabajara) getServiceVersions() []string {
	return generateItems("backend-v", gen.Entropy.ServiceVersionCount)
}

func (gen *Tabajara) getAppVersions() []string {
	return generateItems("v", gen.Entropy.AppVersionCount)
}

func (gen *Tabajara) getDevices() []string {
	return generateItems("-", gen.Entropy.DeviceCount)
}

func (gen *Tabajara) getValueAccident(accidentType string, defaultValue float64, resourceName string) float64 {
	if accident, ok := gen.Accidents[GetAccidentKey(resourceName, accidentType)]; ok {
		return accident.Value
	}
	return defaultValue
}
