package auth

import (
	"fmt"

	"github.com/sev-2/raiden"
	"github.com/sev-2/raiden/pkg/connector/pgmeta"
	"github.com/sev-2/raiden/pkg/supabase/drivers/local/meta"
	"github.com/sev-2/raiden/pkg/supabase/objects"
	"github.com/sev-2/raiden/pkg/supabase/query/sql"
)

var GetUserByEmailQuery = "SELECT email FROM auth.users WHERE email = %s"

func generateGetUserQuery(email string) string {
	return fmt.Sprintf(GetUserByEmailQuery, sql.Literal(email))
}

var UpdateUserRecoveryTokenQuery = "UPDATE auth.users SET recovery_token = %s, recovery_sent_at = now() WHERE email = %s"

func generateUpdateUserRecoveryTokenQuery(token string, email string) string {
	return fmt.Sprintf(UpdateUserRecoveryTokenQuery, sql.Literal(token), sql.Literal(email))
}

func getBaseUrl(cfg *raiden.Config) string {
	return fmt.Sprintf("%s%s", cfg.SupabaseApiUrl, cfg.SupabaseApiBasePath)
}

func GetUserByEmail(cfg *raiden.Config, email string) (result objects.User, err error) {
	q := generateGetUserQuery(email)
	rs, err := pgmeta.ExecuteQuery[[]objects.User](getBaseUrl(cfg), q, nil, meta.DefaultInterceptor(cfg), nil)
	if err != nil {
		err = fmt.Errorf("get email from error : %s", err)
		return
	}

	if len(rs) == 0 {
		err = fmt.Errorf("get email %s is not found", email)
		return
	}
	return rs[0], nil
}

func UpdateUserRecoveryToken(cfg *raiden.Config, email string, token string) error {
	sql := generateUpdateUserRecoveryTokenQuery(token, email)
	_, err := pgmeta.ExecuteQuery[any](getBaseUrl(cfg), sql, nil, meta.DefaultInterceptor(cfg), nil)
	if err != nil {
		return fmt.Errorf("update user recovery token %s error : %s", email, err)
	}
	return nil
}
