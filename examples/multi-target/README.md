# multi-target example

`$ spindeploy prod`

Deploys the "myapp" application to the `prod` target using the local 
`../_templates/simple.yml` DCD pipeline template, 

```yaml
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
