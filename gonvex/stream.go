package gonvex

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/synadia-io/orbit.go/jetstreamext"
)

// var nc *nats.Conn
var js jetstream.JetStream

func init() {

	// var err error
	// opts := &server.Options{ServerName: "nats1"}

	// // Initialize new server with options
	// ns, err := server.NewServer(opts)

	// if err != nil {
	// 	panic(err)
	// }

	// // Start the server via goroutine
	// go ns.Start()

	// // Wait for server to be ready for connections
	// if !ns.ReadyForConnections(4 * time.Second) {
	// 	panic("not ready for connection")
	// }

	// nc, err = nats.Connect(ns.ClientURL(), nats.InProcessServer(ns))
	// if err != nil {
	// 	panic(err)
	// }
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		panic(err)
	}
	js, err = jetstream.New(nc)
	if err != nil {
		panic(err)
	}

	//ns.WaitForShutdown()
}

func (s *Storage[T, PT]) runSnapshot(ctx context.Context) error {
	// ent := PT(new(T))
	// cons, err := s.stream.CreateOrUpdateConsumer(ctx, jetstream.ConsumerConfig{
	// 	Durable:       fmt.Sprintf("snapshoter-%s-%s", s.atype, s.boundedCtx),
	// 	MaxAckPending: 1,
	// 	AckPolicy:     jetstream.AckExplicitPolicy,
	// })
	// if err != nil {
	// 	panic(err)
	// }
	// kv, err := js.CreateOrUpdateKeyValue(ctx, jetstream.KeyValueConfig{
	// 	Bucket:  fmt.Sprintf("%s-%s", s.atype, s.boundedCtx),
	// 	Storage: jetstream.FileStorage,
	// })
	// if err != nil {
	// 	panic(err)
	// }
	// ct, err := cons.Consume(func(msg jetstream.Msg) {
	// 	etype := msg.Headers().Get(EVENT_HEADER)
	// 	ct, ok := EventRegistry[etype]
	// 	if !ok {
	// 		panic("build func no header found")
	// 	}
	// 	ev, err := ct(msg.Data())
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	kv.Get(ctx)
	// 	ent.Reduce(event)
	// 	msg.Ack()

	// })
	// if err != nil {
	// 	return fmt.Errorf("runsnapshot: %w", err)
	// }
	// go func() {
	// 	<-ctx.Done()
	// 	ct.Drain()
	// }()

	return nil
}

func NewEventStream[T any, PT Reducible2[T]](ctx context.Context) EventStream[PT] {
	t := reflect.TypeFor[T]()
	aname := t.Name()
	sep := strings.Split(t.PkgPath(), "/")
	bcname := sep[len(sep)-1]

	//ent := PT(new(T))
	// for _, v := range ent.RegisterEvents() {
	// 	etype := reflect.TypeOf(v)
	// 	if etype.Kind() == reflect.Ptr {
	// 		panic("RegisterEvents return type must be a slice of values")
	// 	}

	// 	eventDefaultRegistry[fmt.Sprintf("%s_%s", aname, etype.Name())] = v
	// 	eventNamesRegistry[v] = etype.Name()

	// }
	store := &Storage[T, PT]{atype: aname, boundedCtx: bcname}

	//	st, ok := streams[ag.Domain().Type]
	//if !ok {
	st, err := js.CreateOrUpdateStream(ctx, jetstream.StreamConfig{
		Subjects:    []string{fmt.Sprintf("%s.*", streamName(store.atype, store.boundedCtx))},
		Name:        streamName(store.atype, store.boundedCtx),
		Storage:     jetstream.FileStorage,
		AllowDirect: true,
	})
	if err != nil {
		panic(err)
	}

	//	streams[ag.Domain().Type] = st
	//}
	store.stream = st
	// if err := store.runSnapshot(ctx); err != nil {
	// 	panic(err)
	// }
	return store
}

type Reducible2[T any] interface {
	*T
	Reducible
}

type EventStream[T Reducible] interface {
	GetAggregate(context.Context, AggregateID) (*Aggregate[T], error)
	NewAggregate() *Aggregate[T]
}

type Storage[T any, PT Reducible2[T]] struct {
	stream     jetstream.Stream
	atype      string
	boundedCtx string
}

func (s *Storage[T, PT]) NewAggregate() *Aggregate[PT] {
	var ent = PT(new(T))
	agr := &Aggregate[PT]{ID: NewAggregateID(), Type: s.atype, BoundedCtx: s.boundedCtx, Root: ent}

	return agr
}

func subjectName(agtype string, bcname string, agrid string) string {
	return fmt.Sprintf("%s-%s.%s", agtype, bcname, agrid)
}
func streamName(agtype string, bcname string) string {
	return fmt.Sprintf("%s-%s", agtype, bcname)
}

func (a *Storage[T, PT]) GetAggregate(ctx context.Context, id AggregateID) (*Aggregate[PT], error) {

	var ent = a.NewAggregate()
	ent.ID = id
	subj := subjectName(a.atype, a.boundedCtx, id.String())

	//start := time.Now()

	msgs, err := jetstreamext.GetBatch(ctx, js, streamName(a.atype, a.boundedCtx), 30, jetstreamext.GetBatchSubject(subj))
	//fmt.Println(time.Since(start))
	if err != nil {
		return nil, fmt.Errorf("build func can't get msg batch: %w", err)
	}

	for msg, err := range msgs {

		if err != nil {
			if errors.Is(err, jetstreamext.ErrNoMessages) {
				return nil, ErrNoAggregate
			}
			return nil, fmt.Errorf("build func can't get msg batch: %w", err)
		}
		etype := msg.Header.Get(EVENT_HEADER)
		ct, ok := EventRegistry[etype]
		if !ok {
			return nil, fmt.Errorf("build func no header found")
		}
		ev, err := ct(msg.Data)
		if err != nil {
			return nil, fmt.Errorf("build func %w", err)
		}

		if err := ent.Root.Reduce(ev); err != nil {
			return nil, fmt.Errorf("build func: %w", err)
		}

		ent.Version = msg.Sequence
	}

	return ent, nil
}
