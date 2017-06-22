package myrouter

import (
	"fmt"
	"strings"
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

func TestGenerateRegExpFromPathNoArgs(t *testing.T) {
	var path = "/api/user"
	var url = "/api/user"
	var expect = "/api/user"

	var regexp = GenerateRegExpFromPath(path, map[string]string{})
	if regexp.String() != expect {
		fmt.Print(regexp.String())
		fmt.Print("\n")
		t.Fail()
	}
	if !regexp.MatchString(url) {
		t.Fail()
	}
}

func TestGenerateRegExpFromPathSingleArg(t *testing.T) {
	var path = "/api/user/{id}"
	var url = "/api/user/5"
	var expect = "/api/user/([^/]+)"

	var regexp = GenerateRegExpFromPath(path, map[string]string{})
	if regexp.String() != expect {
		fmt.Print(regexp.String())
		fmt.Print("\n")
		t.Fail()
	}
	if !regexp.MatchString(url) {
		t.Fail()
	}
}

func TestGenerateRegExpFromPathSingleCopiedArg(t *testing.T) {
	var path = "/api/user/{id}/user-{id}"
	var url = "/api/user/5/user-5"
	var expect = "/api/user/([^/]+)/user-([^/]+)"

	var regexp = GenerateRegExpFromPath(path, map[string]string{})
	if regexp.String() != expect {
		fmt.Print(regexp.String())
		fmt.Print("\n")
		t.Fail()
	}
	if !regexp.MatchString(url) {
		t.Fail()
	}
}

func TestGenerateRegExpFromPathTwoArgs(t *testing.T) {
	var path = "/api/user/{id}/client-{slug}"
	var url = "/api/user/5/client-poltorak-dariusz"
	var expect = "/api/user/([^/]+)/client-([^/]+)"

	var regexp = GenerateRegExpFromPath(path, map[string]string{})
	if regexp.String() != expect {
		fmt.Print(regexp.String())
		fmt.Print("\n")
		t.Fail()
	}
	if !regexp.MatchString(url) {
		t.Fail()
	}
}

func TestGenerateRegExpFromPathNoValues(t *testing.T) {
	var path = "/api/user/5/client-poltorak-dariusz"
	var url = "/api/user/5/client-poltorak-dariusz"
	var expect = []string{}

	var regexp = GenerateRegExpFromPath(path, map[string]string{})
	var result = regexp.FindAllStringSubmatch(url, -1)
	var parameters []string
	for i := 1; i < len(result[0]); i++ {
		parameters = append(parameters, result[0][i])
	}

	if !regexp.MatchString(url) {
		t.Fail()
	}

	if !arrayCompareString(expect, parameters) {
		fmt.Print(strings.Join(parameters, ", "))
		fmt.Print("\n")
		t.Fail()
	}
}

func TestGenerateRegExpFromPathOneValue(t *testing.T) {
	var path = "/api/user/{id}"
	var url = "/api/user/5"
	var expect = []string{"5"}

	var regexp = GenerateRegExpFromPath(path, map[string]string{})
	var result = regexp.FindAllStringSubmatch(url, -1)
	var parameters []string
	for i := 1; i < len(result[0]); i++ {
		parameters = append(parameters, result[0][i])
	}

	if !regexp.MatchString(url) {
		t.Fail()
	}

	if !arrayCompareString(expect, parameters) {
		fmt.Print(strings.Join(parameters, ", "))
		fmt.Print("\n")
		t.Fail()
	}
}

func TestGenerateRegExpFromPathTwoValues(t *testing.T) {
	var path = "/api/user/{id}/client-{slug}"
	var url = "/api/user/5/client-poltorak-dariusz"
	var expect = []string{"5", "poltorak-dariusz"}

	var regexp = GenerateRegExpFromPath(path, map[string]string{})
	var result = regexp.FindAllStringSubmatch(url, -1)
	var parameters []string
	for i := 1; i < len(result[0]); i++ {
		parameters = append(parameters, result[0][i])
	}

	if !regexp.MatchString(url) {
		t.Fail()
	}

	if !arrayCompareString(expect, parameters) {
		fmt.Print(strings.Join(parameters, ", "))
		fmt.Print("\n")
		t.Fail()
	}
}
