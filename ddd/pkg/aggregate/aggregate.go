package aggregate

import (
	"context"
	"log/slog"
	"sync"

	"ddd/internal/registry"
	"ddd/internal/serde"
	"ddd/pkg/store"

	"errors"
	"fmt"

	"github.com/google/uuid"
)

// var nc *nats.Conn

type messageCount uint

const (
	snapshotSize messageCount = 100
)

func init() {
	// js, _ := jetstream.New(js)
	// c, _ := js.CreateOrUpdateConsumer(context.Background(), "d", jetstream.ConsumerConfig{})
	// c.Consume(func(msg jetstream.Msg) {},jetstream.PullMaxMessages())
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

	//ns.WaitForShutdown()
}

type ID[T any] string

func NewID[T any]() ID[T] {
	a := uuid.New()
	return ID[T](a.String())
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

func New[T any](ctx context.Context, es EventStream[T], ss SnapshotStore[T], opts ...Option[T]) *Aggregate[T] {

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
		stream: es,
		snap:   ss,
		serder: serde.DefaultSerder{},

		//serder:          &DefaultSerder[T]{},
	}
	for _, o := range opts {
		o(aggr)
	}

	aggr.eventRegistry = registry.New[Event[T]](aggr.serder)
	aggr.commandRegistry = registry.New[Command[T]](aggr.serder)
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

type EventStore[T any] interface {
	CreateStream(ctx context.Context, aggName string, bCtx string) *EventStream[T]
}

type Envelope[T any] struct {
	AggrID  ID[T]
	Version uint64
	Body    []byte
}

type EventStream[T any] interface {
	StoreEvent(ctx context.Context, events []Envelope[T]) error
	GetEvents(ctx context.Context, AggrID ID[T], fromSeq uint64, handler func(event []byte) error) (uint64, error)
	Subscribe(ctx context.Context, name string, handler func(event []byte) error, ordered bool)
}

type SnapshotStore[T any] interface {
	Store(ctx context.Context, AggrID ID[T], snap []byte) error
	Get(ctx context.Context, AggrID ID[T]) ([]byte, error)
}

//	type Registry[T Reducible[T]] interface {
//		Register(Applyable[T])
//	}
type typeRegistry[T any] interface {
	Exists(tname string) bool
	Get(tname string, b []byte) (T, error)
	Register(item T)
}

type Aggregate[T any] struct {
	stream          EventStream[T]
	snap            SnapshotStore[T]
	serder          serde.Serder
	ermu            sync.RWMutex
	eventRegistry   typeRegistry[Event[T]]
	crmu            sync.RWMutex
	commandRegistry typeRegistry[Command[T]]
}

type snapshot[T any] struct {
	MsgCount messageCount
	Version  uint64
	Body     *T
}

func (a *Aggregate[T]) build(ctx context.Context, id ID[T]) (*snapshot[T], error) {

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
		a.ermu.RLock()
		ev, err := a.eventRegistry.Get(rec.Type, rec.Body)
		if err != nil {
			panic(fmt.Sprintf("event not registered: %s", rec.Type))
		}
		a.ermu.RUnlock()
		ev.Apply(&ent)
		totalMsgs++

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("buid %w", err)
	}
	sn := &snapshot[T]{Version: last, Body: &ent, MsgCount: totalMsgs}

	return sn, nil
}

type CommandFunc[T any] func(*T) Event[T]

func (f CommandFunc[T]) Execute(t *T) Event[T] {
	return f(t)
}

func (a *Aggregate[T]) RegisterEvent(event Event[T]) {
	a.ermu.Lock()
	defer a.ermu.Unlock()
	a.eventRegistry.Register(event)
}
func (a *Aggregate[T]) RegisterCommand(command Command[T]) {
	a.commandRegistry.Register(command)
}

func (a *Aggregate[T]) CommandFunc(ctx context.Context, id ID[T], command func(*T) Event[T]) error {
	return a.Command(ctx, id, CommandFunc[T](command))
}

func (a *Aggregate[T]) Command(ctx context.Context, id ID[T], command Command[T]) error {

	var err error

	snap := &snapshot[T]{}

	snap, err = a.build(ctx, id)
	if err != nil {
		if !errors.Is(err, store.ErrNoAggregate) {
			return fmt.Errorf("build aggrigate: %w", err)
		}
		snap = &snapshot[T]{}
	}

	evt := command.Execute(snap.Body)

	if e, ok := evt.(EventError[T]); ok {
		return fmt.Errorf("command: %w", e)
	}
	b, err := a.serder.Serialize(evt)
	if err != nil {
		return fmt.Errorf("command: %w", err)
	}

	tname := registry.TypeNameFrom(evt)
	if !a.eventRegistry.Exists(tname) {
		slog.Error("event type not registered", "event", tname)
		panic(fmt.Errorf("event type not registered: %s", tname))
	}
	rec := StoreRecord{Body: b, Type: tname}

	r, err := a.serder.Serialize(rec)
	if err != nil {
		return fmt.Errorf("command: %w", err)
	}

	if err := a.stream.StoreEvent(ctx, []Envelope[T]{{AggrID: id, Version: snap.Version, Body: r}}); err != nil {
		return fmt.Errorf("command: %w", err)
	}
	slog.Info("EventStored", "type", tname)
	// Save snapshot if aggregate has more than snapshotSize messages
	if snap != nil {
		if snap.MsgCount >= snapshotSize {
			go func() {
				b, err := a.serder.Serialize(snap)
				if err != nil {
					slog.Warn(err.Error())
				}
				a.snap.Store(ctx, id, b)
			}()

		}
	}

	return nil
}
