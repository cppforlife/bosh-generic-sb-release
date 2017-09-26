#!/bin/bash

set -e

echo "-----> `date`: Upload stemcell"
bosh -n upload-stemcell "https://bosh.io/d/stemcells/bosh-warden-boshlite-ubuntu-trusty-go_agent?v=3363.31" \
  --sha1 03c6981d894e5c53554643ea4316c56a724ac8f7 \
  --name bosh-warden-boshlite-ubuntu-trusty-go_agent \
  --version 3363.31

echo "-----> `date`: Delete previous deployment"
bosh -n -d cf-mysql-broker delete-deployment --force
rm -f broker-creds.yml

echo "-----> `date`: Deploy"
( set -e; cd ./../..;
  bosh -n -d cf-mysql-broker deploy ./manifests/broker.yml -o ./manifests/dev.yml \
  -v director_ip=192.168.50.6 \
  -v director_client=admin \
  -v director_client_secret=$(bosh int ~/workspace/deployments/vbox/creds.yml --path /admin_password) \
  --var-file director_ssl.ca=<(bosh int ~/workspace/deployments/vbox/creds.yml --path /director_ssl/ca) \
  -v broker_name=cf-mysql-broker \
  -v service_name="CF MySQL" \
  -v service_description="CF MySQL" \
  --var-file si_manifest=<(wget -O- https://raw.githubusercontent.com/cloudfoundry/cf-mysql-deployment/develop/cf-mysql-deployment.yml|bosh int - -o examples/cf-mysql/fixes.yml|base64) \
  -v si_params=null \
  -v sb_manifest=null \
  -v sb_params=null \
  --vars-store ./examples/cf-mysql/creds.yml )

echo "-----> `date`: Use SB CLI"
export SB_BROKER_URL=http://$(bosh -d cf-mysql-broker is --column ips|head -1|tr -d '[:space:]'):8080
export SB_BROKER_USERNAME=broker
export SB_BROKER_PASSWORD=$(bosh int creds.yml --path /broker_password)

sb-cli services

echo "-----> `date`: Delete old service instances"
sb-cli delete-service-instance test1

echo "-----> `date`: Create service instances"
sb-cli create-service-instance test1

echo "-----> `date`: Check on service instances"
bosh -d service-instance-test1 manifest

echo "-----> `date`: Delete service instances"
sb-cli delete-service-instance test1

echo "-----> `date`: Delete deployments"
bosh -n -d cf-mysql-broker delete-deployment
rm -f broker-creds.yml

echo "-----> `date`: Done"
