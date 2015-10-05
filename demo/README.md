So I'm going to go ahead and tell you the story of a company you may have heard of before - ACME Inc. ACME Inc has been around for years, but what you may not know is that just this year they decided it was time to get online. Even more  importantly, they figured it'd be best to come online for the first time today! True story.

That's right, they figured all sorts of cool stuff was going to be announced at re:Invent, so they contacted me and asked that I bring them into the cloud.

After gathering some requirements it became pretty obvious that while they wanted a web presence, they weren't yet sure what they wanted to do with it. Still, it's important that every company have a webpage, so I figured we might as well make a placeholder.

# Setup

Before using Empire. We have to do two things.

1. We need to provision an ECS cluster and install Empire.
2. We also need to install the Empire CLI.

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

Great! We have an empty Empire environment all ready for Acme Inc's applications. Let's go ahead and build a placeholder site for Acme Inc.

_Walk through www_

```console
$ docker build -t acmeinc/www .
$ docker run -p 8080:8080 acmeinc/www
$ open http://$(docker-machine ip default)
```

So let's go ahead and deploy this to Empire. The only pre-requisite is that we have the Docker image hosted on a Docker registry somewhere, so let's go ahead and push this docker image.

```console
$ docker push acmeinc/www
$ emp create www
$ emp domain-add acmeinc.com -a www
$ emp deploy acmeinc/www
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
$ docker build -t acmeinc/www .
$ docker push acmeinc/www
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

Wait a sec...turns out I just got a text from the CEO and he says they're going to provide Dropping Anvils as a Service and needs the app running right now for the Hacker News announcement.

Turns out I already have an application built, and we just need to deploy it and then expose the API through our nginx application.

_Walk through anvils app_

```console
$ docker build -t acmeinc/anvils .
$ docker run -p 8080:80 acmeinc/anvils
$ curl http://$(docker-machine ip default):8080/drop -d '{"Target": "Road Runner"}' -i
```

```console
$ docker push acmeinc/anvils:latest
$ emp deploy acmeinc/anvils
```

And now let's mount this API in our nginx application:

```nginx
    location /api {
        rewrite /api/(.*) /$1 break;
        proxy_pass "http://anvils.empire";
    }
```

```console
$ make && make push
$ emp deploy acmeinc/www
$ curl $ELB/api/drop -d '{"Target": "Road Runner"}' -i
```

Ok. So it's launch day and we hit #1 on HN. We're getting slammed and it doesn't help that the Road Runner is trying to DoS us.

```console
$ ./beepbeep $ELB
```

Let's go ahead and scale up our web process to account for the new load.

```console
$ emp scale web=5 -a www
```
