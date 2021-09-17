package ping

import (
	"testing"
	"time"
)

func Test_pingIP(t *testing.T) {
	type args struct {
		pingRecord PingRecord
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "ping 1.1.1.1",
			args: args{
				pingRecord: PingRecord{
					IPAddress: "1.1.1.1",
					PingRTT:   time.Duration(1),
				},
			},
		},
		{
			name: "ping 8.8.8.8",
			args: args{
				pingRecord: PingRecord{
					IPAddress: "8.8.8.8",
					PingRTT:   time.Duration(1),
				},
			},
		},
		{
			name: "ping 114.114.114.114",
			args: args{
				pingRecord: PingRecord{
					IPAddress: "114.114.114.114",
					PingRTT:   time.Duration(1),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pingIP(tt.args.pingRecord)
		})
	}

}

func Test_outputParameter(t *testing.T) {
	type args struct {
		pingRecords []PingRecord
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "output to parameters",
			args: args{
				pingRecords: []PingRecord{
					{"1.1.1.1", 2 * time.Microsecond},
					{"8.8.8.8", 88 * time.Microsecond},
					{"8.8.4.4", 99 * time.Microsecond},
					{"114.114.114.114", 0 * time.Microsecond},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			outputParameter(tt.args.pingRecords)
		})
	}
}

func Test_outputTime(t *testing.T) {
	type args struct {
		pingRecords []PingRecord
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "sort by time",
			args: args{
				pingRecords: []PingRecord{
					{"1.1.1.1", 1 * time.Microsecond},
					{"8.8.8.8", 28 * time.Microsecond},
					{"8.8.4.4", 33 * time.Microsecond},
					{"114.114.114.114", 0 * time.Microsecond},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			outputTime(tt.args.pingRecords)
		})
	}
}
