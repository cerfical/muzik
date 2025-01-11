package api

type errorResponse struct {
	Errors []responseError `json:"errors"`
}

type responseError struct {
	Status string `json:"status"`
	Title  string `json:"title"`
}

type dataItemsResponse[T any] struct {
	Data []T `json:"data"`
}

type dataItemResponse[T any] struct {
	Data T `json:"data"`
}
