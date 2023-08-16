package rest

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/xopxe23/news-server/internal/domain"
)

type ArticlesService interface {
	CreateAuthor(ctx context.Context, author domain.Author) (int, error)
	GetAllAuthors(ctx context.Context) ([]domain.Author, error)
	GetAuthorById(ctx context.Context, authorId int) (domain.Author, error)
	UpdateAuthor(ctx context.Context, authorId int, input domain.UpdateAuthorInput) error
	DeleteAuthor(ctx context.Context, authorId int) error
}

// @Summary Get All Authors
// @Tags Articles
// @ID get-all-articles
// @Accept json
// @Produce json
// @Success 200
// @Failure 400
// @Failure 500
// @Router /authors [get]
func (h *Handler) getAllAuthors(w http.ResponseWriter, r *http.Request) {
	authors, err := h.articlesService.GetAllAuthors(r.Context())
	if err != nil {
		logError("getAllAuthors", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(map[string]interface{}{
		"data": authors,
	})
	if err != nil {
		logError("getAllAuthors", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(response)
}

// @Summary Create Author
// @Tags Articles
// @ID create-author
// @Accept json
// @Produce json
// @Param input body domain.Author true "Author input"
// @Success 200
// @Failure 400
// @Failure 500
// @Router /authors [post]
func (h *Handler) createAuthor(w http.ResponseWriter, r *http.Request) {
	reqBytes, err := io.ReadAll(r.Body)
	if err != nil {
		logError("createAuthor", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var author domain.Author
	if err := json.Unmarshal(reqBytes, &author); err != nil {
		logError("createAuthor", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	auhorId, err := h.articlesService.CreateAuthor(r.Context(), author)
	if err != nil {
		logError("createAuthor", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	response, err := json.Marshal(map[string]int{
		"author id": auhorId,
	})
	if err != nil {
		logError("createAuthor", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(response)
}

// @Summary Get Author By Id
// @Tags Articles
// @ID get-author-by-id
// @Accept json
// @Produce json
// @Param id path int true "Author ID"
// @Success 200
// @Failure 400
// @Failure 500
// @Router /authors/{id} [get]
func (h *Handler) getAuthorById(w http.ResponseWriter, r *http.Request) {
	authorId, err := getIdFromRequest(r)
	if err != nil {
		logError("getAuthorById", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	author, err := h.articlesService.GetAuthorById(r.Context(), authorId)
	if err != nil {
		logError("getAuthorById", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(author)
	if err != nil {
		logError("getAuthorById", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(response)
}

// @Summary Update Author
// @Tags Articles
// @ID update-author
// @Accept json
// @Produce json
// @Param id path int true "Author ID"
// @Param input body domain.UpdateAuthorInput true "Update Author input"
// @Success 200
// @Failure 400
// @Failure 500
// @Router /authors/{id} [put]
func (h *Handler) updateAuthor(w http.ResponseWriter, r *http.Request) {
	authorId, err := getIdFromRequest(r)
	if err != nil {
		logError("updateAuthor", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	reqBytes, err := io.ReadAll(r.Body)
	if err != nil {
		logError("updateAuthor", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var input domain.UpdateAuthorInput
	if err := json.Unmarshal(reqBytes, &input); err != nil {
		logError("updateAuthor", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err = h.articlesService.UpdateAuthor(r.Context(), authorId, input); err != nil {
		logError("updateAuthor", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	response, err := json.Marshal(map[string]string{
		"status": "author updated",
	})
	if err != nil {
		logError("updateAuthor", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(response)
}


// @Summary Delete Author
// @Tags Articles
// @ID delete-author
// @Accept json
// @Produce json
// @Param id path int true "Author ID"
// @Success 200
// @Failure 400
// @Failure 500
// @Router /authors/{id} [delete]
func (h *Handler) deleteAuthor(w http.ResponseWriter, r *http.Request) {
	authorId, err := getIdFromRequest(r)
	if err != nil {
		logError("deleteAuthor", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = h.articlesService.DeleteAuthor(r.Context(), authorId)
	if err != nil {
		logError("deleteAuthor", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	response, err := json.Marshal(map[string]string{
		"status": "author deleted",
	})
	if err != nil {
		logError("deleteAuthor", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(response)
}

func (h *Handler) getAllArticles(w http.ResponseWriter, r *http.Request) {}
func (h *Handler) createArticle(w http.ResponseWriter, r *http.Request)  {}
func (h *Handler) getArticleById(w http.ResponseWriter, r *http.Request) {}
func (h *Handler) updateArticle(w http.ResponseWriter, r *http.Request)  {}
func (h *Handler) deleteArticle(w http.ResponseWriter, r *http.Request)  {}

func getIdFromRequest(r *http.Request) (int, error) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		return 0, err
	}
	if id == 0 {
		return 0, errors.New("id can't be zero")
	}
	return id, nil
}
