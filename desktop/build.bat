rsrc -manifest exe.manifest -ico static/main.ico
rice embed-go
set GOARCH=amd64
go build -ldflags="-H windowsgui -w -s" -o tcpproxy_64bit.exe

set GOARCH=386
go build -ldflags="-H windowsgui -w -s" -o tcpproxy_32bit.exe

zip windows_x64.zip tcpproxy_64bit.exe
zip windows_x32.zip tcpproxy_32bit.exe