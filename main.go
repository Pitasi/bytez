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
			ASCII{},
			Binary{},
			Decimal{},
			Hex{},
			Base64{},
			Bech32{},
		}
		widgetsWithValue := []HTMLer{
			ASCII{},
			Binary{},
			Decimal{},
			Hex{},
			Base64{},
			Bech32{},
			Protobuf{},
		}

		var input []byte

		q := r.URL.Query()

		var widgetID string
		for _, w1 := range q["w"] {
			if len(w1) > 0 {
				widgetID = w1
				break
			}
		}

		if widgetID != "" {
			strInput := q.Get("input-" + widgetID)

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
        <meta name="viewport" content="width=device-width, initial-scale=1" />
		<link rel="stylesheet" href="/static/tailwind.min.css">
    </head>
    <body>
        <script src="/static/htmx-1.9.10.min.js"></script>

		<main class="flex flex-col gap-6 mt-10 px-4 pb-64 md:px-0 max-w-lg mx-auto font-mono">
			<div class="flex flex-col">
				<h1 class="text-3xl text-gray-300 font-bold">Bytez</h1>
				<p class="text-lg text-gray-500">Convert bytes to different formats</p>
			</div>

			<script>
				function updateInput(value) {
					document.getElementById("w-input").value = value;
					document.getElementById("w-input-fallback").value = value;
				}
			</script>
			<form id="form" hx-get="/" hx-target="body" hx-swap="outerHTML" hx-push-url="true" hx-trigger="submit" hx-sync="this:replace">
				<button id="w-input" type="submit" name="w" value="" class="hidden" tabindex="-1"></button>
				<input type="hidden" id="w-input-fallback" name="w" value="" />
				<div class="flex flex-col mt-10 gap-10">
					{{range .Widgets}}
						{{.HTML $.Input $.Params}}
					{{end}}
				</div>
			</form>
		</main>

		<footer class="px-4 pb-12 md:px-0 max-w-lg mx-auto font-mono text-gray-400">
			made by <a class="text-gray-300 hover:underline" target="_blank" href="https://anto.pt">💡 anto.pt</a>
			<br />
			<a class="text-gray-300 hover:underline" target="_blank" href="https://github.com/Pitasi/bytez">⭐️ source code</a>
		</footer>
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
