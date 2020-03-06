package secparse

import (
	"reflect"
	"testing"
)

func TestParseSecrets(t *testing.T) {
	type args struct {
		secfilePtr string
	}
	tests := []struct {
		name string
		args args
		want Secret
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ParseSecrets(tt.args.secfilePtr); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseSecrets() = %v, want %v", got, tt.want)
			}
		})
	}
}