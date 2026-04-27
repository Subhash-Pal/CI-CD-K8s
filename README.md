# CI-CD-K8s

This repository is a complete CI/CD + Kubernetes demo that anyone can recreate from scratch.

- `api` service on port `8080`
- `quote-service` service on port `8081`
- Docker images for both services
- Kubernetes manifests under `k8s/`
- GitHub Actions workflow in `.github/workflows/ci-cd.yaml`

## What You Need

Install these tools first:

- `Git`
- `Go 1.22+`
- `Docker Desktop`
- `kubectl`
- `kind`
- `gh` for GitHub repo creation from terminal

## Project Layout

- `cmd/api` - API service that calls the quote service
- `cmd/quote-service` - simple quote service
- `Dockerfile.api` - container build for the API
- `Dockerfile.quote-service` - container build for the quote service
- `k8s/` - Kubernetes manifests and namespace/config map
- `.github/workflows/ci-cd.yaml` - GitHub Actions CI/CD pipeline

## Recreate The Repo

### 1. Create The GitHub Repository

Create a GitHub repository named `CI-CD-K8s` under your account.

If you are using the terminal:

```powershell
gh auth login
gh repo create Subhash-Pal/CI-CD-K8s --private --confirm
```

If you are using the GitHub website, create a new private repository named `CI-CD-K8s` and do not add a README or `.gitignore`.

If the repository already exists, use the existing remote URL:

```powershell
git remote set-url origin https://github.com/Subhash-Pal/CI-CD-K8s.git
```

### 2. Clone Or Prepare The Project

If you already have the local folder, open it:

```powershell
Set-Location "D:\training_golang\Module 9 (Hour 65-72) Docker, Kubernetes & CMCD Deep Dive\CI-CD-K8s"
```

If you are starting fresh, clone the repository:

```powershell
git clone https://github.com/Subhash-Pal/CI-CD-K8s.git
Set-Location CI-CD-K8s
```

### 3. Initialize Git

If the folder is not yet a git repo:

```powershell
git init
git branch -M main
git config --global --add safe.directory "D:/training_golang/Module 9 (Hour 65-72) Docker, Kubernetes & CMCD Deep Dive/CI-CD-K8s"
```

Set your identity if needed:

```powershell
git config user.name "Subhash Pal"
git config user.email "you@example.com"
```

### 4. Add The Remote

```powershell
git remote add origin https://github.com/Subhash-Pal/CI-CD-K8s.git
```

If `origin` already exists:

```powershell
git remote set-url origin https://github.com/Subhash-Pal/CI-CD-K8s.git
```

### 5. Commit And Push

```powershell
git add .
git commit -m "Initial CI/CD K8s demo"
git push -u origin main
```

## Test Locally

Run these checks before you push or re-run Actions.

### One-command local demo

Run the full local demo from PowerShell:

```powershell
.\run-local-demo.ps1
```

If you want the script to delete the kind cluster at the end:

```powershell
.\run-local-demo.ps1 -Cleanup
```

Keep Docker Desktop running while the script is executing. If Docker stops or sleeps, `kind` and `kubectl` will fail because the cluster runs inside Docker.

The script also updates the Kubernetes deployments to use the locally built images before it waits for rollout, which avoids Docker Hub pull errors.

If the script fails because the cluster is in a bad state, clean it up and run it again:

```powershell
.\run-local-demo.ps1 -Cleanup
.\run-local-demo.ps1
```

If you want to run the steps manually instead of using the script:

```powershell
go test ./...
docker build -f Dockerfile.api -t ci-cd-k8s-api:test .
docker build -f Dockerfile.quote-service -t ci-cd-k8s-quote:test .
kind create cluster --name go-cd-demo
kind load docker-image ci-cd-k8s-api:test --name go-cd-demo
kind load docker-image ci-cd-k8s-quote:test --name go-cd-demo
kubectl apply -k k8s
kubectl -n go-cd-demo get all
kubectl -n go-cd-demo port-forward service/api 18080:80
curl http://127.0.0.1:18080/api
```

Run the app tests:

```powershell
$env:GOCACHE = Join-Path $PWD ".gocache"
go test ./...
```

Build both Docker images:

```powershell
docker build -f Dockerfile.api -t ci-cd-k8s-api:test .
docker build -f Dockerfile.quote-service -t ci-cd-k8s-quote:test .
```

Validate Kubernetes manifests:

```powershell
kubectl apply --dry-run=client -k k8s
```

Run the stack on a local cluster:

```powershell
kubectl apply -k k8s
kubectl -n go-cd-demo get all
```

Smoke test the API on your local cluster:

```powershell
kubectl -n go-cd-demo port-forward service/api 18080:80
```

Open a second terminal and run:

```powershell
curl http://127.0.0.1:18080/api
```

The script already runs the same smoke test for you and prints the API response.

## How CI/CD Works

Push to `main` and open the **Actions** tab in GitHub.

The workflow has 2 main jobs:

1. `test` runs `go test ./...` and `go build`
2. `deploy-kind` creates a temporary `kind` cluster in GitHub Actions, builds Docker images, loads them into `kind`, deploys the manifests, waits for rollout, and smoke-tests the API

The workflow will:

- run `go test ./...`
- build both Docker images
- create a temporary `kind` cluster in GitHub Actions
- load the images into the cluster
- deploy the Kubernetes manifests with `kubectl apply -k k8s`
- update the running images
- wait for rollout completion
- smoke test the API endpoint through a port-forward

No `KUBE_CONFIG_DATA` secret is needed for this demo because the cluster is created inside the workflow itself.

## Verify The Deployment

After the workflow succeeds, open the workflow run and confirm these steps are green:

- `Test`
- `deploy-kind`
- `Deploy manifests`
- `Wait for rollout`
- `Smoke test API`

If you are testing locally, verify the cluster with:

```powershell
kubectl config current-context
kubectl get nodes
kubectl -n go-cd-demo get pods
kubectl -n go-cd-demo get svc
```

Open the API service through NodePort `30007` if you are using your own local cluster:

```text
http://localhost:30007/
```

If you want to use your own remote Kubernetes cluster instead of `kind`, then you would create a `KUBE_CONFIG_DATA` secret and change the workflow to use that cluster. This repository does not need it for the default demo path.

## What To Change For Your Own Demo

If you want to rename the repo or use another GitHub account:

- update the GitHub repo name in the clone and remote URLs
- update the owner name in the `gh repo create` command
- keep the workflow the same if you want the same CI/CD demo behavior
