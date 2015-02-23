package main

import (
	"fmt"
	"os"

	"github.com/fsouza/go-dockerclient"
)

func main() {
	//todo get this information from boot2docker or shell environment
	endpoint := "tcp://192.168.59.103:2376"
	client, err := docker.NewTLSClient(endpoint,
		"/Users/joerg/.boot2docker/certs/boot2docker-vm/cert.pem",
		"/Users/joerg/.boot2docker/certs/boot2docker-vm/key.pem",
		"/Users/joerg/.boot2docker/certs/boot2docker-vm/ca.pem")
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	containers, err := client.ListContainers(docker.ListContainersOptions{All: true})
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	fmt.Println("Current public ports:")
	for _, container := range containers {
		for _, port := range container.Ports {
			if port.PublicPort != 0 {
				fmt.Printf("%s : %d\n", container.Names[0], port.PublicPort)
			}
		}
	}

	//next listen for events regarding new ports

	//then open proxy
}
