cd client
fyne-cross windows -arch=amd64
cp fyne-cross/bin/windows-amd64/client.exe client.exe
cd ..
rm -rf client/fyne-cross
rm client/Icon.png
cd server
env GOOS=windows GOARCH=amd64 go build main.go
