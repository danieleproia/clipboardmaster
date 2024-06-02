# generate icon.syso file for windows
rsrc -ico assets/icon.ico -o icon.syso
rsrc -manifest app.manifest -ico assets/icon.ico -o icon.syso