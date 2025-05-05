package handler

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/rshelekhov/msa-messenger/gateway/internal/domain"
)

func (h *Handler) GetChats() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handler.chat.GetChats"

		ctx := r.Context()
		log := h.logWithReqID(ctx, op)

		params := &GetChatsParams{}
		if err := render.Decode(r, params); err != nil {
			err = fmt.Errorf("failed to decode request: %w", err)
			handleBadRequestError(w, r, log, err)
			return
		}

		chats, nextPageToken, err := h.chatUsecase.GetChats(ctx, params.PageToken)
		if err != nil {
			handleInternalError(w, r, log, err)
			return
		}

		resp := toGetChatsResponse(chats, nextPageToken)

		render.Status(r, http.StatusOK)
		render.JSON(w, r, resp)
	}
}

func (h *Handler) GetChat() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handler.chat.GetChat"

		ctx := r.Context()
		log := h.logWithReqID(ctx, op)

		params := &GetChatsIdParams{}
		if err := render.Decode(r, params); err != nil {
			err = fmt.Errorf("failed to decode request: %w", err)
			handleBadRequestError(w, r, log, err)
			return
		}

		chatID := chi.URLParam(r, "id")
		if chatID == "" {
			handleBadRequestError(w, r, log, errors.New("chat ID is required"))
			return
		}

		chat, nextPageToken, err := h.chatUsecase.GetChat(ctx, chatID, params.PageToken)
		if err != nil {
			if errors.Is(err, domain.ErrChatNotFound) {
				handleError(w, r, log, err, http.StatusNotFound)
				return
			}
			handleInternalError(w, r, log, err)
			return
		}

		resp := toGetChatResponse(chat, nextPageToken)

		render.Status(r, http.StatusOK)
		render.JSON(w, r, resp)
	}
}

func (h *Handler) SendMessage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handler.chat.SendMessage"

		ctx := r.Context()
		log := h.logWithReqID(ctx, op)

		chatID := chi.URLParam(r, "id")
		if chatID == "" {
			handleBadRequestError(w, r, log, errors.New("chat ID is required"))
			return
		}

		req := &MessageRequest{}
		if err := h.decodeAndValidateRequest(w, r, log, req); err != nil {
			return
		}

		message, err := h.chatUsecase.SendMessage(ctx, chatID, req.Content)
		if err != nil {
			if errors.Is(err, domain.ErrChatNotFound) {
				handleError(w, r, log, err, http.StatusNotFound)
				return
			}
			handleInternalError(w, r, log, err)
			return
		}

		resp := toSendMessageResponse(message)

		render.Status(r, http.StatusOK)
		render.JSON(w, r, resp)
	}
}

func (h *Handler) DeleteChat() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handler.chat.DeleteChat"

		ctx := r.Context()
		log := h.logWithReqID(ctx, op)

		chatID := chi.URLParam(r, "id")
		if chatID == "" {
			handleBadRequestError(w, r, log, errors.New("chat ID is required"))
			return
		}

		if err := h.chatUsecase.DeleteChat(ctx, chatID); err != nil {
			if errors.Is(err, domain.ErrChatNotFound) {
				handleError(w, r, log, err, http.StatusNotFound)
				return
			}
			handleInternalError(w, r, log, err)
			return
		}

		render.Status(r, http.StatusOK)
	}
}
