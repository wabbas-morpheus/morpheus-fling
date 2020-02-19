package encrypttext

import (
	"crypto/rsa"
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


func Test_encryptKey(t *testing.T) {
	type args struct {
		publicKey  *rsa.PublicKey
		sourceText []byte
		label      []byte
	}
	tests := []struct {
		name              string
		args              args
		wantEncryptedText []byte
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotEncryptedText := encryptKey(tt.args.publicKey, tt.args.sourceText, tt.args.label); !reflect.DeepEqual(gotEncryptedText, tt.wantEncryptedText) {
				t.Errorf("encryptKey() = %v, want %v", gotEncryptedText, tt.wantEncryptedText)
			}
		})
	}
}

func Test_encryptText(t *testing.T) {
	type args struct {
		plaintext []byte
		key       []byte
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := encryptText(tt.args.plaintext, tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("encryptText() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("encryptText() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_genRandom(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := genRandom(); got != tt.want {
				t.Errorf("genRandom() = %v, want %v", got, tt.want)
			}
		})
	}
}