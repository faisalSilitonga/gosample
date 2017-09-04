package ws

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

type WebService struct {
	Message     []string
	Status      string
	Config      []string
	Data        map[string]string
	Server_time string
}

func HelloWebService(w http.ResponseWriter, r *http.Request) {

	start := time.Now()

	if r.Method != "GET" {
		fmt.Fprintf(w, "hanya dapat diakses dengan method GET")
	}

	r.ParseForm()

	ws := WebService{
		Message: []string{fmt.Sprint("Hallo nama saya ", r.FormValue("name"))},
		Status:  "OK",
		Config:  nil,
		Data:    map[string]string{"name": r.FormValue("name")},
	}

	elapsed := time.Since(start)

	hasilJson, err := json.Marshal(map[string]interface{}{
		"message_status":      ws.Message,
		"status":              ws.Status,
		"config":              ws.Config,
		"data":                ws.Data,
		"server_process_time": strconv.FormatFloat(elapsed.Seconds(), 'f', -1, 64),
	})
	if err != nil {
		fmt.Fprintf(w, "ada error ketika parsing json: %s", err.Error())
	}

	fmt.Fprint(w, string(hasilJson))
	return
}
