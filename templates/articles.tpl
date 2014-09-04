{{define "content"}}

{{if .Documents}}
    {{with .Documents}}
    {{range .}}
    <div class="text">
        <h2 class="title">
            <a href="{{buildUrl "ViewArticle" "year" "2012" "month" "12" "title" "hello"}}">{{.Meta.Title}}</a>
        </h2>
        {>article-meta}
        <div class="article-content">
        {{if .Preview}}
                {{.Preview | safeHtml}}
                <div class="more">&rarr; <a href="/articles/{created.year}/{created.month}/{slug}.html#cut">{preview_link}</a></div>
        {{else}}
            {{.Content | safeHtml}}
        {{end}}
            <div class="comment-counter pull-right"><a href="/articles/{created.year}/{created.month}/{slug}.html#disqus_thread">{title}</a></div>
            <div style="clear:both"></div>
        </div>
    </div>
    {{end}}
    {{end}}
{{else}}
    <div class="text center">
        Nothing here yet :(
    </div>
{{end}}

{>index-pager}

{{end}}
