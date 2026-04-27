# CI-CD-K8s

This repository contains a complete CI/CD demo with Kubernetes deployment:

- `api` service on port `8080`
- `quote-service` on port `8081`
- Docker images for both services
- Kubernetes manifests under `k8s/`
- GitHub Actions workflow in `.github/workflows/ci-cd.yaml`

## 1. Create The Repository

Create a GitHub repository named `CI-CD-K8s` under your GitHub account `Subhash-Pal`.

If you are using the terminal:

```powershell
gh auth login
gh repo create Subhash-Pal/CI-CD-K8s --private --confirm
```

If the repository already exists, skip creation and use the existing remote URL:

```powershell
git remote set-url origin https://github.com/Subhash-Pal/CI-CD-K8s.git
```

## 2. Clone Or Prepare The Project

If you already have the local folder, open it:

```powershell
Set-Location "D:\training_golang\Module 9 (Hour 65-72) Docker, Kubernetes & CMCD Deep Dive\CI-CD-K8s"
```

If you are starting fresh, clone the repository:

```powershell
git clone https://github.com/Subhash-Pal/CI-CD-K8s.git
Set-Location CI-CD-K8s
```

## 3. Initialize Git

If the folder is not yet a git repo:

```powershell
git init
git branch -M main
git config --global --add safe.directory "D:/training_golang/Module 9 (Hour 65-72) Docker, Kubernetes & CMCD Deep Dive/CI-CD-K8s"
```

Set your identity if needed:

```powershell
git config user.name "Your Name"
git config user.email "you@example.com"
```

## 4. Add The Remote

```powershell
git remote add origin https://github.com/Subhash-Pal/CI-CD-K8s.git
```

If `origin` already exists:

```powershell
git remote set-url origin https://github.com/Subhash-Pal/CI-CD-K8s.git
```

## 5. Commit And Push

```powershell
git add .
git commit -m "Initial CI/CD K8s demo"
git push -u origin main
```

## 6. Test Locally

Before you check GitHub Actions, run the project locally to confirm the app and manifests work:

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

Then open a second terminal and run:

```powershell
curl http://127.0.0.1:18080/api
```

## 7. Check CI/CD In GitHub

Push to `main` and open the **Actions** tab in GitHub.

The workflow has 3 stages:

1. `test` runs `go test ./...` and `go build`
2. `deploy-kind` creates a temporary `kind` cluster in GitHub Actions
3. The workflow builds Docker images, loads them into `kind`, deploys the manifests, waits for rollout, and smoke-tests the API

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

## 8. Verify The Deployment

After the workflow succeeds, open the failed or completed workflow run and confirm these steps are green:

- `Test`
- `Build and push` or `deploy-kind` depending on the run type
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

## 9. What To Change For Your Own Demo

If you want to rename the repo or use another GitHub account:

- update the GitHub repo name in the clone and remote URLs
- update the owner name in the `gh repo create` command
- keep the workflow the same if you want the same CI/CD demo behavior
