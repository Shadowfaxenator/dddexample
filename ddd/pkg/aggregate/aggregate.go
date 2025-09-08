package aggregate

import (
	"context"
	"log/slog"

	"ddd/internal/registry"
	"ddd/internal/serde"
	"ddd/pkg/store"
	nstore "ddd/pkg/store/nats"

	"errors"
	"fmt"

	"reflect"
	"strings"

	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

// var nc *nats.Conn
var js jetstream.JetStream

type messageCount uint

const (
	snapshotSize messageCount = 100
)

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

type ID = string

func NewID() ID {
	a := uuid.New()
	return a.String()
}

// type serdeAggregate[T any] struct {
// 	body  T
// 	atype string
// }

type Option[T any] func(a *Aggregate[T])

func WithSerde[T any](s serde.Serder) Option[T] {
	return func(a *Aggregate[T]) {
		a.serder = s
	}
}

type StoreRecord struct {
	Body []byte
	Type string
}

func New[T any](ctx context.Context, opts ...Option[T]) *Aggregate[T] {
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

	aggr := &Aggregate[T]{
		stream: nstore.NewEventStream(ctx, js, bcname, aname),
		snap:   nstore.NewSnapshotStore(ctx, js, bcname, aname),
		serder: serde.DefaultSerder{},

		//serder:          &DefaultSerder[T]{},
	}
	for _, o := range opts {
		o(aggr)
	}

	aggr.eventRegistry = registry.New(aggr.serder)
	aggr.commandRegistry = registry.New(aggr.serder)
	//var ent T
	// ent.Events(func(e Applyable[T]) {

	// 	store.eventRegistry.Add(e)
	// })
	// for _, v := range ent.Events() {

	// }
	//	st, ok := streams[ag.Domain().Type]
	//if !ok {

	//	streams[ag.Domain().Type] = st
	//}

	return aggr
}

type EventStore interface {
	CreateStream(ctx context.Context, aggName string, bCtx string) *EventStream
}

type EventStream interface {
	StoreEvent(ctx context.Context, ID string, version uint64, event []byte) error
	GetEvents(ctx context.Context, ID string, version uint64, handler func(event []byte) error) (uint64, error)
	Subscribe(ctx context.Context, name string, handler func(event []byte) error, ordered bool)
	Name() string
}

type SnapshotStore interface {
	Store(ctx context.Context, ID string, snap []byte) error
	Get(ctx context.Context, ID string) ([]byte, error)
}

// type Registry[T Reducible[T]] interface {
// 	Register(Applyable[T])
// }

type Aggregate[T any] struct {
	stream          EventStream
	snap            SnapshotStore
	serder          serde.Serder
	eventRegistry   *registry.RegistryStore
	commandRegistry *registry.RegistryStore
}

type Snapshot[T any] struct {
	MsgCount messageCount
	Version  uint64
	Body     *T
}

func (a *Aggregate[T]) build(ctx context.Context, id ID) (*Snapshot[T], error) {

	var ent T
	//var snap Snapshot[T]
	// rec, err := a.snap.Get(ctx, id)
	// if err != nil {
	// 	if !errors.Is(err, store.ErrNoSnapshot) {
	// 		return nil, fmt.Errorf("build: %w", err)
	// 	}
	// 	if err := a.serder.Deserialize(rec, &snap); err != nil {
	// 		return nil, fmt.Errorf("build: %w", err)
	// 	}

	// } else {
	// 	if err := a.serder.Deserialize(rec, &snap); err != nil {
	// 		return nil, fmt.Errorf("build: %w", err)
	// 	}
	// }

	var totalMsgs messageCount

	last, err := a.stream.GetEvents(ctx, id, 0, func(b []byte) error {

		var rec StoreRecord

		if err := a.serder.Deserialize(b, &rec); err != nil {

			return fmt.Errorf("build : %w", err)
		}

		ev, err := a.eventRegistry.Get(rec.Type, rec.Body)
		if err != nil {
			panic(fmt.Sprintf("event not registered: %s", rec.Type))
		}

		ev.(Applyer[T]).Apply(&ent)
		totalMsgs++

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("buid %w", err)
	}
	sn := &Snapshot[T]{Version: last, Body: &ent, MsgCount: totalMsgs}

	return sn, nil
}

type CommandFunc[T any] func(*T) (*Event[T], error)

func (f CommandFunc[T]) Execute(t *T) (*Event[T], error) {
	return f(t)
}

func (a *Aggregate[T]) RegisterEvent(event Applyer[T]) {
	a.eventRegistry.Register(event)
}
func (a *Aggregate[T]) RegisterCommand(command Executer[T]) {
	a.commandRegistry.Register(command)
}

func (a *Aggregate[T]) CommandFunc(ctx context.Context, id ID, command func(*T) (*Event[T], error)) error {
	return a.Command(ctx, id, CommandFunc[T](command))
}

func (a *Aggregate[T]) Command(ctx context.Context, id ID, command Executer[T]) error {

	var err error

	snapshot := &Snapshot[T]{}

	snapshot, err = a.build(ctx, id)
	if err != nil {
		if !errors.Is(err, store.ErrNoAggregate) {
			return fmt.Errorf("build aggrigate: %w", err)
		}
		snapshot = &Snapshot[T]{}
	}

	evt, err := command.Execute(snapshot.Body)
	if err != nil {
		return fmt.Errorf("command: %w", err)
	}

	b, err := a.serder.Serialize(evt.Applyer)
	if err != nil {
		return fmt.Errorf("command: %w", err)
	}

	rec := StoreRecord{Body: b, Type: evt.Type}

	r, err := a.serder.Serialize(rec)
	if err != nil {
		return fmt.Errorf("command: %w", err)
	}

	if err := a.stream.StoreEvent(ctx, id, snapshot.Version, r); err != nil {
		return fmt.Errorf("command: %w", err)
	}
	slog.Info("EventStored", "type", evt.Type, "stream", a.stream.Name())
	// Save snapshot if aggregate has more than snapshotSize messages
	if snapshot != nil {
		if snapshot.MsgCount >= snapshotSize {
			go func() {
				b, err := a.serder.Serialize(snapshot)
				if err != nil {
					slog.Warn(err.Error())
				}
				a.snap.Store(ctx, id, b)
			}()

		}
	}

	return nil
}
