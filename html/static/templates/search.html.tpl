{{  define "content" }}

        <script>
            var queryText = '{{ .QueryText }}'
        </script>

        <h1> <i class="fa fa-search fa-fw" aria-hidden="true"></i> {{ gettext "Search results for" }} "{{ .QueryText }}" </h1>
        <h5>
            {{ if eq .SearchSource "global search" }}
                [{{ gettext "global search" }}]
            {{ else }}
                [{{ gettext "in" }} {{ .SearchSource }}]
            {{ end }}
        </h5>
        <h6>
            {{ if eq .NbResults 0 }}
                {{ gettext "No result" }}
            {{ else }}
                {{ .NbResults }} {{ gettext "results" }} 
            {{ end }}
            ({{ .TimeTook }})
        </h6>


            {{ if .Results }}

                <div class="justify-content-md-center row" id="ContentContainer">

                    {{ range .Results }}

                        <div class="col-xl-4 col-lg-6 col-md-6 col-sm-12">
                            <div class="card mb-3 box-shadow shadow-sm">
                                <div class="card-body header">
                                    <img src="{{ .Img }}" style="width:auto;height:25px"/>
                                    <h4 class="card-title" style="margin-bottom:0">{{ .Title }}</h4>
                                </div>
                                <div class="card-body">
                                    <a class="btn btn-outline-secondary" href="{{ .Link }}"><i class="fa fa-book fa-fw" aria-hidden="true"></i> {{ gettext "Read" }}</a>
                                </div>
                            </div>
                        </div>

                    {{ end }}

                </div>

                <div class="container position-absolute bottom-0 start-50 translate-middle-x py-2 d-flex justify-content-between bd-highlight mb-3" style="margin-bottom: 0 !important; padding: 0 !important;">
                    
                    <div class="p-2 bd-highlight">
                        {{ if .PreviousPage }}
                            <form method="POST" action="/search/{{ .Wiki }}/" id="previousPageForm">
                                <input type="hidden" id="sQuery" name="sQuery" value="{{ .QueryText }}">
                                <input type="hidden" id="pageNumber" name="pageNumber" value="{{ .PreviousPage }}">
                                <button class="btn btn-outline-secondary" type="submit"><i class="fa fa-arrow-circle-left fa-fw" aria-hidden="true"></i> {{ gettext "Previous page" }}</button>
                            </form>
                        {{ end }}
                        
                        {{ if not .PreviousPage }}
                            <button class="btn btn-outline-secondary disabled" type="button" aria-disabled="true"><i class="fa fa-arrow-circle-left fa-fw" aria-hidden="true"></i> {{ gettext "Previous page" }}</button>
                        {{ end }}
                    </div>

                    <div class="p-2 bd-highlight text-center align-self-center">
                        {{ gettext "Page" }} {{ .ActualPageNb }}/{{ .NbPages }}
                    </div>

                    <div class="p-2 bd-highlight">
                        {{ if .NextPage }}
                            <form method="POST" action="/search/{{ .Wiki }}/" id="nextPageForm">
                                <input type="hidden" id="sQuery" name="sQuery" value="{{ .QueryText }}">
                                <input type="hidden" id="pageNumber" name="pageNumber" value="{{ .NextPage }}">
                                <button class="btn btn-outline-secondary" type="submit">{{ gettext "Next page" }} <i class="fa fa-arrow-circle-right fa-fw" aria-hidden="true"></i></button>
                            </form>
                        {{ end }}
                        
                        {{ if not .NextPage }}
                            <button class="btn btn-outline-secondary disabled" type="button" aria-disabled="true">{{ gettext "Next page" }} <i class="fa fa-arrow-circle-right fa-fw" aria-hidden="true"></i></button>
                        {{ end }}
                    </div>

            </div>        
                
        {{ end }}

{{ end }}