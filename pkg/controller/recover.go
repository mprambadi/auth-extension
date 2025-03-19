package controller

import (
	"fmt"

	"github.com/mprambadi/raiden-auth-module/pkg/auth"

	"github.com/sev-2/raiden"
)

type RecoverPayload struct {
	Email string `json:"email"`
}

type RecoverResponse struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
}

type RecoverController struct {
	raiden.ControllerBase
	Payload *RecoverPayload
	Result  RecoverResponse
}

func (c *RecoverController) Post(ctx raiden.Context) error {
	raiden.Info("/auth/v1/recover params", ctx.RequestContext().QueryArgs().String())
	authExtension := auth.AuthExtension{}

	if err := ctx.ResolveLibrary(&authExtension); err != nil {
		return err
	}
	webUrl := ctx.Config().GetString("WEB_URL")
	redirectUri := fmt.Sprintf("%s/auth/verify", webUrl)

	if err := authExtension.Recover(c.Payload.Email, redirectUri); err != nil {
		return err
	}
	raiden.Info("Password recovery email sent successfully:")

	c.Result.Status = true
	c.Result.Message = "Password recovery email sent successfully"

	return ctx.SendJson(&c.Result)
}
