package sysgatherer

import (
	"reflect"
	"testing"
)

func TestSysGather(t *testing.T) {
	tests := []struct {
		name string
		want []byte
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SysGather(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SysGather() = %v, want %v", got, tt.want)
			}
		})
	}
}