package main // import "randomint"

import (
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"time"
)

var templ = `
<!DOCTYPE html>
<head>
	<title>random int</title>
</head>
<body>
		<center>
		<p><h1>Get a random Int</h1></p>
		<h2>{{ .Random }}</h2>
		</center>
</body>
</html>
`

func main() {
	http.HandleFunc("/", randomInt())
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func randomInt() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
		t := template.New("random")
		tm, err := t.Parse(templ)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		tm.Execute(w, struct {
			Random int
		}{
			Random: rnd.Int(),
		})
	}
}
