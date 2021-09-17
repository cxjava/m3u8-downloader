package utils

import (
	"testing"
)

func TestExecUnixShell(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "echo hello successfully",
			args: args{
				s: "echo hello",
			},
			wantErr: false,
		},
		{
			name: "exec failed",
			args: args{
				s: "xxx bbb",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ExecUnixShell(tt.args.s); (err != nil) != tt.wantErr {
				t.Errorf("ExecUnixShell() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
