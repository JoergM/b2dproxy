# b2dproxy

A simple proxy that watches the Docker host in the boot2docker virtual machine and creates a Proxy at localhost for every open public port.

## Why?

It's just inconvinient when using boot2docker, that you can't refer to localhost but instead have to use the internal network address of the boot2docker VM.

This little tool makes using boot2docker just a little more convinient.

## Installation

Using standard go build process:

```
$> go get github.com/joergm/b2dproxy
$> go install github.com/joergm/b2dproxy
``` 

The proxy will now be installed at __$GOPATH/bin__.

## Usage

Your shell environment should be set up for use with boot2docker. The fastest way is:

```
$> $(boot2docker shellinit)
```

The proxy can be started using:

```
$> b2dproxy
```

This will start watching Docker and open a Proxy on your _localhost_ for every public port. Proxys shut down automatically when closed in docker.

### Example

Open a webserver using the docker training image, then test it via curl. 

```
$> docker run -d -p 5000:5000 training/webapp
$> curl http://localhost:5000
Hello World!
```

## Status

This tool is still under heavy development. You should not use it in any critical environment. There will probably be bugs. If you find any you are welcome to file an issue or even better a pull request.