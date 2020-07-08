{{ define "content" }}

<iframe type="text/html" src="{{ .Source }}" style="height:90%;position:absolute;width:100%;right:0" target="_parent" id="wikiContent" frameBorder="0" onload="fixURLs()">
</iframe>

<script src="https://cdnjs.cloudflare.com/ajax/libs/jquery/3.5.1/jquery.min.js"></script>

<script>
// Don't toggle sidebar on small devices, since it would open it
if (width > 1430) {
   
}

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
	   }

	   // Set new href
       $(this).attr("href", newlink)
    })
}
</script>

{{ end }}