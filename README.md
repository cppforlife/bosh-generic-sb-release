# BOSH Generic Service Broker

Following service broker release can be deployed and configured with a BOSH manifest representing a service instance template. It also allows to configure parameters for service instance and binding creation.

Each SB API call maps to one or more BOSH commands:

- create service instance
  - runs `bosh deploy service-instance.yml -v ... -o ...`
- update service instance
  - runs `bosh deploy service-instance.yml -v ... -o ...`
- delete service instance
  - runs `bosh delete-deployment`
- create service binding
  - runs `bosh deploy service-binding.yml -v ... -o ...`
  - then `bosh run-errand create-*` capturing binding credentials from errand stdout
- delete service binding
  - runs `bosh deploy service-binding.yml -v ... -o ...`
  - then `bosh run-errand delete-*`

## Deploy

Configure and deploy service broker:

```
$ bosh -n -d cockroachdb-broker deploy ./manifests/broker.yml -o ./manifests/dev.yml \
  -v director_ip=192.168.50.6 \
  -v director_client=admin \
  -v director_client_secret=$(bosh int ~/vbox/creds.yml --path /admin_password) \
  --var-file director_ssl.ca=<(bosh int ~/vbox/creds.yml --path /director_ssl/ca) \
  -v broker_name=cockroachdb-broker \
  -v service_name=CockroachDB \
  -v service_description=CockroachDB \
  --var-file si_manifest=<(cat examples/cockroachdb/service-instance.yml|base64) \
  -v si_params=null \
  --var-file sb_manifest=<(cat examples/cockroachdb/service-binding.yml|base64) \
  --var-file sb_params=<(cat examples/cockroachdb/service-binding-params.yml|base64) \
  --vars-store ./examples/cockroachdb/creds.yml
```

## Usage

Use [`sb-cli`](https://github.com/cppforlife/sb-cli) to talk to issue SB commands:

```
$ export SB_BROKER_URL=http://$(bosh -d cockroachdb-broker is --column ips|head -1|tr -d '[:space:]'):8080
$ export SB_BROKER_USERNAME=broker
$ export SB_BROKER_PASSWORD=$(bosh int creds.yml --path /broker_password)

# List available services and their plans
$ sb-cli services

# Create a new instance with test1 ID
$ sb-cli create-service-instance test1

# Create a new binding with binding2 ID
$ sb-cli create-service-binding test1 --id=binding2 -p read_only=true
```

See `tests/run.sh` and `examples/` directory for more details.

## Service plans configuration format

...

## Parameters configuration format

During creation of service instances and bindings SB API allows to provide customizable parameters. `manifests/broker.yml` allows to provide `si_params` and `sb_params` variables with the following format:

```
- name: nodes        # <--- name of a param
  type: integer      # <--- type of a param
  ops:               # <--- list of operations to be applied to either
  - type: replace    #      service instance or binding manifest
    path: /instance_groups/name=zookeeper/instances
    value: ((value)) # <--- inserted user provided value

- name: include-status-errand
  type: boolean
  ops:
  - type: replace
    path: /instance_groups/name=zookeeper/jobs/-
    value:
    - name: status
      release: zookeeper
      properties: {}
```

Assuming above example is encoded with base64 and provided to `si_params` variable, user should be able to set `nodes` param to `3` to scale down Zookeeper service to 3 nodes:

```
$ sb-cli create-service-instance test1 -p nodes=3
```

Note that service instance create and update parameters are will be configured to be same because from the perspective of a SB create and update is the same operation to converge service instance to a particular state.

# Todo

- plan level properties...
- versions for development
- parameter types
- https broker url
- uaa integration
