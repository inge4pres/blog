# Cloud-native applications: Operator-Framework


__Cloud-native__: sounds attractive right? What does it even mean? [Wikipedia](https://en.wikipedia.org/wiki/Cloud-native) has no page on it already so anyone can give its own definition... Here's mine:

_A Cloud-native application has only concern on the functionalities that it has to deliver as it is completely decoupled from the infrastructure it runs on_

So how can software delivery be cloud-native? Isn't software delivery supposed to "install" software onto some infrastructure? Well if your infrastructure provider _is_ cloud-native, you can transitively deliver software on it in a cloud-native way (counts of cloud-native is over 9000, so stopping here)!

Recently RedHat acquired CoreOS, bold move if you ask me. CoreOS since the M&A has been very quite until a week ago when the [operator-framework](https://github.com/operator-framework) was announced through a [blog post](https://coreos.com/blog/introducing-operator-framework); this is a huge step forward for everyone as this new toolkit will empower the average developer with the ability to run operators on Kubernetes and package their applications as extensions to the Kubernetes API. 

Never heard of [Custom Resource Definitions](https://kubernetes.io/docs/concepts/api-extension/custom-resources/#customresourcedefinitions)? You'd better get on track as this will be driving the next-gen wave of applications that will run _as part_ of the platform that delivers them, with the ability to automate their management and simplify dramatically their operations which will be tightly integrated with the cluster management itself.

And as usual I'm eager to try out this new toy and see what can be done with it: I want to build a cloud-native software delivery application that will enable CI/CD jobs to be running inside the cluster and managed by the same API server! Using a CRD for the "Pipeline" kind I can control the build/test/release flow of my application and moreover monitor the whole thing with the same tools with which I will monitor my application.

### Creating CD³

I decided to start a new project called [CD³ (cd kube)](https://github.com/inge4pres/cdkube): it will be a Continuous-* software that will run in Kubernetes and will be dedicated to deliver software _through_ Kubernetes. I found Weaveworks did something similar with [Flux](https://github.com/weaveworks/flux/) and since I don't want to reinvent the wheel I'll just try and replicate some functionality of running a continer in an operator.

First things first: I installed the [operator-sdk](https://github.com/operator-framework/operator-sdk) and I have the binary in my `$GOPATH/bin`.
Running

```
$GOPATH/bin/operator-sdk new cdkube --api-version=delivery.inge.4pr.es/v1alpha1 --kind=Pipeline
```

I am resulting in an auto-generated project with some scaffold code, as in [this commit](https://github.com/inge4pres/cdkube/commit/bd7a1cf8790cf02169057f759e40272993673181).
Next I add some details which I think are at the core of a pipeline: the repository with code and configurations, the image to build the software and what version/name give to the resulting application artifact. When adding new items to the Spec and Status of CRD I will modify the [pkg/apis/delivery/v1alpha1/types.go](https://github.com/inge4pres/cdkube/blob/master/pkg/apis/delivery/v1alpha1/types.go) file, then the guide suggests to run `operator-sdk generate k8s` to update the code, in fact the `zz_generated_deepcopy.go` is updated, but I don't see the [deploy/cr.yaml](https://github.com/inge4pres/cdkube/commit/5a17429d206475bcef85077d1f9aeb6fdda50357#diff-5cd8fbaa4e3fbb473c3e57c6fabca2f9) changed as it should so probably there's a bug... Moving forward I have my operator logic to be defined now: I add some check whether the building pod is succeded and build the operator

```
operator-sdk build inge4pres/cdkube:v0.0.1
docker push inge4pres/cdkube:v0.0.1
```

My operator is built, pushed to dockerhub and ready to be kicked in a k8s cluster, so I fire up one in GKE

```
gcloud beta container --project "inge4pres-gcp" clusters create "operator-framework-test" --zone "europe-west1-b" --machine-type "g1-small" --image-type "COS" --disk-size "25"

gcloud container clusters get-credentials operator-framework-test --zone europe-west1-b --project inge4pres-gcp
```

and when it's ready I can deploy the auto-generated resources to the cluster. 
An amazing result is that once the depoloyment is done I can create my CRD with the `kubectl` command and see my custom-type declared, running `kubectl create -f deploy/crd.yaml` I have a new object created in k8s

{{< highlight shell >}}
➜  ~ kubectl get pipeline
NAME            AGE
doing-nothing   16s
➜  ~ kubectl describe pipeline
Name:         doing-nothing
Namespace:    default
Labels:       <none>
Annotations:  <none>
API Version:  delivery.inge.4pr.es/v1alpha1
Kind:         Pipeline
Metadata:
  Cluster Name:
  Creation Timestamp:             2018-05-05T18:08:31Z
  Deletion Grace Period Seconds:  <nil>
  Deletion Timestamp:             <nil>
  Initializers:                   <nil>
  Resource Version:               1528
  Self Link:                      /apis/delivery.inge.4pr.es/v1alpha1/namespaces/default/pipelines/doing-nothing
  UID:                            5160a990-508f-11e8-8f06-42010a8401f8
Spec:
Spec:
  Build Arguments:
    building something...
    now I'm done
  Build Commands:
    echo
  Build Image:     busybox
  Repo:            http://github.com/inge4pres/just-a-test
  Target Name:     tesApp
  Target Version:  0.1
Events:            <none>
➜  ~
{{< /highlight >}}

When I deploy my operator and create the custom resource applying the manifest for a pipeline, the operator picks it up and starts a container with the name of the pipeline

{{< highlight shell >}}
➜  cdkube git:(master) ✗ kubectl apply -f deploy/cr.yaml
➜  cdkube git:(master) ✗ kubectl get po
NAME                      READY     STATUS    RESTARTS   AGE
cdkube-57bcf65584-sbvd6   1/1       Running   0          12s
➜  cdkube git:(master) ✗ kubectl apply -f deploy/cr.yaml
pipeline "doing-nothing" created
➜  cdkube git:(master) ✗ kubectl get po
NAME                      READY     STATUS              RESTARTS   AGE
cdkube-57bcf65584-sbvd6   1/1       Running             0          29s
testapp                   0/1       ContainerCreating   0          2s
➜  cdkube git:(master) ✗ kubectl logs cdkube-57bcf65584-sbvd6
time="2018-05-05T18:42:59Z" level=info msg="Go Version: go1.10.2"
time="2018-05-05T18:42:59Z" level=info msg="Go OS/Arch: linux/amd64"
time="2018-05-05T18:42:59Z" level=info msg="operator-sdk Version: 0.0.5+git"
time="2018-05-05T18:42:59Z" level=info msg="starting pipelines controller"
time="2018-05-05T18:43:21Z" level=info msg="build still running, status of builder container: Pending"
➜  cdkube git:(master) ✗ kubectl get po
NAME                      READY     STATUS      RESTARTS   AGE
cdkube-57bcf65584-sbvd6   1/1       Running     0          49s
testapp                   0/1       Completed   2          22s
{{< /highlight >}}

and if I get the `testapp` logs I can see the pipeline execution

{{< highlight shell >}}
➜  cdkube git:(master) ✗ kubectl logs testapp
building something... now I'm done
{{< /highlight >}}

Bonus: if the custom resource is deleted with `kubectl delete -f deploy/cr.yaml` the pod `testapp` gets terminated automatically! I don't know if this magic is done by Kubernetes or by the Operator Framework, but sure I love it as it will enable my resources to be managed with the k8s API just like any other.

I still need to tune the logic of handling the container as I didn't expect the restarts to happen, but hey this looks amazing! In a couple of hours I have an operator that can spawn containers when instructed to do so by custom manifests representing a pipeline!

Stay tuned 😎

