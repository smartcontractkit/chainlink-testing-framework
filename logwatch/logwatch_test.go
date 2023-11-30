package logwatch_test

import (
	"context"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"

	"github.com/smartcontractkit/chainlink-testing-framework/logwatch"
	"github.com/smartcontractkit/chainlink-testing-framework/utils/testcontext"
)

type TestCase struct {
	name                  string
	containers            int
	msg                   string
	msgsAmount            int
	msgsIntervalSeconds   float64
	exitEarly             bool
	mustNotifyList        map[string][]*regexp.Regexp
	expectedNotifications map[string][]*logwatch.LogNotification
}

func getNotificationsAmount(m map[string][]*regexp.Regexp) int {
	if m == nil {
		return 1
	}
	notificationsToAwait := 0
	for _, v := range m {
		notificationsToAwait += len(v)
	}
	return notificationsToAwait
}

// replaceContainerNamePlaceholders this function is used to replace container names with dynamic values
// so we can run tests in parallel
func replaceContainerNamePlaceholders(tc TestCase) []string {
	dynamicContainerNames := make([]string, 0)
	for i := 0; i < tc.containers; i++ {
		staticSortedIndex := strconv.Itoa(i)
		containerName := uuid.NewString()
		dynamicContainerNames = append(dynamicContainerNames, containerName)
		if tc.mustNotifyList != nil {
			tc.mustNotifyList[containerName] = tc.mustNotifyList[staticSortedIndex]
			delete(tc.mustNotifyList, staticSortedIndex)
			for _, log := range tc.expectedNotifications[staticSortedIndex] {
				log.Container = containerName
				log.Prefix = containerName
			}
			tc.expectedNotifications[containerName] = tc.expectedNotifications[staticSortedIndex]
			delete(tc.expectedNotifications, staticSortedIndex)
		}
	}
	return dynamicContainerNames
}

// startTestContainer with custom streams emitted
func startTestContainer(ctx context.Context, containerName string, msg string, amount int, intervalSeconds float64, exitEarly bool) (testcontainers.Container, error) {
	var cmd []string
	if exitEarly {
		cmd = []string{"bash", "-c",
			fmt.Sprintf(
				"for i in {1..%d}; do sleep %.2f; echo '%s'; done",
				amount,
				intervalSeconds,
				msg,
			)}
	} else {
		cmd = []string{"bash", "-c",
			fmt.Sprintf(
				"for i in {1..%d}; do sleep %.2f; echo \"%s-$i\"; done; while true; do sleep 1; done",
				amount,
				intervalSeconds,
				msg,
			)}
	}
	req := testcontainers.ContainerRequest{
		Name:  containerName,
		Image: "ubuntu:latest",
		Cmd:   cmd,
	}
	return testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
}

