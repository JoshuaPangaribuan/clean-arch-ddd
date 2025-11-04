@echo off
REM ci-check.bat - Run tests to verify the project is in a good state
REM This script is designed to be run in CI/CD pipelines

setlocal enabledelayedexpansion

echo 🔍 Running CI checks...

REM Check if Go is installed
where go >nul 2>&1
if %ERRORLEVEL% neq 0 (
    echo ❌ Error: Go is not installed
    exit /b 1
)

echo ✅ Go version:
go version

REM Generate mocks
echo 📦 Generating mocks...
where mockery >nul 2>&1
if %ERRORLEVEL% neq 0 (
    echo Installing mockery...
    go install github.com/vektra/mockery/v2@latest
)
mockery
if %ERRORLEVEL% neq 0 (
    echo ❌ Error: Failed to generate mocks
    exit /b 1
)

REM Run tests
echo 🧪 Running tests...
go test -v -cover ./...
if %ERRORLEVEL% neq 0 (
    echo ❌ Error: Tests failed
    exit /b 1
)

REM Run tests with coverage report
echo 📊 Generating coverage report...
go test -coverprofile=coverage.out ./...
if %ERRORLEVEL% neq 0 (
    echo ❌ Error: Failed to generate coverage report
    exit /b 1
)

echo ✅ All CI checks passed!
exit /b 0
