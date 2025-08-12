package http

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/uptrace/bun"

	"github.com/huntercenter1/backend-test/product-service/internal/middleware"
	"github.com/huntercenter1/backend-test/product-service/internal/models"
	"github.com/huntercenter1/backend-test/product-service/internal/repo"
)

type Router struct {
	db   *bun.DB
	repo repo.ProductRepo
}

func New(db *bun.DB) *Router {
	return &Router{db: db, repo: repo.New(db)}
}

func (rt *Router) Register(r *gin.Engine) {
	r.Use(gin.Recovery(), middleware.RequestID(), middleware.Timeout(5*time.Second))

	r.GET("/health", func(c *gin.Context){ c.JSON(http.StatusOK, gin.H{"status":"ok"}) })

	r.GET("/products", rt.list)
	r.POST("/products", rt.create)
	r.GET("/products/:id", rt.get)
	r.PUT("/products/:id", rt.update)
	r.DELETE("/products/:id", rt.delete)
	r.GET("/products/search", rt.search)
	r.PUT("/products/:id/stock", rt.updateStock)
}

func (rt *Router) list(c *gin.Context) {
	limit, offset := parsePag(c)
	items, total, err := rt.repo.List(c.Request.Context(), limit, offset)
	if err != nil { c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}); return }
	c.JSON(http.StatusOK, gin.H{"items": items, "total": total, "limit": limit, "offset": offset})
}

func (rt *Router) create(c *gin.Context) {
	var p models.Product
	if err := c.ShouldBindJSON(&p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"}); return
	}
	if p.Name == "" || p.Price <= 0 { c.JSON(http.StatusBadRequest, gin.H{"error":"name/price required"}); return }
	res, err := rt.repo.Create(c.Request.Context(), &p)
	if err != nil { c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}); return }
	c.JSON(http.StatusCreated, res)
}

func (rt *Router) get(c *gin.Context) {
	id := c.Param("id")
	p, err := rt.repo.GetByID(c.Request.Context(), id)
	if err != nil { c.JSON(http.StatusNotFound, gin.H{"error":"not found"}); return }
	c.JSON(http.StatusOK, p)
}

func (rt *Router) update(c *gin.Context) {
	id := c.Param("id")
	var body models.Product
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"}); return
	}
	p, err := rt.repo.GetByID(c.Request.Context(), id)
	if err != nil { c.JSON(http.StatusNotFound, gin.H{"error":"not found"}); return }
	if body.Name != "" { p.Name = body.Name }
	if body.Description != "" { p.Description = body.Description }
	if body.Price > 0 { p.Price = body.Price }
	if body.Stock >= 0 { p.Stock = body.Stock }
	res, err := rt.repo.Update(c.Request.Context(), p)
	if err != nil { c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}); return }
	c.JSON(http.StatusOK, res)
}

func (rt *Router) delete(c *gin.Context) {
	id := c.Param("id")
	if err := rt.repo.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error":"not found"}); return
	}
	c.Status(http.StatusNoContent)
}

func (rt *Router) search(c *gin.Context) {
	q := c.Query("q")
	limit, offset := parsePag(c)
	items, total, err := rt.repo.Search(c.Request.Context(), q, limit, offset)
	if err != nil { c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}); return }
	c.JSON(http.StatusOK, gin.H{"items": items, "total": total, "limit": limit, "offset": offset})
}

type stockBody struct { Delta int `json:"delta"` }

func (rt *Router) updateStock(c *gin.Context) {
	id := c.Param("id")
	var body stockBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error":"invalid body"}); return
	}
	res, err := rt.repo.UpdateStock(c.Request.Context(), id, body.Delta)
	if err != nil { c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}); return }
	c.JSON(http.StatusOK, res)
}

func parsePag(c *gin.Context) (int, int) {
	limit := 20; offset := 0
	if v := c.Query("limit"); v != "" { if n, err := strconv.Atoi(v); err==nil && n>0 && n<=100 { limit = n } }
	if v := c.Query("offset"); v != "" { if n, err := strconv.Atoi(v); err==nil && n>=0 { offset = n } }
	return limit, offset
}
