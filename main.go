package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

func getQueryParams(req *http.Request) (string, int, int, bool) {
	termArr, ok := req.URL.Query()["q"]
	if !ok {
		return "", 0, 0, false
	}
	term := termArr[0]

	fromArr, ok := req.URL.Query()["from"]
	if !ok {
		return "", 0, 0, false
	}
	fromStr := fromArr[0]
	from, err := strconv.Atoi(fromStr)
	if err != nil {
		return "", 0, 0, false
	}

	sizeArr, ok := req.URL.Query()["size"]
	if !ok {
		return "", 0, 0, false
	}
	sizeStr := sizeArr[0]
	size, err := strconv.Atoi(sizeStr)
	if err != nil {
		return "", 0, 0, false
	}

	return term, from, size, true
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/populate", func(w http.ResponseWriter, req *http.Request) {
		numberArr, ok := req.URL.Query()["number"]
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Attach proper parameters"))
			return
		}
		numberStr := numberArr[0]
		number, err := strconv.Atoi(numberStr)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Attach proper parameters"))
			return
		}
		err = Populate(number)
		if err != nil || !ok {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
	})

	mux.HandleFunc("/search", func(w http.ResponseWriter, req *http.Request) {
		term, from, size, ok := getQueryParams(req)
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Attach proper parameters"))
			return
		}
		res, err := Search(term, from, size)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Error searching"))
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(res)
	})

	log.Fatal(http.ListenAndServe(":8000", mux))
}
