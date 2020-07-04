# ZimServer
A modern zim fileserver which can handle multiple zim files by serving a beautiful Wiki website. It is a lightweight and performant relpacement for kiwix-serve and can handle many big wiki archives (zim files).

# Installation
- Install go and compile it using `go build -o zimserver`
or
- Download the latest release

# Usage
Create a folder `library` and place your .zim files inside it, or link them using symlinks.<br>
Run the binary and go to `https://localhost:8080`

# Features
- [x] Read/Handle multiple Zim files
- [x] Read Wikis
- [ ] Replace absolute links with relative ones
- [ ] Search