package v1

import "net/http"

func HealthCheck() http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte(`{"status":"OK"}`))
		if err != nil {
			return
		}
	}
}

func ReadyCheck() http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte(`{"status":"OK"}`))
		if err != nil {
			return
		}
	}
}
