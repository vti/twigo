{{define "article-meta"}}

    <div class="article-meta">
        {{dateFmt .Meta.Created}} by {{or .Meta.Author .Conf.Author}}
        <div class="tags">
        </div>
    </div>

{{end}}
