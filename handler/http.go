package handler

import (
	"alibaba/shortener/domain"
	"alibaba/shortener/store"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"net/http"
)

type HttpHandler struct {
	Store  store.Store
	Logger *zap.Logger
}

func (h *HttpHandler) GetOriginal(w http.ResponseWriter, r *http.Request) {
	var req domain.GetOriginalReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		_ = json.NewEncoder(w).Encode(h.makeResp("get_original", nil, domain.ErrBadParamInput))
		return
	}
	resp, err := h.Store.GetOriginal(r.Context(), req.Short)
	_ = json.NewEncoder(w).Encode(h.makeResp("get_original", resp, err))
}

func (h *HttpHandler) GetInfo(w http.ResponseWriter, r *http.Request) {
	var req domain.GetInfoReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		_ = json.NewEncoder(w).Encode(h.makeResp("get_info", nil, domain.ErrBadParamInput))
		return
	}
	resp, err := h.Store.GetInfo(r.Context(), req.Short)
	_ = json.NewEncoder(w).Encode(h.makeResp("get_info", resp, err))
}

func (h *HttpHandler) Save(w http.ResponseWriter, r *http.Request) {
	var req domain.SaveReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		_ = json.NewEncoder(w).Encode(h.makeResp("save", nil, domain.ErrBadParamInput))
		return
	}
	resp, err := h.Store.Save(r.Context(), req.Original)
	_ = json.NewEncoder(w).Encode(h.makeResp("save", resp, err))
}

func (h *HttpHandler) Redirect(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	org, err := h.Store.GetOriginal(r.Context(), vars["short"])
	if err != nil {
		_ = json.NewEncoder(w).Encode(h.makeResp("Redirect", nil, err))
		return
	}
	http.Redirect(w, r, org, http.StatusPermanentRedirect)
}

func (h *HttpHandler) makeResp(endpoint string, resp interface{}, err error) *domain.HttpResponse {
	switch err {
	case nil:
		h.Logger.Info(fmt.Sprintf("%s succeed", endpoint))
		return &domain.HttpResponse{
			StatusCode: 200,
			StatusDesc: "succeed",
			Error:      "",
			Payload:    resp,
		}
	case domain.ErrInternalServerError:
		h.Logger.Error(fmt.Sprintf("%s failed", endpoint), zap.Error(err))
		return &domain.HttpResponse{
			StatusCode: 500,
			StatusDesc: "internal error",
			Error:      err.Error(),
			Payload:    resp,
		}
	case domain.ErrBadParamInput:
		h.Logger.Error(fmt.Sprintf("%s failed", endpoint), zap.Error(err))
		return &domain.HttpResponse{
			StatusCode: 400,
			StatusDesc: "bad request",
			Error:      err.Error(),
			Payload:    resp,
		}
	case domain.ErrNotFound:
		h.Logger.Error(fmt.Sprintf("%s failed", endpoint), zap.Error(err))
		return &domain.HttpResponse{
			StatusCode: 404,
			StatusDesc: "not found",
			Error:      err.Error(),
			Payload:    resp,
		}
	default:
		h.Logger.Error(fmt.Sprintf("%s failed", endpoint), zap.Error(err))
		return &domain.HttpResponse{
			StatusCode: 500,
			StatusDesc: "internal",
			Error:      err.Error(),
			Payload:    resp,
		}
	}
}
