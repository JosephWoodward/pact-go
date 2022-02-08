package graphql

import (
	"fmt"
	"regexp"

	"github.com/pact-foundation/pact-go/dsl"
)

// Variables represents values to be substituted at runtime
type Variables map[string]interface{}

// Request is the main implementation of the Pact interface.
type Request struct {
	// HTTP Headers
	Headers dsl.MapMatcher

	// Path to GraphQL endpoint
	Path dsl.Matcher

	// HTTP Query String
	QueryString dsl.MapMatcher

	// GraphQL Query
	Query string

	// GraphQL Variables
	Variables Variables

	// GraphQL Operation
	Operation string

	// GraphQL method (usually POST, but can be get with a query string)
	// NOTE: for query string users, the standard HTTP interaction should suffice
	Method string
}

// Specify the operation (if any)
func (r *Request) WithOperation(operation string) *Request {
	r.Operation = operation

	return r
}

// Specify the method (defaults to POST)
func (r *Request) WithMethod(method string) *Request {
	r.Method = method

	return r
}

// Given specifies a provider state. Optional.
func (r *Request) WithQuery(query string) *Request {
	r.Query = query

	// TODO: should/can we validate this?

	return r
}

// Given specifies a provider state. Optional.
func (r *Request) WithVariables(variables Variables) *Request {
	r.Variables = variables

	return r
}

func GraphQL(request Request) *dsl.Request {

	// TODO: set defaults and validate aspects of the setup

	return &dsl.Request{
		Method: request.Method,
		Path:   request.Path,
		Query:  request.QueryString,
		Body: graphQLQueryBody{
			Operation: request.Operation,
			Query:     dsl.Regex(request.Query, escapeGraphQlQuery(request.Query)),
			Variables: request.Variables,
		},
		Headers: request.Headers,
	}

}

type graphQLQueryBody struct {
	Operation string      `json:"operationName,omitempty"`
	Query     dsl.Matcher `json:"query"`
	Variables Variables   `json:"variables,omitempty"`
}

func escapeSpace(s string) string {
	r := regexp.MustCompile(`\s+`)
	return r.ReplaceAllString(s, `\s*`)
}

func escapeRegexChars(s string) string {
	r := regexp.MustCompile(`(?m)[\-\[\]\/\{\}\(\)\*\+\?\.\\\^\$\|]`)

	f := func(s string) string {
		return fmt.Sprintf(`\%s`, s)
	}
	return r.ReplaceAllStringFunc(s, f)
}

func escapeGraphQlQuery(s string) string {
	return escapeSpace(escapeRegexChars(s))
}
