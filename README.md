# Albert
I'd be pretty pissed off too if you blasted me off into space and left me to suffocate.

[![Codacy Badge](https://api.codacy.com/project/badge/Grade/c57760a288af415083e6c1331143d368)](https://www.codacy.com/app/codingconcepts/albert?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=codingconcepts/albert&amp;utm_campaign=Badge_Grade)
[![Go Report Card](https://goreportcard.com/badge/github.com/codingconcepts/albert)](https://goreportcard.com/report/github.com/codingconcepts/albert)

## Disclaimer

Albert is the result of a 2-day project I undertook for my company's last hackathon.  As such, it shouldn't be considered reliable, finished or well-tested.  I'd love this to receive some input from the open source community, so if you're in need of a chaos monkey, or just want to make an existing one more chaotic, I'd love to hear from you!

## Overview

Albert is a simple, scalable and platform agnostic chaos monkey.

There are only two components, an orchestrator and an agent.  The orchestrator publishes commands to the agents and the agents execute kill requests against processes or Docker images etc. on the machines they're running on.  Both are hooked together using the wonderful NATS messaging platform.

### Dependencies

The only external dependency you'll need to run Albert is [NATS](http://nats.io/), a FOSS messaging system.  Either install the [binary](http://nats.io/download/nats-io/gnatsd/) or build from [source](https://github.com/nats-io/gnatsd).  Make sure to add the resulting `gnatsd` executable into your PATH.

### Install and run

``` bash
$ go get -u github.com/codingconcepts/albert
```

Out of habit, I'm using Rake, so if you've got it available (it's really easy to install), just run the following at the project root directory:

``` bash
$ rake run
```

This will start the gnatsd node, 2 agents and 1 orchestrator.  The agents are configured to run in `dummy` mode, just in case you're on a Windows machine and happen to have Notepad running :)

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

An agent is configured to kill a process, a machine, a Docker container (not tested), or a network interface (not developed).

As with the orchestrator, here's some example configuration (with the obvious bits omitted) to help you make sense of things:

``` json
{
    "application": "MarketDataAPI",
    "applicationType": "process",
    "identifier": "market.exe"
}
```

The contrived MarketDataAPI example might be running different agents on different OS's to do its job.  For example, there might be a WebAPI running in IIS serving customers and a Go web server running inside a Docker container on a Linux machine serving internal interconnected systems.  To achieve this, just run an agent on each machine with config to suit the job it's performing.  Note, they're both running as part of the MarketDataAPI application which needs to match up to the orchestrator's application name.

``` json
{
    "application": "MarketDataAPI",
    "applicationType": "docker",
    "identifier": "market-data"
}
```

Agents can also be configured perform a custom kill job with the following configuration:

``` json
{
    "application": "MarketDataAPI",
    "applicationType": "custom",
    "customInstructions": [
        "iisreset"
    ],
}
```

### Todo

- [ ] Behavioural unit tests
- [ ] Test cross-platform
- [x] SDK to turn any application into an agent
- [x] Agent overrides for kill procedure (in config)