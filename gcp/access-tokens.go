package gcp

import (
	"fmt"
	"strings"
	"time"
	"github.com/fatih/color"
	"github.com/BishopFox/cloudfox/internal/gcp"
	"github.com/BishopFox/cloudfox/internal"
	"github.com/kyokomi/emoji"
	"github.com/BishopFox/cloudfox/globals"
)

type AccessTokensModule struct {
	// Filtering data
	Organizations []string
	Folders []string
	Projects []string
}

func (m *AccessTokensModule) PrintAccessTokens(version string, outputFormat string, outputDirectory string, verbosity int) error {
	fmt.Printf("[%s][%s] Enumerating gcloud access tokens (%s, %s)...\n", color.CyanString(emoji.Sprintf(":fox:cloudfox %s :fox:", version)), color.CyanString(globals.GCP_WHOAMI_MODULE_NAME), color.CyanString("default user token"), color.RedString("application-default token"))
	tokens := gcp.ReadRefreshTokens()
	accessTokens := gcp.ReadAccessTokens()
	applicationdefaulthash := gcp.GetDefaultApplicationHash()
	activeAccount := gcp.GetActiveAccount()

	var tableHead = []string{"Account", "Credential Type", "Validity", "Refresh Token", "Refresh Token needs password"}
	var tableBody [][]string
	for _, accessToken := range accessTokens {
		// prepare token's associated ID (email or hash)
		var emailinfo string
		if (accessToken.AccountID == activeAccount) {
			emailinfo = color.BlueString(accessToken.AccountID)
		} else {
			emailinfo = accessToken.AccountID
		}

		// format expiration time
		Exp, _ := time.Parse(time.RFC3339, accessToken.TokenExpiry)
		var timeinfo string
		if Exp.After(time.Now()) {
			timeinfo = Exp.Format(time.RFC1123)
		} else {
			timeinfo = "EXPIRED"
		}

		// display token type, user or application credential
		var tokentype string
		if (strings.Contains(accessToken.AccountID, "@")) {
			tokentype = "User"
		} else {
			tokentype = "Application"
		}
		// is there an associated refres htoken ?
		var refreshtokeninfo string
		refreshtokeninfo = "No"
		for _, refreshtoken := range tokens {
			if (accessToken.AccountID == refreshtoken.Email) {
				refreshtokeninfo = "Yes"
			} else if (accessToken.AccountID == applicationdefaulthash) {
				refreshtokeninfo = "Yes"
				emailinfo = color.RedString(emailinfo)
			}
		}

		// is there a proof of reauthentication to use the refresh
		// token ? If not, password is needed
		var raptinfo string
		if refreshtokeninfo == "No" {
			raptinfo = "-"
		} else if accessToken.RaptToken.Valid {
			raptinfo = "No"
		} else {
			if (tokentype == "Application") {
				raptinfo = "TBD"
			} else {
				raptinfo = "Yes"
			}
		}
		tableBody = append(
			tableBody,
			[]string{
				emailinfo,
				tokentype,
				timeinfo,
				refreshtokeninfo,
				raptinfo,
			})
	}
	internal.PrintTableToScreen(tableHead, tableBody, false)
	return nil
}
