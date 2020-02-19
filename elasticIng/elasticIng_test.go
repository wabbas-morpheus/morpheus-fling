package elasticing

import (
	"reflect"
	"testing"
)

func TestElasticHealth(t *testing.T) {
	tests := []struct {
		name string
		want *Esstats
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ElasticHealth(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ElasticHealth() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestElasticIndices(t *testing.T) {
	tests := []struct {
		name string
		want []Esindices
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ElasticIndices(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ElasticIndices() = %v, want %v", got, tt.want)
			}
		})
	}
}