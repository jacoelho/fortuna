package api

import (
	"bytes"
	"net/http"
	"strconv"
)

type IntegersFunction func(min, max, count int) ([]int, error)

func IntegersHandler(f IntegersFunction) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		minKey := r.FormValue("min")
		maxKey := r.FormValue("max")
		countKey := r.FormValue("count")

		min, err := strconv.Atoi(minKey)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		max, err := strconv.Atoi(maxKey)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		count, err := strconv.Atoi(countKey)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if max < min {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		n, err := f(min, max, count)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
		}

		s, err := intToString(n, "\n")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
		}

		w.Write([]byte(s))
	}
}

func intToString(s []int, delimitator string) (string, error) {
	var b bytes.Buffer

	for _, v := range s {
		if _, err := b.WriteString(strconv.Itoa(v)); err != nil {
			return "", err
		}
		if _, err := b.WriteString(delimitator); err != nil {
			return "", err
		}
	}
	return b.String(), nil
}
