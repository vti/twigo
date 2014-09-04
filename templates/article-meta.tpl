{{define "article-meta"}}

    <div class="article-meta">
        {{.Document.Meta.Created}} by {{.Document.Meta.Author}}
        <div class="tags">
        </div>
    </div>

{{end}}
