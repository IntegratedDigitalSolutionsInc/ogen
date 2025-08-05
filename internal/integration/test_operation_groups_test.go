package integration

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	api "github.com/ogen-go/ogen/internal/integration/test_operation_groups"
)

// TestOperationGroups tests the server/per-operation-group feature.
func TestOperationGroups(t *testing.T) {
	// Create handlers for each group
	usersHandler := &testUsersHandler{}
	imagesHandler := &testImagesHandler{}

	// Create servers for each group
	usersServer, err := api.NewUsersServer(usersHandler)
	if err != nil {
		t.Fatalf("Failed to create users server: %v", err)
	}

	imagesServer, err := api.NewImagesServer(imagesHandler)
	if err != nil {
		t.Fatalf("Failed to create images server: %v", err)
	}

	// Test that each server only handles its own routes
	t.Run("UsersServer", func(t *testing.T) {
		// Should handle /users
		req := httptest.NewRequest("GET", "/users", nil)
		w := httptest.NewRecorder()
		usersServer.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}

		// Should NOT handle /images
		req = httptest.NewRequest("GET", "/images", nil)
		w = httptest.NewRecorder()
		usersServer.ServeHTTP(w, req)
		if w.Code != http.StatusNotFound {
			t.Errorf("Expected status 404 for /images, got %d", w.Code)
		}
	})

	t.Run("ImagesServer", func(t *testing.T) {
		// Should handle /images
		req := httptest.NewRequest("GET", "/images", nil)
		w := httptest.NewRecorder()
		imagesServer.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}

		// Should NOT handle /users
		req = httptest.NewRequest("GET", "/users", nil)
		w = httptest.NewRecorder()
		imagesServer.ServeHTTP(w, req)
		if w.Code != http.StatusNotFound {
			t.Errorf("Expected status 404 for /users, got %d", w.Code)
		}
	})
}

// testUsersHandler implements api.UsersHandler for testing
type testUsersHandler struct{}

func (h *testUsersHandler) CreateUser(ctx context.Context, req *api.User) (*api.User, error) {
	return req, nil
}

func (h *testUsersHandler) GetUser(ctx context.Context, params api.GetUserParams) (*api.User, error) {
	return &api.User{
		ID:    params.ID,
		Name:  "Test User",
		Email: "test@example.com",
	}, nil
}

func (h *testUsersHandler) ListUsers(ctx context.Context) ([]api.User, error) {
	return []api.User{
		{ID: "1", Name: "User 1", Email: "user1@example.com"},
		{ID: "2", Name: "User 2", Email: "user2@example.com"},
	}, nil
}

// testImagesHandler implements api.ImagesHandler for testing
type testImagesHandler struct{}

func (h *testImagesHandler) GetImage(ctx context.Context, params api.GetImageParams) (*api.Image, error) {
	return &api.Image{
		ID:   params.ID,
		URL:  "https://example.com/image.jpg",
		Size: 1024,
	}, nil
}

func (h *testImagesHandler) ListImages(ctx context.Context) ([]api.Image, error) {
	return []api.Image{
		{ID: "1", URL: "https://example.com/1.jpg", Size: 1024},
		{ID: "2", URL: "https://example.com/2.jpg", Size: 2048},
	}, nil
}

func (h *testImagesHandler) UploadImage(ctx context.Context, req *api.UploadImageReq) (*api.Image, error) {
	return &api.Image{
		ID:   "new",
		URL:  "https://example.com/new.jpg",
		Size: 4096,
	}, nil
}
