package hook

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/abilioesteves/goh/gohserver"
	"github.com/abilioesteves/goh/gohtypes"
	"github.com/abilioesteves/metrics-generator-tabajara/src/generator"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

// WebHook defines the contracts of the methods that should be implemented by concrete hooks
type WebHook interface {
	CreateAccident(w http.ResponseWriter, r *http.Request)
	DeleteAccident(w http.ResponseWriter, r *http.Request)
	IncreaseEntropy(w http.ResponseWriter, r *http.Request)
	DecreaseEntropy(w http.ResponseWriter, r *http.Request)
}

// DefaultHook defines the gernerator structure that implements the API
type DefaultHook struct {
	Generator generator.Generator
}

// Init initializes a new default webhook
func Init(generator generator.Generator) *DefaultHook {
	hook := &DefaultHook{
		Generator: generator,
	}

	router := mux.NewRouter()
	router.HandleFunc("/accidents/{accidentType}/{resourceName}", hook.DeleteAccident).Methods("DELETE")
	router.HandleFunc("/accidents", hook.DeleteAccidents).Methods("DELETE")
	router.HandleFunc("/accidents", hook.CreateAccident).Methods("POST")
	router.HandleFunc("/entropy/set", hook.CreateAccident).Methods("POST")

	logrus.Info("Initialized Metrics Generator Tabajara Webhook")
	err := http.ListenAndServe("0.0.0.0:32865", router)
	if err != nil {
		logrus.Errorf("Error initializing the Metrics Generator Tabajara: %v", err)
	}

	return hook
}

// CreateAccident creates observation accidents to an specific resource
func (hook *DefaultHook) CreateAccident(w http.ResponseWriter, r *http.Request) {
	defer gohserver.HandleError(w)

	var accident generator.Accident
	err := json.NewDecoder(r.Body).Decode(&accident)
	gohtypes.PanicIfError("Unable to decode the request body", 406, err)

	err = hook.Generator.CreateAccident(accident)
	gohtypes.PanicIfError("Unable to create an Accident", 500, err)

	gohserver.WriteJSONResponse(true, 200, w)
}

// DeleteAccident deletes observation accidents to an specific resource
func (hook *DefaultHook) DeleteAccident(w http.ResponseWriter, r *http.Request) {
	defer gohserver.HandleError(w)

	vars := mux.Vars(r)
	resourceName := vars["resourceName"]
	accidentType := vars["accidentType"]
	if strings.Trim(resourceName, " ") != "" && strings.Trim(accidentType, " ") != "" {
		err := hook.Generator.DeleteAccident(accidentType, resourceName)
		gohtypes.PanicIfError(fmt.Sprintf("Not possible to delete the accident '%s'", vars["resourceName"]), 500, err)

		gohserver.WriteJSONResponse(true, 200, w)
	}

	gohtypes.Panic("Accident type and resource name must be provided in order to identify the accident", 406)
}

// DeleteAccidents deletes all observation accidents
func (hook *DefaultHook) DeleteAccidents(w http.ResponseWriter, r *http.Request) {
	defer gohserver.HandleError(w)

	err := hook.Generator.DeleteAccidents()
	gohtypes.PanicIfError(fmt.Sprintf("Not possible to remove all accidents"), 500, err)

	gohserver.WriteJSONResponse(true, 200, w)
}

// SetEntropy increases the number of returned time-series
func (hook *DefaultHook) SetEntropy(w http.ResponseWriter, r *http.Request) {
	defer gohserver.HandleError(w)

	var entropy generator.Entropy
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&entropy)
	gohtypes.PanicIfError("Unable to decode the payload", 406, err)

	err = hook.Generator.SetEntropy(entropy)
	gohtypes.PanicIfError("Not possible to adjust the entropy of generated metrics", 500, err)

	gohserver.WriteJSONResponse(true, 200, w)
}
