---
date: "2018-01-21T12:31:34+01:00"
title: "Serverless on Kubernetes"
description: "An overview of the options to run functions as a service on k8s"
author: "inge4pres"

categories:
  - kubernetes
  - serverless

tags:
  - serverless
  - kubernetes
  - platform
  - development

draft: false
---

Kubernetes is the _de facto_ platform for running modern applications: its broad adoption in 2017 and the velocity of the project made it so and it's been accepted as the standard for many companies, from small to planet scale. It was impossible that such an extensible platform would be left out the serverless party, so here are the 4 main players offering FaaS to be run via k8s.

#### A premise
If you're new to serverless and FaaS and all the previous buzzwords sound like cacophony to your ears, I really recommend reading [this post](https://martinfowler.com/articles/serverless.html) and watching [this talk](https://www.youtube.com/watch?v=LAWjdZYrUgI). You could also notice how I put FaaS and serverless under the same hat here, this is just a personal opinion: historically I approached the topic using AWS Lambda, and I really tied the idea of writing functions and let someone else manage the infrastructure to the _serverless_ concept. Someone pointed out that the term was confusing because, you know, *there are actual servers* running those functions, so the acronym FaaS (Function as a Service) was born.

#### Why serverless on k8s
It seems like a natural evolution for distributed systems to be composed by smaller and smaller parts. When moving from SOA to microservices the size of the service was reduced to enable development of more fine-grained functionalities into smaller and more maintainable components; taken to the extreme, you can reduce a microservice to be dedicated to just one task or to be  made of just one function, that's where FaaS fits into. Kubernetes is a great activator for such modularity as it creates a very powerful abstraction over infrastructure, so when developing a function as a separate module of a distributed system you can scale both vertically and horizontally any building block, each one independently from another, or you could even let Kubernetes manage that (think [Horizontal Pod Autoscaler](https://kubernetes.io/docs/tasks/run-application/horizontal-pod-autoscale/)).

#### Four players offering FaaS on k8s

