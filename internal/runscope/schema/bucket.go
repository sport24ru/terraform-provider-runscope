package schema

type Bucket struct {
	Key  string     `json:"key"`
	Name string     `json:"name"`
	Team BucketTeam `json:"team"`
}

type BucketTeam struct {
	Name string `json:"name"`
	Id   string `json:"id"`
}

type BucketCreateResponse struct {
	Bucket `json:"data"`
}

type BucketGetResponse struct {
	Bucket `json:"data"`
}

type BucketListResponse struct {
	Buckets []Bucket `json:"data"`
}
