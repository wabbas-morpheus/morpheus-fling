package portscanner

import (
	"reflect"
	"testing"
	"time"
)

func TestScanPort(t *testing.T) {
	type args struct {
		ip      string
		port    int
		timeout time.Duration
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ScanPort(tt.args.ip, tt.args.port, tt.args.timeout); got != tt.want {
				t.Errorf("ScanPort() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStart(t *testing.T) {
	type args struct {
		psEntity []*PortScanner
		timeout  time.Duration
	}
	tests := []struct {
		name string
		args args
		want []ScanResult
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Start(tt.args.psEntity, tt.args.timeout); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Start() = %v, want %v", got, tt.want)
			}
		})
	}
}