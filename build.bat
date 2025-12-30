@echo off
REM Batch file wrapper for build.ps1
REM This allows running the PowerShell build script from cmd.exe

if "%~1"=="" (
    powershell -ExecutionPolicy Bypass -File "%~dp0build.ps1" help
) else (
    powershell -ExecutionPolicy Bypass -File "%~dp0build.ps1" %*
)
