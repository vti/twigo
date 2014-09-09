{{define "content"}}
<div class="text">
    <h1>Archive</h1>
    <br />
    {{range .Years}}
        <h2>{{.Name}}</h2>
        <div style="margin-left:2em">
            {{range .Months}}
            {{.Name}}
            <div style="margin-left:2em">
            {{range .Documents}}
                <a href="">{{.Meta.Title}}</a>
                {{template "article-meta" .}}
            {{end}}
            </div>
            {{end}}
        </div>
    {{end}}
</div>
{{end}}
