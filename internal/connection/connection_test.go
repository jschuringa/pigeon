package connection_test

import (
	"fmt"
	"testing"

	"github.com/jschuringa/pigeon/internal/connection"
)

func TestConnection_Send(t *testing.T) {
	t.Parallel()
	type testCase struct {
		name    string
		conn    mockNetConnection
		wantErr bool
	}

	for _, c := range []testCase{
		{
			name: "no error",
			conn: mockNetConnection{},
		},
		{
			name:    "err returns",
			conn:    mockNetConnection{err: fmt.Errorf("oh no")},
			wantErr: true,
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			connWrapper := connection.NewConnection(&c.conn)
			err := connWrapper.Send([]byte{})
			if err != nil {
				if !c.wantErr {
					t.Fatal("should not have error")
				}
			}
		})
	}
}

func TestConnection_Close(t *testing.T) {
	t.Parallel()
	type testCase struct {
		name    string
		conn    mockNetConnection
		wantErr bool
	}

	for _, c := range []testCase{
		{
			name: "no error",
			conn: mockNetConnection{},
		},
		{
			name:    "err returns",
			conn:    mockNetConnection{err: fmt.Errorf("oh no")},
			wantErr: true,
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			connWrapper := connection.NewConnection(&c.conn)
			err := connWrapper.Close()
			if err != nil {
				if !c.wantErr {
					t.Fatal("should not have error")
				}
			}
		})
	}
}
