package main

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"math/rand"
	"net/http"
)

type entry struct {
	name string
	id   byte
}

func has_id(es []entry, b byte) bool {
	for _, e := range es {
		if e.id == b {
			return true
		}
	}
	return false
}

func has_name(es []entry, name string) bool {
	for _, e := range es {
		if e.name == name {
			return true
		}
	}
	return false
}

func permutate_rand(es []entry) []entry {
	es_copy := make([]entry, len(es))
	copy(es_copy, es)
	rand.Shuffle(len(es_copy), func(i, j int) {
		es_copy[i], es_copy[j] = es_copy[j], es_copy[i]
	})
	return es_copy
}

func main() {
	objs := make([]entry, 0)

	http.HandleFunc("/add_name", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		name := r.Form.Get("name")
		if name == "" {
			return
		}

		for _, e := range objs {
			if e.name == name {
				obs := make([]byte, 2)
				hex.Encode(obs, []byte{e.id})
				w.Write(obs)
				return
			}
		}

		bs := make([]byte, 1)
		for {
			rand.Read(bs)
			if !has_id(objs, bs[0]) {
				break
			}
		}
		objs = append(objs, entry{name, bs[0]})
		fmt.Println(name)
		obs := make([]byte, 2)
		hex.Encode(obs, bs)
		w.Write(obs)
	})

	http.HandleFunc("/result", func(w http.ResponseWriter, r *http.Request) {
		es := permutate_rand(objs)
		es_copy := make([]entry, len(es))
		copy(es_copy, es)
		for i := range es {
			es_copy[i].id = es[(i+1)%len(es)].id
			es_copy[i].name = es[i].name
		}
		permutate_rand(es_copy)
		str := ""
		for _, e := range es_copy {
			str += fmt.Sprintf("%s: %x\n", e.name, e.id)
		}
		w.Write(bytes.NewBufferString(str).Bytes())
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./main.html")
	})

	http.ListenAndServe(":4200", nil)
}
