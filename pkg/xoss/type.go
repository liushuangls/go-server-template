package xoss

type PolicyToken struct {
	Host     string            `json:"host"`
	FormData map[string]string `json:"form_data"`
}

type Credentials struct {
	AccessKeyID     string
	SecretAccessKey string
	SessionToken    string
	Bucket          string
	ExpireAt        int64
	Region          string
}
