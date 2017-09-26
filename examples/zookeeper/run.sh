#!/bin/bash

set -e

which sb-cli

echo "-----> `date`: Upload stemcell"
bosh -n upload-stemcell "https://bosh.io/d/stemcells/bosh-warden-boshlite-ubuntu-trusty-go_agent?v=3445.7" \
  --sha1 4c0670b318ca4c394e72037e05f49cc14d369636 \
  --name bosh-warden-boshlite-ubuntu-trusty-go_agent \
  --version 3445.7

echo "-----> `date`: Delete previous deployment"
bosh -n -d zookeeper-broker delete-deployment --force
rm -f creds.yml

echo "-----> `date`: Deploy"
( set -e; cd ./../..;
  bosh -n -d zookeeper-broker deploy ./manifests/broker.yml -o ./manifests/dev.yml \
  -v director_ip=192.168.50.6 \
  -v director_client=admin \
  -v director_client_secret=$(bosh int ~/workspace/deployments/vbox/creds.yml --path /admin_password) \
  --var-file director_ssl.ca=<(bosh int ~/workspace/deployments/vbox/creds.yml --path /director_ssl/ca) \
  -v broker_name=zookeeper-broker \
  -v srv_id=zookeeper \
  -v srv_name=Zookeeper \
  -v srv_description=Zookeeper \
  --var-file si_manifest=<(wget -O- https://raw.githubusercontent.com/cppforlife/zookeeper-release/master/manifests/zookeeper.yml|bosh int - -o examples/zookeeper/fixes.yml|base64) \
  --var-file si_params=<(cat examples/zookeeper/service-instance-params.yml|base64) \
  -v sb_manifest=null \
  -v sb_params=null \
  --vars-store ./examples/zookeeper/creds.yml )

echo "-----> `date`: Use SB CLI"
export SB_BROKER_URL=http://$(bosh -d zookeeper-broker is --column ips|head -1|tr -d '[:space:]'):8080
export SB_BROKER_USERNAME=broker
export SB_BROKER_PASSWORD=$(bosh int creds.yml --path /broker_password)

sb-cli services

echo "-----> `date`: Delete old service instances"
sb-cli delete-service-instance test1
sb-cli delete-service-instance test2

echo "-----> `date`: Create service instances"
sb-cli create-service-instance test1
sb-cli create-service-instance test2 -p nodes=3

echo "-----> `date`: Check on service instances"
bosh -d service-instance-test1 manifest
bosh -d service-instance-test2 manifest

if [ "x$(bosh -d service-instance-test1 vms|wc -l|tr -d '[:space:]')" != "x5" ]; then
  echo "Expected test1 service instance to have 5 nodes"
  exit 1
fi

if [ "x$(bosh -d service-instance-test2 vms|wc -l|tr -d '[:space:]')" != "x3" ]; then
  echo "Expected test2 service instance to have 3 nodes"
  exit 1
fi

echo "-----> `date`: Update service instances"
sb-cli update-service-instance test1 -p nodes=3

if [ "x$(bosh -d service-instance-test2 vms|wc -l|tr -d '[:space:]')" != "x3" ]; then
  echo "Expected test1 service instance to have 3 nodes after scale down"
  exit 1
fi

echo "-----> `date`: Delete service instances"
sb-cli delete-service-instance test1
sb-cli delete-service-instance test2

echo "-----> `date`: Delete deployments"
bosh -n -d zookeeper-broker delete-deployment
rm -f creds.yml

echo "-----> `date`: Done"
