// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See License.txt for license information.

package oauthoidc

import (
	"encoding/json"
	"github.com/mattermost/mattermost-server/einterfaces"
	"github.com/mattermost/mattermost-server/model"
	"io"
	"strconv"
	"strings"
)

type OidcProvider struct {
}

type OidcUser struct {
	Id      	int64  `json:"sub"`
	Name       	string `json:"name"`
	GivenName  	string `json:"given_name"`
	FamilyName 	string `json:"family_name"`
	Email      	string `json:"email"`
}

func init() {
	provider := &OidcProvider{}
	einterfaces.RegisterOauthProvider(model.USER_AUTH_SERVICE_OIDC, provider)
}


func userFromOidcUser(glu *OidcUser) *model.User {
	user := &model.User{}
	username := glu.Name
	if username == "" {
		username = glu.Email
	}
	user.Username = model.CleanUsername(username)
	user.FirstName = glu.GivenName
	user.LastName = glu.FamilyName
	strings.TrimSpace(user.Email)
	user.Email = glu.Email
	userId := strconv.FormatInt(glu.Id, 10)
	user.AuthData = &userId
	user.AuthService = model.USER_AUTH_SERVICE_OIDC

	return user
}

func oidcUserFromJson(data io.Reader) *OidcUser {
	decoder := json.NewDecoder(data)
	var glu OidcUser
	err := decoder.Decode(&glu)
	if err == nil {
		return &glu
	} else {
		return nil
	}
}

func (glu *OidcUser) ToJson() string {
	b, err := json.Marshal(glu)
	if err != nil {
		return ""
	} else {
		return string(b)
	}
}

func (glu *OidcUser) IsValid() bool {
	if glu.Id == 0 {
		return false
	}

	if len(glu.Email) == 0 {
		return false
	}

	return true
}

func (glu *OidcUser) getAuthData() string {
	return strconv.FormatInt(glu.Id, 10)
}

func (m *OidcProvider) GetIdentifier() string {
	return model.USER_AUTH_SERVICE_OIDC
}

func (m *OidcProvider) GetUserFromJson(data io.Reader) *model.User {
	glu := oidcUserFromJson(data)
	if glu.IsValid() {
		return userFromOidcUser(glu)
	}

	return &model.User{}
}

func (m *OidcProvider) GetAuthDataFromJson(data io.Reader) string {
	glu := oidcUserFromJson(data)

	if glu.IsValid() {
		return glu.getAuthData()
	}

	return ""
}
