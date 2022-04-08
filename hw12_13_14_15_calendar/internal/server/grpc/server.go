package internalgrpc

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/golang/protobuf/ptypes/empty"
	pb "github.com/gregss/otus/hw12_13_14_15_calendar/api"
	"github.com/gregss/otus/hw12_13_14_15_calendar/internal/app"
	"github.com/gregss/otus/hw12_13_14_15_calendar/internal/logger"
	"github.com/gregss/otus/hw12_13_14_15_calendar/internal/storage"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Service struct {
	pb.UnimplementedCalendarServer
	logger logger.Logger
	app    app.App
}

func (s Service) CreateEvent(ctx context.Context, e *pb.Event) (*empty.Empty, error) {
	event := storage.Event{
		Title:       e.Title,
		Time:        e.Time.AsTime(),
		Duration:    time.Duration(e.Duration),
		Description: e.Description,
		UserID:      int(e.UserID),
		NotifyTime:  e.NotifyTime.AsTime(),
	}

	s.app.CreateEvent(ctx, event)

	return nil, nil
}

func (s Service) ChangeEvent(ctx context.Context, e *pb.Event) (*empty.Empty, error) {
	event := storage.Event{
		Title:       e.Title,
		Time:        e.Time.AsTime(),
		Duration:    time.Duration(e.Duration),
		Description: e.Description,
		UserID:      int(e.UserID),
		NotifyTime:  e.NotifyTime.AsTime(),
	}

	s.app.ChangeEvent(ctx, int(e.ID), event)

	return nil, nil
}

func (s Service) RemoveEvent(ctx context.Context, e *pb.Event) (*empty.Empty, error) {
	s.app.RemoveEvent(ctx, int(e.ID))

	return nil, nil
}

func (s Service) DayEvents(ctx context.Context, e *pb.Event) (*pb.EventResponse, error) {
	events := make([]*pb.Event, 0)
	for _, v := range s.app.DayEvents(ctx, e.Time.AsTime()) {
		events = append(events, &pb.Event{
			ID:          uint32(v.ID),
			Title:       v.Title,
			Time:        timestamppb.New(v.Time),
			Duration:    uint32(v.Duration),
			Description: v.Description,
			UserID:      uint32(v.UserID),
			NotifyTime:  timestamppb.New(v.NotifyTime),
		})
	}

	return &pb.EventResponse{
		Events: events,
	}, nil
}

func (s Service) WeekEvents(ctx context.Context, e *pb.Event) (*pb.EventResponse, error) {
	events := make([]*pb.Event, 0)
	for _, v := range s.app.WeekEvents(ctx, e.Time.AsTime()) {
		events = append(events, &pb.Event{
			ID:          uint32(v.ID),
			Title:       v.Title,
			Time:        timestamppb.New(v.Time),
			Duration:    uint32(v.Duration),
			Description: v.Description,
			UserID:      uint32(v.UserID),
			NotifyTime:  timestamppb.New(v.NotifyTime),
		})
	}

	return &pb.EventResponse{
		Events: events,
	}, nil
}

func (s Service) MonthEvents(ctx context.Context, e *pb.Event) (*pb.EventResponse, error) {
	events := make([]*pb.Event, 0)
	for _, v := range s.app.MonthEvents(ctx, e.Time.AsTime()) {
		events = append(events, &pb.Event{
			ID:          uint32(v.ID),
			Title:       v.Title,
			Time:        timestamppb.New(v.Time),
			Duration:    uint32(v.Duration),
			Description: v.Description,
			UserID:      uint32(v.UserID),
			NotifyTime:  timestamppb.New(v.NotifyTime),
		})
	}

	return &pb.EventResponse{
		Events: events,
	}, nil
}

func StartServer(logger logger.Logger, app app.App, port string) {
	lsn, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		logger.Error(err.Error())
	}

	server := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			func(
				ctx context.Context,
				req interface{},
				info *grpc.UnaryServerInfo,
				handler grpc.UnaryHandler) (resp interface{}, err error) {
				logger.Info(fmt.Sprintf("requested: %v", info.FullMethod))
				return nil, nil
			},
		),
	)
	pb.RegisterCalendarServer(server, Service{logger: logger, app: app})

	logger.Info(fmt.Sprintf("starting server on %s", lsn.Addr().String()))
	if err := server.Serve(lsn); err != nil {
		logger.Error(err.Error())
	}
}
