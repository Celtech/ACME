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
- Managing our certifcate store and backing it up to S3
  - Issuing new certificates for new domains against the challenge server
  - Automatically renewing certificates that will expire soon
  - Archiving old certificates

## Usage

### Modes

Mode is controlled via the `ACME_ENV` environment variable. Possible values are:

- development
- production
- testing

### Pre-requistes

- A valid top level domain name (`/etc/hosts` entries will not work for this)
- A linux VM with docker installed on it and ports `80/tcp`, `443/tcp`, `9022/tcp` 
  exposed to the outside world

This project will not be something you can run locally on Docker for Mac, you will 
need linux based server that can properly expose ports to the outside world.

### Quick start

Git clone this repsoitroy to your VM. Then run the following command:

```shell
docker stack deploy -c docker-compose.testing.yml playground
```

### Configuration

| ENV Variable                | Default                                        | Required | Description                                                                                                                                                                                     |
|-----------------------------|------------------------------------------------|----------|-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| ACCOUNT_EMAIL               | NULL                                           | Yes      | Email to associate with a LetsEncrypt account to issue certificates from                                                                                                                        |
| DATA_PATH                   | /data                                          | No       | Base path to store our LetsEncrypt account data, and certificates, this will be created if it doesn't exist                                                                                     |
| HTTP_CHALLENGE_HOST         | 0.0.0.0                                        | No       | The host that our HTTP Challenge server will listen on, for docker `0.0.0.0` is correct.                                                                                                        |
| HTTP_CHALLENGE_PORT         | 80                                             | No       | The port our HTTP challenge server listens on, by default this is 80 and assumes you will use a reverse proxy                                                                                   |
| HTTP_CHALLENGE_PROXY_HEADER | NULL                                           | No       | When using a reverse proxy, you need to have a header that contains the original host. In most cases this will be `X-Forwarded-Host`                                                            |
| TLS_CHALLENGE_HOST          | 0.0.0.0                                        | No       | The host that our TLS Challenge server will listen on, for docker `0.0.0.0` is correct.                                                                                                         |
| TLS_CHALLENGE_PORT          | 443                                            | No       | The port our TLS challenge server listens on, by default this is 443 and assumes you will use a reverse proxy                                                                                   |
| ACME_HOST                   | https://acme-v02.api.letsencrypt.org/directory | No       | The ACME host that will issue certificates. The default is LetsEncrypt's ACME server and is subject to rate limiting. Use `https://acme-staging-v02.api.letsencrypt.org/directory` when testing |