# b2dproxy

A simple proxy that watches the Docker host in the boot2docker virtual machine and creates a Proxy at localhost for every open public port.

## Why?

It's just inconvinient when using boot2docker, that you can't refer to localhost but instead have to use the internal network address of the boot2docker VM.

This tool makes using boot2docker just a little more convinient.

## Installation

Using standard go build process:

```
$> go get github.com/joergm/b2dproxy
``` 

The binary will now be installed at __$GOPATH/bin__.

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

If you want to use low ports (<1024) you have to start b2dproxy with root privileges:

```
$> sudo -E b2dproxy
```

The -E option keeps the environment variables.

### Example

Once you started b2dproxy try to open a webserver using the docker training image, then test it via curl. 

```
$> docker run -d -p 5000:5000 training/webapp
$> curl http://localhost:5000
Hello World!
```

## Status

This tool is still under development. You should not use it in any critical environment. There will probably be bugs. If you find any you are welcome to file an issue or even better a pull request.

### Known Open issues

- NON TLS Connection on localhost
- Passing on connection refused to client.
- UDP does not forward the origin IP address
