# spin-deploy

A command-line deployment script that wraps the [Spinnaker](www.spinnaker.io)
DCD ([declarative continuous delivery](https://github.com/spinnaker/dcd-spec))
file format. With this deployment tool, you can define `.spin.yml` files in
your repositories to define deployment configurations alongside your app.

# example

```yml
# .spin.yml
---
schema: "1"
metadata:
  app: myapp

targets:
  loadtest:
    # Look for an example of such a template in `../_templates/simple.yml`
    template: s3://mybucket-dcd-templates/simple.yml
    notifications:
      on.error:
      - type: slack
        channel: '#det'

  prod:
    template: s3://mybucket-dcd-templates/prod.yml
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