func TestLogWatchDocker(t *testing.T) {
	tests := []TestCase{
		{
			name:                "should read exactly 10 streams (1 container)",
			containers:          1,
			msg:                 "hello!",
			msgsAmount:          10,
			msgsIntervalSeconds: 0.1,
		},
		{
			name:                "should read exactly 10 streams even if container exits (1 container)",
			containers:          1,
			msg:                 "hello!",
			msgsAmount:          10,
			msgsIntervalSeconds: 0.1,
			exitEarly:           true,
		},
		{
			name:                "should read exactly 100 streams fast (1 container)",
			containers:          1,
			msg:                 "hello!",
			msgsAmount:          100,
			msgsIntervalSeconds: 0.01,
		},
		{
			name:                "should read exactly 100 streams fast even if container exits (1 container)",
			containers:          1,
			msg:                 "hello!",
			msgsAmount:          100,
			msgsIntervalSeconds: 0.01,
			exitEarly:           true,
		},
		{
			name:                "should read exactly 10 streams and notify 4 times (2 containers)",
			msg:                 "A\nB\nC\nD",
			containers:          2,
			msgsAmount:          1,
			msgsIntervalSeconds: 0.1,
			mustNotifyList: map[string][]*regexp.Regexp{
				"0": {
					regexp.MustCompile("A"),
					regexp.MustCompile("B"),
				},
				"1": {
					regexp.MustCompile("C"),
					regexp.MustCompile("D"),
				},
			},
			expectedNotifications: map[string][]*logwatch.LogNotification{
				"0": {
					&logwatch.LogNotification{Container: "0", Log: "A\n"},
					&logwatch.LogNotification{Container: "0", Log: "B\n"},
				},
				"1": {
					&logwatch.LogNotification{Container: "1", Log: "C\n"},
					&logwatch.LogNotification{Container: "1", Log: "D\n"},
				},
			},
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := testcontext.Get(t)
			dynamicContainerNames := replaceContainerNamePlaceholders(tc)
			lw, err := logwatch.NewLogWatch(t, tc.mustNotifyList)
			require.NoError(t, err)
			containers := make([]testcontainers.Container, 0)
			for _, cn := range dynamicContainerNames {
				container, err := startTestContainer(ctx, cn, tc.msg, tc.msgsAmount, tc.msgsIntervalSeconds, tc.exitEarly)
				require.NoError(t, err)
				name, err := container.Name(ctx)
				require.NoError(t, err)
				err = lw.ConnectContainer(context.Background(), container, name)
				require.NoError(t, err)
			}

			// streams should be there with a gap of 1 second
			time.Sleep(time.Duration(int(tc.msgsIntervalSeconds*float64(tc.msgsAmount)))*time.Second + 1*time.Second)
			lw.PrintAll()

			// all streams should be recorded
			for _, cn := range dynamicContainerNames {
				require.Len(t, lw.ContainerLogs(cn), tc.msgsAmount*getNotificationsAmount(tc.mustNotifyList))
			}

			// client must receive notifications if mustNotifyList is set
			// each container must have notifications according to their match patterns
			if tc.mustNotifyList != nil {
				notifications := make(map[string][]*logwatch.LogNotification)
				for i := 0; i < getNotificationsAmount(tc.mustNotifyList); i++ {
					msg := lw.Listen()
					if notifications[msg.Container] == nil {
						notifications[msg.Container] = make([]*logwatch.LogNotification, 0)
					}
					notifications[msg.Container] = append(notifications[msg.Container], msg)
				}
				t.Logf("notifications: %v", spew.Sdump(notifications))
				t.Logf("expectations: %v", spew.Sdump(tc.expectedNotifications))
				if !reflect.DeepEqual(tc.expectedNotifications, notifications) {
					t.Fatalf("expected: %v, got: %v", tc.expectedNotifications, notifications)
				}
			}

			defer func() {
				// testcontainers/ryuk:v0.5.1 will handle the shutdown automatically if container exited
				// container.IsReady() is inconsistent and not always showing that container has exited
				// ontainer.Terminate() and container.StopLogProducer() has known bugs, if you call them they can hang
				// forever if container is already exited
				// https://github.com/testcontainers/testcontainers-go/pull/1085
				// tried latest branch with a fix, but no luck
				// this code terminates the containers properly
				for _, c := range containers {
					if !tc.exitEarly {
						if err := lw.DisconnectContainer(c); err != nil {
							t.Fatalf("failed to disconnect container: %s", err.Error())
						}
						if err := c.Terminate(ctx); err != nil {
							t.Fatalf("failed to terminate container: %s", err.Error())
						}
					}
				}
			}()
		})
	}
}

