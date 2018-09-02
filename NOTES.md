## Compiling for Windows

I had to do the following to get things cross compiling to Windows:

    brew install mingw-w64

    GOOS=windows GOARCH=386 CGO_ENABLED=1 CC=i686-w64-mingw32-gcc go build -v github.com/chrismytton/csvquery/cmd/csvquery
