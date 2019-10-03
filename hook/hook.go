package hook

import (
	"net/http"
)

// WebHook defines the contracts of the methods that should be implemented by concrete hooks
type WebHook interface {
	Init()
	CreateAccident(w http.ResponseWriter, r *http.Request)
	DeleteAccident(w http.ResponseWriter, r *http.Request)
	IncreaseEntropy(w http.ResponseWriter, r *http.Request)
	DecreaseEntropy(w http.ResponseWriter, r *http.Request)
}
