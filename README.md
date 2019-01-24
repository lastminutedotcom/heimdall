# Heimdall
Heimdall, all-seeing and all-hearing

## General behaviour
**_Heimdall_** is a scheduler application based on cron and http calls to Cloudflare's API.

Setting your organizationId, **_Heimdall_** retrieve all the zones configured in you account and get the metrics to collect.
 
 

## Setting up the environment

#### Environment variables
```bash
export CLOUDFLARE_ORG_ID=<YOUR ORGANIZATION ID>\
export CLOUDFLARE_EMAIL=<YOUR EMAIL>\  
export CLOUDFLARE_TOKEN=<YOUR TOKEN>\
export CONFIG_PATH=<CONFIGURATION FILE PATH>
```
this variables are required for the correct working of *Heimdall*

#### Configuration file [configuration example](./test/config.json)
```json
{
  "collect_every_minutes" : "5", 
  "graphite_config": {
    "host": "<IP ADDRESS>",
    "port": 2113
  }
}
```
collect_every_minutes: mean that the metrics will be taken every interval defined. 

eg.: 5 meaning: every 5th minute for the last 5 minute

graphite_config: is the configuration required to connect to your graphite host in order to push the metrics.


