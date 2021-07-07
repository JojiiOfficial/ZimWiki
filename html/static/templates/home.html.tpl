{{ define "content" }}

    {{ if not .Cards }}
        <h1>No wiki available</h1>
    {{ else if eq (len .Cards) 1 }}
        <h1>{{ len .Cards }} Available Wiki</h1>
    {{ else }}
        <h1>{{ len .Cards }} Available Wikis</h1>
    {{ end }}

    <div class="justify-content-md-center row" id="ContentContainer">

            {{/* Add all available cards */}}
            {{ range .Cards }}
                <div class="col-xl-4 col-lg-6 col-md-6 col-sm-12">
                    <div class="card mb-3 box-shadow shadow-sm">
                        <div class="card-body">
                            <div class="row">
                                <div class="col-sm-2">
                                    <img src="{{ .Image }}">
                                </div>
                                <div class="col col-xxl-8">
                                    <h4>{{ .Title }}</h4>
                                    <p>{{ .Text }}</p>
                                </div>
                            </div>
                            <a class="btn btn-outline-secondary" href="{{ .Link }}"><i class="fa fa-book fa-fw" aria-hidden="true"></i> Read</a>
                        </div>
                    </div>
                </div>
            {{ end }}

    </div>

{{ end }}