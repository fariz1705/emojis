

package main

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ServiceWeaver/weaver"
)


var indexHtml string 

func main() {
	if err := weaver.Run(context.Background(), run); err != nil {
		panic(err)
	}
}

type app struct {
	weaver.Implements[weaver.Main]
	searcher weaver.Ref[Searcher]
	lis      weaver.Listener `weaver:"emojis"`
}


func run(ctx context.Context, a *app) error {
	a.Logger(ctx).Info("emojis listener available.", "addr", a.lis)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		fmt.Fprint(w, indexHtml)
	})
	http.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
		a.handleSearch(a.searcher.Get().Search, w, r)
	})
	http.HandleFunc("/search_chatgpt", func(w http.ResponseWriter, r *http.Request) {
		a.handleSearch(a.searcher.Get().SearchChatGPT, w, r)
	})
	return http.Serve(a.lis, nil)
}


func (a *app) handleSearch(search func(context.Context, string) ([]string, error), w http.ResponseWriter, r *http.Request) {
	
	query := r.URL.Query().Get("q")
	emojis, err := search(r.Context(), query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	bytes, err := json.Marshal(emojis)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintln(w, string(bytes))
}
