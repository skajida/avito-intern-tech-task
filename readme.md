# Segments service

## Build

```bash
podman build -f build/service-segs/db/Containerfile -t segments-db
podman build -f build/service-segs/app/Containerfile -t segments-app .
podman kube play deploy/segments-pod.yml
```

## Connect

```bash
curl -v localhost:8080/segs \
-H 'Content-Type: application/json' \
-d '{"seg_id": "AVITO_TRAINEE_CONNECT_SEGMENT"}'
```
