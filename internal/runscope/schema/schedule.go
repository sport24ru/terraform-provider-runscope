package schema

type ScheduleBase struct {
	EnvironmentId string `json:"environment_id"`
	Interval      string `json:"interval"`
	Note          string `json:"note"`
}

type Schedule struct {
	ScheduleBase
	Id         string `json:"id"`
	ExportedAt int64  `json:"exported_at"`
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

type ScheduleUpdateRequest struct {
	ScheduleBase
}

type ScheduleUpdateResponse struct {
	Schedule `json:"data"`
}
