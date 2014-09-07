{{define "content"}}

<div class="text">
<h1>Tag {{.Tag}}
<sup>
<a href="/tags/{{.Tag}}.rss"><img src="/static/images/rss.png" alt="RSS" /></a>
</sup>
</h1>
<br />
{{range .Documents}}
    <a href="{{buildViewArticleUrl .}}">{{.Meta.Title}}</a>
    {{template "article-meta" .}}
{{end}}
</div>

{{end}}
