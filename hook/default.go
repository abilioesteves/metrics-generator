package hook

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/labbsr0x/goh/gohserver"
	"github.com/labbsr0x/goh/gohtypes"
	"github.com/abilioesteves/metrics-generator/generator"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
)

// DefaultHook defines the generator structure that implements the API
type DefaultHook struct {
	Generator generator.Generator
	router    *mux.Router
}

// NewDefaultHook instantiates a new default hook structure
func NewDefaultHook(generator generator.Generator) *DefaultHook {
	hook := &DefaultHook{
		Generator: generator,
		router:    mux.NewRouter(),
	}

	hook.router.HandleFunc("/accidents/{accidentType}/{resourceName}", hook.DeleteAccident).Methods("DELETE")
	hook.router.HandleFunc("/accidents", hook.DeleteAccidents).Methods("DELETE")
	hook.router.HandleFunc("/accidents", hook.CreateAccident).Methods("POST")
	hook.router.HandleFunc("/entropy/set", hook.SetEntropy).Methods("POST")
	hook.router.Handle("/metrics", promhttp.Handler())

	return hook
}

// Init Initializes a new default webhook
func (hook *DefaultHook) Init() {
	logrus.Info("Initializing the default webhook...")
	err := http.ListenAndServe("0.0.0.0:32865", hook.router)
	if err != nil {
		logrus.Errorf("Error initializing the Metrics Generator Tabajara: %v", err)
		os.Exit(1)
	}
	logrus.Infof("Default webhook initialized!")
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
