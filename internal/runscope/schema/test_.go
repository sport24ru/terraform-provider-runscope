package schema

type TestMinimal struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type TestBase struct {
	TestMinimal
	DefaultEnvironmentId string `json:"default_environment_id"`
}

type Test struct {
	TestBase
	Id    string     `json:"id"`
	Steps []TestStep `json:"steps"`
}

type TestStep struct {
	Id string `json:"id"`
}

type TestGetResponse struct {
	Test `json:"data"`
}

type TestCreateRequest struct {
	TestMinimal
}

type TestCreateResponse struct {
	Test `json:"data"`
}

type TestUpdateRequest struct {
	TestBase
}

type TestUpdateResponse struct {
	Test `json:"data"`
}
