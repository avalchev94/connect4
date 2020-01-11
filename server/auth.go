package server

import (
	"errors"
	"net/http"

	"github.com/gorilla/securecookie"

	"github.com/avalchev94/tarantula"
	"github.com/avalchev94/tarantula/games"
)

var (
	hashKey      = []byte("fd2a388a983c529183e323178364474f13ac273619337fa09b9291f7b59dbba8")
	secureCookie = securecookie.New(hashKey, nil)
	allowedHosts = []string{"http://192.168.1.100:8081", "http://localhost:8081"}
)

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		originHost := r.Header["Origin"][0]
		for i := range allowedHosts {
			if originHost == allowedHosts[i] {
				w.Header().Set("Access-Control-Allow-Origin", originHost)
			}
		}
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		next.ServeHTTP(w, r)
	})
}

type authCookie struct {
	UUID   tarantula.UUID `json:"uuid"`
	Player games.PlayerID `json:"player_id"`
}

func encodeCookie(w http.ResponseWriter, name string, data authCookie) error {
	value, err := secureCookie.Encode(name, data)
	if err != nil {
		return err
	}

	http.SetCookie(w, &http.Cookie{
		Name:  name,
		Value: value,
	})
	return nil
}

func decodeCookie(r *http.Request, name string) (authCookie, error) {
	httpCookie, err := r.Cookie(name)
	if err != nil {
		return authCookie{}, errors.New("Auth cookie not found")
	}

	cookie := authCookie{}
	return cookie, secureCookie.Decode(name, httpCookie.Value, &cookie)
}
