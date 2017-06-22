package myrouter

import (
	"testing"
)

func TestGenerateEmptyPath(t *testing.T) {
	var path, err = generatePath("", make(map[string]string))
	if err != nil {
		t.Fail()
	}
	if path != "" {
		t.Fail()
	}
}

func TestGeneratePathWithMissingParameter(t *testing.T) {
	var path, err = generatePath("/{id}", make(map[string]string))
	if err == nil {
		t.Fail()
	}
	if path != "" {
		t.Fail()
	}
}

func TestGeneratePathWithParameter(t *testing.T) {
	var path, err = generatePath("/{id}", map[string]string{"id": "test"})
	if err != nil {
		t.Fail()
	}
	if path != "/test" {
		t.Fail()
	}
}

func TestGeneratePathWithExtraParameter(t *testing.T) {
	var path, err = generatePath("/{id}", map[string]string{"id": "test", "slug": "poltorak-dariusz"})
	if err != nil {
		t.Fail()
	}
	if path != "/test" {
		t.Fail()
	}
}

func TestGeneratePathWithTwoParameters(t *testing.T) {
	var path, err = generatePath("/{id}/{slug}", map[string]string{"id": "test", "slug": "poltorak-dariusz"})
	if err != nil {
		t.Fail()
	}
	if path != "/test/poltorak-dariusz" {
		t.Fail()
	}
}
