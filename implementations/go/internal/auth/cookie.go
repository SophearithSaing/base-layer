package auth

import (
	"baselayer/internal/config"
	"net/http"
	"time"
)

func setAuthCookie(w http.ResponseWriter, accessToken, refreshToken string) {
	now := time.Now()
	isProd := config.GetAppEnv() == "prod"
	sameSite := http.SameSiteLaxMode
	if isProd {
		sameSite = http.SameSiteNoneMode
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    accessToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   isProd,
		SameSite: sameSite,
		Expires:  now.Add(15 * time.Minute),
	})
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Path:     "/auth/refresh",
		HttpOnly: true,
		Secure:   isProd,
		SameSite: sameSite,
		Expires:  now.AddDate(0, 0, 30),
	})
}
