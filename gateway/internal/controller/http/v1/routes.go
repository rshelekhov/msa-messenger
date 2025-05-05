package v1

import (
	"github.com/go-chi/chi/v5"
)

func (ar *appRouter) initRoutes() (*chi.Mux, error) {
	ar.mux.Route("/v1", func(r chi.Router) {
		// Public routes
		r.Group(func(r chi.Router) {
			r.Head("/health", HealthCheck())
			r.Head("/ready", ReadyCheck())

			r.Get("/.well-known/jwks.json", ar.h.GetJWKS())

			r.Route("/auth", func(r chi.Router) {
				r.Post("/login", ar.h.Login())
				r.Post("/register", ar.h.Register())
				r.Post("/refresh", ar.h.RefreshToken())
			})
		})

		// Protected routes (access token required)
		r.Group(func(r chi.Router) {
			r.Use(ar.jwtMgr.HTTPMiddleware)

			r.Post("/auth/logout", ar.h.Logout())

			r.Route("/users", func(r chi.Router) {
				r.Get("/me", ar.h.GetOwnProfile())
				r.Patch("/me", ar.h.UpdateOwnProfile())
				r.Delete("/me", ar.h.DeleteOwnProfile())
				r.Get("/search", ar.h.SearchUsers())
			})

			r.Route("/friends", func(r chi.Router) {
				r.Get("/", ar.h.GetFriends())
				r.Get("/{id}", ar.h.GetFriend())
				r.Delete("/{id}", ar.h.RemoveFriend())
				r.Post("/{id}/invite", ar.h.InviteFriend())
			})

			r.Route("/friend-invites", func(r chi.Router) {
				r.Get("/", ar.h.GetFriendInvites())
				r.Patch("/{id}/accept", ar.h.AcceptFriendInvite())
				r.Patch("/{id}/decline", ar.h.DeclineFriendInvite())
			})

			r.Route("/chats", func(r chi.Router) {
				r.Get("/", ar.h.GetChats())
				r.Get("/{id}", ar.h.GetChat())
				r.Post("/{id}", ar.h.SendMessage())
				r.Delete("/{id}", ar.h.DeleteChat())
			})
		})
	})

	return ar.mux, nil
}
