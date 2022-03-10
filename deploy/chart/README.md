# Deploying TFUI with Helm

## Configuration

The default configuration does not make the service available over an Ingress, and does not configure persistent volumes.
It can be adapted to your needs with overrides of the default values, similar to the below YAML:

```yaml
server:
  ingress:
    enabled: true
    host: tfui.local
    ingressClassName: nginx
    annotations:
      ingress.kubernetes.io/ssl-redirect: "true"
      kubernetes.io/ingress.allow-http: "false"
      nginx.ingress.kubernetes.io/whitelist-source-range: "192.168.0.0/24"
    tls:
      - hosts:
          - tfui.local
        secretName: tls-local
  persistence:
    enabled: true
    storageClass: csi-sc-cinderplugin
    capacity: 100Mi
```
