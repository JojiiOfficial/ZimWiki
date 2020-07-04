{{ define "content" }}

<embed type="text/html" src="{{ .Source }}" style="height:100%;position:absolute;width:100%;right:0">


{{/* Hide sidebar */}}
<script>
    setTimeout(function(){
        $("#wrapper").toggleClass("toggled")
    }, 50)
</script>

{{ end }}