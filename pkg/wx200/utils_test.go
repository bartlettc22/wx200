package wx200

import "testing"

func TestSubDecimal(t *testing.T) {

	b := byte('\xff')
	val, _ := subDecimal(b, 0, 0)
	if val != 1 {
		t.Errorf("Decimal value was not correct, got: %d, want: %d.", val, 1)
	}

	b = byte('\xff')
	val, _ = subDecimal(b, 0, 7)
	if val != 255 {
		t.Errorf("Decimal value was not correct, got: %d, want: %d.", val, 255)
	}

	b = byte('\xff')
	val, _ = subDecimal(b, 6, 6)
	if val != 1 {
		t.Errorf("Decimal value was not correct, got: %d, want: %d.", val, 1)
	}

	b = byte('\xaa')
	val, _ = subDecimal(b, 1, 6)
	if val != 21 {
		t.Errorf("Decimal value was not correct, got: %d, want: %d.", val, 21)
	}

	b = byte('\xff')
	_, err := subDecimal(b, 6, 5)
	if err == nil {
		t.Errorf("Should have gotten out of order error")
	}

	b = byte('\xff')
	_, err = subDecimal(b, 0, 8)
	if err == nil {
		t.Errorf("Should have gotten out of bounds error")
	}
}
