package esnats

import (
	"context"
	"ddd/pkg/aggregate"
	"ddd/pkg/store"
	"errors"
	"fmt"

	"github.com/nats-io/nats.go/jetstream"
)

type snapshotStore[T any] struct {
	tname      string
	boundedCtx string
	kv         jetstream.KeyValue
}

func NewSnapshotStore[T any](ctx context.Context, js jetstream.JetStream) *snapshotStore[T] {
	aname, bname := metaFromType[T]()
	store := &snapshotStore[T]{
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

func (s *snapshotStore[T]) snapshotBucketName() string {
	return fmt.Sprintf("snapshot-%s-%s", s.boundedCtx, s.tname)
}

func (s *snapshotStore[T]) Store(ctx context.Context, id aggregate.ID[T], snap []byte) error {
	_, err := s.kv.Put(ctx, string(id), snap)
	return err
}

func (s *snapshotStore[T]) Get(ctx context.Context, id aggregate.ID[T]) ([]byte, error) {
	v, err := s.kv.Get(ctx, string(id))
	if err != nil {
		if errors.Is(err, jetstream.ErrKeyNotFound) {
			return nil, store.ErrNoSnapshot
		}
		return nil, err
	}

	return v.Value(), nil
}
