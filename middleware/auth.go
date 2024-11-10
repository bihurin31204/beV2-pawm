package middleware

import (
    "context"
    "net/http"

    "github.com/auth0-community/go-auth0"
    "gopkg.in/square/go-jose.v2"
)

var auth0Domain string
var auth0Audience string
var jwtValidator *auth0.JWTValidator

func init() {
    auth0Domain = "dev-snipo7mrvybz2qyd.us.auth0.com"
    auth0Audience = "https://dev-snipo7mrvybz2qyd.us.auth0.com/api/v2/" // Use your custom API's Identifier

    if auth0Domain == "" || auth0Audience == "" {
        panic("AUTH0_DOMAIN or AUTH0_AUDIENCE not set in environment variables")
    }

    issuerURL := "https://" + auth0Domain + "/"

    jwksURI := issuerURL + ".well-known/jwks.json"
    provider := auth0.NewJWKClient(auth0.JWKClientOptions{URI: jwksURI}, nil)

    configuration := auth0.NewConfiguration(
        provider,
        []string{auth0Audience},
        issuerURL,
        jose.RS256,
    )

    jwtValidator = auth0.NewValidator(configuration, nil)
}

func AuthMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        token, err := jwtValidator.ValidateRequest(r)
        if err != nil {
            http.Error(w, "Unauthorized: "+err.Error(), http.StatusUnauthorized)
            return
        }

        // Store token claims in the request context
        ctx := context.WithValue(r.Context(), "user", token)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}