* [OpenFaaS](https://github.com/openfaas/faas)
* [Fission](https://github.com/fission/fission)
* [Kubeless](https://github.com/kubeless/kubeless)
* [Fn Project](https://github.com/fnproject/fn)

Now there might be others, but this 4 are the ones I mostly heard of in the last 6 months, so they must be the right ones 😁.

#### Comparison criteria
This is not a technical benchmark on the capabilities of this 4 frameworks: it's a _"look Ma, I can serverless on k8s"_ post where I try and highlight the pros and cons of adopting one or the other; the criteria will be ease of installation (client and server), languages support, cluster interoperability and developer experience, voted from 0 to 5, the higher the better.
I will use Kubernetes 1.8.6 that is, at the moment of writing, the latest available stable version.

The target function to deploy will be a super-serious analytics and business intelligence tool that will read the incoming HTTP request body and save it in a JSON document alongside with a timestamp. The JSON will be stored on a REDIS using a random UUIDv4 as key. All the code that will be deployed as functions is in [Github](https://github.com/inge4pres/blog/tree/master/serverless-on-kubernetes), while for installing the GCP cluster and REDIS I used the following

```
gcloud beta container --project "${GCP_PROJECT}" clusters create "serverless-k8s" \
  --zone "europe-west2-c" --username "admin" --cluster-version "1.8.6-gke.0" \
  --machine-type "g1-small" --image-type "COS" --disk-size "100" --num-nodes "3"

gcloud container clusters get-credentials serverless-k8s \
  --zone europe-west2-a --project "${GCP_PROJECT}"

helm init

helm install stable/redis
```

#### Fn Project
###### Features

* function configuration via YAML
* local development server via `fn start` and `fn run`
* uses DockerHub to store functions as containers
* web UI with function monitoring

###### Client installation: 4
The [installation instructions](https://github.com/fnproject/fn#install-cli-tool) are easy to read and execute, multiple platforms supported out of the box. User is required to set an environment variable with a DockerHub handle

```
export FN_REGISTRY=<DOCKERHUB_USERNAME>
```

###### Server installation: 3
A Helm chart is provided under [fn-helm](https://github.com/fnproject/fn-helm) but it's not immediately linked to the project's page. The installation requires the user to export an environment variable with the command

```
export FN_API_URL=http://$(kubectl get svc --namespace default fn-release-fn-api -o jsonpath='{.status.loadBalancer.ingress[0].ip}'):80
```

###### Language support: 5
Built-in support for Java, Ruby, Go, Python. Runs any docker container as a function.

###### Cluster interoperability: 2
It requires a LoadBalancer resource, so you won't be able to run it on `minikube` out of the box. It has MySQL and REDIS as dependency services and uses a [DaemonSet](https://kubernetes.io/docs/concepts/workloads/controllers/daemonset/) for the API controller, which might impact node's performance. No monitoring provided for the in-cluster components.

###### Developer experience: 3
Very extensive CLI interface. Functions are pipes: they should read Stdin and write to Stdout; some environment variables are injected to the running code to detect request URL and other configurations. I was able to complete my function deployment in roughly 1 hour after digging the docs a while to find out how to add custom configurations to the functions via environment variables.

#### OpenFaaS
###### Features

* sponsored by CNCF
* function grouping configuration via YAML ([stack file](https://github.com/openfaas/faas-cli#use-a-yaml-stack-file))
* public function repository
* Web UI with function monitoring
* runs any docker container as a function

###### Client installation: 2
The [CLI installation](https://github.com/openfaas/faas-cli#get-started-install-the-cli) is straight-forward for Linux and Mac users but it's not immediately available for Windows. I cannot find an easy way to set the cluster address to point the CLI to.

###### Server installation: 1
Helm chart provided under [faas-netes/helm](https://github.com/openfaas/faas-netes/blob/master/HELM.md) but it's failing the first time because of RBAC property not set and not rolling back, so I'm forced to delete and recreate the release.
Even when installation is completed I cannot connect to the FaaS gateway as the service NodePort 31112, and the LoadBalancer creation errors out with
```
Warning: kubectl apply should be used on resource created by either kubectl create --save-config or kubectl apply
The Service "gateway" is invalid: spec.ports[0].nodePort: Invalid value: 31112: provided port is already allocated
```

###### Language support: 5
Built-in support for NodeJS, Ruby, Go, Python, C#, generic Dockerfile. Runs any docker container as a function.

###### Cluster interoperability: 4
OpenFaas provides Prometheus monitoring with alerting out of the box, plus the architecture is very lean: just 2 pods that serve the API gateway and the function runner. Each function runs as a `Deployment` object and therefore can be scaled independently.

###### Developer experience: 2
After 2 hours trying to complete the setup on Kubernetes I'm still not able to run any function on my cluster; I run a port-froward to the OpenFaaS gateway `kubectl port-forward gateway-pod-id 31112:8080` so I can run `faas-cli deploy -f samples.yml --gateway http://127.0.0.1:31112` inside the just cloned `faas-cli` repository and see some function in action.

#### Fission
###### Features

* only 100ms function cold start (more on the topic [here](https://serverless.com/blog/keep-your-lambdas-warm/))
* natively built for Kubernetes

###### Client installation: 3
[Guide](http://fission.io/docs/0.4.0/installation/#install-the-fission-cli) suggests to download a binary distribution via `curl` and place the binary§ under `/usr/local/bin`, very straightforward for Linux and OSX. Windows support is via WSL or using a binary `fission.exe` with download link provided. Some environment variables need to be setup to point to the cluster, but the [instructions](http://fission.io/docs/0.4.0/installation/#cloud-setups) are very well written.

###### Server installation: 5
[Also guide](http://fission.io/docs/0.4.0/installation/#cloud-hosted-clusters-gke-aws-azure-etc): a single `helm` command installs all the components in a dedicated namespace `fission`.

###### Language support: 4
Built-in support for Linux binaries, Go, .NET, NodeJS, Perl, PHP 7, Python 3, Ruby as reported in the [concepts section](https://github.com/fission/fission#fission-concepts). Custom environment can be built and pushed to the cluster as containers.

###### Cluster interoperability: 2
No monitoring at all and no UI provided to verify the functions state of execution or list, CLI is the only source of truth I can get and it's not easy to understand the architecture.

###### Developer experience: 2
Setup is very straightforward but then the development looks cumbersome (at least for Go): environment variables [cannot be set](https://github.com/fission/fission/issues/356) yet, so this means hard-coded values in the code to connect to external services. Plus logging and function debugging is really hard, after 1 hour digging in documentation and trying to understand a cryptic `Internal server error (fission)` message, I am not able to run my Go function, and it's tough to tell why.

#### Kubeless
###### Features

* natively built for kubernetes
* web UI with function monitoring
* [serverless framework plugin](https://github.com/serverless/serverless-kubeless) available

###### Client installation: 5
Binary distribution available for Linux, OSX and Windows in [Github release page](https://github.com/kubeless/kubeless/releases). No fuss, no hassle: download and run.

###### Server installation: 3
[Documentation](https://github.com/kubeless/kubeless#installation) warns to create a namespace `kubeless` and use that. Same release page offers 3 options to install: Kubernetes with or without RBAC and Openshift. Applying the YAML with `kubectl apply -f kubeless-...` gets the server part up and running, but if I'd like to install in a different namespace I would need to change the whole file. Providing a Helm chart is the standard Kubernetes packaging way so why not having one?

###### Language support: 3
Currently only Python, NodeJS, Ruby and .NET Core are supported. [Custom runtimes](https://github.com/kubeless/kubeless/blob/master/docs/runtimes.md#custom-runtime-alpha) in the form of docker containers need to be built to run other languages, a feature in alpha which I'm forced to explore to run my Go app.

###### Cluster interoperability: 4
Monitoring provided via Prometheus integration; all backend services are run into dedicated namespace while functions are exposed using regular namespaces. It uses a StatefulSet to host Kafka and Zookeeper, they keep the functions state and there is a controller talking with Kubernetes API. I really like that it leverages Kubernetes-native primitives such as `ConfigMaps` and `Secrets` to manage function environment. It also uses a `CustomResourceDefinition` called `functions` so you can `kubectl get functions -o yaml`. What I don't like instead is the use of the cluster's etcd to store functions code when deploying from file.

###### Developer experience: 5
Everything worked great even when running an experimental feature such as custom environment: I was able to inject configuration via environment variables, get the function logs either via `kubeless` CLI or `kubectl` and debug my way out of the configuration error I put into my first image. Second deploy I did I was able to call my function and I validated the results connecting to REDIS afterwards.

### Wrapping up
There are lots of people investing in serverless right now, almost as many as there are for Kubernetes and the integrations between the two will bring a new exciting technology scenario for the next years. To my experience building this post none of the previously listed framework is ready for production usage, to be honest most of them are not even ready for the average developer weekend project usage. You need to know Kubernetes quite a bit to troubleshoot issues happening during deployment and/or execution of your functions and this makes the whole serverless idea crumble, as you'll be forced to dig into infrastructure details to have your code running.

If I'd have to bet on one of the project I tried so far, I'd do so on Kubeless; it's been definitely the smoothest setup  among all and the tight Kubernetes integration makes it a perfect candidate for community-driven growth.
