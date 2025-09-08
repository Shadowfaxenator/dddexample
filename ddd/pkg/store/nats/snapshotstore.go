package nats

import (
	"context"
	"ddd/pkg/store"
	"errors"
	"fmt"

	"github.com/nats-io/nats.go/jetstream"
)

type SnapshotStore struct {
	tname      string
	boundedCtx string
	kv         jetstream.KeyValue
}

func NewSnapshotStore(ctx context.Context, js jetstream.JetStream, bname string, aname string) *SnapshotStore {
	store := &SnapshotStore{
		tname:      aname,
		boundedCtx: bname,
	}
	kv, err := js.CreateOrUpdateKeyValue(ctx, jetstream.KeyValueConfig{
		Bucket:  store.snapshotBucketName(),
		Storage: jetstream.FileStorage,
	})
	if err != nil {
		panic(err)
	}

	store.kv = kv
	return store
}

func (s *SnapshotStore) snapshotBucketName() string {
	return fmt.Sprintf("snapshot-%s-%s", s.boundedCtx, s.tname)
}

func (s *SnapshotStore) Store(ctx context.Context, id string, snap []byte) error {
	_, err := s.kv.Put(ctx, id, snap)
	return err
}

func (s *SnapshotStore) Get(ctx context.Context, id string) ([]byte, error) {
	v, err := s.kv.Get(ctx, id)
	if err != nil {
		if errors.Is(err, jetstream.ErrKeyNotFound) {
			return nil, store.ErrNoSnapshot
		}
		return nil, err
	}

	return v.Value(), nil
}
