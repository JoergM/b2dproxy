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

	//initial read of ports
	updateports(client)

	//now listen for events to update ports
	events := make(chan *docker.APIEvents)
	err = client.AddEventListener(events)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	for {
		<-events //it does't matter what kind of update, just refresh the ports
		updateports(client)
	}

}

func updateports(client *docker.Client) {
	containers, err := client.ListContainers(docker.ListContainersOptions{All: true})
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	fmt.Print("Forwarded ports:")
	for _, container := range containers {
		for _, port := range container.Ports {
			if port.PublicPort != 0 {
				//todo check if ports are changed
				fmt.Printf("%d ", port.PublicPort)
				//todo start a proxy
			}
		}
	}
	fmt.Print("\n")
}
