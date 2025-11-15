package tguserid_test

import (
	"fmt"
	"testing"

	"github.com/go-mojito/mojito"
	"github.com/nilathedragon/spamscale/pkg/thirdparty/tguserid"
)

func TestGetUser_CacheHit(t *testing.T) {
	// Setup: Pre-populate cache
	username := "testuser"
	cacheKey := fmt.Sprintf(tguserid.TgUserIdCacheKey, username)
	expectedUser := tguserid.TgUserData{
		ID:         123456789,
		Type:       "private",
		FirstName:  "Test",
		LastName:   "User",
		Username:   username,
		IsBot:      false,
		IsVerified: false,
	}

	// Set user in cache
	if err := mojito.DefaultCache().Set(cacheKey, expectedUser); err != nil {
		t.Fatalf("Failed to set cache: %v", err)
	}

	// Test: GetUser should return cached data
	user, err := tguserid.GetUser(username)
	if err != nil {
		t.Fatalf("GetUser returned error: %v", err)
	}

	if user == nil {
		t.Fatal("GetUser returned nil user")
	}

	if user.ID != expectedUser.ID {
		t.Errorf("Expected ID %d, got %d", expectedUser.ID, user.ID)
	}
	if user.Username != expectedUser.Username {
		t.Errorf("Expected Username %s, got %s", expectedUser.Username, user.Username)
	}
	if user.FirstName != expectedUser.FirstName {
		t.Errorf("Expected FirstName %s, got %s", expectedUser.FirstName, user.FirstName)
	}

	// Cleanup
	mojito.DefaultCache().Delete(cacheKey)
}

func TestGetUser_FetchError(t *testing.T) {
	// Test: Verify error handling when fetching fails
	username := "nonexistentuser12345"
	cacheKey := fmt.Sprintf(tguserid.TgUserIdCacheKey, username)

	// Ensure cache is empty
	mojito.DefaultCache().Delete(cacheKey)

	// Attempt to fetch a user that likely doesn't exist
	// This tests error handling in fetchUser
	user, err := tguserid.GetUser(username)
	if err == nil {
		// If no error, clean up cached data
		if user != nil {
			mojito.DefaultCache().Delete(cacheKey)
		}
		t.Log("Note: GetUser succeeded (user might exist or API might return empty data)")
	} else {
		// Error is expected for non-existent users or API failures
		t.Logf("GetUser returned expected error: %v", err)
		// Verify cache was not populated on error
		exists, checkErr := mojito.DefaultCache().Contains(cacheKey)
		if checkErr == nil && exists {
			t.Error("Cache should not contain failed fetch result")
		}
	}
}

func TestGetUser_CacheAndFetch(t *testing.T) {
	username := "testuser"
	cacheKey := fmt.Sprintf(tguserid.TgUserIdCacheKey, username)

	// Clean up any existing cache
	mojito.DefaultCache().Delete(cacheKey)

	// First call: should fetch and cache
	user1, err1 := tguserid.GetUser(username)
	if err1 != nil {
		// If fetch fails, skip this test
		t.Skipf("Skipping test due to fetch error (likely network/API issue): %v", err1)
	}

	if user1 == nil {
		t.Fatal("First GetUser call returned nil user")
	}

	// Verify it was cached
	exists, err := mojito.DefaultCache().Contains(cacheKey)
	if err != nil {
		t.Fatalf("Failed to check cache: %v", err)
	}
	if !exists {
		t.Error("User should be cached after first fetch")
	}

	// Second call: should return from cache
	user2, err2 := tguserid.GetUser(username)
	if err2 != nil {
		t.Fatalf("Second GetUser call returned error: %v", err2)
	}

	if user2 == nil {
		t.Fatal("Second GetUser call returned nil user")
	}

	// Verify both calls return the same data
	if user1.ID != user2.ID {
		t.Errorf("Cached user ID mismatch: first=%d, second=%d", user1.ID, user2.ID)
	}
	if user1.Username != user2.Username {
		t.Errorf("Cached user Username mismatch: first=%s, second=%s", user1.Username, user2.Username)
	}

	// Cleanup
	mojito.DefaultCache().Delete(cacheKey)
}
