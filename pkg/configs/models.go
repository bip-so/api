package configs

type PGConnectionInfo struct {
	User     string
	Password string
	Host     string
	Port     string
	Name     string
}

type RedisConnectionInfo struct {
	Password string
	Host     string
	Port     string
}

type AWSS3Info struct {
	AccessKeyID     string
	AccessSecretKey string
	BucketName      string
	Region          string
	CloudFrontURL   string
}

type DiscordBotInfo struct {
	AppID        string
	PublicKey    string
	ClientID     string
	ClientSecret string
	Token        string
	Permission   string
}

type KakfaConnectionInfo struct {
	Hosts string
}

type AlgoliaInfo struct {
	AppID       string
	AdminAPIKey string
}

type SlackInfo struct {
	AppID         string
	ClientID      string
	ClientSecret  string
	SigningSecret string
	Nonce         string
}

type APIHost struct {
	APIHost string
}

type SiteRoot struct {
	SiteRoot string
}
type GitInfo struct {
	Hosts  []string
	Secret string
}

type GetStreamConnectionInfo struct {
	ApiKey    string
	ApiSecret string
}

type AppInfo struct {
	FrontendHost string
	BackendHost  string
}

type SupabaseInfo struct {
	SupabaseBaseurl string
	SupabaseToken   string
}
