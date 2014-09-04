{{define "content"}}

{{if .Documents}}
    {{range .Documents}}

    {{$url := buildUrl "ViewArticle" "year" .Created.Year "month" .Created.Month "title" .Slug}}
    <div class="text">
        <h2 class="title">
            <a href="{{$url}}">{{.Meta.Title}}</a>
        </h2>
        {partial "article-meta.tpl" .}
        <div class="article-content">
        {{if .Preview}}
                {{.Preview | safeHtml}}
                <div class="more">&rarr; <a href="{{$url}}#cut">{preview_link}</a></div>
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

{>index-pager}

{{end}}
