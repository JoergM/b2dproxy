package main

import (
	"fmt"
	"log"
	"net"
	"net/url"
	"os"
	"path"

	"github.com/fsouza/go-dockerclient"
)

type Stoppable interface {
	Stop()
}

var proxiedPorts map[int]Stoppable
var b2dhost string

func main() {
	endpoint := os.Getenv("DOCKER_HOST")
	if endpoint == "" {
		fmt.Println("Could not find DOCKER_HOST. You have to set the environment.")
		os.Exit(2)
	}
	certpath := os.Getenv("DOCKER_CERT_PATH")
	if certpath == "" {
		fmt.Println("Could not find DOCKER_CERT_PATH. You have to set the environment.")
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

	//todo handle non tls connections (ENV DOCKER_TLS_VERIFY)
	client, err := docker.NewTLSClient(endpoint,
		path.Join(certpath, "cert.pem"),
		path.Join(certpath, "key.pem"),
		path.Join(certpath, "ca.pem"))
	if err != nil {
		panic(err)
	}

	//initial read of ports
	proxiedPorts = make(map[int]Stoppable)
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
			proxy.Stop()
			delete(proxiedPorts, port)
		}
	}
}

func addNewPorts(currentports map[int]string) {
	for port, ptype := range currentports {
		if proxiedPorts[port] == nil {
			if ptype == "udp" {
				forwarder := NewUDPForwarder(b2dhost, port)
				err := forwarder.start()
				if err == nil {
					proxiedPorts[port] = forwarder
				}
			} else {
				newPort, err := NewSinglePortProxy(b2dhost, port)
				if err == nil {
					proxiedPorts[port] = newPort
				}
			}
		}
	}
}
