{{define "content"}}

<h1>Tags</h1>
    <br />

{{if .Tags}}
    <div class="tags">
    {{range $key, $value := .Tags}}
        <a href="{{buildUrl "ListArticlesByTag" "tag" $key}}">{{$key}}</a>
        <sub>({{$value}})</sub>
    {{end}}
    </div>
{{else}}
    <div class="text center">
        Nothing here yet :(
    </div>
{{end}}

{{end}}
