normal:
	go build -o gorender render.go

install:
	cp gorender ~/.local/bin/

clean:
	rm gorender

uninstall:
	rm ~/.local/bin/gorender

full: normal install clean
