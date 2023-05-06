package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go/rest-ws/models"
	"github.com/go/rest-ws/repository"
	"github.com/go/rest-ws/server"
	"github.com/gorilla/mux"
	"github.com/segmentio/ksuid"
)

type UpsertPostRequest struct {
	PostContent string `json:"post_content"`
}

type PostResponse struct {
	Id          string `json:"id"`
	PostContent string `json:"post_content"`
}

type PostDeleteUpdateResponse struct {
	Message string `json:"message"`
}

func InsertPostHandler(s server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userId := r.Header.Get("userId")
		user, err := repository.GetUserById(r.Context(), userId)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var postRequest = UpsertPostRequest{}
		if err := json.NewDecoder(r.Body).Decode(&postRequest); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		id, err := ksuid.NewRandom()

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		post := models.Post{
			Id:          id.String(),
			PostContent: postRequest.PostContent,
			UserId:      user.Id,
		}

		err = repository.InsertPost(r.Context(), &post)

		var postMessage = models.WebsocketMessage{
			Type:    "post_create",
			Payload: post,
		}

		s.Hub().Broadcast(postMessage, nil)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(PostResponse{
			Id:          post.Id,
			PostContent: post.PostContent,
		})

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func GetPostByIdHandler(s server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		post, err := repository.GetPostById(r.Context(), params["id"])

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(post)
	}
}

func UpdatePostHandler(s server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userId := r.Header.Get("userId")
		user, err := repository.GetUserById(r.Context(), userId)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		params := mux.Vars(r)

		var postRequest = UpsertPostRequest{}
		if err := json.NewDecoder(r.Body).Decode(&postRequest); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		post := models.Post{
			Id:          params["id"],
			PostContent: postRequest.PostContent,
			UserId:      user.Id,
		}

		err = repository.UpdatePost(r.Context(), &post)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(PostDeleteUpdateResponse{
			Message: "Update post successfully",
		})
	}
}

func DeletePostDHandler(s server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userId := r.Header.Get("userId")
		user, err := repository.GetUserById(r.Context(), userId)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		params := mux.Vars(r)

		err = repository.DeletePost(r.Context(), params["id"], user.Id)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(PostDeleteUpdateResponse{
			Message: "Post deleted successfully",
		})

	}
}

func ListPosts(s server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var err error
		pageStr := r.URL.Query().Get("page")
		var page = uint64(0)

		if pageStr != "" {
			page, err = strconv.ParseUint(pageStr, 10, 64)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		}

		posts, err := repository.PostList(r.Context(), page)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(posts)

	}
}
