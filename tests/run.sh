#!/bin/bash

set -e # -x

echo "-----> `date`: Upload stemcell"
bosh -n upload-stemcell "https://bosh.io/d/stemcells/bosh-warden-boshlite-ubuntu-trusty-go_agent?v=3421.4" \
  --sha1 e7c440fc20bb4bea302d4bfdc2369367d1a3666e \
  --name bosh-warden-boshlite-ubuntu-trusty-go_agent \
  --version 3421.4

echo "-----> `date`: Delete previous deployment"
bosh -n -d zookeeper-broker delete-deployment --force
rm -f broker-creds.yml

echo "-----> `date`: Deploy"
( set -e; cd ./..;
  bosh -n -d zookeeper-broker deploy ./manifests/broker.yml -o ./manifests/dev.yml \
  -v director_ip=192.168.50.6 \
  -v director_client=admin \
  -v director_client_secret=$(bosh int ~/workspace/deployments/vbox/creds.yml --path /admin_password) \
  --var-file uaa_ssl_ca=<(bosh int ~/workspace/deployments/vbox/creds.yml --path /uaa_ssl/ca) \
  --var-file director_ssl.ca=<(bosh int ~/workspace/deployments/vbox/creds.yml --path /director_ssl/ca) \
  --var-file si_manifest=manifests/zookeeper.yml \
  --vars-store ./tests/creds.yml )

echo "-----> `date`: Use SB CLI"
export SB_BROKER_URL=http://$(bosh -d zookeeper-broker is --column ips|head -1|tr -d '[:space:]'):8080
export SB_BROKER_USERNAME=zookeeper-broker
export SB_BROKER_PASSWORD=wyra3zehvvr9agtqsc8s

sb-cli services

sb-cli delete-service-instance test1
sb-cli delete-service-instance test2

sb-cli create-service-instance test1
sb-cli create-service-instance test2

echo "-----> `date`: Check on service instances"
bosh -d service-instance_test1 manifest
bosh -d service-instance_test2 manifest

sb-cli delete-service-instance test1
sb-cli delete-service-instance test2

echo "-----> `date`: Delete deployments"
bosh -n -d zookeeper-broker delete-deployment
rm -f broker-creds.yml

echo "-----> `date`: Done"
