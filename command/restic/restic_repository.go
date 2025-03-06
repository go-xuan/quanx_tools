package restic

import (
	"fmt"
	"runtime"
	"strings"
)

// Repository 存储库
type Repository struct {
	Uri         string `json:"uri" yaml:"uri"`                   // 存储库uri（必填）
	Password    string `json:"password" yaml:"password"`         // 存储库密码（必填）
	AwsUser     string `json:"aws_user" yaml:"aws_user"`         // AmazonS3秘钥ID（本地存储库可忽略）
	AwsPassword string `json:"aws_password" yaml:"aws_password"` // AmazonS3秘钥（本地存储库可忽略）
}

// SetEnv 设置环境变量
func (r *Repository) SetEnv() string {
	var repoEnv, awsEnv string
	switch runtime.GOOS {
	case "windows":
		repoEnv = fmt.Sprintf(`set RESTIC_REPOSITORY=%s&&set RESTIC_PASSWORD=%s&&`, r.Uri, r.Password)
		awsEnv = fmt.Sprintf(`set AWS_ACCESS_KEY_ID=%s&&set AWS_SECRET_ACCESS_KEY=%s&&`, r.AwsUser, r.AwsPassword)
	default:
		repoEnv = fmt.Sprintf(`export RESTIC_REPOSITORY="%s"; export RESTIC_PASSWORD="%s"; `, r.Uri, r.Password)
		awsEnv = fmt.Sprintf(`export AWS_ACCESS_KEY_ID="%s"; export AWS_SECRET_ACCESS_KEY="%s"; `, r.AwsUser, r.AwsPassword)
	}
	if strings.HasPrefix(r.Uri, "s3:") {
		return repoEnv + awsEnv
	} else {
		return repoEnv
	}
}
