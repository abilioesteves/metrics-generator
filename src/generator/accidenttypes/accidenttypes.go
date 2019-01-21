package accidenttypes

const (

	// Calls the number of times a specific resource gets called
	Calls = "calls"

	// Latency the latency of a request to a specific resource
	Latency = "latency"

	// ErrorRate the rate at which error occurs to a specific resource [0,1]. Defaults to 0,33
	ErrorRate = "errorrate"
)

const (
	// DefaultNumberOfCalls defines the default number of calls a single resource gets called.
	DefaultNumberOfCalls = 1.0

	// DefaultLatency defines the default call latency in seconds
	DefaultLatency = 0.3

	// DefaultErrorRate defines the default error rate applicable to a single resource
	DefaultErrorRate = 0.33
)
