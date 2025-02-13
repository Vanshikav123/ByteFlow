package main

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/go-playground/form/v4"
)

/*
Logging the Error: It logs the error details, including a stack trace, to help developers debug the issue.
debug.Stack(): This retrieves the current call stack as a byte slice.

	The stack trace helps developers identify where the error occurred in the code.
	Writes the error trace (error message + stack trace) to the log. This helps developers diagnose the issue later.
*/
func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errorLog.Print(trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}

func (app *application) render(w http.ResponseWriter, status int, page string, data *templateData) {
	ts, ok := app.templateCache[page]

	if !ok {
		err := fmt.Errorf("the template %s does not exist", page)
		app.serverError(w, err)
		return
	}
	/*A new buffer (bytes.Buffer) is created to temporarily store the rendered HTML.
	This is done to avoid writing incomplete or erroneous HTML directly to the response.*/
	buf := new(bytes.Buffer)

	err := ts.ExecuteTemplate(buf, "base", data)
	if err != nil {
		app.serverError(w, err)
		return
	}

	w.WriteHeader(status)
	/*buf.WriteTo(w): The contents of the buffer (the rendered HTML) are written to the http.ResponseWriter,
	which sends the response to the client.*/
	buf.WriteTo(w)
}

func (app *application) decodePostForm(r *http.Request, dst any) error {
	err := r.ParseForm()
	if err != nil {
		return err
	}
	/*Parsing the Form Data:
	     Before you can access the form data using r.PostForm, you need to call r.ParseForm().
		 This method parses the raw query from the URL and the body of the request, and populates r.Form and r.PostForm.
	     r.ParseForm() must be called before accessing r.PostForm; otherwise, r.PostForm will be empty.

		Accessing Form Data:

		r.PostForm is a url.Values type, which is a map of string slices (map[string][]string).
	    It contains the form data from the POST request body.
		You can access individual form fields using the Get method, which returns the first value associated with the given key.
	    If the key does not exist, it returns an empty string.*/
	err = app.formDecoder.Decode(dst, r.PostForm)
	if err != nil {
		var invalidDecoderError *form.InvalidDecoderError
		if errors.As(err, &invalidDecoderError) {
			panic(err)
		}

		return err
	}
	return nil
}
func (app *application) isAuthenticated(r *http.Request) bool {
	/*The function retrieves a value stored in the request context using r.Context().Value(isAuthenticatedContextKey).
	isAuthenticatedContextKey is a constant key of type contextKey (type contextKey string), ensuring unique identification.
	The retrieved value is then asserted as a boolean (bool).*/
	isAuthenticated, ok := r.Context().Value(isAuthenticatedContextKey).(bool)
	if !ok {
		return false
	}
	return isAuthenticated
}
