# clone-washoku.ps1
# GitHubから bonsai/washoku をcloneしてRAG CLIを起動する

$ErrorActionPreference = "Stop"

$Repo = "https://github.com/bonsai/washoku.git"
$Name = "washoku"
$Root = Get-Location
$Target = Join-Path $Root $Name

Write-Host ""
Write-Host "=== washoku RAG ===" -ForegroundColor Cyan
Write-Host ""

if (-not (Get-Command git -ErrorAction SilentlyContinue)) {
    Write-Error "git が見つかりません。"
    exit 1
}

if (-not (Test-Path $Target)) {
    Write-Host "Clone: $Repo" -ForegroundColor Yellow
    git clone $Repo $Target
}
else {
    Write-Host "既存の $Target を使用します。" -ForegroundColor Yellow
    Push-Location $Target
    try {
        git pull --ff-only
    }
    finally {
        Pop-Location
    }
}

$Rag = Join-Path $Target "rag"

if (-not (Test-Path (Join-Path $Rag "cli.py"))) {
    Write-Error "rag/cli.py が見つかりません。"
    exit 1
}

Write-Host ""
Write-Host "RAG ready: $Rag" -ForegroundColor Green
Write-Host ""

if ($args.Count -gt 0) {
    $Question = $args -join " "
    Push-Location $Rag
    try {
        python cli.py ask $Question
    }
    finally {
        Pop-Location
    }
    exit $LASTEXITCODE
}

Write-Host "質問を入力してください。" -ForegroundColor Cyan
Write-Host "終了: exit" -ForegroundColor DarkGray
Write-Host ""

while ($true) {
    $Question = Read-Host "washoku>"
    if ([string]::IsNullOrWhiteSpace($Question)) { continue }
    if ($Question -eq "exit") { break }

    Push-Location $Rag
    try {
        python cli.py ask $Question
    }
    finally {
        Pop-Location
    }
    Write-Host ""
}
