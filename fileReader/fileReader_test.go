package filereader

import (
	portscanner "github.com/wabbas-morpheus/morpheus-fling/portScanner"
	"reflect"
	"testing"
	
)

func TestFileToStructArray(t *testing.T) {
	type args struct {
		fn     string
		uLimit int64
	}
	tests := []struct {
		name string
		args args
		want []*portscanner.PortScanner
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FileToStructArray(tt.args.fn, tt.args.uLimit); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FileToStructArray() = %v, want %v", got, tt.want)
			}
		})
	}
}