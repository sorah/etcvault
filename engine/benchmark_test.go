package engine

import (
	"testing"
)

func BenchmarkEncryptV1Short(b *testing.B) {
	engine := NewEngine(testKeychain)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = engine.Transform("ETCVAULT::plain1:the-key:this text should be encrypted::ETCVAULT")
	}
}

func BenchmarkEncryptV1Long(b *testing.B) {
	engine := NewEngine(testKeychain)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = engine.Transform("ETCVAULT::plain1:the-key:this text should be encrypted aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa::ETCVAULT")
	}
}

func BenchmarkDecryptV1Short(b *testing.B) {
	engine := NewEngine(testKeychain)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = engine.Transform("ETCVAULT::1:the-key::oXKv3edU7AjUXK1+7+Ng7y5tjByLzMe8MRL2lCxlsE03pHS2AXnd3mvar5dkbgeTU4dY8lcMPYAqRGXi2y9YJ7MD+8vKpkORczLYOBTiSXY8cuttvWY+ffjeJMSsLiHn0tDdtjvCtshSBTe9vLz75yyW8J91DUm9CriHWtQhaXw=::ETCVAULT")
	}
}

func BenchmarkDecryptV1Long(b *testing.B) {
	engine := NewEngine(testKeychain)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = engine.Transform("ETCVAULT::1:the-key:long:JRrn3XxO/HJEu/xYblTkxooOGvFkvnHz4AyinTceZMI2ybRbS2TyoOS+fTGZTTdUMnQ0gKhqH/KsCBjtvW/lw+CXEXVooCmpRCRyVYJIu/FH+oarHIGkpDTeJruEVaL1Jlvo0gb9Ea4zeZuKSiabY+puoTHVCEm1sEN8pHE48xA=,6LaTIBRfKOMBfHq/2JaF/ooeVe97GLGe5gJB8DBYMI30q8mynk9DoMgDKX4ROoiUXatFhSS20hvIIZEUwt62qN7ksivXSb9OybZwU22h6Kw=::ETCVAULT")
	}
}

func BenchmarkPlain(b *testing.B) {
	engine := NewEngine(testKeychain)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = engine.Transform("i'm plain text.")
	}
}