func TestLogWatchConnectWithDelayDocker(t *testing.T) {
	t.Parallel()
	ctx := testcontext.Get(t)
	containerName := fmt.Sprintf("%s-container-%s", "TestLogWatchConnectRetryDocker", uuid.NewString())
	message := "message"
	interval := float64(1)
	amount := 10

	//set initial timeout to 0 so that it retries to connect using fibonacci backoff
	lw, err := logwatch.NewLogWatch(t, nil)
	require.NoError(t, err)
	container, err := startTestContainer(ctx, containerName, message, amount, interval, false)
	require.NoError(t, err)
	name, err := container.Name(ctx)
	require.NoError(t, err)

	time.Sleep(5 * time.Second)

	err = lw.ConnectContainer(context.Background(), container, name)
	require.NoError(t, err)

	time.Sleep(time.Duration(int(interval*float64(amount)))*time.Second + 5*time.Second)
	lw.PrintAll()

	require.Len(t, lw.ContainerLogs(containerName), amount)

	t.Cleanup(func() {
		if err := lw.Shutdown(ctx); err != nil {
			t.Fatalf("failed to shutodwn logwatch: %s", err.Error())
		}
		if err := container.Terminate(ctx); err != nil {
			t.Fatalf("failed to terminate container: %s", err.Error())
		}
	})
}

type MockedLogProducingContainer struct {
	name       string
	id         string
	isRunning  bool
	consumer   testcontainers.LogConsumer
	startError error
	startSleep time.Duration
	// acceptsLogs       bool
	stopError         error
	errorChannelError error
	startCounter      int
	messages          []string
	logMutex          sync.Mutex
	errorCh           chan error
}

func (m *MockedLogProducingContainer) Name(ctx context.Context) (string, error) {
	return m.name, nil
}

func (m *MockedLogProducingContainer) FollowOutput(consumer testcontainers.LogConsumer) {
	m.consumer = consumer
}

func (m *MockedLogProducingContainer) StartLogProducer(ctx context.Context, timeout time.Duration) error {
	m.startCounter++
	m.errorCh = make(chan error, 1)

	if m.startError != nil {
		// m.acceptsLogs = false
		return m.startError
	}

	if m.startSleep > 0 {
		time.Sleep(m.startSleep)
	}

	// m.logMutex.Lock()
	// defer m.logMutex.Unlock()
	// m.acceptsLogs = true

	// store index of last processed log
	// iterate over m.messages in a goroutine and accept new logs
	go func() {
		fmt.Println("starting log producer loop")
		lastProcessedLogIndex := -1
		for {
			time.Sleep(200 * time.Millisecond)
			{
				// m.lock("loop")
				m.errorCh <- m.errorChannelError
				if m.errorChannelError != nil {
					fmt.Println("stopping log producer loop")
					// m.unlock("loop")
					return
				}
				// m.unlock("loop")
			}
			for i, msg := range m.messages {
				time.Sleep(200 * time.Millisecond)
				// fmt.Printf("lastProcessedLogIndex: %d, i: %d\n", lastProcessedLogIndex, i)
				if i <= lastProcessedLogIndex {
					fmt.Println("skipping log")
					continue
				}
				lastProcessedLogIndex = i
				fmt.Println("processing log")
				m.consumer.Accept(testcontainers.Log{
					LogType: testcontainers.StdoutLog,
					Content: []byte(msg),
				})
			}
		}
	}()

	return nil
}

func (m *MockedLogProducingContainer) StopLogProducer() error {
	return m.stopError
}

func (m *MockedLogProducingContainer) GetLogProducerErrorChannel() <-chan error {
	return m.errorCh
}

func (m *MockedLogProducingContainer) IsRunning() bool {
	return m.isRunning
}

func (m *MockedLogProducingContainer) GetContainerID() string {
	return m.id
}

func (m *MockedLogProducingContainer) SendLog(msg string) {
	m.messages = append(m.messages, msg)
	fmt.Println("new log sent")
}

func (m *MockedLogProducingContainer) lock(msg string) {
	m.logMutex.Lock()
	fmt.Printf("lock acquired: %s\n", msg)
}

func (m *MockedLogProducingContainer) unlock(msg string) {
	m.logMutex.Unlock()
	fmt.Printf("lock released: %s\n", msg)
}

