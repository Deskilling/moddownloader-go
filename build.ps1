$BINARY_NAME = "moddownloader-go"
$OUTPUT_DIR = "./builds"

if (!(Test-Path -Path $OUTPUT_DIR)) {
    New-Item -ItemType Directory -Path $OUTPUT_DIR | Out-Null
}

$PLATFORMS = @(
    "linux/amd64",    # Linux 64-bit
    "darwin/amd64",   # macOS 64-bit (Intel)
    "darwin/arm64",   # macOS 64-bit (Apple Silicon)
    "windows/amd64"   # Windows 64-bit
)

foreach ($PLATFORM in $PLATFORMS) {
    $parts = $PLATFORM -split "/"
    $GOOS = $parts[0]
    $GOARCH = $parts[1]
	
	$PLATFORMNAME = $GOOS
	$ARCHNAME = $GOARCH
    if ($PLATFORMNAME -eq "darwin") {
        $PLATFORMNAME = "macos"
		
        if ($ARCHNAME -eq "amd64") {
            $ARCHNAME = "intel"
        } 
		
        if ($ARCHNAME -eq "arm64") {
            $ARCHNAME = "silicon"
        }
    }
	
	if ($ARCHNAME -eq "amd64") {
		$ARCHNAME = "x64"
	}
	    	
    $OUTPUT_NAME = "$OUTPUT_DIR/$BINARY_NAME-$PLATFORMNAME-$ARCHNAME"
	
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
