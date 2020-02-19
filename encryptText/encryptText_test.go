package encrypttext

import (
	"reflect"
	"testing"
)

func TestEncryptItAll(t *testing.T) {
	type args struct {
		pubKeyFile string
		plaintext  string
	}
	tests := []struct {
		name string
		args args
		want EncryptResult
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := EncryptItAll(tt.args.pubKeyFile, tt.args.plaintext); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("EncryptItAll() = %v, want %v", got, tt.want)
			}
		})
	}
}