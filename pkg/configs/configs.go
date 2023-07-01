package configs

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

func InitConfig(configName, configPath string) {
	viper.SetConfigType("dotenv")
	viper.SetConfigFile(configName)
	viper.AddConfigPath(configPath)

	//Overrides is available.
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
}

func GetConfigString(key string) string {
	value := viper.GetString(key)
	return value
}

func GetConfigInt(key string) int {
	value := viper.GetInt(key)
	return value
}

func GetConfigBool(key string) bool {
	value := viper.GetBool(key)
	return value
}

func IsLive() bool {
	return GetConfigString(ENV_KEY) == ENV_PRODUCTION
}

func IsDev() bool {
	return GetConfigString(ENV_KEY) == ENV_DEV
}

func IsLocal() bool {
	return GetConfigString(ENV_KEY) == ENV_LOCAL
}

func GetPGConfig() *PGConnectionInfo {
	return &PGConnectionInfo{
		User:     GetConfigString(PG_USER_KEY),
		Password: GetConfigString(PG_PASS_KEY),
		Host:     GetConfigString(PG_HOST_KEY),
		Port:     GetConfigString(PG_PORT_KEY),
		Name:     GetConfigString(PG_NAME_KEY),
	}
}

// to be removed after this code goes in production
func GetOldPGConfig() *PGConnectionInfo {
	return &PGConnectionInfo{
		User:     GetConfigString("OLD_" + PG_USER_KEY),
		Password: GetConfigString("OLD_" + PG_PASS_KEY),
		Host:     GetConfigString("OLD_" + PG_HOST_KEY),
		Port:     GetConfigString("OLD_" + PG_PORT_KEY),
		Name:     GetConfigString("OLD_" + PG_NAME_KEY),
	}
}

func GetRedisConfig() *RedisConnectionInfo {
	return &RedisConnectionInfo{
		Password: GetConfigString(REDIS_PASS_KEY),
		Host:     GetConfigString(REDIS_HOST_KEY),
		Port:     GetConfigString(REDIS_PORT_KEY),
	}
}

func GetAWSS3Config() *AWSS3Info {
	return &AWSS3Info{
		AccessKeyID:     GetConfigString(S3_ACCESS_KEY_ID),
		AccessSecretKey: GetConfigString(S3_SECRET_ACCESS_KEY),
		BucketName:      GetConfigString(S3_BUCKET_NAME),
		Region:          GetConfigString(S3_REGION),
		CloudFrontURL:   GetConfigString(S3_CLOUDFRONT_URL),
	}
}

func GetDiscordBotConfig() *DiscordBotInfo {
	return &DiscordBotInfo{
		AppID:        GetConfigString(DISCORD_APP_ID),
		PublicKey:    GetConfigString(DISCORD_PUBLIC_KEY),
		ClientID:     GetConfigString(DISCORD_CLIENT_ID),
		ClientSecret: GetConfigString(DISCORD_CLIENT_SECRET),
		Token:        GetConfigString(DISCORD_TOKEN),
		Permission:   GetConfigString(DISCORD_PERMISSION),
	}
}

func GetKafkaConfig() *KakfaConnectionInfo {
	return &KakfaConnectionInfo{
		Hosts: GetConfigString(KAFKA_HOSTS_KEY),
	}
}

func GetAlgoliaConfig() *AlgoliaInfo {
	return &AlgoliaInfo{
		AppID:       GetConfigString(ALGOLIA_APP_ID),
		AdminAPIKey: GetConfigString(ALGOLIA_ADMIN_API_KEY),
	}
}

func GetSlackConfig() *SlackInfo {
	return &SlackInfo{
		AppID:         GetConfigString(SLACK_APP_ID),
		ClientID:      GetConfigString(SLACK_CLIENT_ID),
		ClientSecret:  GetConfigString(SLACK_CLIENT_SECRET),
		SigningSecret: GetConfigString(SLACK_SIGNING_SECRET),
		Nonce:         GetConfigString(SLACK_NONCE),
	}
}

func GetSecretShortKey() string {
	return GetConfigString(SECRET_SHORT_KEY)
}
func GetGitConfig() *GitInfo {
	hostsStr := GetConfigString(GIT_HOSTS)
	hosts := strings.Split(hostsStr, ",")
	return &GitInfo{
		Hosts:  hosts,
		Secret: GetConfigString(GIT_SECRET_KEY),
	}
}

func GetStreamConfig() *GetStreamConnectionInfo {
	return &GetStreamConnectionInfo{
		ApiKey:    GetConfigString(GET_STREAM_API_KEY),
		ApiSecret: GetConfigString(GET_STREAM_API_SECRET),
	}
}

func GetDiscordLambdaServerEndpoint() string {
	return GetConfigString(DISCORD_LAMBDA_SERVER_ENDPOINT)
}

func GetSlackLambdaServerEndpoint() string {
	return GetConfigString(SLACK_LAMBDA_SERVER_ENDPOINT)
}

func GetAppInfoConfig() *AppInfo {
	return &AppInfo{
		FrontendHost: GetConfigString(FRONTEND_HOST),
		BackendHost:  GetConfigString(BACKEND_HOST),
	}
}

func GetSupabaseConfig() *SupabaseInfo {
	return &SupabaseInfo{
		SupabaseBaseurl: GetConfigString(SUPABASE_BASEURL),
		SupabaseToken:   GetConfigString(SUPABASE_TOKEN),
	}
}

// Returns right email based on ENV
func GetCurrentSystemEmail() string {
	if IsLive() {
		return PRODUCTION_EMAIL
	}
	return STAGE_EMAIL
}
