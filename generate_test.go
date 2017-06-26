package myrouter

import (
	"fmt"
	"regexp"
	"strings"
	"testing"
)

func testGenerateRegExpFromPathVerification(t *testing.T, path string, url string, expect string, requirements map[string]string) {
	var regexp, err = generateRegExpFromPath(path, requirements)
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}

	if regexp.String() != expect {
		fmt.Print(regexp.String())
		fmt.Print("\n")
		t.Fail()
	}

	if !regexp.MatchString(url) {
		t.Fail()
	}
}

func testGenerateRegExpFromPathValueVerification(t *testing.T, path string, url string, expect []string, requirements map[string]string) {
	var regexp, err = generateRegExpFromPath(path, requirements)
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}

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

func TestGenerateEmptyPath(t *testing.T) {
	var path, err = generatePath("", make(map[string]string), make(map[string]*regexp.Regexp))
	if err != nil {
		t.Fail()
	}
	if path != "" {
		t.Fail()
	}
}

func TestGeneratePathWithMissingParameter(t *testing.T) {
	var path, err = generatePath("/{id}", make(map[string]string), make(map[string]*regexp.Regexp))
	if err == nil {
		t.Fail()
	}

	if path != "" {
		t.Fail()
	}
}

func TestGeneratePathWithParameter(t *testing.T) {
	var path, err = generatePath("/{id}", map[string]string{"id": "test"}, make(map[string]*regexp.Regexp))
	if err != nil {
		t.Fail()
	}

	if path != "/test" {
		t.Fail()
	}
}

func TestGeneratePathWithExtraParameter(t *testing.T) {
	var path, err = generatePath("/{id}", map[string]string{"id": "test", "slug": "poltorak-dariusz"}, make(map[string]*regexp.Regexp))
	if err != nil {
		t.Fail()
	}

	if path != "/test" {
		t.Fail()
	}
}

func TestGeneratePathWithTwoParameters(t *testing.T) {
	var path, err = generatePath("/{id}/{slug}", map[string]string{"id": "test", "slug": "poltorak-dariusz"}, make(map[string]*regexp.Regexp))
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

	testGenerateRegExpFromPathVerification(t, path, url, expect, map[string]string{})
}

func TestGenerateRegExpFromPathSingleArg(t *testing.T) {
	var path = "/api/user/{id}"
	var url = "/api/user/5"
	var expect = "/api/user/([^/]+)"

	testGenerateRegExpFromPathVerification(t, path, url, expect, map[string]string{})
}

func TestGenerateRegExpFromPathSingleCopiedArg(t *testing.T) {
	var path = "/api/user/{id}/user-{id}"
	var url = "/api/user/5/user-5"
	var expect = "/api/user/([^/]+)/user-([^/]+)"

	testGenerateRegExpFromPathVerification(t, path, url, expect, map[string]string{})
}

func TestGenerateRegExpFromPathTwoArgs(t *testing.T) {
	var path = "/api/user/{id}/client-{slug}"
	var url = "/api/user/5/client-poltorak-dariusz"
	var expect = "/api/user/([^/]+)/client-([^/]+)"

	testGenerateRegExpFromPathVerification(t, path, url, expect, map[string]string{})
}

func TestGenerateRegExpFromPathNoValues(t *testing.T) {
	var path = "/api/user/5/client-poltorak-dariusz"
	var url = "/api/user/5/client-poltorak-dariusz"
	var expect = []string{}

	testGenerateRegExpFromPathValueVerification(t, path, url, expect, map[string]string{})
}

func TestGenerateRegExpFromPathOneValue(t *testing.T) {
	var path = "/api/user/{id}"
	var url = "/api/user/5"
	var expect = []string{"5"}

	testGenerateRegExpFromPathValueVerification(t, path, url, expect, map[string]string{})
}

func TestGenerateRegExpFromPathTwoValues(t *testing.T) {
	var path = "/api/user/{id}/client-{slug}"
	var url = "/api/user/5/client-poltorak-dariusz"
	var expect = []string{"5", "poltorak-dariusz"}

	testGenerateRegExpFromPathValueVerification(t, path, url, expect, map[string]string{})
}

func TestGenerateRegExpFromPathWithRequirement(t *testing.T) {
	var path = "/api/user/{id}/client-{slug}"
	var url = "/api/user/5/client-poltorak-dariusz"
	var expect = "/api/user/([^/]+)/client-([a-z\\-]+)"
	var requirements = map[string]string{"slug": "[a-z\\-]+"}

	testGenerateRegExpFromPathVerification(t, path, url, expect, requirements)
}

func TestGenerateRegExpFromPathWithTwoRequirements(t *testing.T) {
	var path = "/api/user/{id}/client-{slug}"
	var url = "/api/user/5/client-poltorak-dariusz"
	var expect = "/api/user/([1-9]+[0-9]*)/client-([a-z\\-]+)"
	var requirements = map[string]string{"id": "[1-9]+[0-9]*", "slug": "[a-z\\-]+"}

	testGenerateRegExpFromPathVerification(t, path, url, expect, requirements)
}

func TestGenerateUrl(t *testing.T) {
	var expect = "https://example.com/api/user/5"
	var schema = "https"
	var host = "example.com"
	var port int
	var path = "/api/user/{id}"
	var parameters = map[string]string{"id": "5"}

	var url, err = generateURL(schema, host, port, path, parameters, make(map[string]*regexp.Regexp))
	if err != nil {
		t.Fail()
	}
	if url != expect {
		t.Fail()
	}
}

func TestGenerateUrlTwoParameters(t *testing.T) {
	var expect = "https://example.com/api/user/5/dariusz-poltorak"
	var schema = "https"
	var host = "example.com"
	var port int
	var path = "/api/user/{id}/{slug}"
	var parameters = map[string]string{"slug": "dariusz-poltorak", "id": "5"}

	var url, err = generateURL(schema, host, port, path, parameters, make(map[string]*regexp.Regexp))
	if err != nil {
		t.Fail()
	}
	if url != expect {
		t.Fail()
	}
}

func TestGenerateUrlMissingParameters(t *testing.T) {
	var expect = "https://example.com/api/user/5/dariusz-poltorak"
	var schema = "https"
	var host = "example.com"
	var port int
	var path = "/api/user/{id}/{slug}"
	var parameters = map[string]string{"id": "5"}

	var url, err = generateURL(schema, host, port, path, parameters, make(map[string]*regexp.Regexp))
	if err == nil {
		t.Fail()
	}
	if url == expect {
		t.Fail()
	}
}

func TestGenerateUrlTwoParametersWithPort(t *testing.T) {
	var expect = "https://example.com:80/api/user/5/dariusz-poltorak"
	var schema = "https"
	var host = "example.com"
	var port = 80
	var path = "/api/user/{id}/{slug}"
	var parameters = map[string]string{"slug": "dariusz-poltorak", "id": "5"}

	var url, err = generateURL(schema, host, port, path, parameters, make(map[string]*regexp.Regexp))
	if err != nil {
		t.Fail()
	}
	if url != expect {
		t.Fail()
	}
}
