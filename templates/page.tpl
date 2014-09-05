{{define "content"}}

    <div class="text">
        <h1 class="title">{{.Document.Meta.Title}}</h1>
        <div class="article-content">
            {{.Document.Content | safeHtml}}
        </div>
    </div>

{{end}}
