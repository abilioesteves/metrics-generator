package hook

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/abilioesteves/goh/gohserver"
	"github.com/abilioesteves/goh/gohtypes"
	"github.com/abilioesteves/metrics-generator-tabajara/src/generator"
	"github.com/gorilla/mux"
	"github.com/labbsr0x/bindman-dns-webhook/src/types"
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
	router.HandleFunc("/accidents/{resourceName}", hook.DeleteAccident).Methods("DELETE")
	router.HandleFunc("/accidents/{resourceName}", hook.CreateAccident).Methods("POST")
	router.HandleFunc("/entropy/increase/{n}", hook.CreateAccident).Methods("POST")
	router.HandleFunc("/entropy/decrease/{n}", hook.DecreaseEntropy).Methods("POST")

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
	if strings.Trim(resourceName, " ") != "" {
		err := hook.Generator.DeleteAccident(vars["resourceName"])

		types.PanicIfError(types.Error{Message: fmt.Sprintf("Not possible to delete the accident '%s'", vars["resourceName"]), Code: 500, Err: err})

		gohserver.WriteJSONResponse(true, 200, w)
	}

	gohtypes.Panic("No resource name given.", 406)
}

// IncreaseEntropy increases the number of returned time-series
func (hook *DefaultHook) IncreaseEntropy(w http.ResponseWriter, r *http.Request) {
	hook.increaseOrDecreaseEntropy(w, r, hook.Generator.IncreaseEntropy)
}

// DecreaseEntropy descreaes the number of returned time-series
func (hook *DefaultHook) DecreaseEntropy(w http.ResponseWriter, r *http.Request) {
	hook.increaseOrDecreaseEntropy(w, r, hook.Generator.DecreaseEntropy)
}

func (hook *DefaultHook) increaseOrDecreaseEntropy(w http.ResponseWriter, r *http.Request, do func(n int) error) {
	defer gohserver.HandleError(w)

	vars := mux.Vars(r)
	nstr := vars["n"]
	if strings.Trim(nstr, " ") != "" {
		n, err := strconv.Atoi(nstr)
		gohtypes.PanicIfError("The parameter n is not a valid number", 406, err)

		err = do(n)
		gohtypes.PanicIfError("Not possible to adjust the entropy of generated metrics", 500, err)

		gohserver.WriteJSONResponse(true, 200, w)
	}

	gohtypes.Panic("No number of time-series to adjust provided", 406)
}
