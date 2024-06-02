$appname="clipboardMaster.exe"

if ($args -contains "--debug") {
    Write-Host "Building in debug mode"
    go build -o dist\\$appname # for debugging
} else {
    Write-Host "Building in release mode"
    go build -ldflags "-H=windowsgui" -o dist\\$appname # for release
}