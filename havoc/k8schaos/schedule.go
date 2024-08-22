package k8schaos

import (
	"context"
	"time"

	"github.com/chaos-mesh/chaos-mesh/api/v1alpha1"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type ScheduleStatus string

const (
	ScheduleStatusCreated ScheduleStatus = "created"
	ScheduleStatusDeleted ScheduleStatus = "deleted"
	ScheduleStatusUnknown ScheduleStatus = "unknown" // For any state that doesn't match the above
)

type Schedule struct {
	Object        *v1alpha1.Schedule
	Description   string
	DelayCreate   time.Duration // Delay before creating the chaos object
	Duration      time.Duration // Duration for which the chaos object should exist
	Status        ChaosStatus
	Client        client.Client
	listeners     []ChaosListener
	cancelMonitor context.CancelFunc
	startTime     time.Time
	endTime       time.Time
	logger        *zerolog.Logger
}

type ScheduleOpts struct {
	Object      *v1alpha1.Schedule
	Description string
	DelayCreate time.Duration
	Duration    time.Duration
	Client      client.Client
	Listeners   []ChaosListener
	Logger      *zerolog.Logger
}

func NewSchedule(opts ScheduleOpts) (*Schedule, error) {
	if opts.Client == nil {
		return nil, errors.New("client is required")
	}
	if opts.Object == nil {
		return nil, errors.New("chaos object is required")
	}
	if opts.Logger == nil {
		return nil, errors.New("logger is required")
	}

	return &Schedule{
		Object:      opts.Object,
		Description: opts.Description,
		DelayCreate: opts.DelayCreate,
		Duration:    opts.Duration,
		Client:      opts.Client,
		listeners:   opts.Listeners,
		logger:      opts.Logger,
	}, nil
}

// Create initiates a delayed creation of a chaos object, respecting context cancellation and deletion requests.
// It uses a timer based on `DelayCreate` and calls `create` method upon expiration unless preempted by deletion.
func (s *Schedule) Create(ctx context.Context) {
	done := make(chan struct{})

	// Create the timer with the delay to create the chaos object
	timer := time.NewTimer(s.DelayCreate)

	go func() {
		select {
		case <-ctx.Done():
			// If the context is canceled, stop the timer and exit
			if !timer.Stop() {
				<-timer.C // If the timer already expired, drain the channel
			}
			close(done) // Signal that the operation was canceled
		case <-timer.C:
			// Timer expired, check if deletion was not requested
			if s.Status != StatusDeleted {
				s.createNow(ctx)
			}
			close(done) // Signal that the creation process is either done or skipped
		}
	}()
}

func (s *Schedule) Delete(ctx context.Context) error {
	if err := s.Client.Delete(ctx, s.Object); err != nil {
		return errors.Wrap(err, "failed to delete chaos object")
	}

	// Cancel the monitoring goroutine
	if s.cancelMonitor != nil {
		s.cancelMonitor()
	}

	s.endTime = time.Now()
	s.Status = StatusDeleted
	s.notifyListeners(string(ScheduleStatusDeleted))

	return nil
}

func (s *Schedule) AddListener(listener ChaosListener) {
	s.listeners = append(s.listeners, listener)
}

func (s *Schedule) GetObject() client.Object {
	return s.Object
}

func (s *Schedule) GetChaosName() string {
	return s.Object.GetName()
}

func (s *Schedule) GetChaosDescription() string {
	return s.Description
}

func (s *Schedule) GetChaosSpec() interface{} {
	return s.Object.Spec.ScheduleItem
}

func (s *Schedule) GetChaosDuration() (time.Duration, error) {
	return s.Duration, nil
}

func (s *Schedule) GetStartTime() time.Time {
	return s.startTime
}

func (s *Schedule) GetEndTime() time.Time {
	return s.endTime
}

func (s *Schedule) GetExpectedEndTime() (time.Time, error) {
	duration, err := s.GetChaosDuration()
	if err != nil {
		return time.Time{}, err
	}
	return s.startTime.Add(duration), nil
}

func (s *Schedule) createNow(ctx context.Context) {
	if err := s.Client.Create(ctx, s.Object); err != nil {
		Logger.Error().Err(err).Interface("chaos", s).Msg("failed to create chaos object")
		return
	}
	s.startTime = time.Now()
	s.Status = StatusCreated
	s.notifyListeners(string(ScheduleStatusCreated))

	// Create a cancellable context for monitorStatus
	monitorCtx, cancel := context.WithCancel(ctx)
	s.cancelMonitor = cancel
	go s.monitorStatus(monitorCtx)

	// Start a deletion timer to delete the chaos object after the specified duration
	done := make(chan struct{})
	deleteTimer := time.NewTimer(s.Duration)
	go func() {
		select {
		case <-ctx.Done():
			// Context was canceled, ensure chaos object is deleted
			if !deleteTimer.Stop() {
				<-deleteTimer.C // Drain the timer if it already fired
			}
			err := s.Delete(context.Background())
			if err != nil {
				s.logger.Error().Err(err).Msg("failed to delete chaos object")
			}
			close(done)
		case <-deleteTimer.C:
			// Duration elapsed, delete the chaos object
			err := s.Delete(context.Background())
			if err != nil {
				s.logger.Error().Err(err).Msg("failed to delete chaos object")
			}
			close(done)
		}
	}()
}

func (s *Schedule) notifyListeners(event string) {
	for _, listener := range s.listeners {
		switch event {
		case string(ScheduleStatusCreated):
			listener.OnScheduleCreated(*s)
		case string(ScheduleStatusDeleted):
			listener.OnScheduleDeleted(*s)
		}
	}
}

func (s *Schedule) monitorStatus(ctx context.Context) {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			// Context canceled, stop monitoring
			return
		case <-ticker.C:
			// Fetch the latest state of the Schedule object
			var schedule v1alpha1.Schedule
			if err := s.Client.Get(ctx, client.ObjectKey{
				Namespace: s.Object.GetNamespace(),
				Name:      s.Object.GetName(),
			}, &schedule); err != nil {
				Logger.Error().Err(err).Msg("Failed to get Schedule object")
				continue
			}

			// Log or process the schedule's current status
			// Loggerog.Info().Interface("status", schedule.Status).Msg("Current Schedule Status")
		}
	}
}
