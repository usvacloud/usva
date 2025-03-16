package api

type APIHeaders struct {
	APIKeyHeaderName string
}

var Headers = APIHeaders{
	APIKeyHeaderName: "X-Usva-Api-Key",
}
