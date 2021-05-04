package authRepo

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/borscht/backend/config"
	"github.com/borscht/backend/internal/models"
	"github.com/gomodule/redigo/redis"
	"github.com/stretchr/testify/require"
)

func TestNewSessionRepo(t *testing.T) {
	redisConn, err := redis.Dial("tcp", config.RedisHost)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}
	sessionRepo := NewSessionRepo(redisConn)
	if sessionRepo != nil {
		return
	}
}

func TestSessionRepo_Create(t *testing.T) {
	serverRedis, err := miniredis.Run()
	if err != nil {
		t.Fatalf("cant create miniredis: %s", err)
	}
	defer serverRedis.Close()

	redisConn, err := redis.Dial("tcp", serverRedis.Addr())
	sessionRepo := NewSessionRepo(redisConn)

	c := context.Background()
	ctx := context.WithValue(c, "request_id", 1)

	sessionData := models.SessionData{
		Session: "session",
		Id:      1,
		Role:    config.RoleUser,
	}
	expected, err := json.Marshal(models.SessionInfo{
		Id:   sessionData.Id,
		Role: sessionData.Role,
	})

	err = sessionRepo.Create(ctx, sessionData)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	got, err := serverRedis.Get(headKey + sessionData.Session)
	if err != nil {
		require.NoError(t, err)
	}
	require.EqualValues(t, expected, got)
}

func TestSessionRepo_CheckTrue(t *testing.T) {
	serverRedis, err := miniredis.Run()
	if err != nil {
		t.Fatalf("cant create miniredis: %s", err)
	}
	defer serverRedis.Close()

	sessionData := models.SessionData{
		Session: "session",
		Id:      1,
		Role:    config.RoleUser,
	}
	expected := models.SessionInfo{
		Id:   sessionData.Id,
		Role: sessionData.Role,
	}
	jsonData, err := json.Marshal(expected)

	redisConn, err := redis.Dial("tcp", serverRedis.Addr())
	sessionRepo := NewSessionRepo(redisConn)

	redisConn.Do("SET", headKey+sessionData.Session, jsonData)

	c := context.Background()
	ctx := context.WithValue(c, "request_id", 1)

	got, ok, err := sessionRepo.Check(ctx, sessionData.Session)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	require.EqualValues(t, expected, got)
	require.True(t, ok)
}

func TestSessionRepo_CheckFalse(t *testing.T) {
	serverRedis, err := miniredis.Run()
	if err != nil {
		t.Fatalf("cant create miniredis: %s", err)
	}
	defer serverRedis.Close()

	redisConn, err := redis.Dial("tcp", serverRedis.Addr())
	sessionRepo := NewSessionRepo(redisConn)

	c := context.Background()
	ctx := context.WithValue(c, "request_id", 1)

	_, ok, err := sessionRepo.Check(ctx, "session")

	require.False(t, ok)
}

func TestSessionRepo_DeleteTrue(t *testing.T) {
	serverRedis, err := miniredis.Run()
	if err != nil {
		t.Fatalf("cant create miniredis: %s", err)
	}
	defer serverRedis.Close()

	sessionData := models.SessionData{
		Session: "session",
		Id:      1,
		Role:    config.RoleUser,
	}
	expected := models.SessionInfo{
		Id:   sessionData.Id,
		Role: sessionData.Role,
	}
	jsonData, err := json.Marshal(expected)

	redisConn, err := redis.Dial("tcp", serverRedis.Addr())
	sessionRepo := NewSessionRepo(redisConn)

	redisConn.Do("SET", headKey+sessionData.Session, jsonData)

	c := context.Background()
	ctx := context.WithValue(c, "request_id", 1)

	err = sessionRepo.Delete(ctx, sessionData.Session)

	require.Nil(t, err)
}
