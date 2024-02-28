package main

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"

	"github.com/akshat-kaushik/goPriceFetcher/types"
)

type APIfunc func(ctx context.Context, w http.ResponseWriter, r *http.Request) error

type JSONApiServer struct {
	svc        PriceFetcher
	listenAddr string
}

func NewJSONApiServer(svc PriceFetcher) *JSONApiServer {
	return &JSONApiServer{
		svc:        svc,
		listenAddr: ":8080",
	}
}

func (s *JSONApiServer) Run() {
	http.HandleFunc("/", makeHttpHandler(s.handleFetchPrice))
	http.ListenAndServe(s.listenAddr, nil)
	fmt.Println("Server is running on ", s.listenAddr)
}

type contextKey string

func makeHttpHandler(fn APIfunc) http.HandlerFunc {
	ctx := context.Background()
	ctx = context.WithValue(ctx, contextKey("reqID"), rand.Intn(100000000))
	return func(w http.ResponseWriter, r *http.Request) {
		if err := fn(ctx, w, r); err != nil {
			writeJson(w, http.StatusInternalServerError, err.Error())
		}
	}
}
func (s *JSONApiServer) handleFetchPrice(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	ticker := r.URL.Query().Get("ticker")

	price, err := s.svc.FetchPrice(ctx, ticker)
	if err != nil {
		return err
	}
	priceResp := types.PriceResponse{
		Ticker: ticker,
		Price:  price,
	}

	return writeJson(w, http.StatusOK, priceResp)
}

func writeJson(w http.ResponseWriter, s int, v any) error {
	w.WriteHeader(s)
	return json.NewEncoder(w).Encode(v)
}
