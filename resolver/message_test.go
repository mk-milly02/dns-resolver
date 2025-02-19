package resolver

import (
	"testing"
)

func TestEncodeURL(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			args: args{
				name: "dns.google.com",
			},
			want: "03646e7306676f6f676c6503636f6d00",
		},
		{
			args: args{
				name: "www.regent.edu.gh",
			},
			want: "0377777706726567656e740365647502676800",
		},
		{
			args: args{
				name: "bbc.co.uk",
			},
			want: "0362626302636f02756b00",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := encodeURL(tt.args.name); got != tt.want {
				t.Errorf("EncodeURL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_decodeURL(t *testing.T) {
	type args struct {
		encoded string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			args: args{
				encoded: "03646e7306676f6f676c6503636f6d00",
			},
			want: "dns.google.com",
		},
		{
			args: args{
				encoded: "0377777706726567656e740365647502676800",
			},
			want: "www.regent.edu.gh",
		},
		{
			args: args{
				encoded: "0362626302636f02756b00",
			},
			want: "bbc.co.uk",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := decodeURL(tt.args.encoded); got != tt.want {
				t.Errorf("decodeURL() = %v, want %v", got, tt.want)
			}
		})
	}
}
