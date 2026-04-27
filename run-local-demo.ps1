param(
  [switch]$Cleanup
)

$ErrorActionPreference = 'Stop'

function Write-Step {
  param([string]$Message)
  Write-Host ""
  Write-Host "==> $Message"
}

function Assert-Command {
  param([string]$Name)
  if (-not (Get-Command $Name -ErrorAction SilentlyContinue)) {
    throw "Required command '$Name' was not found in PATH."
  }
}

$repoRoot = Split-Path -Parent $MyInvocation.MyCommand.Path
Set-Location $repoRoot

Assert-Command go
Assert-Command docker
Assert-Command kubectl
Assert-Command kind

try {
  docker info | Out-Null
} catch {
  throw "Docker Desktop must be running before you start this script."
}

$clusterName = 'go-cd-demo'
$clusterContext = "kind-$clusterName"
$namespace = 'go-cd-demo'
$apiImage = 'ci-cd-k8s-api:test'
$quoteImage = 'ci-cd-k8s-quote:test'
$portForwardProcess = $null

try {
  Write-Step 'Running Go tests'
  $env:GOCACHE = Join-Path $repoRoot '.gocache'
  go test ./...

  # Build the two container images used by the demo.
  Write-Step 'Building Docker images'
  docker build -f Dockerfile.api -t $apiImage .
  docker build -f Dockerfile.quote-service -t $quoteImage .

  # Create a local Kubernetes cluster if one does not already exist.
  Write-Step 'Creating kind cluster'
  $existingCluster = kind get clusters | Where-Object { $_ -eq $clusterName }
  if (-not $existingCluster) {
    kind create cluster --name $clusterName
  } else {
    Write-Host "Cluster '$clusterName' already exists, reusing it."
  }

  kubectl config use-context $clusterContext | Out-Null

  # Load the freshly built images into kind so Kubernetes can run them locally.
  Write-Step 'Loading Docker images into kind'
  kind load docker-image $apiImage --name $clusterName
  kind load docker-image $quoteImage --name $clusterName

  # Apply the Kubernetes manifests and then point the deployments at the local images.
  Write-Step 'Deploying Kubernetes manifests'
  kubectl apply -k k8s --validate=false

  Write-Step 'Updating deployment images'
  kubectl -n $namespace set image deployment/api api=$apiImage
  kubectl -n $namespace set image deployment/quote-service quote-service=$quoteImage

  Write-Step 'Waiting for rollouts'
  kubectl -n $namespace rollout status deployment/api --timeout=180s
  kubectl -n $namespace rollout status deployment/quote-service --timeout=180s

  Write-Step 'Checking deployed resources'
  kubectl -n $namespace get pods
  kubectl -n $namespace get svc

  Write-Step 'Starting API smoke test'
  $portForwardProcess = Start-Process `
    -FilePath 'kubectl' `
    -WindowStyle Hidden `
    -PassThru `
    -ArgumentList @(
      '-n',
      $namespace,
      'port-forward',
      'service/api',
      '18080:80'
    )

  Start-Sleep -Seconds 5
  $response = $null
  for ($attempt = 1; $attempt -le 12; $attempt++) {
    try {
      $response = Invoke-WebRequest -UseBasicParsing -Uri 'http://127.0.0.1:18080/api'
      break
    } catch {
      if ($attempt -eq 12) {
        throw
      }
      Start-Sleep -Seconds 2
    }
  }
  Write-Host $response.Content

  Write-Step 'Demo completed successfully'
}
finally {
  if ($portForwardProcess -and -not $portForwardProcess.HasExited) {
    Stop-Process -Id $portForwardProcess.Id -Force
  }

  if ($Cleanup) {
    Write-Step 'Cleaning up kind cluster'
    kind delete cluster --name $clusterName
  }
}
