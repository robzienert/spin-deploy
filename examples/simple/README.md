# simple example

`$ spindeploy test`

Deploys the "simple" application to the `test` target using the local
`../_templates/simple.yml` DCD pipeline template file.

```yaml
---
schema: "1"
targets:
  test:
    template: ../_templates/simple.yml
```

## notes 

* The `simple` app name is derived from the directory that the `.spin.yml` file
  is a child of. It will be processed to match [Frigga](https://github.com/netflix/frigga)
  naming conventions. This value can be set manually with the `metadata.app` YAML
  config.

