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
        $previousCi = $env:CI
        try {
            $env:CI = 'true'
            pnpm install --frozen-lockfile
            $installExitCode = $LASTEXITCODE
        }
        finally {
            $env:CI = $previousCi
        }
        if ($installExitCode -ne 0) { throw "pnpm install failed for ifs-$Name" }

        if ($Name -eq 'admin') {
            pnpm run build:prod
        } elseif ($Name -eq 'ledger') {
            pnpm run build:h5
        } else {
            pnpm run build
        }
        if ($LASTEXITCODE -ne 0) { throw "Build failed for ifs-$Name" }

        $distDir = if ($Name -eq 'ledger') {
            Join-Path $projectDir 'dist/build/h5'
        } else {
            Join-Path $projectDir 'dist'
        }
        if (-not (Test-Path (Join-Path $distDir 'index.html'))) {
            throw "Build output is missing: $distDir/index.html"
        }

        $archivePath = Join-Path $projectDir "ifs-$Name.zip"
        if (Test-Path $archivePath) { Remove-Item -LiteralPath $archivePath -Force }
        Add-Type -AssemblyName System.IO.Compression
        $archiveStream = [System.IO.File]::Open($archivePath, [System.IO.FileMode]::CreateNew)
        $archive = [System.IO.Compression.ZipArchive]::new(
            $archiveStream,
            [System.IO.Compression.ZipArchiveMode]::Create
        )
        try {
            $rootPrefix = $distDir.TrimEnd([System.IO.Path]::DirectorySeparatorChar) + [System.IO.Path]::DirectorySeparatorChar
            Get-ChildItem -LiteralPath $distDir -Recurse -File -Force | ForEach-Object {
                $entryName = $_.FullName.Substring($rootPrefix.Length).Replace('\', '/')
                $entry = $archive.CreateEntry($entryName, [System.IO.Compression.CompressionLevel]::Optimal)
                $sourceStream = $_.OpenRead()
                $entryStream = $entry.Open()
                try {
                    $sourceStream.CopyTo($entryStream)
                }
                finally {
                    $entryStream.Dispose()
                    $sourceStream.Dispose()
                }
            }
        }
        finally {
            $archive.Dispose()
            $archiveStream.Dispose()
        }
        Write-Host "Packaged: $archivePath"
    }
    finally {
        Pop-Location
    }
}

if ($Target -eq 'all' -or $Target -eq 'admin') { Build-Frontend 'admin' }
if ($Target -eq 'all' -or $Target -eq 'portal') { Build-Frontend 'portal' }
if ($Target -eq 'all' -or $Target -eq 'ledger') { Build-Frontend 'ledger' }
