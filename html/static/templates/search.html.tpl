{{  define "content" }}

    <center>
    <div>
        <h1> Search results for "{{ .QueryText }}" </h1>
        <span style="width:100%;text-align:center">
            [in {{ .SearchSource }}]
        </span>
    </div>
    </center>

    <br><br>

    {{ range .Results }}
        <div class="row" style="padding-bottom: 15px;">
                <div class="col">
                    <div class="card shadow-sm" style="cursor:pointer" onclick="window.location.href='{{ .Link }}'">
                        <div class="card-body header">
                            <img src="{{ .Img }}" style="width:auto;height:25px"/>
                            <h4 class="card-title" style="margin-bottom:0">{{ .Title }}</h4>
                        </div>
                    </div>
                </div>
            </div>
    {{ end }}

{{ end }}