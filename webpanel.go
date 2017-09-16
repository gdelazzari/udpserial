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

	router.HandleFunc("/ports", handlerPortsIndex).Methods("GET")
	router.HandleFunc("/ports/{portName}", handlerPortGet).Methods("GET")
	router.HandleFunc("/ports", handlerPortPost).Methods("POST")
	router.HandleFunc("/ports/{portName}", handlerPortPut).Methods("PUT")
	router.HandleFunc("/ports/{portName}", handlerPortDelete).Methods("DELETE")
	router.HandleFunc("/freePortNames", handlerFreePortNames).Methods("GET")
	router.HandleFunc("/baudrates", handlerBaudrates).Methods("GET")
	router.HandleFunc("/listenIPs", handlerListenIPs).Methods("GET")
	router.HandleFunc("/restartThreads", handlerRestartThreads).Methods("GET")

	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./panel/")))

	http.Handle("/", router)

	logger("webpanel", LogFatal, http.ListenAndServe(":8080", nil))
}

func handlerPortsIndex(w http.ResponseWriter, r *http.Request) {
	logger("webpanel", LogInfo, "requested ports index")

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	var data []string
	for _, portConfig := range config.Ports {
		data = append(data, portConfig.Name)
	}
	json.NewEncoder(w).Encode(data)
}

func handlerPortGet(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	portName := vars["portName"]

	logger("webpanel", LogInfo, "requested config for port "+portName)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	data, err := getPortConfig(config, portName)
	if err != nil {
		json.NewEncoder(w).Encode(nil)
		return
	}

	json.NewEncoder(w).Encode(data)
}

func handlerPortPost(w http.ResponseWriter, r *http.Request) {
	var portConfig PortConfig

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		json.NewEncoder(w).Encode(nil)
		return
	}
	if err := r.Body.Close(); err != nil {
		json.NewEncoder(w).Encode(nil)
		return
	}
	if err := json.Unmarshal(body, &portConfig); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		json.NewEncoder(w).Encode(nil)
		return
	}

	logger("webpanel", LogInfo, "posted config for port "+portConfig.Name)

	_, getPortError := getPortConfig(config, portConfig.Name)

	if getPortError != nil {
		config.Ports = append(config.Ports, portConfig)

		onConfigChange()

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusCreated)

		json.NewEncoder(w).Encode(portConfig)
	}

	json.NewEncoder(w).Encode(nil)
	return
}

func handlerPortPut(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	portName := vars["portName"]

	var portConfig PortConfig

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		json.NewEncoder(w).Encode(nil)
		return
	}
	if err := r.Body.Close(); err != nil {
		json.NewEncoder(w).Encode(nil)
		return
	}
	if err := json.Unmarshal(body, &portConfig); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		json.NewEncoder(w).Encode(nil)
		return
	}
	if portConfig.Name != portName {
		json.NewEncoder(w).Encode(nil)
		return
	}

	for i := 0; i < len(config.Ports); i++ {
		if config.Ports[i].Name == portName {
			config.Ports[i] = portConfig
			onConfigChange()
			json.NewEncoder(w).Encode(portConfig)
			return
		}
	}

	json.NewEncoder(w).Encode(nil)
	return
}

func handlerPortDelete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	portName := vars["portName"]

	foundIdx := -1
	foundPortConfig := PortConfig{}

	for i := 0; i < len(config.Ports); i++ {
		if config.Ports[i].Name == portName {
			foundPortConfig = config.Ports[i]
			foundIdx = i
			break
		}
	}

	if foundIdx >= 0 {
		config.Ports = append(config.Ports[:foundIdx], config.Ports[foundIdx+1:]...)
		onConfigChange()
		json.NewEncoder(w).Encode(foundPortConfig)
		return
	}

	json.NewEncoder(w).Encode(nil)
	return
}

func handlerFreePortNames(w http.ResponseWriter, r *http.Request) {
	logger("webpanel", LogInfo, "requested ports index")

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	var data []string

	for _, portDefinition := range definitions.PortDefinitions {
		_, getPortError := getPortConfig(config, portDefinition.PortName)
		if getPortError != nil {
			data = append(data, portDefinition.PortName)
		}
	}

	json.NewEncoder(w).Encode(data)
}

func handlerBaudrates(w http.ResponseWriter, r *http.Request) {
	logger("webpanel", LogInfo, "requested ports index")

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(definitions.BaudRates)
}

func handlerListenIPs(w http.ResponseWriter, r *http.Request) {
	logger("webpanel", LogInfo, "requested ports index")

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

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

	json.NewEncoder(w).Encode(data)
}

func handlerRestartThreads(w http.ResponseWriter, r *http.Request) {
	logger("webpanel", LogInfo, "requested threads restart")

	restartAllThreads()

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(nil)
}

func onConfigChange() {
	config.saveToFile(configFilename)
}
