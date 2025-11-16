$OutName = "server"
$OS = $PSVersionTable.OS

if ($OS -like "*Windows*") {
    $exe = "builds/${AppName}-windows-amd64.exe"
} else {
    Write-Output "Unsupported OS: $OS"
    exit 1
}

if (Test-Path $exe) {
    Write-Output "Running $exe..."
    & $exe
} else {
    Write-Output "Executable not found: $exe"
    exit 1
}

