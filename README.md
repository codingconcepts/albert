# Albert
I'd be pretty pissed off too if you blasted me off into space and left me to suffocate.

[![Build Status](https://travis-ci.org/codingconcepts/albert.svg?branch=master)](https://travis-ci.org/codingconcepts/albert)
[![Go Report Card](https://goreportcard.com/badge/github.com/codingconcepts/albert)](https://goreportcard.com/report/github.com/codingconcepts/albert)
[![Exago](https://api.exago.io:443/badge/rank/github.com/codingconcepts/albert)](https://exago.io/project/github.com/codingconcepts/albert)
[![Exago](https://api.exago.io:443/badge/cov/github.com/codingconcepts/albert)](https://exago.io/project/github.com/codingconcepts/albert)

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

Agents can be configured to kill anything you like (see "Kill Instructions")

As with the orchestrator, here's some example configuration (with the obvious bits omitted) to help you make sense of things:

``` json
{
    "application": "MarketDataAPI",
    "instructions": [ "taskkill", "/f", "/t", "/im", "market.exe" ]
}
```

This contrived MarketDataAPI example might be running different agents on different OS's to do its job.  For example, there might be a WebAPI running in IIS serving customers and a Go web server running inside a Docker container on a Linux machine serving internal interconnected systems.  To achieve this, just run an agent on each machine with config to suit the job it's performing.  Note, they're both running as part of the MarketDataAPI application which needs to match up to the orchestrator's application name.


### Kill Instructions

<table>
    <tr>
        <th>&nbsp;</th>
        <th>Windows</th>
        <th>Linux</th>
    </tr>
    <tr>
        <td>Machine</td>
        <td>"shutdown", "-t", "0", "-r", "-f"</td>
        <td>"reboot", "-f"</td>
    </tr>
    <tr>
        <td>Process</td>
        <td>"taskkill", "/f", "/t", "/im", "PROCNAME.exe"</td>
        <td>"kill", "-KILL", "'pgrep PROCNAME'"</td>
    </tr>
    <tr>
        <td>Docker image</td>
        <td colspan="2">$(docker stop $(docker ps -a -q --filter ancestor=IMAGENAME --format="{{.ID}}"))</td>
    </tr>
</table>

### Todo

- [x] Behavioural unit tests
- [ ] Test cross-platform
- [x] SDK to turn any application into an agent
- [x] Agent overrides for kill procedure (in config)
- [ ] Test that the channel pipe goroutines in the agent's natsProcessor end
- [ ] Remove fmt.Scanln() because that's just filthy