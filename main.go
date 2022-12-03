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
		str := ""
		for i, e := range es {
			str += e.name + " " + hex.EncodeToString([]byte{es[(i+1)%len(es)].id}) + "\n"
		}
		w.Write(bytes.NewBufferString(str).Bytes())
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./main.html")
	})

	http.ListenAndServe(":4200", nil)
}
