# Operator Build/Deployment

## Deploying with helm

```bash
$> helm install android-farm-operator deploy/charts/android-farm-operator
```

### Helm Values


| Variable                  | Description                                     | Default                                                              |
|:-------------------------:|:------------------------------------------------|:--------------------------------------------------------------------:|
| `nameOverride`            | Name override for resources                     | `---`                                                                |
| `fullnameOverride`        | Full name override for resources                | `---`                                                                |
| `replicaCount`            | The number of operator instances to run         | `1`                                                                  |
| `image.repository`        | The docker registry to pull the operator from   | `quay.io/tinyzimmer/android-farm-operator` |
| `image.pullPolicy`        | The image pull policy                           | `IfNotPresent`                                                       |
| `imagePullSecrets`        | Image pull secrets for the operator             | `[]`                               |
| `serviceAccount.create`   | Whether to create service account and roles     | `true`                                                               |
| `serviceAccount.name`     | A name override for the service account         | `---`                                                                |
| `podSecurityContext`      | A pod security context to apply to the operator | `{}`                                                                 |
| `securityContext`         | Security context to apply to the operator pod   | `{}`                                                                 |
| `resources`               | Resource limits/requests for the operator       | `(see values.yaml)`                                                  |
| `nodeSelector`            | Node selector for the operator pod(s)           | `{}`                                                                 |
| `tolerations`             | Tolerations for the operator pod(s)             | `[]`                                                                 |
| `affinity`                | Affinity for the operator pod(s)                | `{}`                                                                 |
