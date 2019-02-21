package main

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"
)

var containerMap map[string]string

func home(w http.ResponseWriter, r *http.Request){
	// title & description for homepage
	fmt.Fprintf(w, "Modular Finance Project: FaaS")
}

func dockerContainersMapCreate() map[string]string {
	// initiate map
	containerMap := make(map[string]string)

	// initiate dockerContainersMapCreate client
	cli, err := client.NewClientWithOpts(client.WithVersion("1.38"))
	if err != nil {
		panic(err)
	}

	// get information of docker containers
	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		panic(err)
	}

	// populate map with containers' details
	for _, container := range containers {

		serviceName := container.Labels["faas.name"]
		containerPort := container.Labels["faas.port"]

		// container name in network (not faas.name) eg: /faas_modular_fibonacci-service_1 (get rid of /)
		containerNameInNetwork := strings.Split(container.Names[0], "/")[1]

		// get information of docker network
		network, err := cli.NetworkInspect(context.Background(), "faas_modular_default", types.NetworkInspectOptions{})
		if err != nil {
			panic(err)
		}

		// search for container and its IP in the network
		for _, containerInNetwork := range network.Containers {

			// we don't include faas-gateway (with no service name)
			if serviceName != "" && containerNameInNetwork == containerInNetwork.Name {

				// CIDR to IP (x.x.x.x/y -> x.x.x.x)
				containerIP := strings.Split(containerInNetwork.IPv4Address, "/")[0]

				//populate map with information
				containerMap[serviceName] = containerIP+":"+containerPort
			}
		}
	}

	return containerMap
}

func dockerContainersMapCreateCron(seconds int) {
	for t := range time.NewTicker(time.Duration(seconds) * time.Second).C {
		containerMap = dockerContainersMapCreate()
		_ = t
	}
}

func proxyRequest(w http.ResponseWriter, r *http.Request) {
	// read service from url
	service := mux.Vars(r)["service"]

	// proxy request
	serviceExists := false

	for serviceName, ip := range containerMap {
		if serviceName == service {

			host := ip
			urlProxy := url.URL{Scheme: "http", Host: host,}

			proxy := httputil.NewSingleHostReverseProxy(&urlProxy)

			r.URL.Path = "/"
			proxy.ServeHTTP(w, r)

			serviceExists = true
			break
		}
	}

	// return 400 for bad service name
	if !serviceExists {
		fmt.Fprintf(w, "400 Bad Request: No such service")
	}
}

func main() {
	// create map (needed for the first 10 seconds goroutine won't run)
	containerMap = dockerContainersMapCreate()

	// update map every t seconds
	go dockerContainersMapCreateCron(10)

	// confirm gateway is running
	fmt.Println("Gateway On")

	// init router for building API
	r := mux.NewRouter()

	// route handles & endpoints
	r.HandleFunc("/", home).Methods("GET")
	r.HandleFunc("/lambda/{service}", proxyRequest).Methods("GET")

	// start server that listens to 8080. wrap with log for error checking
	log.Fatal(http.ListenAndServe(":8080", r))
}