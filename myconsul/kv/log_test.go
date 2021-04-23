package kv

import "testing"

func Test_kvLog_Info(t *testing.T) {
	type args struct {
		msg  string
		args []interface{}
	}
	tests := []struct {
		name string
		l    *kvLog
		args args
	}{
		{
			name: "test-1",
			args: args{
				msg: "test-1",
				args: []interface{}{
					"k1",
					"v1",
					"k2",
					2,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewLogger()
			l.Info(tt.args.msg, tt.args.args...)
		})
	}
}
