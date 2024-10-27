package server

import (
	"context"
	"fmt"
	"os"

	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
)

type Authenticator struct {
	*oidc.Provider
	oauth2.Config
}

func New() (*Authenticator, error) {
	provider, err := oidc.NewProvider(context.Background(), "https://"+os.Getenv("AUTH0_DOMAIN")+"/")
	if err != nil {
		return nil, err
	}

	return &Authenticator{
		Provider: provider,
		Config: oauth2.Config{
			ClientID:     os.Getenv("AUTH0_RA_CLIENT_ID"),
			ClientSecret: os.Getenv("AUTH0_RA_CLIENT_SECRET"),
			RedirectURL:  os.Getenv("APP_URL") + "/callback",
			Endpoint:     provider.Endpoint(),
			Scopes:       []string{oidc.ScopeOpenID, "profile"},
		},
	}, nil
}

func (a *Authenticator) VerifyIDToken(ctx context.Context, token *oauth2.Token) (*oidc.IDToken, error) {
	rawIDToken, ok := token.Extra("id_token").(string)
	if !ok {
		return nil, fmt.Errorf("no id_token in token")
	}

	oidcConfig := &oidc.Config{
		ClientID: a.ClientID,
	}

	return a.Verifier(oidcConfig).Verify(ctx, rawIDToken)
}
