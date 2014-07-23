package main

import (
	"net/http"
	"html/template"
	"github.com/hoisie/redis"
)

var Redis redis.Client

var page = `
<!DOCTYPE html>
<html lang="en">
<body>
<h1>counting online users with redis and go</h1>
<p>{{.count}} online users:</p>
<ul>
{{range .users}}
<li>{{.}} is online.</li>
{{end}}
</ul>
</body>
</html>
`
type FakeApp struct {
	Name string
	Env string
}

var App = &FakeApp{Name: "roc_test", Env: "dev"}

func main() {
	var tmpl = template.New("test")
	tmpl, _ = tmpl.Parse(page)
	var data = map[string]interface{}{}

	// /?u=username1...
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		username := r.FormValue("u")
		if username != "" {
			add_online_username(username)
		}

		data["count"] = online_users_count()
		data["users"] = online_usernames()
		tmpl.Execute(w, data)
	})
	http.ListenAndServe(":8080", nil)
}
