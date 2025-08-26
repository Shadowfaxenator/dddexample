package gonvex

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"math"
	"reflect"
	"strings"

	"github.com/google/uuid"
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

func NewAggregateRoot[T any, PT Reducible2[T]](ctx context.Context) *NatsStorage[T, PT] {
	t := reflect.TypeFor[T]()
	if t.Kind() != reflect.Struct {
		panic("T must be a struct")
	}
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
	store := &NatsStorage[T, PT]{atype: aname, boundedCtx: bcname}
	ent := PT(new(T))
	for _, v := range ent.Events() {
		registerEvent(v, bcname)
	}
	//	st, ok := streams[ag.Domain().Type]
	//if !ok {
	st, err := js.CreateOrUpdateStream(ctx, jetstream.StreamConfig{
		Subjects:    []string{store.allSubjects()},
		Name:        store.streamName(),
		Storage:     jetstream.FileStorage,
		AllowDirect: true,
	})
	if err != nil {
		panic(err)
	}

	//	streams[ag.Domain().Type] = st
	//}
	store.stream = st
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

type Reducible2[T any] interface {
	*T
	Reducible
}

type AggregateRoot[AggrRoot Reducible] interface {
	Command(context.Context, *ID, func(AggrRoot) (Event, error)) (ID, error)
	StartProjection(ctx context.Context, name string, handler func(Event) any) error
	Type() string
	BoundedContext() string
}

type NatsStorage[T any, PT Reducible2[T]] struct {
	stream     jetstream.Stream
	kv         jetstream.KeyValue
	atype      string
	boundedCtx string
	proj       *sql.DB
}

func (s *NatsStorage[T, PT]) Type() string {
	return s.atype
}
func (s *NatsStorage[T, PT]) BoundedContext() string {
	return s.boundedCtx
}

func (a *NatsStorage[T, PT]) subjectNameForID(agrid string) string {
	return fmt.Sprintf("%s:%s.%s", a.boundedCtx, a.atype, agrid)
}

func (a *NatsStorage[T, PT]) streamName() string {
	return fmt.Sprintf("%s:%s", a.boundedCtx, a.atype)
}

func (a *NatsStorage[T, PT]) allSubjects() string {
	return fmt.Sprintf("%s.*", a.streamName())
}

func (a *NatsStorage[T, PT]) snapshotBucketName() string {
	return fmt.Sprintf("snapshot-%s-%s", a.boundedCtx, a.atype)
}

func eventName(etype string, bcname string) string {
	return fmt.Sprintf("%s:%s", bcname, etype)
}

func (a *NatsStorage[T, PT]) StoreEvent(ctx context.Context, id ID, version uint64, event Event) error {

	data, err := Serialize(event)
	if err != nil {
		return fmt.Errorf("send event: %w", err)
	}
	ae := &AggregateEvent{
		ID:       uuid.New(),
		BContext: a.boundedCtx,
		AggrID:   uuid.UUID(id),
		AggrType: a.atype,
		Event:    data,
	}

	msg := nats.NewMsg(a.subjectNameForID(id.String()))

	msg.Header.Add(EVENT_TYPE_HEADER, eventName(event.Type(), a.boundedCtx))
	msg.Header.Add(AGGREGATE_ID_HEADER, id.String())

	p, err := Serialize(ae)
	if err != nil {
		return fmt.Errorf("send event func: %w", err)
	}

	msg.Data = p

	_, err = js.PublishMsg(ctx, msg, jetstream.WithExpectLastSequencePerSubject(version))
	if err != nil {
		return fmt.Errorf("send event func: %w", err)
	}
	return nil

}
func (a *NatsStorage[T, PT]) setKV(ctx context.Context, id ID, data any) error {
	b, err := Serialize(data)
	if err != nil {
		return fmt.Errorf("set KV can't serialize: %w", err)
	}

	_, err = a.kv.Put(ctx, id.String(), b)
	return err
}

func (a *NatsStorage[T, PT]) getKV(ctx context.Context, id ID) (*Aggregate[PT], error) {

	entry, err := a.kv.Get(ctx, id.String())
	if err != nil {
		return nil, fmt.Errorf("get KV can't get msg batch: %w", err)
	}
	var aggrKV Aggregate[PT]
	if err := Deserialize(entry.Value(), &aggrKV); err != nil {
		return nil, fmt.Errorf("get KVfunc can't deserialize: %w", err)
	}

	// var ent = PT(new(T))
	// if err := Deserialize(entry.Value(), &ent); err != nil {
	// 	return nil, fmt.Errorf("get KVfunc can't deserialize: %w", err)
	// }
	return &aggrKV, nil
}

func (a *NatsStorage[T, PT]) build(ctx context.Context, aid ID) (PT, uint64, error) {
	var lastseq uint64
	var ent = PT(new(T))
	e, err := a.getKV(ctx, aid)

	if err == nil {
		lastseq = e.Version
		*ent = *e.Data
	}

	subj := a.subjectNameForID(aid.String())

	msgs, err := jetstreamext.GetBatch(ctx, js, a.streamName(), math.MaxInt, jetstreamext.GetBatchSubject(subj), jetstreamext.GetBatchSeq(lastseq+1))
	//fmt.Println(time.Since(start))
	if err != nil {
		return nil, lastseq, fmt.Errorf("build func can't get msg batch: %w", err)
	}
	var totalMsgs uint
	for msg, err := range msgs {

		if err != nil {
			if errors.Is(err, jetstreamext.ErrNoMessages) {
				return nil, lastseq, ErrNoAggregate
			}
			return nil, lastseq, fmt.Errorf("build func can't get msg batch: %w", err)
		}
		etype := msg.Header.Get(EVENT_TYPE_HEADER)

		ct, ok := EventRegistry[etype]
		if !ok {
			return nil, lastseq, fmt.Errorf("build func no header found")
		}
		ev, err := ct(msg.Data)
		if err != nil {
			return nil, lastseq, fmt.Errorf("build func %w", err)
		}
		//fmt.Printf("ev: %v\n", ev)
		if err := ent.Reduce(ev); err != nil {
			return nil, lastseq, fmt.Errorf("build func: %w", err)
		}
		totalMsgs++

		lastseq = msg.Sequence
	}

	if totalMsgs >= MAX_MSGS {
		go func() {

			ag := Aggregate[PT]{
				Version: lastseq,
				Data:    ent,
			}
			a.setKV(ctx, aid, ag)

		}()
	}
	return ent, lastseq, nil
}

func (a *NatsStorage[T, PT]) Command(ctx context.Context, id *ID, handler func(PT) (Event, error)) (ID, error) {
	var lastseq uint64
	var ent = PT(new(T))
	var err error
	aid := NewID()
	if id != nil {
		aid = *id
	}

	//start := time.Now()
	if id != nil {
		ent, lastseq, err = a.build(ctx, aid)
		if err != nil {
			return aid, fmt.Errorf("build aggrigate: %w", err)
		}

	}
	evt, err := handler(ent)
	if err != nil {
		return aid, fmt.Errorf("handle aggrigate: %w", err)
	}

	return aid, a.StoreEvent(ctx, aid, lastseq, evt)
}
