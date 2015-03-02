package main

import (
	"fmt"
	"net/url"
	"os"
	"sort"

	"github.com/fsouza/go-dockerclient"
)

var oldports []int
var b2dhost string

func main() {
	//todo get this information from boot2docker or shell environment
	endpoint := "tcp://192.168.59.103:2376"
	parsed, err := url.ParseRequestURI(endpoint)
	if err != nil {
		panic(err)
	}

	b2dhost = parsed.Host

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

	currentports := findcurrentports(containers)
	removeoldports(currentports)
	addnewports(currentports)
	oldports = currentports
}

func findcurrentports(containers []docker.APIContainers) []int {

	currentports := make([]int, 0)

	for _, container := range containers {
		for _, port := range container.Ports {
			if port.PublicPort != 0 {
				currentports = append(currentports, int(port.PublicPort))
			}
		}
	}

	sort.Ints(currentports)

	return currentports
}

func removeoldports(currentports []int) {
	for _, port := range oldports {
		i := sort.SearchInts(currentports, port)
		if i == len(currentports) || currentports[i] != port {
			fmt.Printf("Removing Port: %d\n", port)
			//todo remove proxy
		}
	}
}

func addnewports(currentports []int) {
	for _, port := range currentports {
		i := sort.SearchInts(oldports, port)
		if i == len(oldports) || oldports[i] != port {
			fmt.Printf("Adding Port: %d\n", port)
			//todo start a proxy
		}
	}
}
