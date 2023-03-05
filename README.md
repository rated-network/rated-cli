# Rated CLI

A tool to interact with the [Rated Network API](https://api.rated.network/docs).

```
Beacon Chain validator ratings from the command line.

Usage:
  rated-cli [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  watch       Watch performances of Ethereum validator keys

Flags:
      --config string   config file (default is $HOME/.rated-cli.yaml)
  -h, --help            help for rated-cli

Use "rated-cli [command] --help" for more information about a command.
```

## Configuration 

Example of a configuration file:

```
rated:
  # Endpoint to the Rated Network API
  apiEndpoint: https://api.rated.network

  # JWT Access token to the Rated Network API
  accessToken: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c

  # Optional: network to query, currently supports mainnet and prater
  network: mainnet

  # Interface to listen on to expose Prometheus metrics and health probes handlers
  listenOn: :8080
  
  # Configuration of the watcher service
  watcher:
    # Can be `day` or `hour`
    # Note: Days run on 24h intervals (225 epochs) from genesis. This means noon UTC to noon UTC the following day.
    # Note: "hours" are 57.6 minute intervals (9 epochs). A day contains exactly 25 such intervals.
    granularity: hour

    # List of validation keys to monitor
    validationKeys:
    - a26c2959d31d3f143eedd9f7377e33d842be77dffadaa4b12f8be2ccee725e797b2d7e93d5c55dbc1cb35177181401f4
    - a4a962fa6a13b57cb04c24821abc04dc219d4f3ffc65308bb483792e12402a983d68eb304a3af87e95ea9f9d8163ef46
    - a8e15d4cd973953be4f33d1107c571023a2f6330892eb087e36dd8c9386e2035d28e4d7fb57959df4862fce75908d12d
```

## Building & Running

```
make
bin/rated-cli --config config.yaml watch
```

## Features

### Grafana

If you are running it in an environment with Prometheus & Grafana, you can
install the ready-to-use dashboard in Grafana by:

- clicking on the `+` icon on the left menu and selecting `Import`,
- pasting the [JSON](grafana/dashboard.json) and clicking the `Load` button.

You will end up with a similar-looking dashboard:

![](grafana/rated-dashboard.png)

### Watcher

The watcher watches a list of validation keys and exposes their metrics to an
endpoint to be polled by Prometheus. This allows to implement blackbox
monitoring on an Ethereum Validation infrastructure without relying on internal
signals, independant of the Ethereum clients used.
