//go:generate go run github.com/deepmap/oapi-codegen/v2/cmd/oapi-codegen --config=../../../../../pkg/api/v1/models.cfg.yaml ../../../../../pkg/api/v1/api.yaml
//go:generate go run github.com/deepmap/oapi-codegen/v2/cmd/oapi-codegen --config=../../../../../pkg/api/v1/cfg.yaml ../../../../../pkg/api/v1/api.yaml

package api

import (
	"context"
	"errors"

	"github.com/maybecoding/keep-it-safe/internal/server/adapters/api/v1/models"
	"github.com/maybecoding/keep-it-safe/internal/server/core/entity"
	"github.com/maybecoding/keep-it-safe/internal/server/core/services/secret"
	"github.com/maybecoding/keep-it-safe/internal/server/core/services/user"
)

type API struct {
	user   *user.Service
	secret *secret.Service
}

func New(user *user.Service, secret *secret.Service) *API {
	return &API{user: user, secret: secret}
}

var _ StrictServerInterface = (*API)(nil)

// Login user
// (POST /login)
func (a *API) Login(ctx context.Context, request LoginRequestObject) (LoginResponseObject, error) {
	if request.Body == nil {
		return Login400Response{}, nil
	}
	token, err := a.user.Login(ctx, entity.UserLogin(request.Body.Login), entity.UserPassword(request.Body.Password))
	if err != nil {
		if errors.Is(err, entity.ErrUserNotFound) || errors.Is(err, entity.ErrIncorrectPassword) {
			return Login401JSONResponse{}, nil
		}
		return Login500JSONResponse{}, nil
	}
	cookie := authCookie(token)
	return Login200Response{Login200ResponseHeaders{SetCookie: cookie}}, nil
}

// (POST /register)
// Register new user
// (POST /register)
func (a *API) Register(ctx context.Context, request RegisterRequestObject) (RegisterResponseObject, error) {
	if request.Body == nil {
		return Register400Response{}, nil
	}
	token, err := a.user.Register(ctx, entity.UserLogin(request.Body.Login), entity.UserPassword(request.Body.Password))
	if err != nil {
		if errors.Is(err, entity.ErrUserNotAvailable) {
			return Register409Response{}, nil
		}
		return Register500Response{}, nil

	}
	cookie := authCookie(token)
	return Register200Response{Register200ResponseHeaders{SetCookie: cookie}}, nil
}

// Get list of secrets of user
// (GET /secrets)
func (a *API) SecretList(ctx context.Context, request SecretListRequestObject) (SecretListResponseObject, error) {
	userID, mErr := a.getUserID(request.Params.Authorization)
	if mErr != nil {
		return SecretList401JSONResponse(*mErr), nil
	}

	list, err := a.secret.List(ctx, userID)
	if err != nil {
		return SecretList500JSONResponse(models.Error{Error: err.Error()}), nil
	}
	resp := make(SecretList200JSONResponse, 0, len(list))
	for _, s := range list {
		resp = append(resp, models.Secret{
			Created: s.Created,
			Id:      int32(s.ID),
			Name:    string(s.Name),
			Type:    int32(s.Type),
			Updated: s.Updated,
			UserId:  int32(s.UserID),
		})
	}

	return resp, nil
}

// Creates new secret of user
// (POST /secrets)
func (a *API) SecretSet(ctx context.Context, request SecretSetRequestObject) (SecretSetResponseObject, error) {
	userID, mErr := a.getUserID(request.Params.Authorization)
	if mErr != nil {
		return SecretSet401JSONResponse(*mErr), nil
	}
	b := request.Body
	if b == nil ||
		b.SecretType == int32(entity.SecretTypeCredentials) && b.Credentials == nil ||
		b.SecretType == int32(entity.SecretTypeText) && b.Text == nil ||
		b.SecretType == int32(entity.SecretTypeBinary) && b.Binary == nil ||
		b.SecretType == int32(entity.SecretTypeBankCard) && b.BankCard == nil {
		return SecretSet400Response{}, nil
	}

	data := entity.Data{
		SecretName: entity.SecretName(b.SecretName),
		SecretType: entity.SecretType(b.SecretType),
	}
	if b.SecretMeta != nil {
		for _, m := range *b.SecretMeta {
			data.SecretMeta = append(data.SecretMeta, entity.SecretAttr{Attr: m.Attr, Value: m.Value})
		}
	}
	if b.Credentials != nil {
		data.Credentials = &entity.DataCredentials{Login: b.Credentials.Login, Password: b.Credentials.Password}
	}

	if b.Text != nil {
		dt := entity.DataText(*b.Text)
		data.Text = &dt
	}
	if b.Binary != nil {
		db := entity.DataBinary(*b.Binary)
		data.Binary = &db
	}

	if b.BankCard != nil {
		data.BankCard = &entity.DataBankCard{
			Number:         b.BankCard.Number,
			Valid:          b.BankCard.Valid,
			Holder:         b.BankCard.Holder,
			ValidationCode: b.BankCard.ValidationCode,
		}
	}

	_, err := a.secret.Set(ctx, userID, data)
	if err != nil {
		return SecretSet500JSONResponse(models.Error{Error: err.Error()}), nil
	}
	return SecretSet200Response{}, nil
}

// Get secret by id
// (GET /secrets/{secret_id})
func (a *API) SecretGetByID(ctx context.Context, request SecretGetByIDRequestObject) (SecretGetByIDResponseObject, error) {
	userID, mErr := a.getUserID(request.Params.Authorization)
	if mErr != nil {
		return SecretGetByID401JSONResponse(*mErr), nil
	}

	data, err := a.secret.GetByID(ctx, userID, entity.SecretID(request.SecretId))
	if err != nil {
		if errors.Is(err, entity.ErrSecretNotFound) || errors.Is(err, entity.ErrSecretForbiden) {
			return SecretGetByID404Response{}, nil
		}
		return SecretGetByID500JSONResponse(models.Error{Error: err.Error()}), nil
	}
	resp := models.Data{
		SecretName: string(data.SecretName),
		SecretType: int32(data.SecretType),
	}
	if data.Credentials != nil {
		resp.Credentials = &models.DataCredentials{Login: data.Credentials.Login, Password: data.Credentials.Password}
	}
	if data.Text != nil {
		txt := string(*data.Text)
		resp.Text = &txt
	}
	if data.Binary != nil {
		b := []byte(*data.Binary)
		resp.Binary = &b
	}
	if data.BankCard != nil {
		resp.BankCard = &models.DataBankCard{
			Number:         data.BankCard.Number,
			Valid:          data.BankCard.Valid,
			Holder:         data.BankCard.Holder,
			ValidationCode: data.BankCard.ValidationCode,
		}
	}
	if data.SecretMeta != nil && len(data.SecretMeta) > 0 {
		meta := make([]models.SecretAttr, 0, len(data.SecretMeta))
		for _, attr := range data.SecretMeta {
			meta = append(meta, models.SecretAttr{Attr: attr.Attr, Value: attr.Value})
		}
		resp.SecretMeta = &meta
	}

	return SecretGetByID200JSONResponse(resp), err
}
