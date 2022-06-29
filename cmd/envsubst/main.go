package main

import (
	"io"
	"os"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig/v3"
)

func main() {
	buf, err := io.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}

	env := env()

	t := template.Must(template.New("tmpl").Funcs(sprig.TxtFuncMap()).Parse(string(buf)))
	t.Execute(os.Stdout, env)
}

func env() map[string]string {
	env := make(map[string]string)
	for _, kv := range os.Environ() {
		kv := strings.SplitN(kv, "=", 2)
		if len(kv) != 2 {
			continue
		}
		env[kv[0]] = kv[1]
	}
	return env
}
