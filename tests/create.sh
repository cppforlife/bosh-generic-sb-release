#!/bin/bash

set -e

# Assumes bosh and sb-cli are on path

export SB_BROKER_URL=http://10.244.0.8:8080
export SB_BROKER_USERNAME=zookeeper-broker
export SB_BROKER_PASSWORD=wyra3zehvvr9agtqsc8s

sb-cli delete-service-instance test1
sb-cli delete-service-instance test2

sb-cli create-service-instance test1
sb-cli create-service-instance test2

bosh -d service-instance_test1 manifest
bosh -d service-instance_test2 manifest

sb-cli delete-service-instance test1
sb-cli delete-service-instance test2
