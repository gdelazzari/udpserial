package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
)

func serveWebPanel(wg *sync.WaitGroup) {
	defer wg.Done()

	logger("webpanel", LogInfo, "starting web server")

	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/api/ports", handlerPortsIndex).Methods("GET")
	router.HandleFunc("/api/ports/{portName}", handlerPortGet).Methods("GET")
	router.HandleFunc("/api/ports", handlerPortPost).Methods("POST")
	router.HandleFunc("/api/ports/{portName}", handlerPortPut).Methods("PUT")
	router.HandleFunc("/api/ports/{portName}", handlerPortDelete).Methods("DELETE")
	router.HandleFunc("/api/statistics", handlerStatistics).Methods("GET")
	router.HandleFunc("/api/freePortNames", handlerFreePortNames).Methods("GET")
	router.HandleFunc("/api/baudrates", handlerBaudrates).Methods("GET")
	router.HandleFunc("/api/listenIPs", handlerListenIPs).Methods("GET")
	router.HandleFunc("/api/reloadConfigAndRestartThreads", handlerReloadConfigAndRestartThreads).Methods("GET")

	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./panel/dist")))

	http.Handle("/", router)

	logger("webpanel", LogFatal, http.ListenAndServe(":8080", nil))
}

func handlerPortsIndex(w http.ResponseWriter, r *http.Request) {
	changingConfig := readConfig(configFilename)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	var data []string
	for _, portConfig := range changingConfig.Ports {
		data = append(data, portConfig.Name)
	}
	json.NewEncoder(w).Encode(data)
}

func handlerPortGet(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	portName := vars["portName"]

	changingConfig := readConfig(configFilename)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	data, err := getPortConfig(changingConfig, portName)
	if err != nil {
		answerError(&w)
		return
	}

	json.NewEncoder(w).Encode(data)
}

func handlerPortPost(w http.ResponseWriter, r *http.Request) {
	var portConfig PortConfig

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		answerError(&w)
		return
	}
	if err := r.Body.Close(); err != nil {
		answerError(&w)
		return
	}
	if err := json.Unmarshal(body, &portConfig); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		json.NewEncoder(w).Encode(nil)
		return
	}

	logger("webpanel", LogInfo, "posted config for port "+portConfig.Name)

	changingConfig := readConfig(configFilename)

	_, getPortError := getPortConfig(changingConfig, portConfig.Name)

	if getPortError != nil {
		changingConfig.Ports = append(changingConfig.Ports, portConfig)

		changingConfig.saveToFile(configFilename)

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusCreated)

		json.NewEncoder(w).Encode(portConfig)
		return
	}

	answerError(&w)
	return
}

func handlerPortPut(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	portName := vars["portName"]

	var portConfig PortConfig

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		answerError(&w)
		return
	}
	if err := r.Body.Close(); err != nil {
		answerError(&w)
		return
	}
	if err := json.Unmarshal(body, &portConfig); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		json.NewEncoder(w).Encode(nil)
		return
	}

	logger("webpanel", LogInfo, "changing config for port "+portConfig.Name)

	changingConfig := readConfig(configFilename)

	if portConfig.Name != portName {
		answerError(&w)
		return
	}

	for i := 0; i < len(changingConfig.Ports); i++ {
		if changingConfig.Ports[i].Name == portName {
			changingConfig.Ports[i] = portConfig

			changingConfig.saveToFile(configFilename)

			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(http.StatusCreated)

			json.NewEncoder(w).Encode(portConfig)

			return
		}
	}

	answerError(&w)
	return
}

func handlerPortDelete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	portName := vars["portName"]

	logger("webpanel", LogInfo, "deleting config for port "+portName)

	changingConfig := readConfig(configFilename)

	foundIdx := -1
	foundPortConfig := PortConfig{}

	for i := 0; i < len(changingConfig.Ports); i++ {
		if changingConfig.Ports[i].Name == portName {
			foundPortConfig = changingConfig.Ports[i]
			foundIdx = i
			break
		}
	}

	if foundIdx >= 0 {
		changingConfig.Ports = append(changingConfig.Ports[:foundIdx], changingConfig.Ports[foundIdx+1:]...)

		changingConfig.saveToFile(configFilename)

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusCreated)

		json.NewEncoder(w).Encode(foundPortConfig)
		return
	}

	answerError(&w)
	return
}

func handlerFreePortNames(w http.ResponseWriter, r *http.Request) {
	changingConfig := readConfig(configFilename)

	var data []string

	for _, portDefinition := range definitions.PortDefinitions {
		_, getPortError := getPortConfig(changingConfig, portDefinition.PortName)
		if getPortError != nil {
			data = append(data, portDefinition.PortName)
		}
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(data)
}

func handlerBaudrates(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(definitions.BaudRates)
}

func handlerListenIPs(w http.ResponseWriter, r *http.Request) {
	type IPDescription struct {
		IP          string `json:"ip"`
		Description string `json:"description"`
	}
	var data []IPDescription

	data = append(data, IPDescription{"0.0.0.0", "any"})
	data = append(data, IPDescription{"localhost", "only local"})

	ifaces, err := net.Interfaces()
	if err == nil {
		for _, i := range ifaces {
			addrs, err := i.Addrs()
			if err == nil {
				for _, addr := range addrs {
					var ip net.IP
					switch v := addr.(type) {
					case *net.IPNet:
						ip = v.IP
					case *net.IPAddr:
						ip = v.IP
					}
					data = append(data, IPDescription{ip.String(), i.Name})
				}
			}
		}
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(data)
}

func handlerReloadConfigAndRestartThreads(w http.ResponseWriter, r *http.Request) {
	logger("webpanel", LogInfo, "requested threads restart")

	config = readConfig(configFilename)

	restartAllThreads()

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(nil)
}

func handlerStatistics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(getPublicStatistics(&statistics))
}

func answerError(w *http.ResponseWriter) {
	(*w).Header().Set("Content-Type", "application/json; charset=UTF-8")
	(*w).WriteHeader(http.StatusCreated)

	json.NewEncoder(*w).Encode(nil)
}
