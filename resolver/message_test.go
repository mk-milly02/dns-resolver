package resolver

import "testing"

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
			want: "3dns6google3com0",
		},
		{
			args: args{
				name: "www.regent.edu.gh",
			},
			want: "3www6regent3edu2gh0",
		},
		{
			args: args{
				name: "bbc.co.uk",
			},
			want: "3bbc2co2uk0",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := EncodeURL(tt.args.name); got != tt.want {
				t.Errorf("EncodeURL() = %v, want %v", got, tt.want)
			}
		})
	}
}
