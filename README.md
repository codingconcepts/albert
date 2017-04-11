# albert
I'd be pretty pissed off too if you blasted me into space and left me to suffocate.

## Disclaimer

Albert is the result of a 2-day project I undertook for my company's last hackathon.  As such, it shouldn't be considered reliable, finished or well-tested.  I'd love for this to get traction, so if you're in need of a chaos monkey, or just want to make an existing one more ... chaosy, I'd love to hear from you!

## Overview

Albert is a simple, scalable and platform agnostic chaos monkey.

There are only two components, both hooked together using the wonderful NATS messaging platform:

### Orchestrator

One or more orchestrators execute configurable tasks to take down zero or more nodes of zero or more applications.  Importantly, orchestrators are completely unaware of the agents executing the tasks.

To clear this up, here's an example which kills 50% of a make-believe `MarketDataApi` every 10 minutes on Mondays to Fridays, between 09:00 and 10:00.  It then begins killing 75% of the same service between 10:00 and 11:00, every 5 minutes.  Just in case you got cocky with how well your service was holding up against the first wave.  I've removed the obvious bits:

``` json
"applications": [
    {
        "name": "MarketDataAPI",
        "schedule": "10 9-10 * * 1-5",
        "percentage": 0.5
    },
    {
        "name": "MarketDataAPI",
        "schedule": "5 10-11 * * 1-5",
        "percentage": 0.75
    }
]
```

### Agent

An agent is configured to kill a process, a machine, a docker container (not tested), or a network interface (not developed).

As with the orchestrator, here's some example configuration (with the obvious bits omitted) to help you make sense of things:

``` json
{
    "application": "MarketDataAPI",
    "applicationType": "process",
    "identifier": "market.exe"
}
```

The contrived MarketDataAPI example might be running different agents on different OS's to do its job.  For example, there might be a WebAPI running in IIS serving customers and a Go web server running inside a Docker container on a Linux machine serving internal interconnected systems.  To acheive this, just run an agent on each machine with config to suit the job it's performing.  Note, they're both running as part of the MarketDataAPI application, this has to match up to the orchestrator's name.

``` json
{
    "application": "MarketDataAPI",
    "applicationType": "docker",
    "identifier": "market-data"
}
```

### Todo

- [ ] TLS for NATS connections (clients and cluster nodes)
- [ ] Test cross-machine
- [ ] Test with Linux machines
- [ ] SDK to turn any application into an agent
- [ ] SDKs for clients in other languages
- [ ] Agent overrides for kill procedure