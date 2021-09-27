# ZimWiki
A modern zim fileserver which can handle multiple zim files by serving a beautiful Wiki website. It is a lightweight and performant replacement for kiwix-serve and can handle many big wiki archives (zim files).

# Screenshots
<table>
<thead>
    <td>
        Desktop
    </td>
    <td>
        Mobile
    </td>
</thead>
<tr>
    <td>
        <img src=".img/home.png" width="auto" height="259px"/>
    </td>
    <td>
        <img src=".img/home_mobile.png" width="auto" height="259px"/>
    </td>
</tr>
<tr>
    <td>
        <img src=".img/wiki.png" width="auto" height="259px"/>
    </td>
    <td>
        <img src=".img/wiki_mobile.png" width="auto" height="259px"/>
    </td>
</tr>
</table>
<br> 

# Installation
- Install Go and compile it using `make build`
or
- Download the latest release

# Configuration
### Example
config.toml:
```toml
[Config]
LibraryPath = "./library"
Address = ":8080"
EnableSearchCache = "true"
SearchCacheDuration = "2"
```

# Usage
Your `LibraryPath` must contain your .zim files inside it, you can also link them using symlinks.  
Run the binary and go to `https://localhost:Port`

# Features
- [x] Read/Handle multiple Zim files
- [x] Read Wikis
- [x] Search (inside a wiki or globally)
- [x] Create wiki index files for faster search
- [x] Caching of search results for smoother navigation
- [x] Send content gzipped
- [x] Use symlinks as .zim link
- [x] Replace absolute links with relative ones
- [X] Config file

***

This program is free software: you can redistribute it and/or modify it under the terms of the GNU Affero General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version.

This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License along with this program. If not, see http://www.gnu.org/licenses/.