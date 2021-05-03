---
date: "2019-11-01T15:56:42+01:00"
title: "CKA exam experience and preparation"
author: "inge4pres"
emoji: true
comment: true

categories:
  - tech
tags:
  - cka
  - kubernetes
---

Yes! Yesterday I received an awesome email stating that I cleared the [Certified Kubernetes Administrator](https://www.cncf.io/certification/cka/) exam! 😎
Here I want to report my experience in preparing and taking the exam, hopefully this info can help others Kubernetes practitioners get the certification too. 

### Preparation
I consider myself lucky because for the past two years I had the chance to use Kubernetes working at [lastminute.com](https://lmgroup.lastminute.com/); on top this on-the-job training I went through a lot of studying and practicing because the exam itself has a lot of content.
Having some expertise of working with Kubernetes is definitely a huge help, as this specific exam is tailored for cluster administrators.
You will need to prove a thorough understanding of the Kubernetes primitives, architecture and operational strategies.

Do not fear though! The [documentation](https://kubernetes.io/docs/home/) is extensive and very detailed on (most of) the topics that are needed to use and operate a cluster so here's my suggestions.

* ##### Exam structure
Understand the structure of the exam: 3 hours, 24 questions, 74% score needed to pass. 
In my experience the first 10 questions will be the easiest, so try to complete them during the first hour; the last 3-4 questions are going to be very long to read and execute so you will need more time in the end.
The exam software will not warn you if the question exercise is completed: *read carefully* the tasks that will mark the exercise complete.  

* ##### Knowledge
Read through the whole documentation carefully, at least twice.
Get familiar with the structure of the docs: concepts, tasks, reference: they are organized in such a way that the content is mixed up so know where to look for.
The following books are recommended read:
    - [Kubernetes up and running](https://www.oreilly.com/library/view/kubernetes-up-and/9781492046523/)
    - [Managing Kubernetes](https://www.oreilly.com/library/view/managing-kubernetes/9781492033905/)
    - [Kubernetes in action](https://www.oreilly.com/library/view/kubernetes-in-action/9781617293726/)


* ##### Practice
If you have free credits for a cloud provider use them to spin up a cluster to play with; you will need to know well how to interact with the cluster using `kubectl` CLI.
Try all the commands listed in the [kubectl cheat sheet](https://kubernetes.io/docs/reference/kubectl/cheatsheet/) and understand their effect on the cluster.
Explore the tutorials and tasks sections of the documentation and run them on the cluster to find caveats and errors in the docs.

Run through the amazing Kelsey Hightower [Kubernetes The Hard Way](https://github.com/kelseyhightower/kubernetes-the-hard-way) at least twice; spin up VMs on your PC using Vagrant and try building a cluster from scratch.
Push it even further trying to break the cluster in every possible way and observe the effects of failure: systemd logs, kubectl errors, etc...

### Taking the exam
I have to say the exam infrastructure is well done: in a single browser window you will have a left menu with the questions and a central section with a tmux terminal.
Before starting the exam a proctor will ask to confirm your identity and will validate the exam conditions are met: clean desk (no headset, no phone, no paper), empty and silent room.
During the exam you can have drinks but not food, and you can request a break at any time but the clock will continue to run during the break.
Some tips to meet the requirements.

* #####  Workstation
Remove everything from your exam desk except:

    - valid passport or national ID with picture
    - a glass of water or a juice that can help you stay hydrated 
    - the PC you will run the exam on, the optional external monitor and external mouse/keyboard

* #####  Connectivity
Ensure you have enough bandwidth for a streaming connection (10Mb should be enough), because you will have to share your desktop(s) and show your face through the webcam for the whole time.
Prefer WiFi over cable: before the exam can start you will have to pan your camera around to show the proctor the whole room, probably twice. If you (like me) used the notebook integrated webcam you will have to move the screen left and right and all across the room
If you're at home reboot your router one hour before the exam starts, just to be sure.

* ##### Manage time 
There's a timer in the upper left corner of the exam window: use it to check that you are completing enough questions on time.
Questions vary from 2 to 8 percent weight; on average you will have 7.5 minutes for evey question, but the last 3/4 are much longer to read and complete as they involve troubleshooting and actions to complete: you probably have to recover/install parts of a cluster.
Save time on the first questions to have more for the last ones; when switching from a question to another, check the score percentage of the question so you can evaluate if to do it or skip it with respect to the remaining time.

* ##### Use the resources
There's a staggering feature in the Kubernetes docs: *search*. When you are solving an exercise and need a reference or an example, search for it in the documentation!
As using the whole `*.kubernetes.io` domain is permitted during the exam you can consult examples, copy-paste YAML text and commands and even use snippets from the blog posts!
Use these resources as soon as you are confronted with a question for which `kubectl` commands are not enough, reading carefully a doc page or a blog post can get you out of trouble when something peculiar does not work as expected.

Another nice feature of the exam is a notepad you can open to take notes and copy-paste and edit content. I used it to take notes on the questions I was not able to complete fully and get back on them during the last 15 minutes of time.

### The overall exam experience
Although the exams rules are quite strict and I admit I was intimated reading the [CKA FAQ](https://www.cncf.io/certification/cka/faq/), setting up the exam room was easy and being able to run the test at home is a major plus compared to force you to go to a certification center.
The testing facility ran smoothly thanks to a Chrome extension that is required to run the exam and the instructions on how to complete the exercise were clear enough.

### Last but not least
If you purchased the exam from the Linux Foundation or CNCF Foundation you should have a free retake.
If you fail the first time don't worry! CKA is a difficult exam, very rich of content and all hands-on, definitely one of the most difficult certifications I've ever done. 
Try again after you get back studying and practicing on what you could not solve the first time 💪    
