@echo off
setlocal enabledelayedexpansion

set BINARY_NAME=moddownloader-go
set OUTPUT_DIR=builds
if not exist "%OUTPUT_DIR%" mkdir "%OUTPUT_DIR%"

set PLATFORMS=windows/amd64 linux/amd64

for %%P in (%PLATFORMS%) do (
    for /f "tokens=1,2 delims=/" %%A in ("%%P") do (
        set GOOS=%%A
        set GOARCH=%%B
        set OUTPUT_NAME=%OUTPUT_DIR%\%BINARY_NAME%-!GOOS!-!GOARCH!
        
        if "!GOOS!" == "windows" set OUTPUT_NAME=!OUTPUT_NAME!.exe
        
        echo Building for !GOOS!/!GOARCH!...
        set "CMD=go build -o !OUTPUT_NAME! ."
        call %CMD%
        if errorlevel 1 (
            echo An error occurred while building for !GOOS!/!GOARCH!
            exit /b 1
        )
    )
)

echo Builds completed successfully.
endlocal
