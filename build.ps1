$BINARY_NAME = "moddownloader-go"
$OUTPUT_DIR = "./builds"

if (!(Test-Path -Path $OUTPUT_DIR)) {
    New-Item -ItemType Directory -Path $OUTPUT_DIR | Out-Null
}

$PLATFORMS = @(
    "windows/amd64",
    "linux/amd64"
)

foreach ($PLATFORM in $PLATFORMS) {
    $parts = $PLATFORM -split "/"
    $GOOS = $parts[0]
    $GOARCH = $parts[1]
    
    $OUTPUT_NAME = "$OUTPUT_DIR/$BINARY_NAME-$GOOS-$GOARCH"
    if ($GOOS -eq "windows") {
        $OUTPUT_NAME += ".exe"
    }

    Write-Host "Building for $GOOS/$GOARCH..."
    
    $env:GOOS = $GOOS
    $env:GOARCH = $GOARCH
    go build -o $OUTPUT_NAME .

    if ($LASTEXITCODE -ne 0) {
        Write-Host "An error occurred while building for $GOOS/$GOARCH" -ForegroundColor Red
        exit 1
    }
}

Write-Host "Builds completed successfully!" -ForegroundColor Green
