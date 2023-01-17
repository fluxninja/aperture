package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/mux"

	"github.com/fluxninja/aperture/playground/graphql_demo_app/graph"
	"github.com/fluxninja/aperture/playground/graphql_demo_app/graph/generated"
)

const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	mux := mux.NewRouter()
	mux.Use(authMiddleware)

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{}}))

	mux.Handle("/", playground.Handler("GraphQL playground", "/query"))
	mux.Handle("/query", srv)

	server := &http.Server{
		Addr:              ":" + port,
		Handler:           mux,
		ReadHeaderTimeout: 3 * time.Second,
	}

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		// if auth is not available then proceed to resolver
		if authHeader == "" {
			next.ServeHTTP(w, r)
		} else {
			// if auth is available then verify token
			tokenStr := strings.Replace(authHeader, "Bearer ", "", 1)
			token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
				fmt.Printf("token: %+v", token)
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}
				return []byte("secret"), nil
			})
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}

			// extract claims from jwt token
			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok || !token.Valid {
				http.Error(w, "invalid token", http.StatusUnauthorized)
				return
			}

			userID := claims["userID"].(string)

			// merge userID, userName into request context
			ctx := context.WithValue(r.Context(), graph.UserID{}, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
	})
}
