---
title: "Getting Started With Google Cloud Builder"
date: 2017-10-01T21:31:52+02:00
author: "inge4pres"
description: "GCB is a build tool to automate the creation of containerized applications"
tags:
  - gcp
  - cicd
  - development

categories:
  - tech

draft: false
---

One of the advantages of containerized applications is the standardization, some would say "write it once, runs everywhere" but that's another motto for another product. Anyway with a new packaging technology the same problems are faced: build reproducibility, or the necessity for people doing Ops to know they are going to deploy the same exact piece of code the Dev team used in their tests. So to address this issue the container image needs to be immutable: once it's built, it's not going to be changed, ever. And the same image will be used for testing, QA, beta, preview, presales-demo, whatever environment you need to deploy the app to.

#### Building it
Docker has been around for a few years now, it's mature and stable, but ask anyone using it if they'd allow images built on development workstation to run in production: not gonna happen! The artifact that will serve production traffic will pass through the CI/CD pipelines and pushed to the registry, this is the way of shipping containers. "That's easy", you'll say, "I'll stick my Dockerfile in the repo and let the build system do the magic!", but that's not the whole picture: there are lower layers to pull, tests to be run and _only after they succeed_ you can build the container. So who is going to maintain all of this configurations? Can we store them into the repo too? Yes! With GCB you can write a declarative multi-step workflow that will compile, test and package the code in a container; the container will get immediately pushed to the Google Container Registry too, so it's ready to be consumed (maybe Kubernetes on GKE?).

#### A simple application
There is quite a good number of examples in the [cloud builder repository](https://github.com/GoogleCloudPlatform/cloud-builders "GCB on Github") but I'd like to create a fresh one with Golang: a random number generator. The app will serve a random integer via HTTP, and you can see the code is [very straightforward](https;//github.com/inge4pres/blog/getting-started-with-google-cloud-builder/main.go).
Now I want to build a container to run into GCP with confidence so that every time I make a build I will have a new container image tagged with the version and ready to roll it out.

#### Using GCB
All builds in GCB happen in a container, right now the only engine supported is Docker but more are going to be added. You can leverage Google pre-baked builders or use your own images as builders, in the example I am using Google's `golang-project` image to compile and `docker` to build the final docker image and push it to registry. Note how some environment variables are injected to the container, like `PROJECT_ID` is your GCP running project as configured via
```
gcloud auth login
gcloud config set project your-gcp-project
```

##### _Side note for gophers_ #####
The `golang-project` image does some checks at setup to determine the workspace structure: there need to be a `./src` folder or `GOPATH` must be passed, or the simplest way is to [insert a comment next to `package main` in `main.go`](https://github.com/inge4pres/blog/blob/master/getting-started-with-google-cloud-builder/main.go#L1)

#### Running the build
It's as easy as executing
```
gcloud container builds submit --config cloudbuild.yaml .
```
in the folder holding `main.go` or the root of your project.

See it in action!

<script type="text/javascript" src="https://asciinema.org/a/US7J2mZlGcbCoqlW7zA66A1bh.js" id="asciicast-US7J2mZlGcbCoqlW7zA66A1bh" async></script>

As you can see the current folder (`.` as last parameter) is compressed and shipped to a Cloud Storage random location, then GCB starts the steps listed in the configuration YAML file, running step by step the containers with their arguments.

#### Conclusion
GCB is fast and very easy to use but for what I've been able to test is bound to GCP right now, so if you are willing to deliver a service from Google Cloud and your application is containerized or _"cloud native"_ you have a lightweight build system ready to go, but if you have a private registry or other integrations to do, GCB might still be a too small niche.

I'll make more tests and try to hack a bit GCB in the future so stay tuned!
