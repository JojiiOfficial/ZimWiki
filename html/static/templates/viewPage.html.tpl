{{ define "content" }}

<iframe type="text/html" src="{{ .Source }}" style="position:fixed;height:100%;width:100%;right:0" target="_parent" id="wikiContent" frameBorder="0" onload="fixURLs()">
</iframe>

<script src="/assets/js/jquery.min.js"></script>

<script>
function fixURLs(){
	// Get contents of the iframe
    var wikiContent = $("#wikiContent").contents();

    // Open links in parent tab
    wikiContent.find("head").append($("<base target='_parent'>"));

    // Replace all links to preview links
    wikiContent.find("a").each(function(){
        var oldLink = $(this).attr("href");
        if (oldLink === undefined){
            return
        }

        var newLinkBase = "/wiki/view/"+wiki+"/"+namespace+"/"

        // Append subfolders
        var url = window.location.href;
        if (url.endsWith("/")){
            url = url.substr(0, url.length-1);
        }
        var surl = url.split("/").slice(2);
        if (surl.length > 6){
            newLinkBase += surl.slice(5, surl.length-1).join("/")+"/";
        }

        var newlink = newLinkBase + oldLink+"/";

        if (oldLink.startsWith("http")){
            // If link is a full URL, follow it.
            newlink = oldLink;
        } else if (oldLink.startsWith("/") && !oldLink.startsWith("/wiki/view/"+wiki)){
            var url = oldLink;
            if (url.endsWith("/")){
                url = url.substr(0, url.length-1);
            }

            var splitURL = url.split("/");
            newlink = newLinkBase+splitURL[splitURL.length-1];
        } else if (oldLink.startsWith("#") && oldLink.length > 1 ){
            // If pages are only jumpTo's keep them
            newlink = window.location.href+oldLink;
        } else if (oldLink.startsWith("../") && oldLink.length > 1){
            var goback = countInstances(oldLink,"../");
            if (goback != 0){
                var url = window.location.href;
                if (url.endsWith("/")){
                    url = url.substr(0, url.length-1);
                }

                var splitURL = url.split("/");
                var nurl = splitURL.slice(0,splitURL.length-1-goback).join("/");
                if (!nurl.endsWith("/")){
                    nurl += "/";
                }
                nurl += oldLink.substr(3*goback,oldLink.length-1);
                newlink = nurl;
            }
        } 
        
        if (newlink.startsWith("./")){
            newlink = newlink.substr(2);
        }
        newlink = newlink.replace("./","/");

        // Set new href
        $(this).attr("href", newlink)
    })

    // Scroll to object if url ends with '#<id>'
    if (window.location.href.includes("#")){
        var url = window.location.href;
        if (url.endsWith("/")){
            url = url.substr(0, url.length-1);
        }
        var elemID = url.substr(url.lastIndexOf("#")+1, url.length);
        var elem = wikiContent.find("#"+elemID);
        wikiContent.scrollTop(elem.offset().top);
    }
}

function countInstances(string, word) {
   return string.split(word).length - 1;
}
</script>

{{ end }}