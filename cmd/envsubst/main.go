/*
 * Copyright ADH Partnership
 *
 *  Licensed under the Apache License, Version 2.0 (the "License");
 *  you may not use this file except in compliance with the License.
 *  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 */

package main

import (
	"io"
	"os"
	"strings"
	"text/template"

	"github.com/adh-partnership/sprig/v3"
)

func main() {
	buf, err := io.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}

	env := env()

	t := template.Must(template.New("tmpl").Funcs(sprig.TxtFuncMap()).Parse(string(buf)))
	_ = t.Execute(os.Stdout, env)
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
