$appname="clipboardMaster"

if ($args -contains "--debug") {
    Write-Host "Building in debug mode"
    go build -o "dist\\$appname.exe" # for debugging
} else {
    Write-Host "Building in release mode"
    go build -ldflags "-H=windowsgui" -o "dist\\$appname.exe" # for release
}

# copy languages folder to dist
Remove-Item -Path .\dist\languages -Recurse -Force
Copy-Item -Path .\languages -Destination .\dist\ -Recurse -Force
# copy plugins folder to dist
Remove-Item -Path .\dist\plugins -Recurse -Force
Copy-Item -Path .\plugins -Destination .\dist\ -Recurse -Force
# copy assets\icon.ico to dist
Copy-Item -Path .\assets\icon.ico -Destination .\dist\icon.png -Force

# create zip file for distribution
Remove-Item -Path .\dist\$appname-portable.zip -Force
Compress-Archive -Path .\dist\* -DestinationPath .\dist\$appname-portable.zip -Force