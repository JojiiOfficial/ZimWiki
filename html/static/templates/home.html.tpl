{{ define "content" }}

<h1>Available Wikis</h1>

<div class="justify-content-md-center row" id="ContentContainer">

        {{/* Display message if no wiki was found */}}
        {{if not .Cards}} Nothing here yet {{end}}

        {{/* Add all available cards */}}
        {{ range .Cards }}
            <div class="col-xl-4 col-lg-6 col-md-6 col-sm-12">
                <div class="card mb-3 box-shadow shadow-sm">
                    <div class="card-body">
                        <div class="row">
                            <div class="col-sm-2">
                                <img src="{{ .Image }}">
                            </div>
                            <div class="col col-xxl-5">
                                <h4>{{ .Title }}</h4>
                                <p>{{ .Text }}</p>
                            </div>
                        </div>
                        <a class="btn btn-outline-secondary" href="{{ .Link }}"><i class="fa fa-book"></i> Read</a>
                    </div>
                </div>
            </div>
        {{ end }}

</div>

{{ end }}