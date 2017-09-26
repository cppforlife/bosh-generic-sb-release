#!/bin/bash

set -e

echo "-----> `date`: Testing zookeeper"
cd ../examples/zookeeper/
./run.sh

echo "-----> `date`: Testing cockroachdb"
cd ../cockroachdb/
./run.sh

echo "-----> `date`: Done"
