{{define "layouts/html"}}

<!doctype html>
    <head>
        <meta charset="utf-8">
        <title>{{.Conf.Title}}</title>
        <link rel="stylesheet" href="/static/bootstrap/css/bootstrap.min.css" type="text/css" />
        <link rel="stylesheet" href="/static/bootstrap/css/bootstrap-responsive.min.css" type="text/css" />
        <link rel="stylesheet" href="/static/css/codemirror.css" type="text/css" />
        <link rel="stylesheet" href="/static/css/styles.css" type="text/css" />
        <link rel="alternate" type="application/rss+xml" title="{{.Conf.Title}}" href="/index.rss" />
        <meta name="viewport" content="width=device-width, initial-scale=1.0">
        <meta name="generator" content="{{.Conf.Generator}}" />
    </head>
    <body>
        <div class="container">
        <div class="row">
            <div class="span2">&nbsp;</div>
            <div class="span8">
                <div class="page-header">
                    <div id="header">
                        <div id="title">
                            <a href="/">{{.Conf.Title}}</a>
                            <sup><a href="/index.rss"><img src="/static/images/rss.png" alt="RSS" /></a></sup>
                        </div>
                        <div id="description">{{.Conf.Description}}</div>
                        <span id="author">{{.Conf.Author}}</span>,
                        <span id="about">{{.Conf.About}}</span>
                        <div class="menu">
                            {{range .Conf.Menu}}
                            <a href="{{.Link}}">{{.Title}}</a>
                            {{end}}
                        </div>
                    </div>
                </div>
            </div>
            <div class="span2">&nbsp;</div>
        </div>

        <div class="row">
            <div class="span2">&nbsp;</div>
            <div class="span8">
            {{template "content" .}}
            </div>
            <div class="span2">&nbsp;</div>
        </div>

        <div class="row">
            <div class="span2">&nbsp;</div>
            <div class="span8">
                <div id="footer">
                {{or .Conf.Footer "Powered by <a href=\"http://github.com/vti/twigo\">twigo</a>" | safeHtml}}
                </div>
            </div>
            <div class="span2">&nbsp;</div>
        </div>
        </div>

        <script type="text/javascript" src="/static/javascripts/jquery.js"></script>
        <script type="text/javascript" src="/static/javascripts/codemirror.js"></script>
        <script type="text/javascript" src="/static/javascripts/perl.js"></script>

        <script type="text/javascript">
            $(document).ready(function() {
                var editors = [];
                $('pre.perl').each(function() {
                    $(this).replaceWith('<textarea class="code perl">' + $(this).text() + '</textarea>');
                });
                $('textarea').each(function() {
                    var editor = CodeMirror.fromTextArea(this, {readOnly: true, lineNumbers: true});
                    editors.push(editor);
                });
            });
        </script>

    </body>
</html>

{{end}}
