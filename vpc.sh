#!/usr/bin/env bash

#Create a pulumi stack

pulumi stack init aws-go
pulumi stack select aws-go
pulumi config set name aws-go
pulumi config set aws:region "us-east-1"

#removing default VPC
/usr/bin/remove_vpc

#Executing pulumi
pulumi up -y

