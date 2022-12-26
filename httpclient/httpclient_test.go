package httpclient

import (
	"bytes"
	"testing"
)

func TestRequest(t *testing.T) {
	type args struct {
		url    string
		method string
		body   map[string]interface{}
	}

	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "test1",
			args: args{
				url:    "https://api.agify.io/?name=bella",
				method: "GET",
				body:   nil,
			},
			want:    []byte(`{"age":41,"count":40372,"name":"bella"}`),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Request(tt.args.url, tt.args.method, tt.args.body)
			if (err != nil) != tt.wantErr {
				t.Errorf("Request() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !bytes.Equal(got, tt.want) {
				t.Errorf("Request() = %v, want %v", got, tt.want)
			}
		})
	}
}
