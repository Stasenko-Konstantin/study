#cd client
#fyne-cross windows -arch=amd64
#cd ..
#cp client/fyne-cross/bin/windows-amd64/док_и_серт.exe client.exe
#rm -rf client/fyne-cross
#rm client/Icon.png
cd server
env GOOS=windows GOARCH=amd64 go build main.go
cd ..
cp server/main.exe server.exe