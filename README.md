# CI-CD-K8s

This repository contains a complete CI/CD demo with Kubernetes deployment:

- `api` service on port `8080`
- `quote-service` on port `8081`
- Docker images for both services
- Kubernetes manifests under `k8s/`
- GitHub Actions workflow in `.github/workflows/ci-cd.yaml`

## 1. Create The Repository

Create a GitHub repository named `CI-CD-K8s` under your account `Subhash-Pal`.

If you are using the terminal:

```powershell
gh auth login
gh repo create Subhash-Pal/CI-CD-K8s --private --confirm
```

If the repository already exists, skip creation and use the existing remote.

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

## 6. Local Validation

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

Run the stack on a cluster:

```powershell
kubectl apply -k k8s
kubectl -n go-cd-demo get all
```

## 7. Check CI/CD In GitHub

Push to `main` and open the **Actions** tab in GitHub.

The workflow will:

- run `go test ./...`
- build both Docker images
- push images to GHCR
- deploy the Kubernetes manifests with `kubectl apply -k k8s`
- update the running images
- wait for rollout completion

## 8. GitHub Secrets Needed

Add this secret in the repository settings:

- `KUBE_CONFIG_DATA` = base64 encoded kubeconfig for the target cluster

Example encoding command:

```powershell
[Convert]::ToBase64String([IO.File]::ReadAllBytes("$env:USERPROFILE\.kube\config"))
```

## 9. Verify The Deployment

After the workflow succeeds, check the cluster:

```powershell
kubectl config current-context
kubectl get nodes
kubectl -n go-cd-demo get pods
kubectl -n go-cd-demo get svc
```

Open the API service through NodePort `30007`:

```text
http://localhost:30007/
```
