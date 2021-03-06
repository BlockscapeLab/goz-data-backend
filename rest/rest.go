package rest

import (
	"fmt"
	"log"
	"net/http"

	"github.com/BlockscapeLab/goz-data-backend/rest/types"
	"github.com/gorilla/mux"
)

//StartRestServer start the server at specified port and ip.
func StartRestServer(ip string, port int, dp types.DataProvider) {
	r := registerHandlers(dp)
	if err := http.ListenAndServe(fmt.Sprintf("%s:%d", ip, port), r); err != nil {
		log.Println("[Error] http server failed:", err)
	}

}

func registerHandlers(dp types.DataProvider) *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/scoreboard", createScoreboardHandler(dp)).Methods("GET", "OPTIONS")
	r.HandleFunc("/teams", createTeamsHandler(dp)).Methods("GET", "OPTIONS")
	r.HandleFunc("/teams/{chainID}/details", createTeamDetailHandler(dp)).Methods("GET", "OPTIONS")
	r.HandleFunc("/teams/{chainID}/chart", createTeamChartHandler(dp)).Methods("GET", "OPTIONS")
	r.HandleFunc("/syncstate", createSyncStateHandler(dp)).Methods("GET", "OPTIONS")
	return r
}

func createScoreboardHandler(dp types.DataProvider) http.HandlerFunc {
	h := func(res http.ResponseWriter, req *http.Request) {
		bz, err := dp.GetScoreboardJSON()
		sendResponse(bz, err, res)
	}

	return h
}

func createTeamDetailHandler(dp types.DataProvider) http.HandlerFunc {
	h := func(res http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		bz, err := dp.GetTeamDetailsJSON(vars["chainID"])
		sendResponse(bz, err, res)
	}
	return h
}

func createTeamChartHandler(dp types.DataProvider) http.HandlerFunc {
	h := func(res http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		bz, err := dp.GetTeamChartDataJSON(vars["chainID"])
		sendResponse(bz, err, res)
	}
	return h
}

func createTeamsHandler(dp types.DataProvider) http.HandlerFunc {
	h := func(res http.ResponseWriter, req *http.Request) {
		bz, err := dp.GetTeamsJSON()
		sendResponse(bz, err, res)
	}
	return h
}

func createSyncStateHandler(dp types.DataProvider) http.HandlerFunc {
	h := func(res http.ResponseWriter, req *http.Request) {
		bz, err := dp.GetSyncStateJSON()
		sendResponse(bz, err, res)
	}
	return h
}

func sendResponse(bz []byte, err error, res http.ResponseWriter) {
	res.Header().Set("Access-Control-Allow-Origin", "*")
	if err != nil {
		res.WriteHeader(types.GetStatusCode(err))
		res.Write([]byte(err.Error()))
	} else {
		res.WriteHeader(200)
		res.Write(bz)
	}
}
