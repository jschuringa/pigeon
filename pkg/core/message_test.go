package core_test

import (
	"encoding/json"
	"testing"

	"github.com/jschuringa/pigeon/pkg/core"

	"github.com/google/go-cmp/cmp"
)

func TestMessage_GetContent(t *testing.T) {
	t.Parallel()
	type testCase struct {
		name    string
		content any
		want    any
		wantErr bool
	}

	for _, c := range []testCase{
		{
			name: "Content matches",
			content: &core.BaseModel{
				Val: "foo",
			},
			want: &core.BaseModel{
				Val: "foo",
			},
		},
		{
			name:    "Wrong type fails",
			content: &core.Message{},
			want:    &core.BaseModel{},
			wantErr: true,
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()
			enc, err := json.Marshal(c.content)
			if err != nil {
				t.Fatal("could not marshal content")
			}
			res, err := core.GetContent[core.BaseModel](&core.Message{Content: enc})
			if err != nil {
				if !c.wantErr {
					t.Fatalf("GetContent failed: %v", err)
				}
				return
			}
			if !cmp.Equal(c.want, res) {
				t.Fatal(cmp.Diff(c.want, res))
			}
		})
	}
}
