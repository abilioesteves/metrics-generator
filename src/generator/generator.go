package generator

// Generator defines the methods a Generator should implement
type Generator interface {
	CreateAccident(accident Accident) (err error)
	DeleteAccident(accidentType, resourceName string) (err error)
	DeleteAccidents() (err error)
	SetEntropy(e Entropy) (err error)
}

// Accident accident
type Accident struct {
	ResourceName string  `json:"resourcename,omitempty"`
	Type         string  `json:"type,omitempty"`
	Value        float64 `json:"value,omitempty"`
}

// Entropy defines the available configurations for a Tabajara Generator
type Entropy struct {
	URICount int `json:"uricount,omitempty"`

	ServiceVersionCount int `json:"serviceversioncount,omitempty"`

	AppVersionCount int `json:"apiversioncount,omitempty"`

	DeviceCount int `json:"devicecount,omitempty"`
}

// GetDefaultEntropy returns the default entropy setup for a generator
func GetDefaultEntropy() Entropy {
	return Entropy{
		URICount:            10,
		ServiceVersionCount: 2,
		AppVersionCount:     2,
		DeviceCount:         2,
	}
}

// GetAccidentKey returns the accident key based on the resource name and resource type
func GetAccidentKey(resourceName, accidentType string) string {
	return resourceName + accidentType
}
