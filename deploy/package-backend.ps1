[CmdletBinding()]
param(
    [ValidateSet('amd64', 'arm64')]
    [string]$Architecture = 'amd64'
)

$ErrorActionPreference = 'Stop'
$repoRoot = Split-Path -Parent $PSScriptRoot
$serverDir = Join-Path $repoRoot 'server'
$packageDir = Join-Path $PSScriptRoot 'packages'
$stageDir = Join-Path $packageDir "ifs-api-linux-$Architecture"
$binaryPath = Join-Path $stageDir 'ifs-api'
$archivePath = Join-Path $packageDir "ifs-api-linux-$Architecture.zip"

if (-not (Get-Command go -ErrorAction SilentlyContinue)) {
    throw 'Go is required. Install the Go version specified by server/go.mod and retry.'
}

New-Item -ItemType Directory -Path $packageDir -Force | Out-Null
if (Test-Path $stageDir) { Remove-Item -LiteralPath $stageDir -Recurse -Force }
New-Item -ItemType Directory -Path $stageDir -Force | Out-Null

$previousCgo = $env:CGO_ENABLED
$previousGoos = $env:GOOS
$previousGoarch = $env:GOARCH
try {
    $env:CGO_ENABLED = '0'
    $env:GOOS = 'linux'
    $env:GOARCH = $Architecture
    Push-Location $serverDir
    try {
        go mod download
        if ($LASTEXITCODE -ne 0) { throw 'go mod download failed.' }

        go build -trimpath -ldflags='-s -w' -o $binaryPath .
        if ($LASTEXITCODE -ne 0) { throw 'Linux backend build failed.' }
    }
    finally {
        Pop-Location
    }
}
finally {
    $env:CGO_ENABLED = $previousCgo
    $env:GOOS = $previousGoos
    $env:GOARCH = $previousGoarch
}

if (Test-Path $archivePath) { Remove-Item -LiteralPath $archivePath -Force }
Compress-Archive -Path $binaryPath -DestinationPath $archivePath -CompressionLevel Optimal
Write-Host "Backend package created: $archivePath"
Write-Host 'The archive contains only the Linux executable. Keep the server config.yaml unchanged.'
