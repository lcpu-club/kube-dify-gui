default: pack-windows pack-linux

pack-windows:
    GOOS=windows GOARCH=amd64 CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc fyne package -os windows

pack-linux:
    GOOS=linux GOARCH=amd64 CGO_ENABLED=1 CC=x86_64-linux-gnu-gcc fyne package -os linux

build-windows:
    GOOS=windows GOARCH=amd64 CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc go build -o ./gui.exe -ldflags="-H windowsgui" ./cmd/

build-windows-terminal:
    GOOS=windows GOARCH=amd64 CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc go build -o ./gui-term.exe ./cmd/
    