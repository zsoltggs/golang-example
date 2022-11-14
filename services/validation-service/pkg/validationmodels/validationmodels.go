package validationmodels

const (
	PostSchemaPath   = "/schema/{schemaId}"
	GetSchemaPath    = "/schema/{schemaId}"
	PostValidatePath = "/validate/{schemaId}"
)

type UpsertSchemaRequest struct {
	SchemaID string
	Schema   string
}

type UpsertSchemaResponse struct {
	HttpResponse StatusHttpResponse
}

type GetSchemaRequest struct {
	ID string `json:"id"`
}

type GetSchemaResponse struct {
	ID     string `json:"id"`
	Schema string `json:"schema"`
}

type ValidateRequest struct {
	SchemaID string
	Document string
}

type StatusHttpResponse struct {
	Action string `json:"action"`
	ID     string `json:"id"`
	Status string `json:"status"`
}

type ValidateResponse struct {
	HttpResponse StatusHttpResponse
}
