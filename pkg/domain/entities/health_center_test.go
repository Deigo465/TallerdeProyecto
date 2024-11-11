package entities

import "testing"

func TestNewHealthCenter(t *testing.T) {
	// Given
	name := "New Name"
	district := "New District"
	address := "New Address"

	// When
	healthCenter := NewHealthCenter(1, name, district, address)

	// Then
	if healthCenter.Name != name {
		t.Fatalf("Expecting name to be %s, got %s", name, healthCenter.Name)
	}

	if healthCenter.District != district {
		t.Fatalf("Expecting district to be %s, got %s", district, healthCenter.District)
	}

	if healthCenter.Address != address {
		t.Fatalf("Expecting address to be %s, got %s", address, healthCenter.Address)
	}
}

func TestNewFakeHealthCenter(t *testing.T) {
	// Given
	// When
	healthCenter := NewFakeHealthCenter()

	// Then
	if healthCenter.Name == "" {
		t.Fatalf("Expecting name to not be empty")
	}

	if healthCenter.District == "" {
		t.Fatalf("Expecting name to not be district")
	}

	if healthCenter.Address == "" {
		t.Fatalf("Expecting name to not be address")
	}

}
