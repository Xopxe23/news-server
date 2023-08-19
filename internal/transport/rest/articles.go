package rest

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
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

	CreateArticle(ctx context.Context, input domain.Article) (int, error)
	GetAllArticles(ctx context.Context) ([]domain.ArticleOutput, error)
	GetArticleById(ctx context.Context, articleId int) (domain.ArticleOutput, error)
	UpdateArticle(ctx context.Context, articleId int, input domain.UpdateArticleInput) error
	DeleteArticle(ctx context.Context, articleId int) error
}

// @Summary Get All Authors
// @Security ApiKeyAuth
// @Tags Authors
// @ID get-all-authors
// @Accept json
// @Produce json
// @Success 200 {array} domain.Author
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
// @Security ApiKeyAuth
// @Tags Authors
// @ID create-author
// @Accept json
// @Produce json
// @Param input body domain.Author true "Author input"
// @Success 200 {integer} domain.Author.Id
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
// @Security ApiKeyAuth
// @Tags Authors
// @ID get-author-by-id
// @Accept json
// @Produce json
// @Param id path int true "Author ID"
// @Success 200 {object} domain.Author
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
// @Security ApiKeyAuth
// @Security ApiKeyAuth
// @Tags Authors
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
// @Security ApiKeyAuth
// @Tags Authors
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
		"status": fmt.Sprintf("author №%d deleted", authorId),
	})
	if err != nil {
		logError("deleteAuthor", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(response)
}

// @Summary Get All Articles
// @Security BearerAuth
// @Tags Articles
// @ID get-all-articles
// @Accept json
// @Produce json
// @Success 200 {array} domain.Article
// @Failure 400
// @Failure 500
// @Router /articles [get]
func (h *Handler) getAllArticles(w http.ResponseWriter, r *http.Request) {
	var articles []domain.ArticleOutput
	articles, err := h.articlesService.GetAllArticles(r.Context())
	if err != nil {
		logError("getAllArticles", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	response, err := json.Marshal(map[string][]domain.ArticleOutput{
		"data": articles,
	})
	if err != nil {
		logError("getAllArticles", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(response)
}

// @Summary Create Article
// @Security ApiKeyAuth
// @Tags Articles
// @ID create-articles
// @Accept json
// @Produce json
// @Param input body domain.Article true "Article input"
// @Success 200
// @Failure 400
// @Failure 500
// @Router /articles [post]
func (h *Handler) createArticle(w http.ResponseWriter, r *http.Request) {
	reqBytes, err := io.ReadAll(r.Body)
	if err != nil {
		logError("createArticle", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var input domain.Article
	if err := json.Unmarshal(reqBytes, &input); err != nil {
		logError("createArticle", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	articleId, err := h.articlesService.CreateArticle(r.Context(), input)
	if err != nil {
		logError("createArticle", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	response, err := json.Marshal(map[string]int{
		"article id": articleId,
	})
	if err != nil {
		logError("createArticle", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(response)
}

// @Summary Get Article By Id
// @Security ApiKeyAuth
// @Tags Articles
// @ID get-article-by-id
// @Accept json
// @Produce json
// @Param id path int true "Article ID"
// @Success 200 {object} domain.Article
// @Failure 400
// @Failure 500
// @Router /articles/{id} [get]
func (h *Handler) getArticleById(w http.ResponseWriter, r *http.Request) {
	articleId, err := getIdFromRequest(r)
	if err != nil {
		logError("getArticleById", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	article, err := h.articlesService.GetArticleById(r.Context(), articleId)
	if err != nil {
		logError("getArticleById", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(map[string]domain.ArticleOutput{
		"article": article,
	})
	if err != nil {
		logError("getArticleById", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(response)
}

// @Summary Update Article
// @Security ApiKeyAuth
// @Tags Articles
// @ID update-article
// @Accept json
// @Produce json
// @Param id path int true "Article ID"
// @Param input body domain.UpdateArticleInput true "Update Article input"
// @Success 200
// @Failure 400
// @Failure 500
// @Router /articles/{id} [put]
func (h *Handler) updateArticle(w http.ResponseWriter, r *http.Request) {
	articleId, err := getIdFromRequest(r)
	if err != nil {
		logError("updateArticle", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var input domain.UpdateArticleInput
	reqBytes, err := io.ReadAll(r.Body)
	if err != nil {
		logError("updateArticle", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err := json.Unmarshal(reqBytes, &input); err != nil {
		logError("updateArticle", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = h.articlesService.UpdateArticle(r.Context(), articleId, input)
	if err != nil {
		logError("updateArticle", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(map[string]string{
		"status": fmt.Sprintf("article №%d updated", articleId),
	})
	if err != nil {
		logError("updateArticle", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(response)
}

// @Summary Delete Article
// @Security ApiKeyAuth
// @Tags Articles
// @ID delete-article
// @Accept json
// @Produce json
// @Param id path int true "Article ID"
// @Success 200
// @Failure 400
// @Failure 500
// @Router /articles/{id} [delete]
func (h *Handler) deleteArticle(w http.ResponseWriter, r *http.Request) {
	articleId, err := getIdFromRequest(r)
	if err != nil {
		logError("deleteArticle", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := h.articlesService.DeleteArticle(r.Context(), articleId); err != nil {
		logError("deleteArticle", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(map[string]string{
		"status": fmt.Sprintf("article №%d deleted", articleId),
	})
	if err != nil {
		logError("deleteArticle", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(response)
}

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
