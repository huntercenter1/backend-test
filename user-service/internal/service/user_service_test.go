package service

import (
	"context"
	"testing"

	"github.com/huntercenter1/backend-test/user-service/internal/auth"
	"github.com/huntercenter1/backend-test/user-service/internal/models"
	"github.com/huntercenter1/backend-test/user-service/internal/repo"
)

type fakeRepo struct{
	byID map[string]*models.User
	byU  map[string]*models.User
}
func (f *fakeRepo) Create(ctx context.Context, u *models.User)(*models.User,error){ f.byID[u.ID]=u; f.byU[u.Username]=u; return u,nil }
func (f *fakeRepo) GetByID(ctx context.Context, id string)(*models.User,error){ if u,ok:=f.byID[id]; ok { return u,nil }; return nil, repo.ErrNotFound }
func (f *fakeRepo) GetByUsername(ctx context.Context, un string)(*models.User,error){ if u,ok:=f.byU[un]; ok { return u,nil }; return nil, repo.ErrNotFound }
func (f *fakeRepo) Update(ctx context.Context, u *models.User)(*models.User,error){ f.byID[u.ID]=u; f.byU[u.Username]=u; return u,nil }
func (f *fakeRepo) Delete(ctx context.Context, id string) error { if _,ok:=f.byID[id]; !ok { return repo.ErrNotFound }; delete(f.byID,id); return nil }

func TestAuthenticateAndValidate(t *testing.T){
	f := &fakeRepo{byID:map[string]*models.User{}, byU:map[string]*models.User{}}
	s := New(f)
	hash,_ := auth.HashPassword("123456")
	u := &models.User{ID:"u1", Username:"demo", Email:"d@e.com", PasswordHash: hash}
	f.byID["u1"]=u; f.byU["demo"]=u

	if id, err := s.Authenticate(context.Background(),"demo","123456"); err!=nil || id!="u1" {
		t.Fatalf("auth failed %v %s", err, id)
	}
	if ok, _ := s.Validate(context.Background(),"u1"); !ok {
		t.Fatalf("validate should be true")
	}
	if ok, _ := s.Validate(context.Background(),"nope"); ok {
		t.Fatalf("validate should be false")
	}
}
