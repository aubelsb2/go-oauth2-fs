package go_oauth2_fs

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/officeadminsorted/oauth2/v5"
	"github.com/officeadminsorted/oauth2/v5/models"
	"os"
	"path/filepath"
)

type FSStore interface {
	oauth2.TokenStore
	oauth2.ClientStore
}

func New(fs string) FSStore {
	return struct {
		oauth2.TokenStore
		oauth2.ClientStore
	}{
		TokenStore:  NewTokenStore(fs),
		ClientStore: NewClientStore(fs),
	}
}

func NewClientStore(fs string) oauth2.ClientStore {
	return &ClientStore{fs: fs}
}

func NewTokenStore(fs string) oauth2.TokenStore {
	return &TokenStore{fs: fs}
}

type TokenStore struct {
	fs string
}

type ClientStore struct {
	fs string
}

type Token struct {
	models.Token
	Extra json.RawMessage
}

type Client struct {
	models.Client
	Extra json.RawMessage
}

func (s *TokenStore) Create(ctx context.Context, info oauth2.TokenInfo) error {
	code := info.GetCode()
	fn := s.codeFileName(code)
	f, err := os.Create(fn)
	if err != nil {
		return err
	}
	defer f.Close()
	os.Symlink(fn, s.accessTokenFileName(info.GetAccess()))
	os.Symlink(fn, s.refreshTokenFileName(info.GetRefresh()))
	return json.NewEncoder(f).Encode(info)
}

func (s *TokenStore) codeFileName(code string) string {
	if code == "" {
		code = "BANANANA"
	}
	path := s.tokenPath()
	os.MkdirAll(path, 0755)
	fn := filepath.Join(path, fmt.Sprintf("%s.json", code))
	return fn
}

func (s *TokenStore) accessTokenFileName(code string) string {
	if code == "" {
		code = "BANANANA"
	}
	path := s.accessTokenPath()
	os.MkdirAll(path, 0755)
	fn := filepath.Join(path, fmt.Sprintf("%s.json", code))
	return fn
}

func (s *TokenStore) refreshTokenFileName(code string) string {
	if code == "" {
		code = "BANANANA"
	}
	path := s.refreshTokenPath()
	os.MkdirAll(path, 0755)
	fn := filepath.Join(path, fmt.Sprintf("%s.json", code))
	return fn
}

func (s *TokenStore) RemoveByCode(ctx context.Context, code string) error {
	fn := s.codeFileName(code)
	info, _ := s.GetByCode(ctx, code)
	os.Remove(s.accessTokenFileName(info.GetAccess()))
	os.Remove(s.refreshTokenFileName(info.GetRefresh()))
	return os.Remove(fn)
}

func (s *TokenStore) RemoveByAccess(ctx context.Context, access string) error {
	info, _ := s.GetByAccess(ctx, access)
	os.Remove(s.accessTokenFileName(info.GetAccess()))
	os.Remove(s.refreshTokenFileName(info.GetRefresh()))
	fn := s.codeFileName(info.GetCode())
	return os.Remove(fn)
}

func (s *TokenStore) RemoveByRefresh(ctx context.Context, refresh string) error {
	info, _ := s.GetByRefresh(ctx, refresh)
	os.Remove(s.accessTokenFileName(info.GetAccess()))
	os.Remove(s.refreshTokenFileName(info.GetRefresh()))
	fn := s.codeFileName(info.GetCode())
	return os.Remove(fn)
}

func (s *TokenStore) GetByCode(ctx context.Context, code string) (oauth2.TokenInfo, error) {
	fn := s.codeFileName(code)
	return s.readTokenFile(fn)
}

func (s *TokenStore) readTokenFile(fn string) (oauth2.TokenInfo, error) {
	f, err := os.Open(fn)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var token Token
	if err := json.NewDecoder(f).Decode(&token); err != nil {
		return nil, err
	}
	return &token, nil
}

func (s *TokenStore) GetByAccess(ctx context.Context, access string) (oauth2.TokenInfo, error) {
	fn := s.accessTokenFileName(access)
	return s.readTokenFile(fn)
}

func (s *TokenStore) GetByRefresh(ctx context.Context, refresh string) (oauth2.TokenInfo, error) {
	fn := s.refreshTokenFileName(refresh)
	return s.readTokenFile(fn)
}

func (s *TokenStore) tokenPath() string {
	return filepath.Join(s.fs, "token_codes")
}
func (s *TokenStore) accessTokenPath() string {
	return filepath.Join(s.fs, "access_tokens")
}
func (s *TokenStore) refreshTokenPath() string {
	return filepath.Join(s.fs, "refresh_tokens")
}

func (s *ClientStore) GetByID(ctx context.Context, id string) (oauth2.ClientInfo, error) {
	return s.readClientFile(s.codeFileName(id))
}

func (s *ClientStore) clientPath() string {
	return filepath.Join(s.fs, "clients")
}
func (s *ClientStore) codeFileName(code string) string {
	if code == "" {
		code = "BANANANA"
	}
	path := s.clientPath()
	os.MkdirAll(path, 0755)
	fn := filepath.Join(path, fmt.Sprintf("%s.json", code))
	return fn
}

func (s *ClientStore) readClientFile(fn string) (oauth2.ClientInfo, error) {
	f, err := os.Open(fn)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var client Client
	if err := json.NewDecoder(f).Decode(&client); err != nil {
		return nil, err
	}
	return &client, nil
}

func (s *ClientStore) AddClient(ctx context.Context, client oauth2.ClientInfo) error {
	fn := s.codeFileName(client.GetID())
	f, err := os.Create(fn)
	if err != nil {
		return err
	}
	defer f.Close()
	if err := json.NewEncoder(f).Encode(&client); err != nil {
		return err
	}
	return nil
}
