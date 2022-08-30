tasklist | findstr "go-fly-pro"
if %ERRORLEVEL% EQU 1 (
    go-fly-pro.exe server
)
exit