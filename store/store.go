package store

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/schollz/jsonstore"
	"io"
	"myvgrest/models"
	"os"
	"path/filepath"
	"regexp"
)

type Storage struct {
	path       string
	persistent bool
	store      *jsonstore.JSONStore
}

func OpenStorage(config models.Config) (*Storage, error) {
	var store *jsonstore.JSONStore
	if config.PersistentStore {
		if err := os.MkdirAll(filepath.Dir(config.StoreFile), 0770); err != nil {
			return nil, err
		}

		f, err := os.OpenFile(config.StoreFile, os.O_CREATE|os.O_RDWR, 0644)
		if err != nil {
			return nil, err
		}
		cont, err := io.ReadAll(f)
		if err != nil {
			return nil, err
		}
		if len(cont) <= 1 {
			_, err := f.Write([]byte("{}"))
			if err != nil {
				return nil, err
			}
		}
		err = f.Close()
		if err != nil {
			return nil, err
		}

		store, err = jsonstore.Open(config.StoreFile)
		if err != nil {
			return nil, err
		}
	} else {
		store = new(jsonstore.JSONStore)
	}

	return &Storage{
		path:       config.StoreFile,
		store:      store,
		persistent: config.PersistentStore,
	}, nil
}

func (s *Storage) write() error {
	if s.persistent {
		return jsonstore.Save(s.store, s.path)
	}
	return nil
}

func (s *Storage) SetPeer(peer *models.PeerWrapper) error {
	err := s.store.Set(peer.Id.String(), peer)
	if err != nil {
		return err
	}
	return s.write()
}

func (s *Storage) GetPeer(peerId uuid.UUID) (*models.PeerWrapper, error) {
	peer := new(models.PeerWrapper)
	err := s.store.Get(peerId.String(), peer)
	return peer, err
}

func (s *Storage) ListPeer() ([]*models.PeerWrapper, error) {
	rawMessages := s.store.GetAll(regexp.MustCompile("[\\s\\S]*"))
	res := make([]*models.PeerWrapper, 0)
	for _, msg := range rawMessages {
		peer := new(models.PeerWrapper)
		json.Unmarshal(msg, peer)
		res = append(res, peer)
	}
	return res, nil
}

func (s *Storage) DeletePeer(peerId uuid.UUID) error {
	s.store.Delete(peerId.String())
	return nil
}
