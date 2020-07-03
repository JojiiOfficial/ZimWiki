{{ define "content" }}

<center><h1>Available Wikis</h1></center>
<br>
<br>

<div class="container-fluid" id="ContentContainer" style="margin-top: 18px;">
    <div class="row row-cols-3 float-none d-xl-flex justify-content-xl-center align-items-xl-center">

        {{/* Add all available cards */}}
        {{ range .Cards }}
            <div class="col-lg-4" style="min-width: 400px;">
                <div class="card mb-4 box-shadow">
                    <img class="card-img-top w-100 d-block" src='{{ .Image }}' style="width: 200px;height: 200px;">
                    <div class="card-body">
                        <h4 class="card-title">{{ .Title }}</h4>
                        <p class="card-text">{{ .Text }}</p><a class="card-link" href="{{ .Link }}">Read</a>
                    </div>
                </div>
            </div>
        {{ end }}

    </div>
</div>

{{ end }}