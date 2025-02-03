package api_test

import (
	"fmt"
	"net/url"
	"path/filepath"
	"strconv"

	"github.com/gavv/httpexpect/v2"
)

func isErrorResponse(status int) func(*httpexpect.Response) {
	return func(r *httpexpect.Response) {
		obj := r.Status(status).
			JSON().
			Schema(modelRef("error")).
			Object()

		obj.Path("$.error.status").
			IsEqual(strconv.Itoa(status))
	}
}

func isDataResponse(model string, status int, data any) func(*httpexpect.Response) {
	return func(r *httpexpect.Response) {
		r.Status(status).
			JSON().
			Schema(modelRef(model)).
			Object().
			ContainsValue(data)
	}
}

func modelRef(model string) string {
	p, err := filepath.Abs("../../api/models.json")
	if err != nil {
		panic(err)
	}

	p, err = url.JoinPath("/", p)
	if err != nil {
		panic(err)
	}

	return fmt.Sprintf("file://%s#/$def/%s", p, model)
}
