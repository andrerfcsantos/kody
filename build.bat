@echo off
for /f %%i in ('git describe --tags --always --dirty') do set GIT_VERSION=%%i
set VERSION=dev-%GIT_VERSION%
for /f %%i in ('git rev-parse HEAD') do set COMMIT=%%i
for /f %%i in ('powershell -Command "Get-Date -UFormat '%%Y-%%m-%%dT%%H:%%M:%%SZ'"') do set DATE=%%i

echo Building with version: %VERSION%, commit: %COMMIT%, date: %DATE%
go build -ldflags="-s -w -X main.version=%VERSION% -X main.commit=%COMMIT% -X main.date=%DATE%"
go install -ldflags="-s -w -X main.version=%VERSION% -X main.commit=%COMMIT% -X main.date=%DATE%"