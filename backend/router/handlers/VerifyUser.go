package handlers

import (
	"encoding/json"
	"fmt"
	"meetover/backend/firebase"
	"meetover/backend/user"
	"net/http"

	"github.com/gorilla/mux"
)

// VerifyUser - token exchange and authentication at user login
func VerifyUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	tempUserCode := params["code"]
	redirectURI := r.URL.Query().Get("redirect_uri")
	fmt.Println("[+] Recieved code: " + tempUserCode)
	fmt.Println("[+] Recieved redirect_uri: " + redirectURI)

	aTokenResp, err := user.ExchangeToken(tempUserCode, redirectURI)
	if err != nil {
		respondWithError(w, FailedTokenExchange, err.Error())
		fmt.Println("Sending failed token exchange error")
		return
	}
	fmt.Println("[+] After ExchangeToken: " + aTokenResp.AToken)

	p, err := user.GetProfile(aTokenResp.AToken)
	if err != nil {
		respondWithError(w, FailedTokenExchange, err.Error())
		fmt.Println("Sending failed token exchange error")
		return
	}
	// Updates access token if user exists or adds a new User
	userExists, err := firebase.InitUser(p, aTokenResp)
	if err != nil {
		respondWithError(w, FailedUserInit, err.Error())
		return
	}

	// gets firebase access token for user's IM chat
	customToken, err := firebase.CreateCustomToken(p.ID)
	if err != nil {
		respondWithError(w, FailedTokenExchange, err.Error())
		return
	}

	resp := AuthResponse{
		AccessToken:         aTokenResp,
		Profile:             p,
		FirebaseCustomToken: customToken,
		UserExists:          userExists,
	}

	json.NewEncoder(w).Encode(resp)
}
