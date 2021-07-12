build:
	go mod download
	go build -ldflags "-X github.com/JojiiOfficial/ZimWiki/handlers.version=`git describe --tags` -X github.com/JojiiOfficial/ZimWiki/handlers.buildTime=`date +%FT%T%z`" -o ZimWiki

default: build

upgrade:
	go get -u -v
	go mod download
	go mod tidy
	go mod verify

test:
	go test

man: build
	./ZimWiki --help-man | man -l -

run:
	./ZimWiki

debug: build
	./ZimWiki

install:
	@if ! test -f ZimWiki;then echo 'run "make build" first'; exit 1; fi

ifneq ($(shell id -u), 0)
	@echo "You must be root to perform this action."
	@exit 1
endif
	@mkdir -p /usr/local/share/man/man8
	cp ZimWiki /usr/bin/ZimWiki
	/usr/bin/ZimWiki --help-man > ZimWiki.1
	install -Dm644 ZimWiki.1 /usr/share/man/man8/ZimWiki.8
	@rm ZimWiki.1
	@echo Installed successfully!

uninstall:
ifneq ($(shell id -u), 0)
	@echo "You must be root to perform this action."
	@exit 1
endif
	rm /usr/bin/ZimWiki
	rm -f /usr/share/man/man8/ZimWiki.8
	@echo Uninstalled successfully!

clean:
	go clean
	rm -f ZimWiki.1
	rm -f ZimWiki
	rm -f main
