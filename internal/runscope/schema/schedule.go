package schema

type ScheduleBase struct {
	EnvironmentId string `json:"environment_id"`
	Interval      string `json:"interval"`
	Note          string `json:"note"`
}

type Schedule struct {
	ScheduleBase
	Id string `json:"id"`
}

type ScheduleGetResponse struct {
	Schedule `json:"data"`
}

type ScheduleCreateRequest struct {
	ScheduleBase
}

type ScheduleCreateResponse struct {
	Schedule `json:"data"`
}
