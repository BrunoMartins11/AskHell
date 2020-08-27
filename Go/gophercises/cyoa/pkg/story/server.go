package cyoa

import (
	"html/template"
	"net/http"
	"strings"
)

func init() {
	tpl = template.Must(template.New("").Parse(defaultHandlerTmpl))
}

var tpl *template.Template

type storyHandler struct {
	story       Story
	template    *template.Template
	chapterFunc func(r *http.Request) string
}

func (s storyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := s.chapterFunc(r)

	if st, ok := s.story[path]; ok {
		err := s.template.Execute(w, st)
		if err != nil {
			http.Error(w, "Something went wrong.", http.StatusInternalServerError)
		}
		return
	}

	http.Error(w, "Chapter not found", http.StatusNotFound)
}

type HandlerOption func(h *storyHandler)

func WithTemplate(t *template.Template) HandlerOption {
	return func(h *storyHandler) {
		h.template = t
	}
}

func WithChapterFunc(fn func(r *http.Request) string) HandlerOption {
	return func(h *storyHandler) {
		h.chapterFunc = fn
	}
}

func NewStoryHandler(s Story, opts ...HandlerOption) http.Handler {
	h := storyHandler{s, tpl, defaultPath}

	for _, opt := range opts {
		opt(&h)
	}

	return h
}

func defaultPath(r *http.Request) string {
	path := strings.TrimSpace(r.URL.Path)

	if path == "" || path == "/" {
		path = "/intro"
	}

	return path[1:] // Remove '/'
}

var defaultHandlerTmpl = `
<!DOCTYPE html>
<html>
  <head>
	<meta charset="utf-8">
	<title>Choose Your Own Adventure</title>
  </head>
  <body>
	<section class="page">
	  <h1>{{.Title}}</h1>
	  {{range .Paragraphs}}
		<p>{{.}}</p>
	  {{end}}
	  {{if .Options}}
		<ul>
		{{range .Options}}
		  <li><a href="/{{.Arc}}">{{.Text}}</a></li>
		{{end}}
		</ul>
	  {{else}}
		<h3>The End</h3>
		<h2><a href="/">Back to the start!</a></h2>
	  {{end}}
	</section>
	<style>
	  body {
		font-family: helvetica, arial;
	  }
	  h1 {
		text-align:center;
		position:relative;
	  }
	  .page {
		width: 80%;
		max-width: 500px;
		margin: auto;
		margin-top: 40px;
		margin-bottom: 40px;
		padding: 80px;
		background: #FFFCF6;
		border: 1px solid #eee;
		box-shadow: 0 10px 6px -6px #777;
	  }
	  ul {
		border-top: 1px dotted #ccc;
		padding: 10px 0 0 0;
		-webkit-padding-start: 0;
	  }
	  li {
		padding-top: 10px;
	  }
	  a,
	  a:visited {
		text-decoration: none;
		color: #6295b5;
	  }
	  a:active,
	  a:hover {
		color: #7792a2;
	  }
	  p {
		text-indent: 1em;
	  }
	</style>
  </body>
</html>`
