# Setup

Before using Empire. We have to do two things.

1. We need to provision an ECS cluster and install Empire.
2. We also need to install the Empire CLI.

As I said before, one of our primary design goals was to make Empire really easy to run. To install Empire, we provide a simple CloudFormation stack. In a production installation, you'd probably want to build your own CloudFormation stack, but the one we provide is suitable for a Demo and quickly testing out Empire and ECS.

After we've setup Empire, we just need to install the Empire CLI. So I'll go ahead and do that:

```console
$ go get -u github.com/remind101/emp
```

Now that we've installed it, we can run `emp` and get a list of the available commands:

```console
$ emp
```

The first thing we'll need to do is log in to our Empire environment. In order to do that, we'll need to tell the CLI where our Empire API is located, which I can get from the CloudFormation stack outputs:

```console
$ export EMPIRE_API_URL=http://<elb>
$ emp login
```

And now, let's run emp apps to see what apps we have running in our Empire environment.

```console
$ emp apps
```

Great, no apps, and we can verify that the ECS cluster has no running services:

```console
$ aws ecs list-clusters
$ export CLUSTER=<cluster>
$ aws ecs list-services --cluster $CLUSTER
```

So we have 1 service right now, which is for the Empire API and was setup when we created the cloudformation stack.

## Deploying your first app

Alright, let's go ahead and deploy our first app to Empire.

I'm going to go ahead and create a little Go application that has both a web process that returns some content, and then a worker process that does some background work.

I'm going to go ahead and paste in some code for our application:

```console
$ mkdir app
$ cd app
```

```go
package main

import (
  "flag"
  "fmt"
  "net/http"
  "os"
  "time"
)

func init() {
  flag.Usage = usage
}

func runWeb() {
  port := os.Getenv("PORT")
  fmt.Printf("Starting web server on %s\n", port)

  http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hello, %s\n", os.Getenv("NAME"))
  })

  fmt.Fprintln(os.Stdout, http.ListenAndServe(":"+port, nil))
}

func runWorker() {
  fmt.Println("Starting worker")
  for range time.Tick(1 * time.Second) {
    fmt.Println("Doing hard work...")
  }
}

func main() {
  flag.Parse()

  if len(os.Args) <= 1 {
    usage()
    os.Exit(1)
  }

  cmd := os.Args[1]
  switch cmd {
  case "web":
    runWeb()
  case "worker":
    runWorker()
  default:
    usage()
    os.Exit(1)
  }
}

func usage() {
  fmt.Fprintf(os.Stderr, `Usage of %s:

Commands:

  web  Run the web server
  worker Run the background worker
`, os.Args[0])
  flag.PrintDefaults()
}
```

This app is a simple Go program that has two subcommands; one that runs a web server, and another that runs a background process that does some background processing.

Now let's create the Procfile that will define the web and worker processes that compose this application:

```yaml
web: app web
worker: app worker
```

These map named processes to commands to run, so that we can easily scale these process separately via Empire.

And finaly, let's tie it all together by adding a Dockerfile to build the application.

```dockerfile
# Since this is a Go project, we'll use the base golang image.
FROM golang

# Empire expects web processes to bind to the PORT environment variable, we'll set a sensible default in
ENV PORT=80

# Empire will extract the Procfile from the WORKDIR
WORKDIR /go/src/github.com/remind101/reinvent-2015/demo/app

# Add the source to the image and compile the app.
ADD . /go/src/github.com/remind101/reinvent-2015/demo/app
RUN go install github.com/remind101/reinvent-2015/demo/app

# We'll run the API command, and print the help by default.
CMD ["app", "-h"]
```

And then we'll build it and verify that it works properly.

```console
$ docker build -t ejholmes/app .
$ docker run ejholmes/app
```

----

Empire makes it easy to deploy a Docker image as an application. To do this, we'll need to first push our image to a Docker registry, then deploy it using the Empire CLI.

Let's start by pushing our image to the Docker registry:

```console
$ docker push ejholmes/app:latest
```

The `emp deploy` command is our primary interface to deploying docker images from a docker registry.

```console
$ emp help deploy
```

Let's go ahead and create an app and deploy the latest tag of our Docker image.

```console
$ emp create app
$ emp domain-add app -a app
$ emp deploy ejholmes/app:latest
```

Empire pulled the Docker image, extracted the process types from the Procfile, and created the ECS services.

```console
$ aws ecs list-services --cluster $CLUSTER
$ export SERVICE=
```

We can see that we now have two additional ECS services for our app; one for the web process, and another for the worker process.

We can also see that Empire created an ELB attached to the web process.

```console
$ aws ecs describe-services --service $SERVICE --cluster $CLUSTER --query 'services[0].loadBalancers[0]'
$ export ELB_NAME=
$ aws elb describe-load-balancers --load-balancer-names $ELB_NAME --query 'LoadBalancerDescriptions[0].DNSName'
$ export ELB=
```

It'll take a minute or two for the ELB to begin resolving, so let's take a look at some of the other emp CLI commands that we can use.

We can us the `ps` subcommand to list the running and pending processes for the application.

```console
$ emp ps -a app
```

We can see that we're running 1 instance of the web process.

This just lists all of the ECS tasks for the ECS services that compose this application.

```console
$ aws ecs list-tasks --service $SERVICE --cluster $CLUSTER
```

Let's go ahead and scale up the worker process.

```console
$ emp scale worker=1 -a app
```

In a minute or two, we should see a new process for the worker process.

```console
$ emp ps -a app
```

Empire allows us to set environment variables on the application:

```console
$ emp set NAME=reInvent -a app
```

And also view the environment for the application:

```console
$ emp env -a app
```

Every time we deploy or update environment variables on the app, we create a new release:

```console
$ emp releases -a app
```

And Empire makes it easy to rollback to an existing release. Let's rollback to our first release:

```console
$ emp rollback v1 -a app
$ emp releases -a app
$ emp env -a app
```

Now, let's go ahead and open a browser for our application:

```console
$ open $ELB
```

Woops, it looks like we're missing that NAME environment variable because we rolled back. Let's go ahead and add it again:

```console
$ emp set NAME=reInvent -a app
```

After we reload the browser, you'll notice that nothing has changed. This is because ECS is now spinning up new versions of our application, and then waiting for connections to drain from the existing tasks. It will take about 1-2 minutes for the new version of our web process to start up.

```console
$ emp ps
```

And there we go.
