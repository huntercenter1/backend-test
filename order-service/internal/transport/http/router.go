package http

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/huntercenter1/backend-test/order-service/internal/service"
)

type Router struct {
	svc service.Service
}

func New(svc service.Service) *Router { return &Router{svc: svc} }

func (rt *Router) Register(r *gin.Engine) {
	r.GET("/health", func(c *gin.Context){ c.JSON(http.StatusOK, gin.H{"status":"ok"}) })
	r.POST("/orders", rt.create)
	r.GET("/orders/:id", rt.get)
	r.GET("/orders/:id/items", rt.items)
	r.GET("/orders/user/:user_id", rt.byUser)
	r.PUT("/orders/:id/status", rt.updateStatus)
}

type createReq struct {
	UserID string               `json:"user_id"`
	Items  []service.CreateItem `json:"items"`
}

func (rt *Router) create(c *gin.Context) {
	var req createReq
	if err := c.ShouldBindJSON(&req); err != nil || req.UserID == "" || len(req.Items) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error":"invalid payload"}); return
	}
	o, items, err := rt.svc.Create(c.Request.Context(), req.UserID, req.Items)
	if err != nil { c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}); return }
	c.JSON(http.StatusCreated, gin.H{"order": o, "items": items})
}

func (rt *Router) get(c *gin.Context) {
	o, err := rt.svc.Get(c.Request.Context(), c.Param("id"))
	if err != nil { c.JSON(http.StatusNotFound, gin.H{"error":"not found"}); return }
	c.JSON(http.StatusOK, o)
}

func (rt *Router) items(c *gin.Context) {
	items, err := rt.svc.Items(c.Request.Context(), c.Param("id"))
	if err != nil { c.JSON(http.StatusNotFound, gin.H{"error":"not found"}); return }
	c.JSON(http.StatusOK, gin.H{"items": items})
}

func (rt *Router) byUser(c *gin.Context) {
	list, err := rt.svc.ByUser(c.Request.Context(), c.Param("user_id"))
	if err != nil { c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}); return }
	c.JSON(http.StatusOK, gin.H{"orders": list})
}

type statusReq struct{ Status string `json:"status"` }

func (rt *Router) updateStatus(c *gin.Context) {
	var body statusReq
	if err := c.ShouldBindJSON(&body); err != nil || body.Status == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error":"invalid body"}); return
	}
	o, err := rt.svc.UpdateStatus(c.Request.Context(), c.Param("id"), body.Status)
	if err != nil { c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}); return }
	c.JSON(http.StatusOK, o)
}
