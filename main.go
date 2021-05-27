package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/fasthttp/router"
	"github.com/sethvargo/go-password/password"
	"github.com/valyala/fasthttp"
)

type PasswordInput struct {
	LowerLetters string `json:"lowerLetters"`
	UpperLetters string `json:"upperLetters"`
	Symbols      string `json:"symbols"`
	Digits       string `json:"digits"`
	Length       int    `json:"length"`
	NumDigits    int    `json:"numDigits"`
	NumSymbols   int    `json:"numSymbols"`
	NoUpper      bool   `json:"noUpper"`
	AllowRepeat  bool   `json:"allowRepeat"`
}

func PasswordGenerator(ctx *fasthttp.RequestCtx) {
	passwordInput := &PasswordInput{}
	data := ctx.Request.Body()
	err := json.Unmarshal(data, passwordInput)
	if err != nil {
		fmt.Println(err)
		ctx.Response.SetStatusCode(500)
		return
	}

	g, err := password.NewGenerator(&password.GeneratorInput{
		LowerLetters: passwordInput.LowerLetters,
		UpperLetters: passwordInput.UpperLetters,
		Symbols:      passwordInput.Symbols,
		Digits:       passwordInput.Digits,
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
	res, err := g.Generate(passwordInput.Length, passwordInput.NumDigits, passwordInput.NumSymbols, passwordInput.NoUpper, passwordInput.AllowRepeat)
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
	r.POST("/passwords", PasswordGenerator)
	log.Fatal(fasthttp.ListenAndServe("localhost:8080", r.Handler))
}
