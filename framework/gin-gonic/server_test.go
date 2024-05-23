package gin_gonic

import (
	"testing"
	"time"

	"github.com/charlesbases/logger"

	"github.com/charlesbases/library"
	"github.com/charlesbases/library/framework/gin-gonic/hfwctx"
	"github.com/charlesbases/library/framework/gin-gonic/websocket"
)

func TestRun(t *testing.T) {
	SetModel(RandomModel)

	Run(func(srv *Server) {
		srv.RegisterRouter(func(r *Router) {
			r.GET("/get", func(ctx *hfwctx.Context) {
				ctx.ReturnData(library.NowString())
			})
		})
		srv.RegisterRouter(func(r *Router) {
			r.POST("/post", func(ctx *hfwctx.Context) {
				ctx.ReturnData(library.NowString())
			})
		})
		srv.RegisterRouter(func(r *Router) {
			r.PUT("/put", func(ctx *hfwctx.Context) {
				ctx.ReturnData(library.NowString())
			})
		})
		srv.RegisterRouter(func(r *Router) {
			r.DELETE("/delete", func(ctx *hfwctx.Context) {
				ctx.ReturnData(library.NowString())
			})
		})

		srv.RegisterRouter(func(r *Router) {
			r.GET("/ws", websocket.NewStream(func(o *websocket.Options) {
				o.Action = websockerAction
			}))
		})

		// type Custom struct {
		// 	A struct {
		// 		AA string `json:"aa"`
		// 	} `json:"a"`
		// 	B struct {
		// 		BB []string `json:"bb"`
		// 	} `json:"b"`
		// 	C int  `json:"c"`
		// 	D bool `json:"d"`
		// }
		//
		// var a Custom
		// if err := srv.Unmarshal(&a); err != nil {
		// 	logger.Fatal(err)
		// }
		// fmt.Println(func() string {
		// 	data, _ := json.Marshaler.Marshal(&a)
		// 	return string(data)
		// }())

	})
}

// websockerAction .
func websockerAction(c *hfwctx.Context, session *websocket.Session) {
	logger.Context(c).Debug(time.Now())

	for {
		select {
		case <-time.NewTicker(3 * time.Second).C:
			session.WriteResponse(library.NowString())
		}
	}
}
