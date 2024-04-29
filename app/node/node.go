package node

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"piprim.net/gbcl/app/config"
	"piprim.net/gbcl/app/db"
	"piprim.net/gbcl/app/node/jsonhandler"
)

const (
	DefaultPort = uint16(8180)
)

type httpJSONResp struct {
	Data  any   `json:"data"`
	Error error `json:"error"`
}

func Run() error {
	state, err := db.NewStateFromDisk()
	if err != nil {
		return errors.Wrap(err, "Initializing state from runing node")
	}

	jsonhandler.State = state
	defer state.Close()

	http.HandleFunc("/balances/list", jsonHandlerFunc(jsonhandler.ListBalances))
	http.HandleFunc("/tx/add", jsonHandlerFunc(jsonhandler.TxAdd))

	port := config.Get().GetPort()
	log.Debug().Msg(fmt.Sprintf("Launching GBCL node and its HTTP API on http://localhost:%d", port))

	//nolint:gosec // For simplicity
	err = http.ListenAndServe(fmt.Sprintf(":%d", port), nil)

	return errors.Wrap(err, "Server error")
}

func jsonHandlerFunc[Ti any, To any](handler func(*Ti) (*To, error)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		logger := log.Debug()
		requestID := uuid.New().String()
		logger.Str("http request", r.RequestURI)

		input := new(Ti)
		err := setFromRequestBody(r, input)
		if err != nil {
			logger.Msg("can not handle http request")
			jsonHandle[any](nil, err, w)

			return
		}

		logger.Msg("Handling http request " + requestID)
		jsonOut, err := handler(input)

		jsonHandle(jsonOut, err, w)

		loggerO := log.Debug()
		loggerO.Any("data", jsonOut)
		loggerO.Any("error", err)
		loggerO.Msg("Result of http request" + requestID)
	}
}

func setFromRequestBody(r *http.Request, dest any) error {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return errors.New(fmt.Sprintf("unable to read request body : %s", err.Error()))
	}
	defer r.Body.Close()

	if len(body) > 0 {
		err = json.Unmarshal(body, dest)
		if err != nil {
			return errors.New("unable to unmarshal request body : " + err.Error())
		}
	}

	return nil
}

func jsonHandle[To any](data To, err error, w http.ResponseWriter) {
	resp, _ := json.Marshal(httpJSONResp{
		Data:  data,
		Error: err,
	})
	w.Header().Set("Content-Type", "application/json")

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
	}

	_, _ = w.Write(resp)
}
