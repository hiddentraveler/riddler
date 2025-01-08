run:
	mkdir -p tmp
	go build -o tmp/riddler main.go
	cp -u tmp/riddler ~/files/scripts/cmdline/
