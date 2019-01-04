package generator

import (
	"github.com/abilioesteves/metrics-generator-tabajara/src/metrics"
	"github.com/sirupsen/logrus"
)

// Generator defines the methods a Generator should implement
type Generator interface {
	CreateAccident(accident Accident) (err error)
	DeleteAccident(resourceName string) (err error)
	IncreaseEntropy(n int) (err error)
	DecreaseEntropy(n int) (err error)
}

// Tabajara generates metrics
type Tabajara struct {
	Collector *metrics.Collector
}

// Accident accident
type Accident struct {
	ResourceName string `json:"resourcename,omitempty"`
	Type         string `json:"type,omitempty"`
	Value        string `json:"value,omitempty"`
}

// New instantiates a new
func New(collector *metrics.Collector) *Tabajara {
	return &Tabajara{
		Collector: collector,
	}
}

// Init initializes the generation of the dummy metrics
func (gen *Tabajara) Init() {
	logrus.Infof("Starting requests simulation to generate metrics...")
	// var statuses = []string{"4xx", "2xx", "5xx"}
	// var methods = []string{"POST", "GET", "DELETE", "PUT"}

}

// CreateAccident creates observation accidents to an specific resource
func (gen *Tabajara) CreateAccident(accident Accident) (err error) {
	return
}

// DeleteAccident deletes observation accidents to an specific resource
func (gen *Tabajara) DeleteAccident(resourceName string) (err error) {
	return
}

// IncreaseEntropy increases the number of returned time-series by n
func (gen *Tabajara) IncreaseEntropy(n int) (err error) {
	return
}

// DecreaseEntropy descreases the number of returned time-series by n
func (gen *Tabajara) DecreaseEntropy(n int) (err error) {
	return
}
