{{define "article-meta"}}

    <div class="article-meta">
        {{dateFmt .Document.Created}} by {{or .Document.Meta.Author .Conf.Author}}
        <div class="tags">
        </div>
    </div>

{{end}}
