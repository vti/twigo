{{define "content"}}

<div class="text">
    <h1 class="title"><a href="{{buildViewArticleUrl .Document}}">{{.Document.Meta.Title}}</a></h1>
    {{template "article-meta" .}}
    <div class="article-content">
        {{if .Document.Preview}}
        {{.Document.Preview | safeHtml}}
        <a id="cut"></a>
        {{end}}
        {{.Document.Content | safeHtml}}
    </div>

{{if or .NewerDocument .OlderDocument }}
    <div id="pager">
        <span class="active">
        {{if .NewerDocument}}
            <span class="arrow">&larr; </span><a
            href="{{buildViewArticleUrl .NewerDocument}}">{{.NewerDocument.Meta.Title}}</a> &nbsp;
        {{end}}
    |
        {{if .OlderDocument}}
            &nbsp;<a href="{{buildViewArticleUrl .OlderDocument}}">{{.OlderDocument.Meta.Title}}</a><span class="arrow"> &rarr;</span>
        {{end}}
        </span>
    </div>
{{end}}

    <h2>Comments</h2>

</div>

{{end}}
