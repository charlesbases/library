package websocket

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

var authWithSign = func(r *http.Request) error {
	fmt.Println(r.URL)
	return nil
}

func TestWS(t *testing.T) {
	n := negroni.New()

	r := mux.NewRouter()
	r.Handle("/stream/{sign:[a-z]*}", NewStream(WithAuth(authWithSign)))

	n.UseHandler(r)
	n.Run(":8080")
}
