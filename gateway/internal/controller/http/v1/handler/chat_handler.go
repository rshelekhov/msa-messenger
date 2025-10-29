package handler

import (
	"errors"
	"log/slog"
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
			log.Warn("failed to decode request", slog.String("error", err.Error()))
			handleBadRequestError(w, r, ErrInvalidRequest)
			return
		}

		chats, nextPageToken, err := h.chatUsecase.GetChats(ctx, params.PageToken)
		if err != nil {
			log.Error("failed to get chats", slog.String("error", err.Error()))
			handleInternalError(w, r)
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
			log.Warn("failed to decode request", slog.String("error", err.Error()))
			handleBadRequestError(w, r, ErrInvalidRequest)
			return
		}

		chatID := chi.URLParam(r, "id")
		if chatID == "" {
			log.Warn("chat ID is required")
			handleBadRequestError(w, r, ErrChatIDRequired)
			return
		}

		chat, nextPageToken, err := h.chatUsecase.GetChat(ctx, chatID, params.PageToken)
		if err != nil {
			if errors.Is(err, domain.ErrChatNotFound) {
				log.Warn("failed to get chat", slog.String("error", err.Error()))
				handleError(w, r, domain.ErrChatNotFound, http.StatusNotFound)
				return
			}

			log.Error("failed to get chat", slog.String("error", err.Error()))
			handleInternalError(w, r)
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
			log.Warn("chat ID is required")
			handleBadRequestError(w, r, ErrChatIDRequired)
			return
		}

		req := &MessageRequest{}
		if err := h.decodeAndValidateRequest(w, r, log, req); err != nil {
			return
		}

		message, err := h.chatUsecase.SendMessage(ctx, chatID, req.Content)
		if err != nil {
			if errors.Is(err, domain.ErrChatNotFound) {
				log.Warn("failed to send message", slog.String("error", err.Error()))
				handleError(w, r, domain.ErrChatNotFound, http.StatusNotFound)
				return
			}

			log.Error("failed to send message", slog.String("error", err.Error()))
			handleInternalError(w, r)
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
			log.Warn("chat ID is required")
			handleBadRequestError(w, r, ErrChatIDRequired)
			return
		}

		if err := h.chatUsecase.DeleteChat(ctx, chatID); err != nil {
			if errors.Is(err, domain.ErrChatNotFound) {
				log.Warn("failed to delete chat", slog.String("error", err.Error()))
				handleError(w, r, domain.ErrChatNotFound, http.StatusNotFound)
				return
			}

			log.Error("failed to delete chat", slog.String("error", err.Error()))
			handleInternalError(w, r)
			return
		}

		render.Status(r, http.StatusOK)
	}
}
