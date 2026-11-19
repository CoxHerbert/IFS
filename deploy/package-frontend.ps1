[CmdletBinding()]
param(
    [ValidateSet('admin', 'portal', 'ledger', 'all')]
    [string]$Target = 'all'
)

$ErrorActionPreference = 'Stop'
$repoRoot = Split-Path -Parent $PSScriptRoot

function Build-Frontend([string]$Name) {
    $projectDir = Join-Path $repoRoot "ifs-$Name"
    Push-Location $projectDir
    try {
        if (-not (Get-Command pnpm -ErrorAction SilentlyContinue)) {
            throw 'pnpm is required. Install pnpm and retry.'
        }
        pnpm install --frozen-lockfile
        if ($LASTEXITCODE -ne 0) { throw "pnpm install failed for ifs-$Name" }

        if ($Name -eq 'admin') {
            pnpm run build:prod
        } elseif ($Name -eq 'ledger') {
            pnpm run build:h5
        } else {
            pnpm run build
        }
        if ($LASTEXITCODE -ne 0) { throw "Build failed for ifs-$Name" }

        $distDir = Join-Path $projectDir 'dist'
        if (-not (Test-Path (Join-Path $distDir 'index.html'))) {
            throw "Build output is missing: $distDir/index.html"
        }

        $archivePath = Join-Path $projectDir "ifs-$Name.zip"
        if (Test-Path $archivePath) { Remove-Item -LiteralPath $archivePath -Force }
        Compress-Archive -Path (Join-Path $distDir '*') -DestinationPath $archivePath -CompressionLevel Optimal
        Write-Host "Packaged: $archivePath"
    }
    finally {
        Pop-Location
    }
}

if ($Target -eq 'all' -or $Target -eq 'admin') { Build-Frontend 'admin' }
if ($Target -eq 'all' -or $Target -eq 'portal') { Build-Frontend 'portal' }
if ($Target -eq 'all' -or $Target -eq 'ledger') { Build-Frontend 'ledger' }
