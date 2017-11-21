package app

import (
	"database/sql"
	"net/http"

	"github.com/acoshift/postlist/pkg/view"
)

func MakeHandler(db *sql.DB) http.Handler {
	mux := http.NewServeMux()
	mux.Handle("/", makeIndexHandler(db))
	mux.Handle("/healthz", makeHealthzHandler(db))
	mux.Handle("/create", makeCreateHandler(db))
	return mux
}

func makeIndexHandler(db *sql.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}

		ctx := r.Context()

		rows, err := db.QueryContext(ctx, `
            select
                name, content
            from posts
            order by created_at desc
        `)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var posts []*Post
		for rows.Next() {
			var post Post
			err = rows.Scan(&post.Name, &post.Content)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			posts = append(posts, &post)
		}
		if err = rows.Err(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		data := struct {
			Posts []*Post
		}{posts}

		view.Index(w, &data)
	})
}

func makeCreateHandler(db *sql.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			name := r.FormValue("name")
			content := r.FormValue("content")
			if len(name) == 0 || len(content) == 0 {
				http.Redirect(w, r, "/create", http.StatusSeeOther)
				return
			}

			ctx := r.Context()

			_, err := db.ExecContext(ctx, `
                insert into posts (
                    name, content
                ) values (
                    $1, $2
                )
            `, name, content)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		view.Create(w)
	})
}

func makeHealthzHandler(db *sql.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := db.Ping()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write([]byte("ok"))
	})
}
