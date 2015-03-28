package main

import (
	"fmt"
	"log"
	"net"
	"net/url"
	"os"
	"os/user"
	"path"

	"github.com/fsouza/go-dockerclient"
)

var proxiedPorts map[int]*SinglePortProxy
var b2dhost string

func main() {
	endpoint := os.Getenv("DOCKER_HOST")
	if endpoint == "" {
		fmt.Println("Could not find DOCKER_HOST. You have to set the environment.")
		os.Exit(2)
	}
	parsed, err := url.ParseRequestURI(endpoint)
	if err != nil {
		panic(err)
	}

	b2dhost, _, err = net.SplitHostPort(parsed.Host)
	if err != nil {
		panic(err)
	}

	user, err := user.Current()
	if err != nil {
		panic(err)
	}

	//todo handle non tls connections (ENV DOCKER_TLS_VERIFY)
	client, err := docker.NewTLSClient(endpoint,
		path.Join(user.HomeDir, ".boot2docker/certs/boot2docker-vm/cert.pem"),
		path.Join(user.HomeDir, ".boot2docker/certs/boot2docker-vm/key.pem"),
		path.Join(user.HomeDir, ".boot2docker/certs/boot2docker-vm/ca.pem"))
	if err != nil {
		panic(err)
	}

	//initial read of ports
	proxiedPorts = make(map[int]*SinglePortProxy)
	log.Printf("Started b2dProxy. Now watching Docker for open ports ...\n")
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

func findCurrentPorts(containers []docker.APIContainers) map[int]string {

	currentports := make(map[int]string)

	for _, container := range containers {
		for _, port := range container.Ports {
			if port.PublicPort != 0 {
				currentports[int(port.PublicPort)] = port.Type
			}
		}
	}

	return currentports
}

func removeOldPorts(currentports map[int]string) {
	for port, proxy := range proxiedPorts {
		if currentports[port] == "" {
			proxy.stopListen()
			delete(proxiedPorts, port)
		}
	}
}

func addNewPorts(currentports map[int]string) {
	for port, ptype := range currentports {
		if proxiedPorts[port] == nil {
			newPort, err := NewSinglePortProxy(b2dhost, port, ptype)
			if err == nil {
				proxiedPorts[port] = newPort
			}
		}
	}
}
