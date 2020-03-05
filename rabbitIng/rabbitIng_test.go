package rabbiting

import (
	"reflect"
	"testing"
)

func TestRabbitStats(t *testing.T) {
	tests := []struct {
		name string
		want *RabbitResults
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RabbitStats(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RabbitStats() = %v, want %v", got, tt.want)
			}
		})
	}
}