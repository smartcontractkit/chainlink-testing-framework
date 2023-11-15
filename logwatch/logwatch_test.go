package logwatch_test

import (
	"context"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"testing"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/google/uuid"
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
	pushToLoki            bool
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
				"for i in {1..%d}; do sleep %.2f; echo '%s'; done; while true; do sleep 1; done",
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
				err = lw.ConnectContainer(ctx, container, name, tc.pushToLoki)
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
						lw.DisconnectContainer(c)
						if err := c.Terminate(ctx); err != nil {
							t.Fatalf("failed to terminate container: %s", err.Error())
						}
					}
				}
			}()
		})
	}
}
