@echo off
REM build.bat - Cross-compile το Go engine για Windows, macOS, Linux
REM Τρέξε από τον φάκελο engine\: build.bat

echo Κατεβάζω dependencies...
go mod tidy

echo.
echo Compiling για όλες τις πλατφόρμες...

set BIN_DIR=..\bin
if not exist %BIN_DIR% mkdir %BIN_DIR%

REM Windows (amd64)
echo   Windows (amd64)...
set GOOS=windows
set GOARCH=amd64
go build -ldflags="-s -w" -o %BIN_DIR%\engine-windows.exe .

REM Linux (amd64)
echo   Linux (amd64)...
set GOOS=linux
set GOARCH=amd64
go build -ldflags="-s -w" -o %BIN_DIR%\engine-linux .

REM macOS Apple Silicon (arm64)
echo   macOS arm64...
set GOOS=darwin
set GOARCH=arm64
go build -ldflags="-s -w" -o %BIN_DIR%\engine-mac-arm64 .

REM macOS Intel (amd64)
echo   macOS amd64...
set GOOS=darwin
set GOARCH=amd64
go build -ldflags="-s -w" -o %BIN_DIR%\engine-mac .

echo.
echo Build ολοκληρώθηκε! Τα binaries είναι στο φάκελο: %BIN_DIR%
dir %BIN_DIR%
