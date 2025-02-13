package main

/*Defines a custom type contextKey as a string alias.
Declares a constant isAuthenticatedContextKey of type contextKey with the value "isAuthenticated".
Used to store/retrieve authentication status in Go’s context.Context.*/

type contextKey string

const isAuthenticatedContextKey = contextKey("isAuthenticated")
