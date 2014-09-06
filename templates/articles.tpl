{{define "content"}}

{{if .Documents}}
    {{range .Documents}}

    {{$url := buildViewArticleUrl .}}
    <div class="text">
        <h2 class="title">
            <a href="{{$url}}">{{.Meta.Title}}</a>
        </h2>
        {{template "article-meta" .}}
        <div class="article-content">
        {{if .Preview}}
                {{.Preview | safeHtml}}
                <div class="more">&rarr; <a href="{{$url}}#cut">{{ or .PreviewLink "Read more"}}</a></div>
        {{else}}
            {{.Content | safeHtml}}
        {{end}}
            <div class="comment-counter pull-right"><a href="{{$url}}#disqus_thread">{{.Meta.Title}}</a></div>
            <div style="clear:both"></div>
        </div>
    </div>
    {{end}}
{{else}}
    <div class="text center">
        Nothing here yet :(
    </div>
{{end}}

{{if or .PrevPageOffset .NextPageOffset}}
    <div id="pager">
        {{if .PrevPageOffset}}
        <a href="/?timestamp={{.PrevPageOffset}}"><span class="arrow">&larr; </span>Later</a>
        {{else}}
        <span class="arrow">&larr; </span>Later
        {{end}}

        {{if .NextPageOffset}}
        <a href="/?timestamp={{.NextPageOffset}}">Earlier<span class="arrow"> &rarr;</span></a>
        {{else}}
        Earlier<span class="arrow"> &rarr;</span>
        {{end}}
    </div>
{{end}}

{{end}}
