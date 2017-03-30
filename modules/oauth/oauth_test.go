package oauth

import (
	"testing"

	"github.com/ryankurte/authplz/controllers/datastore"
	"github.com/ryankurte/authplz/test"
)

func TestOauth(t *testing.T) {

	ts, err := test.NewTestServer()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	config := Config{
		TokenSecret: "reasonable-test-secret-here-plz",
	}

	oauthModule, err := NewController(ts.DataStore, config)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	u, err := ts.DataStore.AddUser(test.FakeEmail, test.FakeName, test.FakePass)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	user := u.(*datastore.User)

	scopes := []string{"public.read", "public.write", "private.read", "private.write"}
	redirects := []string{"https://fake-redirect.cows"}
	responses := []string{"token"}
	grants := []string{"client_credential"}

	t.Run("Users can create specified grant types", func(t *testing.T) {
		for _, g := range AdminGrantTypes {
			_, err := oauthModule.CreateClient(user.GetExtID(), scopes, redirects, []string{g}, responses, true)
			if arrayContains(UserGrantTypes, g) && err != nil {
				t.Error(err)
			}
			if !arrayContains(UserGrantTypes, g) && err == nil {
				t.Errorf("Unexpected allowed grant type: %s", g)
			}
		}
	})

	t.Run("Admins can create all grant types", func(t *testing.T) {
		user.SetAdmin(true)
		ts.DataStore.UpdateUser(user)

		for _, g := range AdminGrantTypes {
			c, err := oauthModule.CreateClient(user.GetExtID(), scopes, redirects, []string{g}, responses, true)
			if err != nil {
				t.Error(err)
			} else if c == nil {
				t.Errorf("Nil client returned")
			}
		}

		user.SetAdmin(false)
		ts.DataStore.UpdateUser(user)
	})

	t.Run("Users can only create valid scopes", func(t *testing.T) {
		scopes := []string{"FakeScope"}
		_, err := oauthModule.CreateClient(user.GetExtID(), scopes, redirects, grants, responses, true)
		if err == nil {
			t.Errorf("Unexpected allowed scope: %s", scopes)
		}
	})

}