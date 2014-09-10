{{define "content"}}

<div class="text">
    <h1 class="title"><a href="{{buildViewArticleUrl .Document}}">{{.Document.Meta.Title}}</a></h1>
    {{template "article-meta" .Document}}
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
    {{$conf := conf}}
    {{if conf.Disqus}}
    <div id="disqus_thread"></div>
    <script type="text/javascript">
        /* * * CONFIGURATION VARIABLES: EDIT BEFORE PASTING INTO YOUR WEBPAGE * * */
        var disqus_shortname = '{{conf.Disqus.Shortname}}'; // required: replace example with your forum shortname

        {{if conf.Disqus.Developer}}
        var disqus_developer = 1; // developer mode is on_
        {{end}}

        /* * * DON'T EDIT BELOW THIS LINE * * */
        (function() {
            var dsq = document.createElement('script'); dsq.type = 'text/javascript'; dsq.async = true;
            dsq.src = 'http://' + disqus_shortname + '.disqus.com/embed.js';
            (document.getElementsByTagName('head')[0] || document.getElementsByTagName('body')[0]).appendChild(dsq);
        })();
    </script>
    <noscript>Please enable JavaScript to view the <a href="http://disqus.com/?ref_noscript">comments powered by Disqus.</a></noscript>
    <a href="http://disqus.com" class="dsq-brlink">blog comments powered by <span class="logo-disqus">Disqus</span></a>
    {{end}}

</div>

{{end}}
