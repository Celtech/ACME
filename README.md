# ACME Server

Thin LetsEncrypt ACME client, challenge server, and management API written in go.

## Why?

Working with SSL inside of Docker swarm, with HAProxy, powering it can be a complex
thing. If you add a ton of domains, and dynamic custom domains into the mix, it
seems almost impossible. Current solutions rely on bash scripts and really have
no good way of handling dynamic domains. To solve this, we've created our own ACME
client which will be used inside a docker container. This client has 3
responsibilities:

- Exposing a API that can be hit to issue a new certificate for a domain pointed at
  our ingress servers
- Creating a temporary challenge server per certificate request to validate the domain
- Managing LetsEncrypt user accounts as to not get rate limited
- Managing our certificate store and backing it up to S3
  - Issuing new certificates for new domains against the challenge server
  - Automatically renewing certificates that will expire soon
  - Archiving old certificates

## Usage

### Pre-requisites

- A valid top level domain name (`/etc/hosts` entries will not work for this)
- A linux VM with docker installed on it and ports `80/tcp`, `443/tcp`, `9022/tcp`
  exposed to the outside world

This project will not be something you can run locally on Docker for Mac for a true
end-to-end test, you will need linux based server that can properly expose ports to
the outside world.

### Quick start

1. [Install Taskfile if you haven't already](https://taskfile.dev/installation/)
2. Git clone this repsoitroy to your VM:
3. Start the stack (`task docker:start`)

### Configuration

Most of the configuration is handled through yaml files in the config folder.
By default we include 3 config files:

- testing.yaml
- development.yaml
- production.yaml

Each of these config files correlate to the mode. This allows you to add more
config files than the defaults included, just ensure you set the mode to the
name you gave the file. All files other than `production.yaml` will start the
server in debug mode with additional logging.

### Modes

Mode is controlled via the `ACME_ENV` environment variable. Possible values are:

- development
- production
- testing
- your-custom-modes-here

Custom modes require a corresponding `.yaml` file. For example if you set `ACME_ENV` to `staging`, you would need the
corresponding config file `config/staging.yaml`.

### The CLI

```text
$ go run . --help
Thin Let's Encrypt ACME client and challenge server written in go.

Usage:
  acme [command]

Available Commands:
  add         Adds a new API authorized user to the database
  completion  Generate the autocompletion script for the specified shell
  hash        Returns a hashed version of a plaintext password
  help        Help about any command
  start       Start the web server

Flags:
  -h, --help      help for acme
  -v, --version   version for acme

Use "acme [command] --help" for more information about a command.
```

## Attributions

<p float="left">
<img src="https://plugins.jetbrains.com/static/versions/22143/jetbrains-simple.svg" alt="drawing" width="100"/>
<img src="https://github.com/go-acme/lego/raw/master/docs/static/images/lego-logo.min.svg" alt="Lego ACME Logo" width="200"/>
</p>

This project was made possible by the wonderful developers over at [lego
acme](https://github.com/go-acme/lego). Lego acme was used as the base
ACME client with a few small tweaks to make it easily compatible with the web
which power interfacing with Lets Encrypt.

Also, a special thank you to [JetBrains](https://jb.gg/OpenSourceSupport) for providing 
a free license to GoLand to support the development of this project.
