package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/fasthttp/router"
	"github.com/sethvargo/go-password/password"
	"github.com/valyala/fasthttp"
)

func Index(ctx *fasthttp.RequestCtx) {
	ctx.WriteString("Welcome!")
}

func PasswordGenerator(ctx *fasthttp.RequestCtx) {

	g, err := password.NewGenerator(&password.GeneratorInput{
		LowerLetters: "abcde",
		UpperLetters: "abcde",
		Symbols:      "!@#$%",
		Digits:       "01234",
	})
	if err != nil {
		fmt.Println(err)
		ctx.Response.SetStatusCode(500)
		return
	}

	// Generate generates a password with the given requirements. length is the
	// total number of characters in the password. numDigits is the number of digits
	// to include in the result. numSymbols is the number of symbols to include in
	// the result. noUpper excludes uppercase letters from the results. allowRepeat
	// allows characters to repeat.
	//
	// The algorithm is fast, but it's not designed to be performant; it favors
	// entropy over speed. This function is safe for concurrent use.
	res, err := g.Generate(64, 10, 10, false, true)
	if err != nil {
		fmt.Println(err)
		ctx.Response.SetStatusCode(500)
		return
	}

	body := map[string]interface{}{
		"password": res,
	}

	if err := json.NewEncoder(ctx).Encode(body); err != nil {
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
	} else {

		ctx.Response.Header.SetCanonical([]byte("Content-Type"), []byte("application/json"))
		ctx.Response.SetStatusCode(200)
	}

}

func main() {
	r := router.New()
	r.GET("/passwords", PasswordGenerator)
	log.Fatal(fasthttp.ListenAndServe("localhost:8080", r.Handler))
}
