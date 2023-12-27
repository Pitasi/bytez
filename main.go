package main

import (
	"html/template"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type Parser interface {
	ID() string
	Parse(s string, params url.Values) ([]byte, error)
}

type HTMLer interface {
	HTML(b []byte, params url.Values) template.HTML
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/static") {
			http.FileServer(http.FS(static)).ServeHTTP(w, r)
			return
		}

		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}

		parser := []Parser{
			Hex{},
			ASCII{},
			Base64{},
			Bech32{},
		}
		widgetsWithValue := []HTMLer{
			Hex{},
			ASCII{},
			Base64{},
			Bech32{},
		}

		var input []byte

		q := r.URL.Query()
		if len(q["w"]) > 0 && len(q["input"]) > 0 {
			widgetID := q.Get("w")
			strInput := q.Get("input")

			var found bool
			for _, parser := range parser {
				found = true
				if parser.ID() == widgetID {
					var err error
					input, err = parser.Parse(strInput, q)
					if err != nil {
						log.Printf("parser=%s err: %v", widgetID, err)
						break
					}
				}
			}

			if !found {
				w.Write([]byte("no parser found"))
				return
			}
		}

		err := template.Must(template.New("index").Parse(`
<!DOCTYPE html>
<html class="bg-black">
    <head>
        <title>Bytez</title>
		<link rel="stylesheet" href="/static/tailwind.min.css">
    </head>
    <body>
        <script src="/static/htmx-1.9.10.min.js"></script>
		<main class="flex flex-col gap-6 mt-10 max-w-lg mx-auto font-mono">
			<div class="flex flex-col">
				<h1 class="text-3xl text-gray-300 font-bold">Bytez</h1>
				<p class="text-lg text-gray-500">Convert bytes to different formats</p>
			</div>

			<div class="flex flex-col mt-10 gap-10">
				{{range .Widgets}}
					{{.HTML $.Input $.Params}}
				{{end}}
			</div>
		</main>
    </body>
</html>
        `)).Execute(w, struct {
			Widgets []HTMLer
			Input   []byte
			Params  url.Values
		}{
			Widgets: widgetsWithValue,
			Input:   input,
			Params:  q,
		})
		if err != nil {
			log.Print(err)
		}
	})

	log.Print("listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
