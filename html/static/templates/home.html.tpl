{{ define "content" }}

    {{ if not .Cards }}
        <h1><i class="fa fa-frown-o fa-fw" aria-hidden="true"></i> {{ gettext "No wiki available" }}</h1>
    {{ else if eq (len .Cards) 1 }}
        <h1><i class="fa fa-book fa-fw" aria-hidden="true"></i> {{ len .Cards }} {{ gettext "Available Wiki" }}</h1>
    {{ else }}
        <h1><i class="fa fa-book fa-fw" aria-hidden="true"></i> {{ len .Cards }} {{ gettext "Available Wikis" }}</h1>
    {{ end }}

    <div class="justify-content-md-center row">

            {{/* Add all available cards */}}
            {{ range .Cards }}
                <div class="d-flex col-xl-4 col-lg-6 col-md-6 col-sm-12">
                    <div class="card flex-fill mb-3 box-shadow shadow-sm">
                        <div class="d-flex card-body flex-column">
                            <div class="row">
                                <div class="col-sm-2">
                                    <img src="{{ .Image }}">
                                </div>
                                <div class="col col-xxl-8">
                                    <h4>{{ .Title }}</h4>
                                    <p>{{ .Text }}</p>
                                </div>
                            </div>
                            <a class="btn btn-outline-secondary mt-auto ms-auto" href="{{ .Link }}"><i class="fa fa-book fa-fw" aria-hidden="true"></i> {{ gettext "Read" }}</a>
                        </div>
                    </div>
                </div>
            {{ end }}

    </div>

    {{ if or .Version .BuildTime }}
        <div class="container position-absolute bottom-0 start-50 translate-middle-x py-2 d-flex justify-content-between bd-highlight mb-3" style="margin-bottom: 0 !important;">
            <small class="text-muted">{{ if .Version }}{{ .Version }}{{ end }} {{ if .BuildTime }}({{ .BuildTime }}){{end}}</small>
        </div>
    {{ end }}

{{ end }}