// make sure that before it stopped working it received at least 1 log
// and that when it's started again, then it removes old logs and then receives old ones again and then follow the new ones
func TestLogWatchConnectRetryMockContainer(t *testing.T) {
	t.Parallel()
	ctx := testcontext.Get(t)
	uuid := uuid.NewString()
	amount := 10
	interval := float64(1.12)

	mockedContainer := &MockedLogProducingContainer{
		name: fmt.Sprintf("%s-container-%s", t.Name(), uuid),
		id:   uuid,
		// isRunning:         true,
		startError:        nil,
		stopError:         nil,
		errorChannelError: nil,
		// acceptsLogs:       true,
	}

	lw, err := logwatch.NewLogWatch(t, nil, logwatch.WithLogProducerTimeout(1*time.Second))
	require.NoError(t, err)

	go func() {
		time.Sleep(1 * time.Second)
		{
			// mockedContainer.lock("set error")
			mockedContainer.startSleep = 1 * time.Second
			mockedContainer.errorChannelError = errors.New("test error")
			// mockedContainer.unlock("set error")
		}
		// mockedContainer.startSleep = 1 * time.Second
		time.Sleep(1 * time.Second)
		{
			// mockedContainer.lock("unset error")
			mockedContainer.errorChannelError = nil
			// mockedContainer.unlock("unset errors")
		}
	}()

	go func() {
		// time.Sleep(500 * time.Millisecond)
		for i := 0; i < amount; i++ {
			mockedContainer.SendLog(fmt.Sprintf("message-%d", i))
			time.Sleep(time.Duration(time.Duration(interval) * time.Second))
		}
	}()

	err = lw.ConnectContainer(context.Background(), mockedContainer, mockedContainer.name)
	require.NoError(t, err)

	time.Sleep(time.Duration(int(interval*float64(amount)))*time.Second + 5*time.Second)
	lw.PrintAll()

	require.Len(t, lw.ContainerLogs(mockedContainer.name), 10)
	require.Equal(t, 2, mockedContainer.startCounter)

	t.Cleanup(func() {
		if err := lw.Shutdown(ctx); err != nil {
			t.Fatalf("failed to shutodwn logwatch: %s", err.Error())
		}
	})
}

// as above, but with 3 restarts, 2nd it is still running, 3rd it works
// it should have 10 logs in the end
func TestLogWatchConnectRetryMockContainer_NotStoppedFirstTime(t *testing.T) {
	t.Parallel()
	ctx := testcontext.Get(t)
	uuid := uuid.NewString()
	amount := 10
	interval := float64(1)

	mockedContainer := &MockedLogProducingContainer{
		name:              fmt.Sprintf("%s-container-%s", t.Name(), uuid),
		id:                uuid,
		isRunning:         false,
		startError:        nil,
		stopError:         nil,
		errorChannelError: nil,
	}

	lw, err := logwatch.NewLogWatch(t, nil, logwatch.WithLogProducerTimeout(1*time.Second))
	require.NoError(t, err)

	go func() {
		for i := 0; i < amount; i++ {
			mockedContainer.SendLog(fmt.Sprintf("message-%d", i))
			time.Sleep(time.Duration(time.Duration(interval) * time.Second))
		}
		time.Sleep(500 * time.Millisecond)
	}()

	err = lw.ConnectContainer(context.Background(), mockedContainer, mockedContainer.name)
	require.NoError(t, err)

	go func() {
		mockedContainer.startSleep = 1 * time.Second
		mockedContainer.isRunning = true
		mockedContainer.errorChannelError = errors.New("read error")
		mockedContainer.startError = errors.New("still running")
		time.Sleep(2 * time.Second)
		mockedContainer.startError = nil
		mockedContainer.errorChannelError = nil
	}()

	time.Sleep(time.Duration(int(interval*float64(amount)))*time.Second + 5*time.Second)
	lw.PrintAll()

	require.Len(t, lw.ContainerLogs(mockedContainer.name), 10)
	require.Equal(t, 3, mockedContainer.startCounter)

	t.Cleanup(func() {
		if err := lw.Shutdown(ctx); err != nil {
			t.Fatalf("failed to shutodwn logwatch: %s", err.Error())
		}
	})
}
