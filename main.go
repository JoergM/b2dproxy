package main

import (
	"net"
	"net/url"

	"github.com/fsouza/go-dockerclient"
)

var proxiedPorts map[int]*SinglePortProxy
var b2dhost string

func main() {
	//todo get this information from boot2docker or shell environment
	endpoint := "tcp://192.168.59.103:2376"
	parsed, err := url.ParseRequestURI(endpoint)
	if err != nil {
		panic(err)
	}

	b2dhost, _, err = net.SplitHostPort(parsed.Host)
	if err != nil {
		panic(err)
	}

	client, err := docker.NewTLSClient(endpoint,
		"/Users/joerg/.boot2docker/certs/boot2docker-vm/cert.pem",
		"/Users/joerg/.boot2docker/certs/boot2docker-vm/key.pem",
		"/Users/joerg/.boot2docker/certs/boot2docker-vm/ca.pem")
	if err != nil {
		panic(err)
	}

	//initial read of ports
	proxiedPorts = make(map[int]*SinglePortProxy)
	updateports(client)

	//now listen for events to update ports
	events := make(chan *docker.APIEvents)
	err = client.AddEventListener(events)
	if err != nil {
		panic(err)
	}

	for {
		<-events //it does't matter what kind of update, just refresh the ports
		updateports(client)
	}

}

func updateports(client *docker.Client) {
	containers, err := client.ListContainers(docker.ListContainersOptions{All: true})
	if err != nil {
		panic(err)
	}

	currentports := findCurrentPorts(containers)
	removeOldPorts(currentports)
	addNewPorts(currentports)
}

func findCurrentPorts(containers []docker.APIContainers) map[int]bool {

	currentports := make(map[int]bool)

	for _, container := range containers {
		for _, port := range container.Ports {
			if port.PublicPort != 0 {
				currentports[int(port.PublicPort)] = true
			}
		}
	}

	return currentports
}

func removeOldPorts(currentports map[int]bool) {
	for port, proxy := range proxiedPorts {
		if !currentports[port] {
			proxy.stopListen()
			delete(proxiedPorts, port)
		}
	}
}

func addNewPorts(currentports map[int]bool) {
	for port, _ := range currentports {
		if proxiedPorts[port] == nil {
			proxiedPorts[port] = NewSinglePortProxy(b2dhost, port)
		}
	}
}
