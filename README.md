**DEPRECATED: Use [tiller](https://github.com/robzienert/tiller) instead.**

# spin-deploy

A command-line deployment script that wraps the [Spinnaker](www.spinnaker.io)
DCD ([declarative continuous delivery](https://github.com/spinnaker/dcd-spec))
file format. With this deployment tool, you can define `.spin.yml` files in
your repositories to define deployment configurations alongside your app.

# install

Grab one of the [releases](https://github.com/robzienert/spin-deploy/releases).

```
export SPINNAKER_CLIENT_HOST=https://your-spinnaker-api
export SPINNAKER_CLIENT_X509CERTPATH=/path/to/your/x509/cert
export SPINNAKER_CLIENT_X509KEYPATH=/path/to/your/x509/key
```

# example

```
$ cd myapp-project
$ spindeploy loadtest
```

Starts the `loadtest` target pipeline defined in `.spin.yml`.

```yml
# .spin.yml
---
schema: "1"
metadata:
  app: myapp

targets:
  loadtest:
    # Pretend this has a bunch of deploy-related stages in it.
    template: s3://mybucket-dcd-templates/deploy.yml
    notifications:
      on.error:
      - type: slack
        channel: '#det'

  prod:
    # Pretend that this template exists, and that it has a "deploy" stage in it.
    template: s3://mybucket-dcd-templates/deploy.yml
    notifications:
      on.error:
      - type: slack
        channel: '#det'
    stages:
    - id: failoverWait
      type: wait
      dependsOn:
      - deploy
      name: Wait Before Failover
      config:
        waitTime: 3600
    - id: deployFailover
      type: pipeline
      dependsOn:
      - failoverWait
      name: Deploy to Failover
      config: {}
```
