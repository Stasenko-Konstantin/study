main: main.go
	fyne-cross windows -arch=386
	fyne-cross windows -arch=amd64
	cp fyne-cross/bin/windows-amd64/модуль.exe module_x64.exe
	cp fyne-cross/bin/windows-386/модуль.exe module_x86.exe
	rm -rf fyne-cross
	rm Icon.png
