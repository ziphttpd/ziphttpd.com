package firebase

import "google.golang.org/api/option"

var (
	opt option.ClientOption
)

// GetOption はOptionを生成します。(未設定の場合は admin.json を読みます)
func GetOption() option.ClientOption {
	if opt == nil {
		SetOption("admin.json")
	}
	return opt
}

// SetOption は firebase のオプションを設定します。
func SetOption(credFile string) {
	opt = option.WithCredentialsFile(credFile)
}
