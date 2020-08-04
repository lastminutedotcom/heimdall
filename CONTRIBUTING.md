# Welcome
We are very happy to accept contributions on Heimdall 😍

## Opening an issue

### Types of issue

When opening an issue please use one fo the issue templates:

* [Bug](.github/ISSUE_TEMPLATE/bug.md)
* [Question](.github/ISSUE_TEMPLATE/question.md)
* [Feature request](.github/ISSUE_TEMPLATE/feature_request.md)

The maintainers will pick up the issues and engage the proper conversation with you.

## Local Development

### Dependencies

* Go version 1.13+ is required
* Recent Docker version for manuel testing against Graphite

### Running tests

Just run 
```bash
make test
```

this will not require any external dependency.

### Running the application locally

You can spin up a Graphite+Grafana stack using [docker-compose](./docker-compose.yml).
You will need a docker network named `heimdall` setup beforehand (needed only once):

```bash
docker network create heimdall
``` 

Run the stack with the Graphite backend and Grafana:
```bash
docker-compose up
```

Build and run locally the application provided yuo have all the Cloudflare authentication details:

```bash
make build

CLOUDFLARE_ORG_ID=<ORGID> \
 CLOUDFLARE_EMAIL=me@domain.com CLOUDFLARE_TOKEN=<TOKEN> \
 CONFIG_PATH=./example/docker-compose.json ./heimdall
```

### Verifying metrics manually

You have 2 options:

* check Grafana
* exec a shell into the running `graphite` container and list the whisper files
```bash
docker exec -ti graphite /bin/sh
#->
ls -l /opt/graphite/storage/whisper/cloudflare/
```
