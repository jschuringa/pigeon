package broker_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/jschuringa/pigeon/internal/core"
	"github.com/jschuringa/pigeon/pkg/broker"

	"github.com/google/go-cmp/cmp"
	"github.com/gorilla/websocket"
)

// given how much of a pain this test was to setup, I should
// probably consider refactoring this
func TestQueue_Start(t *testing.T) {
	type testCase struct {
		name     string
		messages []*core.BaseModel
		want     []string
	}

	for _, c := range []testCase{
		{
			name: "one message",
			messages: []*core.BaseModel{
				{
					Val: "foo",
				},
			},
			want: []string{"foo"},
		},
		{
			name: "multiple messages",
			messages: []*core.BaseModel{
				{
					Val: "foo",
				},
				{
					Val: "bar",
				},
			},
			want: []string{"foo", "bar"},
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			res := make([]string, 0)
			srv := testServer(ctx, c.name, c.messages)
			go testClient(ctx, &res)
			// replace sleep with a channel? ugh
			// maybe I can pass in a second context that's like
			// test executed context?
			// and then do <- testCtx.Done() to block, then call cancel
			time.Sleep(10 * time.Second)
			// do I really want to have to use concurrency to test?
			// no haha
			srv.Shutdown(ctx)
			cancel()
			if !cmp.Equal(c.want, res) {
				t.Fatal(cmp.Diff(c.want, res))
			}
		})
	}
}

func TestQueue_StartCancelled(t *testing.T) {

}

func testServer(ctx context.Context, name string, msgs []*core.BaseModel) *http.Server {
	srv := &http.Server{Addr: "localhost:8080"}
	http.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		upgrader := websocket.Upgrader{}
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			fmt.Printf("failed to start websocket")
		}
		q := broker.NewQueue(name, conn)
		defer conn.Close()
		go q.Start(ctx)
		for _, m := range msgs {
			q.Push(m)
		}
		<-ctx.Done()
	})

	go srv.ListenAndServe()
	return srv
}

func testClient(ctx context.Context, res *[]string) {
	u := url.URL{Scheme: "ws", Host: "localhost:8080", Path: "/test"}
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return
	}
	defer c.Close()

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				_, barr, err := c.ReadMessage()
				if err != nil {
					return
				}
				bm := &core.BaseModel{}
				err = json.Unmarshal(barr, &bm)
				if err != nil {
					return
				}
				*res = append(*res, bm.Val)
			}
		}
	}()
	<-ctx.Done()
}
