{{ define "content" }}

<center><h1>Available Wikis</h1></center>
<br>
<br>

<div class="container-fluid" id="ContentContainer" style="margin-top: 18px;">
    <div class="row float-none d-xl-flex justify-content-xl-center align-items-xl-center">

        {{/* Add all available cards */}}
        {{ range .Cards }}
            <div class="col-xl-5" style="max-width:600px">
                <div class="card mb-3 box-shadow">
                    <div class="card-body">
                        <div class="row">
                            <div class="col-sm-4">
                                <img src="{{ .Image }}" style="width: 50px;height: auto;margin-right:10px">
                            </div>
                            <div class="col col-xxl-5">
                                <h4>{{ .Title }}</h4>
                                <p>{{ .Text }}</p>
                            </div>
                        </div>
                        <a class="card-link" style="float:right" href="{{ .Link }}">Read</a>
                    </div>
                </div>
            </div>
        {{ end }}

    </div>
</div>

{{ end }}