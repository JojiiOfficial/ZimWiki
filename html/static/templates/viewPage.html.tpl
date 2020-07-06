{{ define "content" }}

<embed type="text/html" src="{{ .Source }}" style="height:90%;position:absolute;width:100%;right:0">


{{/* Hide sidebar */}}
<script>
    var width = (window.innerWidth > 0) ? window.innerWidth : screen.width;

    if width > 1430 {
        setTimeout(function(){
            $("#wrapper").toggleClass("toggled")
        }, 50)
    }
</script>

{{ end }}