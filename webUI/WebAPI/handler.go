package WebAPI

import (
	"ConfigServer/clientManage"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

func Handler(w http.ResponseWriter, r *http.Request, q url.Values) {
	apiPath := strings.TrimPrefix(r.URL.Path, "/api/")
	pathSegments := strings.Split(apiPath, "/")

	if r.Method == http.MethodPost {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Failed to read request body", http.StatusInternalServerError)
			return
		}

		defer func(Body io.ReadCloser) {
			_ = Body.Close()
		}(r.Body)

		_, err = fmt.Printf("POST request body: %s\n", body)

		if len(pathSegments) > 0 {
			switch pathSegments[0] {
			case "login":

				adminLogin(w, r, q.Get("username"), q.Get("password"))
			case "logout":

				adminLogout(w, r, q.Get("username"))
			case "cfg":
				command(w, r, &body)
			case "func":
				function(w, r, &body)
			case "cliInfo":
				cliInfo(w, r)
			default:

				http.Error(w, "API not found", http.StatusNotFound)
			}
		} else {
			http.Error(w, "Invalid API request", http.StatusBadRequest)
		}
	} else if r.Method == http.MethodGet {
		if len(pathSegments) > 0 {
			switch pathSegments[0] {
			case "cliInfo":
				cliInfo(w, r)
			default:
				http.Error(w, "API not found", http.StatusNotFound)
			}
		} else {
			http.Error(w, "Invalid API request", http.StatusBadRequest)
		}
	}
}

func cliInfo(w http.ResponseWriter, r *http.Request) {
	/*
	* return a list containing all host names if no arguments are included
	* otherwise, return detailed information
	 */

	idList := r.URL.Query()["id"]

	var rsp []byte

	if len(idList) > 0 {
		var clients = make([]clientManage.FriendlyClient, len(idList))
		for i, id := range idList {
			cli, err := clientManage.Get(id)
			if err != nil {
				clients[i] = clientManage.FriendlyClient{}
				println(err.Error())
			} else {
				clients[i] = cli.HumanFriendly()
			}
		}
		rsp, _ = json.Marshal(clients)
	} else {
		rsp, _ = json.Marshal(clientManage.AllHostName())
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err := w.Write(rsp)
	if err != nil {
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
		return
	}
}

func command(w http.ResponseWriter, r *http.Request, body *[]byte) {
	print(body)
	//bodyJson, _ := utils.JsonDecode(*body)
	//com := bodyJson["command"].(string)
	//console.Handler(com)
}

func adminLogin(w http.ResponseWriter, r *http.Request, userName string, password string) {
	if userName == "" || password == "" {
		http.Error(w, "Invalid username or password", http.StatusBadRequest)
	}

}

func adminLogout(w http.ResponseWriter, r *http.Request, userName string) {
	if userName == "" {
		http.Error(w, "Invalid username", http.StatusBadRequest)
	}

}
