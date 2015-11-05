# Intro

Thanks Ben. So I'm gonna go ahead and jump into a quick demo.

I will show you how to set up an Empire environment and how to build and deploy a simple docker image.

And hopefully Murphy is good to us today.

Ok. Is everybody can see ok? Is the font big enough?

# Story

So I'm gonna go ahead and tell you a little story. It's the story of a company you may have heard of before - ACME Inc. ACME Inc has been around for years, but what you may not know is that they just decided it was time to get online. More  importantly, they have asked us to help them achieve it with Empire.

So, after gathering some requirements it became pretty obvious that while they wanted a web presence, they weren't yet sure what they wanted to do with it. So we figured we will just make a placeholder website for them.

# Setup

Before using Empire, we need to set it up, of course. What we need to do is:

1. Install Empire within an ECS cluster.
2. And install the Empire CLI.

We provide a demo Cloudformation stack for Empire. It' a very easy and fast way to try it out.
IF you were to run Empire in production though, you would want to build your own stack.

After we've setup Empire, we just need to install the Empire CLI. So I'll go ahead and do that:

```console
$ go get -u github.com/remind101/emp
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

Great! We have an empty Empire environment all ready for Acme Inc's applications. **We can verify that by going to the ECS console**.

# Build the placeholder website

Ok. So I'm gonna go ahead and deploy the placeholder website for Acme, Inc. I'ts pretty simple, I already built it beforehand.

_Walk through www_

```console
$ docker build -t acmeinc/acme .
$ docker run -p 8080:8080 acmeinc/acme
$ open http://$(docker-machine ip default):8080
```

So let's go ahead and deploy this to Empire. The only pre-requisite is that we have the Docker image hosted on a Docker registry somewhere, so let's go ahead and push this docker image.

**CREATE THE APP BEFORE DEPLOYING**

```console
$ docker push acmeinc/acme
$ emp create www
$ emp domain-add acmeinc.com -a www
$ emp deploy acmeinc/acme
```

So that created an ECS service and attached an ELB to our new application.

_Find ECS service and ELB_

So, it looks like our application isn't running yet, we can see what processes are running using the `ps` command:

```console
$ emp ps -a www
```

Great! We're up and running. But it looks like I made a typo and called the company Acem Inc. Let me go ahead and fix that real quick.

_Fix typo, build and deploy_

```console
$ docker build -t acmeinc/acme .
$ docker push acmeinc/acme
$ emp deploy
```

When we deploy a new version, we'll see that Empire created a new release for us:

```console
$ emp releases -a www
```

When we deploy a new version of a web process, ECS will spin up new version of the app, wait for connections to drain from the existing version and then remove the old processes.

```console
$ emp ps -a www
```

Now that looks better.

# No. 1 on Hackernews!
