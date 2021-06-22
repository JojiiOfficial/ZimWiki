{{  define "content" }}

        <h1> <i class="fa fa-search fa-fw" aria-hidden="true"></i> Search results for "{{ .QueryText }}" </h1>
        <h3>
            [in {{ .SearchSource }}]
        </h3>

   <div class="justify-content-md-center row" id="ContentContainer">

        {{/* Display message if no result was found */}}
        {{if not .Results}} <p id="noResult"><i class="fa fa-frown-o fa-fw" aria-hidden="true"></i> Nothing was found.</p> {{end}}

        {{ range .Results }}
            <div class="col-xl-4 col-lg-6 col-md-6 col-sm-12">
                <div class="card mb-3 box-shadow shadow-sm">
                    <div class="card-body header">
                        <img src="{{ .Img }}" style="width:auto;height:25px"/>
                        <h4 class="card-title" style="margin-bottom:0">{{ .Title }}</h4>
                    </div>
                    <div class="card-body">
                        <a class="btn btn-outline-secondary" href="{{ .Link }}"><i class="fa fa-book fa-fw" aria-hidden="true"></i> Read</a>
                    </div>
                </div>
            </div>
        {{ end }}

    </div>

{{ end }}