# GOATS

Goats is a performance and scale testing tool that takes advantage of the built in monitoring tooling in Kubernetes to measure and visualize the cpu and memory consumption of an app.
It uses hey to drive traffic to any application, as long as the following expectations are met:

1. The application container hosts a webserver when it is run without any additional run commands
2. The application continues to host a webserver on port 8000

An example app for the Go programming language is packaged in the apps/Go/example directory. A convenience script to build applications was provided, but that may be going away
in the future. Please *do not* put any secret values like license keys in your applications or containers. They will get that information from kubernetes in a protected way. This app only supports connecting Agents to the staging servers at the moment.

## Build GOATS

```go
cd runner
go build -o ../goats
cd ..
```

## GOATS Config

In order to run a job, goats requires a yaml config file.

```yaml
version: 1.0

runs:
  - name: example
    kubeconfig: /Users/emiliogarcia/.kube/config
    app:
      image: "quay.io/emiliogarcia_1/example-app:latest"
      environment-variables:
        EXAMPLE_VARIABLE: example value
    traffic-driver:
      startup-delay: 10   # time in seconds the traffic driver waits to send traffic to the application
      service-endpoint: /custom_events
      traffic:
        duration: 500   # time in seconds the traffic driver runs, make sure its at least 4 minutes or the performance metrics may not generate fully
        requests-per-second: 2    # number of requests sent to the server per second
        concurrent-requests: 2    # number of concurrent requests sent each time a request is sent
```

For version 1.0, I would recommend only ever creating one run. The multi-run capability is completely un-tested. The validation built in is really weak, so please respect the following guidelines. A run must have:

- a unique name
- a kubeconfig of the cluster you want to test on
- an app config
- a traffic-driver config

The app config tells the tool what *Public* container image to run as the application server. You can inject additional environment variables into it with the `environment-variables` map.

The traffic-driver config allows you to tune the traffic sent to the server. It accepts the following fields:

| field | type | definition |
| --- | --- | --- |
| startup-delay | uint | time in seconds that the traffic driver will wait to send traffic to the app |
| service-endpoint | string | the http endpoint that the traffic driver will send traffic to:  localhost:8000/\<service-endpoint\> |
| traffic.duration | uint | time in seconds that the driver will send traffic to the service endpoint |
| traffic.requests-per-second | uint | the number of requests the driver will make to the service endpoint per second |
| traffic.concurrent-requests | uint | the number of concurrent requests that are allowed to be sent to the server |


## Running

Once you have your `config.yaml` is ready, all you need to do is run the command:

```go
./goats -c config.yaml
```

It will create a kubernetes job for you, which is in the default namespace of your cluster.
