package provider

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"golang.org/x/oauth2"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/99designs/keyring"
	log "github.com/sirupsen/logrus"
	"github.com/lordbyron/oauth2-auth-cli"
	"github.com/mulesoft-labs/aws-keycloak/provider/saml"
)

const (
	keycloakCookie    = "KEYCLOAK_IDENTITY"
	keycloakSamlPath  = "/auth/realms/Mulesoft/protocol/saml/clients/"
	keycloakAuthPath  = "/auth/realms/Mulesoft/protocol/openid-connect/auth"
	keycloakTokenPath = "/auth/realms/Mulesoft/protocol/openid-connect/token"
)

type KeycloakProviderIf interface {
	RetrieveKeycloakCreds() bool
	BrowserAuth() error
	BasicAuth() error
	GetSamlAssertion() (saml.SAMLStruct, error)
	StoreKeycloakCreds()
}

type KeycloakProvider struct {
	Keyring         keyring.Keyring
	ProfileName     string
	ApiBase         string
	SamlPath        string
	AwsClient       string
	AwsClientSecret string
	kcToken         string
	kcCreds         KeycloakCreds
}

type KeycloakCreds struct {
	Username string
	Password string
}

type KeycloakUserAuthn struct {
	AccessToken           string `json:"access_token"`
	ExpiresIn             int    `json:"expires_in"`
	RefreshTokenExpiresIn int    `json:"refresh_expires_in"`
	RefreshToken          string `json:"refresh_token"`
	TokenType             string `json:"token_type"`
	SessionState          string `json:"session_state"`
}

func NewKeycloakProvider(kr keyring.Keyring, kcprofile string, kcConf map[string]string) (*KeycloakProvider, error) {
	k := KeycloakProvider{
		Keyring:     kr,
		ProfileName: kcprofile,
	}
	if v, e := kcConf["keycloak_base"]; e {
		k.ApiBase = v
	} else {
		return nil, errors.New("Config must specify keycloak_base")
	}
	if v, e := kcConf["aws_client_id"]; e {
		k.AwsClient = v
	} else {
		return nil, errors.New("Config must specify aws_client_id")
	}
	if v, e := kcConf["aws_client_secret"]; e {
		k.AwsClientSecret = v
	} else {
		return nil, errors.New("Config must specify aws_client_secret")
	}
	return &k, nil
}

/**
 * return bool is whether the creds should be stored in keyring if they work
 */
func (k *KeycloakProvider) RetrieveKeycloakCreds() bool {
	var keycloakCreds KeycloakCreds
	keyName := k.keycloakkeyname()

	item, err := k.Keyring.Get(keyName)
	if err == nil {
		log.Debug("found creds in keyring")
		if err = json.Unmarshal(item.Data, &keycloakCreds); err != nil {
			log.Error("could not unmarshal keycloak creds")
		} else {
			k.kcCreds = keycloakCreds
			return false
		}
	} else {
		log.Debugf("couldnt get keycloak creds from keyring: %s", keyName)
		k.kcCreds = k.promptUsernamePassword()
	}
	return true
}

func (k *KeycloakProvider) StoreKeycloakCreds() {
	encoded, err := json.Marshal(k.kcCreds)
	// failure would be surprising, but just dont save
	if err != nil {
		log.Debugf("Couldn't marshal keycloak creds... %s", err)
	} else {
		keyName := k.keycloakkeyname()
		newKeycloakItem := keyring.Item{
			Key:   keyName,
			Data:  encoded,
			Label: keyName + " credentials",
			KeychainNotTrustApplication: false,
		}
		if err := k.Keyring.Set(newKeycloakItem); err != nil {
			log.Debugf("Failed to write keycloak creds to keyring!")
		} else {
			log.Debugf("Successfully stored keycloak creds to keyring!")
		}
	}
}

func (k *KeycloakProvider) promptUsernamePassword() (creds KeycloakCreds) {
	fmt.Fprintf(ProviderOut, "Enter username/password for keycloak (env: %s)\n", k.ProfileName)
	for creds.Username == "" {
		u, err := Prompt("Username", false)
		if err != nil {
			fmt.Fprintf(ProviderOut, "Invalid username: %s\n", creds.Username)
		} else {
			creds.Username = u
		}
	}
	for creds.Password == "" {
		x, err := Prompt("Password", true)
		if err != nil {
			fmt.Fprintf(ProviderOut, "Invalid password: %s\n", creds.Username)
		} else {
			creds.Password = x
		}
	}
	fmt.Fprint(ProviderOut, "\n")
	return
}

func (k *KeycloakProvider) keycloakkeyname() string {
	return "keycloak-creds-" + k.ProfileName
}

/**
 * Initiate OAuth2 Authorization Grant flow
 */
func (k *KeycloakProvider) BrowserAuth() error {
	oauth := &oauth2.Config{
		ClientID:     k.AwsClient,
		ClientSecret: k.AwsClientSecret,
		Endpoint: oauth2.Endpoint{
			AuthURL:  fmt.Sprintf("%s%s", k.ApiBase, keycloakAuthPath),
			TokenURL: fmt.Sprintf("%s%s", k.ApiBase, keycloakTokenPath),
		},
	}
	o := o2cli.Oauth2CLI{
		Conf: oauth,
		Log:  log.StandardLogger(),
	}
	token, err := o.Authorize()
	if err == nil {
		k.kcToken = token.AccessToken
	}
	return err
}

/**
 * Deprecated
 * Must populate kcCreds before calling (eg. by calling RetrieveKeycloakCreds)
 */
func (k *KeycloakProvider) BasicAuth() error {
	payload := url.Values{}
	payload.Set("username", k.kcCreds.Username)
	payload.Set("password", k.kcCreds.Password)
	payload.Set("client_id", k.AwsClient)
	payload.Set("client_secret", k.AwsClientSecret)
	payload.Set("grant_type", "password")

	header := http.Header{
		"Accept":       []string{"application/json"},
		"Content-Type": []string{"application/x-www-form-urlencoded"},
	}

	body, err := k.doHttp("POST", keycloakTokenPath, header, []byte(payload.Encode()))
	if err != nil {
		return nil
	}

	var userAuthn KeycloakUserAuthn
	err = json.Unmarshal(body, &userAuthn)
	if err != nil {
		return err
	}
	log.Debug("successfully authenticated to keycloak")
	k.kcToken = userAuthn.AccessToken
	return nil
}

func (k *KeycloakProvider) GetSamlAssertion() (samlStruct saml.SAMLStruct, err error) {
	header := http.Header{
		"Cookie": []string{fmt.Sprintf("%s=%s", keycloakCookie, k.kcToken)},
	}
	body, err := k.doHttp("GET", keycloakSamlPath+k.AwsClient, header, nil)
	if err != nil {
		return
	}

	if err = saml.Parse(body, &samlStruct); err != nil {
		err = fmt.Errorf("Couldn't access SAML app; is the user %s in a group that has access to AWS? (%s)", k.kcCreds.Username, err)
	}
	return
}

func (k *KeycloakProvider) doHttp(method, path string, header http.Header, data []byte) (body []byte, err error) {
	url, err := url.Parse(fmt.Sprintf("%s/%s", k.ApiBase, path))
	if err != nil {
		return
	}

	req := &http.Request{
		Method: method,
		URL:    url,
		Header: header,
		Body:   ioutil.NopCloser(bytes.NewReader(data)),
	}

	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		err = fmt.Errorf("%s %v: %s", method, url, res.Status)
		return
	}

	body, err = ioutil.ReadAll(res.Body)
	return
}
