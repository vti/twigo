{{define "article-meta"}}

{{ $conf := conf }}

    <div class="article-meta">
        {{dateFmt .Created}} by {{or .Meta.Author $conf.Author}}
        <div class="tags">
        </div>
    </div>

{{end}}
