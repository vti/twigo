{{define "content"}}

<div class="text">
    <h1 class="title"><a href="/articles///.html">{{.Document.Meta.Title}}</a></h1>
    {{template "article-meta" .}}
    <div class="article-content">
        {{if .Document.Preview}}
        {{.Document.Preview | safeHtml}}
        <a id="cut"></a>
        {{end}}
        {{.Document.Content | safeHtml}}
    </div>
    {>article-pager}
    <h2>Comments</h2>

</div>

{{end}}
