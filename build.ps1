param(
    [string]$Version = "00.1"
)

$counterFile = "build_counter.txt"
$date = (Get-Date).ToString("yyyy-MM-dd_HH:mm")

if (Test-Path $counterFile) {
    $build = [int](Get-Content $counterFile) + 1
} else {
    $build = 1
}
Set-Content $counterFile $build

Write-Host "[BTL] sborka v$Version (build #$build) ot $date..." -ForegroundColor Cyan

go build -ldflags "-X main.Version=$Version -X main.BuildNum=$build -X main.BuildDate=$date" -o bltapp.exe

Write-Host "[BTL] binary saved: bltapp.exe" -ForegroundColor Green