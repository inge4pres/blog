#|/bin/bash
yum update -y
yum install -y python-pip golang jq
pip install --upgrade pip
pip install awscli
aws configure set preview.cloudfront true
mkdir $HOME/golang $HOME/.lego
export GOPATH=$HOME/golang
export PATH=$PATH:$GOPATH/bin
go get -u github.com/xenolf/lego

lego --email="your@email.com" --domains="domain1.tld" --domains="domain.san" --dns="route53" --path $HOME/.lego run
