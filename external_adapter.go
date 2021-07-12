package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/rs/zerolog/log"
)

type ExternalAdapter struct {
	LocalAddr        string
	InsideDockerAddr string
}

type ExternalAdapterResponse struct {
	JobRunId string              `json:"id"`
	Data     ExternalAdapterData `json:"data"`
	Error    error               `json:"error"`
}

type ExternalAdapterData struct {
	Result int `json:"result"`
}

const adapterPort = ":6060"

type OkResult struct{}

var variableData int

// starts an external adapter on specified port
func main() {
	router := httprouter.New()
	router.POST("/", index)
	router.POST("/random", randomNumber)
	router.POST("/five", five)
	router.POST("/variable", variable)
	router.POST("/set_variable", setVariable)

	log.Info().Str("Port", adapterPort).Msg("Starting external adapter")
	log.Fatal().AnErr("Error", http.ListenAndServe(adapterPort, router)).Msg("Error occured while running external adapter")
}

// index allows a status check on the adapter
func index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Adapter listening!")
	result := &ExternalAdapterResponse{
		JobRunId: "",
		Data:     ExternalAdapterData{Result: 0},
		Error:    nil,
	}
	_ = json.NewEncoder(w).Encode(result)
}

// RandomNumber returns a random int from 0 to 100
func randomNumber(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	num := rand.Intn(100)
	result := &ExternalAdapterResponse{
		JobRunId: "",
		Data:     ExternalAdapterData{Result: num},
		Error:    nil,
	}
	_ = json.NewEncoder(w).Encode(result)
}

// five returns five
func five(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	result := &ExternalAdapterResponse{
		JobRunId: "",
		Data:     ExternalAdapterData{Result: 5},
		Error:    nil,
	}
	_ = json.NewEncoder(w).Encode(result)
}

func setVariable(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	q := r.URL.Query()
	log.Info().Interface("query", q).Msg("params received")
	v := q.Get("var")
	data, _ := strconv.Atoi(v)
	variableData = data
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	result := &OkResult{}
	_ = json.NewEncoder(w).Encode(result)
}

func variable(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	log.Info().Int("data", variableData).Msg("variable response")
	result := &ExternalAdapterResponse{
		JobRunId: "",
		Data:     ExternalAdapterData{Result: variableData},
		Error:    nil,
	}
	_ = json.NewEncoder(w).Encode(result)
}
