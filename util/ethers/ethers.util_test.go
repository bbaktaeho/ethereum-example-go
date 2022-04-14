package ethers

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
		{"1", args{big.NewInt(1000000000000)}, big.NewFloat(0.000001)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotEther := WeiToEther(tt.args.wei); !reflect.DeepEqual(gotEther.String(), tt.wantEther.String()) {
				t.Errorf("WeiToEther() = %v, want %v", gotEther, tt.wantEther)
			}
		})
	}
}
