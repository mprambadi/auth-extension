package module

import (
	"github.com/sev-2/raiden"
	"github.com/sev-2/raiden/examples/raiden-auth-module/pkg/auth"
	"github.com/sev-2/raiden/examples/raiden-auth-module/pkg/controller"
	"github.com/valyala/fasthttp"
)

type AuthExtentionModule struct {
}

func (m *AuthExtentionModule) Routes() []*raiden.Route {
	return []*raiden.Route{
		{
			Type:       raiden.RouteTypeCustom,
			Path:       "/auth/v1/recover",
			Controller: &controller.RecoverController{},
			Methods:    []string{fasthttp.MethodPost},
		},
	}
}

func (m *AuthExtentionModule) Libs() []func(config *raiden.Config) any {
	return []func(config *raiden.Config) any{
		auth.NewLibrary,
	}
}
