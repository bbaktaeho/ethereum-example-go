package ethersutil

import (
	"math/big"
	"reflect"
	"testing"
)

func TestWeiToEther(t *testing.T) {
	type args struct {
		wei *big.Int
	}
	tests := []struct {
		name      string
		args      args
		wantEther *big.Float
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotEther := WeiToEther(tt.args.wei); !reflect.DeepEqual(gotEther, tt.wantEther) {
				t.Errorf("WeiToEther() = %v, want %v", gotEther, tt.wantEther)
			}
		})
	}
}
