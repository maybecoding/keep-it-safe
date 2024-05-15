package pgstore

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/maybecoding/keep-it-safe/internal/server/config"
	"github.com/maybecoding/keep-it-safe/internal/server/core/entity"
	"github.com/maybecoding/keep-it-safe/pkg/postgres"
	"github.com/stretchr/testify/require"
)

func TestUser(t *testing.T) {
	// Prepare store
	cfg := &config.Config{
		DB: config.DB{Path: "postgres://api:pwd@localhost:5432/keep_it_safe?sslmode=disable"},
	}
	pg, err := postgres.New(cfg.DB.Path)
	require.NoError(t, err)
	pgs := New(pg)

	ctx := context.Background()
	srcUsr := entity.User{UserLogin: "TestUser1", UserPasswordHash: "strong password hash"}
	srcUsr2 := entity.User{UserLogin: "TestUser2", UserPasswordHash: "strong password hash 2"}

	pgs.WithTx(ctx, func(ctx context.Context) error {
		t.Run("creating user", func(t *testing.T) {
			usr, err := pgs.UserNew(ctx, srcUsr.UserLogin, srcUsr.UserPasswordHash)
			require.NoError(t, err)
			require.NotNil(t, usr)
			srcUsr.UserID = usr.UserID
			require.Equal(t, srcUsr, *usr)

			_, err = pgs.UserNew(ctx, srcUsr2.UserLogin, srcUsr2.UserPasswordHash)
			require.NoError(t, err)
		})
		t.Run("creating same user", func(t *testing.T) {
			usrSame, err := pgs.UserNew(ctx, srcUsr.UserLogin, srcUsr.UserPasswordHash)
			require.Nil(t, usrSame)

			require.NotNil(t, err)
			require.ErrorIs(t, err, entity.ErrUserNotAvailable)
		})
		return errors.New("we must rollback transaction")
	})

	pgs.WithTx(ctx, func(ctx context.Context) error {
		t.Run("user available", func(t *testing.T) {
			ok, err := pgs.LoginAvailable(ctx, srcUsr.UserLogin)
			require.NoError(t, err)
			require.Equal(t, ok, true)

			usr, err := pgs.UserNew(ctx, srcUsr.UserLogin, srcUsr.UserPasswordHash)
			require.NoError(t, err)
			srcUsr.UserID = usr.UserID

			ok, err = pgs.LoginAvailable(ctx, srcUsr.UserLogin)
			require.NoError(t, err)
			require.Equal(t, ok, false)
		})

		t.Run("get user", func(t *testing.T) {
			usr, err := pgs.UserGet(ctx, srcUsr.UserLogin)
			require.NoError(t, err)
			require.Equal(t, srcUsr, *usr)
		})

		return errors.New("we must rollback transaction")
	})
}

func TestSecret(t *testing.T) {
	// Prepare store
	cfg := &config.Config{
		DB: config.DB{Path: "postgres://api:pwd@localhost:5432/keep_it_safe?sslmode=disable"},
	}
	pg, err := postgres.New(cfg.DB.Path)
	require.NoError(t, err)
	pgs := New(pg)

	_ = pgs

	ctx := context.Background()
	pgs.WithTx(ctx, func(ctx context.Context) error {
		srcSecret := entity.SecretDetail{
			Secret: entity.Secret{
				Type: 0,
				Name: "Secret1",
			},
			Value:        []byte{1, 2, 3, 4, 5, 6, 7, 8, 9},
			Nonce:        []byte{},
			EncryptionSK: []byte{9, 8, 7, 6, 5, 4, 3, 2, 1},
		}
		t.Run("secret set", func(t *testing.T) {
			// first create user
			usr, err := pgs.UserNew(ctx, "login2", "hash2")
			require.NoError(t, err)
			srcSecret.UserID = usr.UserID
			// Create new secret

			secretID, err := pgs.SecretSet(ctx, srcSecret)
			require.NoError(t, err)
			srcSecret.ID = secretID
		})

		t.Run("secret get", func(t *testing.T) {
			// Check inputed secret
			secret, err := pgs.SecretGet(ctx, srcSecret.ID)
			require.NoError(t, err)
			require.Equal(t, srcSecret.UserID, secret.UserID)
			require.Equal(t, srcSecret.Type, secret.Type)
			require.Equal(t, srcSecret.Name, secret.Name)
			require.Equal(t, srcSecret.Value, secret.Value)
			require.Equal(t, srcSecret.EncryptionSK, secret.EncryptionSK)
		})

		t.Run("secret list", func(t *testing.T) {
			// insert another one secret
			srcSecret2 := srcSecret
			srcSecret2.Name = "Secret2"
			secretID, err := pgs.SecretSet(ctx, srcSecret2)
			require.NoError(t, err)
			srcSecret2.ID = secretID

			// insert same secret
			secretID, err = pgs.SecretSet(ctx, srcSecret2)
			require.NoError(t, err)
			require.Equal(t, srcSecret2.ID, secretID)

			// get list of secrets
			list, err := pgs.SecretList(ctx, srcSecret.UserID)
			require.NoError(t, err)
			require.Len(t, list, 2, "При повторной вставке с тем же usr_id, name не должно быть задвоений")
			for i, exp := range []entity.SecretDetail{srcSecret, srcSecret2} {
				act := list[i]
				require.Equal(t, exp.ID, act.ID)
				require.Equal(t, exp.UserID, act.UserID)
				require.Equal(t, exp.Type, act.Type)
			}
		})

		return errors.New("we must rollback transaction")
	})
}

func TestSecretAttr(t *testing.T) {
	// Prepare store
	cfg := &config.Config{
		DB: config.DB{Path: "postgres://api:pwd@localhost:5432/keep_it_safe?sslmode=disable"},
	}
	pg, err := postgres.New(cfg.DB.Path)
	require.NoError(t, err)
	pgs := New(pg)

	_ = pgs

	ctx := context.Background()
	pgs.WithTx(ctx, func(ctx context.Context) error {
		t.Run("attr insert", func(t *testing.T) {
			// first insert user
			usr, err := pgs.UserNew(ctx, "user_with_secretWithMeta", "hash")
			require.NoError(t, err)

			// first just insert secret
			secret := sampleSecret("secret with meta")
			secret.UserID = usr.UserID
			secretID, err := pgs.SecretSet(ctx, secret)
			require.NoError(t, err)

			metaForInsert := entity.SecretMeta([]entity.SecretAttr{
				{Attr: "Attr1", Value: "Value1"},
				{Attr: "Very very very very very very long attr name", Value: "Long logng logng logng logng logng Value2"},
				{Attr: "Attr3", Value: "Value3"},
			})
			// insert secret meta
			for _, attrValue := range metaForInsert {
				err = pgs.SecretAttrSet(ctx, secretID, attrValue)
				require.NoError(t, err)
			}

			// check insertion
			secretDest, err := pgs.SecretGet(ctx, secretID)
			fmt.Println(secretDest.Meta)
			require.NoError(t, err)
			require.Len(t, secretDest.Meta, len(metaForInsert))
			require.Equal(t, metaForInsert, secretDest.Meta)
		})

		return errors.New("we must rollback transaction")
	})
}

func sampleSecret(name string) entity.SecretDetail {
	return entity.SecretDetail{
		Secret: entity.Secret{
			Type: 0,
			Name: entity.SecretName(name),
		},
		Value:        []byte{1, 2, 3, 4, 5, 6, 7, 8, 9},
		Nonce:        []byte{},
		EncryptionSK: []byte{9, 8, 7, 6, 5, 4, 3, 2, 1},
	}
}
