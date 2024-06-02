$appname="clipboardMaster.exe"

if ($args -contains "--debug") {
    Write-Host "Building in debug mode"
    go build -o dist\\$appname # for debugging
} else {
    Write-Host "Building in release mode"
    go build -ldflags "-H=windowsgui" -o dist\\$appname # for release
}

# copy languages folder to dist
Copy-Item -Path .\languages -Destination .\dist\languages -Recurse -Force
# copy plugins folder to dist
Copy-Item -Path .\plugins -Destination .\dist\plugins -Recurse -Force
# copy assets\icon.ico to dist
Copy-Item -Path .\assets\icon.ico -Destination .\dist\icon.ico -Force