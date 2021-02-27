package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

var (
	listendAddress = flag.String("web.ports", ":8087", "Ports ")
)


func startServer() {
	flag.Parse()
	router := mux.NewRouter()
	router.HandleFunc("/listAllFilesInDir", listAllFilesInDir).Methods("GET")
	router.HandleFunc("/getFileContentJSON/{filename}", getFileContentJSON).Methods("GET")
	router.HandleFunc("/getFileContentJSON/{filename}", postFileContentJSON).Methods("POST")
	router.HandleFunc("/getFileContentJSON/{filename}", deleteFileContentJSON).Methods("DELETE")
	srv := &http.Server{
		Handler: router,
		Addr: *listendAddress,
		WriteTimeout: 15*time.Second,
		ReadTimeout: 15*time.Second,
	}
	log.Fatal(srv.ListenAndServe())
}


func listAllFilesInDir(w http.ResponseWriter, r *http.Request) {
	fileList, err := readFilenamesFromDir(PATH)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, fileList)
}

func getFileContentJSON(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fileContentJSON, err := readFromFileJSON(vars["filename"])
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, fileContentJSON)
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	log.Println(message)
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}


func postFileContentJSON(w http.ResponseWriter, r *http.Request) {
	var data interface{}
	vars := mux.Vars(r)
	decoder := json.NewDecoder(r.Body)
	if err:= decoder.Decode(&data); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()
	if err := saveJSONtoFile(vars["filename"], data); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, "Success. File saved")
}

func deleteFileContentJSON(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	err := deleteFile(vars["filename"])
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, "Success. File deleted")
}