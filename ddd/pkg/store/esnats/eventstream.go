package esnats

import (
	"context"
	"ddd/pkg/aggregate"
	"ddd/pkg/store"
	"errors"
	"fmt"
	"log/slog"
	"reflect"
	"strings"
	"time"

	"math"

	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/synadia-io/orbit.go/jetstreamext"
)

// const (
// 	eventTypeHeader string = "ev_type"
// )

type eventStream[T any] struct {
	aggregate.EventStream[T]
	//commandRegistry: schema.NewRegistry[Executable[T]](),
	tname      string
	boundedCtx string
	js         jetstream.JetStream
}

func metaFromType[T any]() (aname string, bctx string) {
	t := reflect.TypeFor[T]()
	if t.Kind() != reflect.Struct {
		panic("T must be a struct")
	}
	aname = t.Name()
	sep := strings.Split(t.PkgPath(), "/")
	bctx = sep[len(sep)-1]
	return
}

func NewEventStream[T any](ctx context.Context, js jetstream.JetStream) *eventStream[T] {
	aname, bcname := metaFromType[T]()
	stream := &eventStream[T]{js: js, tname: aname, boundedCtx: bcname}

	_, err := stream.js.CreateOrUpdateStream(ctx, jetstream.StreamConfig{
		Subjects:    []string{stream.allSubjects()},
		Name:        stream.streamName(),
		Storage:     jetstream.FileStorage,
		AllowDirect: true,
	})
	if err != nil {
		panic(err)
	}
	return stream
}

func (s *eventStream[T]) subjectNameForID(agrid string) string {
	return fmt.Sprintf("%s:%s.%s", s.boundedCtx, s.tname, agrid)
}

func (s *eventStream[T]) streamName() string {
	return fmt.Sprintf("%s:%s", s.boundedCtx, s.tname)
}

func (s *eventStream[T]) allSubjects() string {
	return fmt.Sprintf("%s.*", s.streamName())
}

func (s *eventStream[T]) Name() string {
	return s.streamName()
}

func (s *eventStream[T]) StoreEvent(ctx context.Context, events []aggregate.Envelope[T]) error {
	for _, event := range events {
		msg := nats.NewMsg(s.subjectNameForID(string(event.AggrID)))
		//msg.Header.Add(eventTypeHeader, event.Type)
		msg.Header.Add(jetstream.MsgIDHeader, uuid.New().String())

		msg.Data = event.Body
		retries := 0
		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
				_, err := s.js.PublishMsg(ctx, msg, jetstream.WithExpectLastSequencePerSubject(event.Version))
				if err != nil {

					var seqerr *jetstream.APIError

					if errors.As(err, &seqerr) {
						if seqerr.ErrorCode == jetstream.JSErrCodeStreamWrongLastSequence {
							slog.Warn("OCC", "version", event.Version, "name", s.subjectNameForID(string(event.AggrID)))
							retries++
							if retries > 50 {
								panic("OCC DeadLock")
							}
							continue
						}
					}
					return fmt.Errorf("store event func: %w", err)

				}
				slog.Info("EventStored", "subject", s.subjectNameForID(string(event.AggrID)), "stream", s.streamName())
				return nil
			}
		}
	}
	return nil
}

func (s *eventStream[T]) GetEvents(ctx context.Context, id aggregate.ID[T], version uint64, handler func(event []byte) error) (uint64, error) {

	subj := s.subjectNameForID(string(id))
	msgs, err := jetstreamext.GetBatch(ctx,
		s.js, s.streamName(), math.MaxInt, jetstreamext.GetBatchSubject(subj),
		jetstreamext.GetBatchSeq(version+1))
	//fmt.Println(time.Since(start))

	if err != nil {
		return 0, fmt.Errorf("get events: %w", err)
	}

	var lastevent uint64
	for msg, err := range msgs {
		if err != nil {
			if errors.Is(err, jetstreamext.ErrNoMessages) {
				return 0, store.ErrNoAggregate
			}
			return 0, fmt.Errorf("build func can't get msg batch: %w", err)
		}

		lastevent = msg.Sequence

		if err := handler(msg.Data); err != nil {
			return 0, fmt.Errorf("get events: %w", err)
		}
	}
	return lastevent, nil
}

func (e *eventStream[T]) Subscribe(ctx context.Context, name string, handler func(event []byte) error, ordered bool) {
	maxpend := 1000
	if ordered {
		maxpend = 1
	}

	cons, err := e.js.CreateOrUpdateConsumer(ctx, e.streamName(), jetstream.ConsumerConfig{
		Durable:       fmt.Sprintf("subscription-%s", name),
		DeliverPolicy: jetstream.DeliverAllPolicy,
		AckPolicy:     jetstream.AckExplicitPolicy,
		MaxAckPending: maxpend,
	})
	if err != nil {
		panic(fmt.Errorf("subscription create consumer: %w", err))
	}
	ct, err := cons.Consume(func(msg jetstream.Msg) {

		if err := handler(msg.Data()); err != nil {
			slog.Warn("Redeliver", "error", err)
			msg.NakWithDelay(1 * time.Second)
			return
		}
		msg.Ack()

	}, jetstream.ConsumeErrHandler(func(consumeCtx jetstream.ConsumeContext, err error) {}))
	if err != nil {
		panic(fmt.Errorf("subscription consume: %w", err))
	}
	go func() {
		<-ctx.Done()
		ct.Drain()
		fmt.Println("CLOSED")
	}()
}
