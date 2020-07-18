build:
	go mod download
	go build -o zimwiki

default: build

upgrade:
	go get -u -v
	go mod download
	go mod tidy
	go mod verify

test:
	go test

man: build
	./zimwiki --help-man | man -l -

run:
	./zimwiki

debug: build
	./zimwiki

install:
	@if ! test -f zimwiki;then echo 'run "make build" first'; exit 1; fi

ifneq ($(shell id -u), 0)
	@echo "You must be root to perform this action."
	@exit 1
endif
	@mkdir -p /usr/local/share/man/man8
	cp zimwiki /usr/bin/zimwiki
	/usr/bin/zimwiki --help-man > zimwiki.1
	install -Dm644 zimwiki.1 /usr/share/man/man8/zimwiki.8
	@rm zimwiki.1
	@echo Installed successfully!

uninstall:
ifneq ($(shell id -u), 0)
	@echo "You must be root to perform this action."
	@exit 1
endif
	rm /usr/bin/zimwiki
	rm -f /usr/share/man/man8/zimwiki.8
	@echo Uninstalled successfully!

clean:
	rm -f zimwiki.1
	rm -f zimwiki
	rm -f main
