package ridsl

import "testing"

func Test_canConvertLeftTypeToRightType(t *testing.T) {
	tests := []struct {
		name      string
		leftType  ComplexType
		rightType ComplexType
		want      bool
	}{{
		name:      "with equal simple types",
		leftType:  QUICConnectionType,
		rightType: QUICConnectionType,
		want:      true,
	}, {
		name:      "with different simple types",
		leftType:  TCPConnectionType,
		rightType: QUICConnectionType,
		want:      false,
	}, {
		name:      "with equal complex types",
		leftType:  SumType(TCPConnectionType, TLSConnectionType, QUICConnectionType),
		rightType: SumType(TCPConnectionType, TLSConnectionType, QUICConnectionType),
		want:      true,
	}, {
		name:      "with equal complex types and swapped types order",
		leftType:  SumType(TCPConnectionType, TLSConnectionType, QUICConnectionType),
		rightType: SumType(TCPConnectionType, QUICConnectionType, TLSConnectionType),
		want:      true,
	}, {
		name:      "with left type being a strict subset of the right type",
		leftType:  SumType(TCPConnectionType, TLSConnectionType),
		rightType: SumType(TCPConnectionType, QUICConnectionType, TLSConnectionType),
		want:      true,
	}, {
		name:      "with left type being simple and the right type being complex",
		leftType:  TCPConnectionType,
		rightType: SumType(TCPConnectionType, QUICConnectionType, TLSConnectionType),
		want:      true,
	}, {
		name:      "with partially overalapped but unrelated types",
		leftType:  SumType(TCPConnectionType, EndpointType),
		rightType: SumType(TLSConnectionType, TCPConnectionType, DNSLookupResultType),
		want:      false,
	}, {
		name:      "with completely unrelated types",
		leftType:  SumType(TCPConnectionType, EndpointType),
		rightType: SumType(TLSConnectionType, DNSLookupResultType),
		want:      false,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := canConvertLeftTypeToRightType(tt.leftType, tt.rightType); got != tt.want {
				t.Fatalf("canConvertLeftTypeToRightType() = %v, want %v", got, tt.want)
			}
		})
	}
}
