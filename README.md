# go-cd-demo

A small Go demo that shows a full CI/CD flow into Kubernetes.

## Services

- `api` on port `8080`
- `quote-service` on port `8081`

## Local run

```bash
go run ./cmd/quote-service
go run ./cmd/api
```

Run the services in separate terminals so the API can reach the quote service.

## Docker build

```bash
docker build -f Dockerfile.api -t go-cd-demo-api .
docker build -f Dockerfile.quote-service -t go-cd-demo-quote-service .
```

## Kubernetes

Apply the manifests:

```bash
kubectl apply -k k8s
```

The API is exposed through `NodePort` `30007`.

## CI/CD

The GitHub Actions workflow:

- runs Go tests
- builds both container images
- pushes them to GHCR on `main`
- deploys the latest images to Kubernetes using a base64-encoded kubeconfig stored as `KUBE_CONFIG_DATA`
