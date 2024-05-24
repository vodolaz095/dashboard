package service

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/go-playground/assert/v2"
	"github.com/rs/zerolog"
	"github.com/vodolaz095/dashboard/model"
	"github.com/vodolaz095/dashboard/sensors"
	"github.com/vodolaz095/dqueue"
)

type testSensor struct {
	sensors.UnimplementedSensor
	mu     *sync.RWMutex
	T      *testing.T
	Alive  bool
	inner  float64
	Closed bool
}

func (ts *testSensor) Init(_ context.Context) error {
	ts.mu = &sync.RWMutex{}
	return nil
}

func (ts *testSensor) Ping(_ context.Context) error {
	if ts.Alive {
		return nil
	}
	return errors.New("sensor is dead")
}

func (ts *testSensor) Close(_ context.Context) error {
	ts.Closed = true
	return nil
}

func (ts *testSensor) GetValue() float64 {
	ts.mu.RLock()
	defer ts.mu.RUnlock()
	return ts.inner
}

func (ts *testSensor) GetUpdatedAt() time.Time {
	ts.mu.RLock()
	defer ts.mu.RUnlock()
	return time.Now()
}

func (ts *testSensor) Update(_ context.Context) error {
	ts.mu.Lock()
	defer ts.mu.Unlock()
	if ts.inner != 0 {
		return nil
	}
	ts.T.Logf("Zero value not allowed")
	return errors.New("zero value not allowed")
}

func (ts *testSensor) Set(val float64) {
	ts.mu.Lock()
	defer ts.mu.Unlock()
	ts.inner = val
}

func TestSensorsServiceKeepUpdated(t *testing.T) {
	t.Skipf("not implemented")
}

func TestSensorsServiceBroadcast(t *testing.T) {
	updates := []model.Update{
		{
			Name:  "testSensor",
			Value: 21.1,
			Error: "",
		},
		{
			Name:  "testSensor",
			Value: 21.2,
			Error: "",
		},
		{
			Name:  "testSensor",
			Value: 21.3,
			Error: "",
		},
		{
			Name:  "testSensor",
			Value: 0,
			Error: "zero value not allowed",
		},
	}

	zerolog.SetGlobalLevel(zerolog.TraceLevel)
	ts := testSensor{}
	ts.Name = "testSensor"
	ts.T = t
	wg := sync.WaitGroup{}
	q := dqueue.New()
	err := ts.Init(context.TODO())
	if err != nil {
		t.Error(err)
		return
	}

	service := SensorsService{
		ListOfSensors: []string{"testSensor"},
		Sensors: map[string]sensors.ISensor{
			"testSensor": &ts,
		},
		UpdateInterval: 100 * time.Millisecond,
		UpdateQueue:    &q,
	}
	err = service.Ping(context.TODO())
	if err != nil {
		if err.Error() != "sensor is dead" {
			t.Error("wrong error")
			return
		}
	} else {
		t.Error("sensor is not alive")
		return
	}
	ts.Alive = true
	err = service.Ping(context.TODO())
	if err != nil {
		t.Error(err)
		return
	}
	ctx, cancel := context.WithCancel(context.TODO())
	sub_chan, err := service.Subscribe(ctx, "test_subscriber")
	if err != nil {
		t.Error(err)
		cancel()
		return
	}
	_, err = service.Subscribe(ctx, "test_subscriber")
	if err != nil {
		if err.Error() != "duplicate subscriber" {
			t.Errorf("wrong duplicate subscribber error: %s", err.Error())
			cancel()
			return
		}
	}
	assert.Equal(t, 1, len(service.subscribers))
	list := service.List()
	assert.Equal(t, 1, len(list))
	assert.Equal(t, "testSensor", list[0].Name)

	wg.Add(1)
	go func() {
		j := 0
		for {
			select {
			case <-ctx.Done():
				t.Logf("Closing subscribber")
				wg.Done()
				return
			case payload := <-sub_chan:
				t.Logf("Update %v - payload: name=%s value=%v error=%s",
					j, payload.Name, payload.Value, payload.Error)
				assert.Equal(t, updates[j].Name, payload.Name)
				assert.Equal(t, updates[j].Error, payload.Error)
				assert.Equal(t, updates[j].Value, payload.Value)
				t.Logf("Update %v is valid", j)
				j++
			}
		}
	}()
	go func() {
		for i := range updates {
			t.Logf("Sending update %v: %s %v...", i, updates[i].Name, updates[i].Value)
			ts.Set(updates[i].Value)
			_, err1 := service.Refresh(ctx, updates[i].Name)
			if err1 != nil {
				if err1.Error() != updates[i].Error {
					t.Errorf("unexpected error %s for update %v", err1, i)
				}
			}
			time.Sleep(time.Second)
		}
		cancel()
	}()
	wg.Wait()
	err = service.Close(context.Background())
	if err != nil {
		t.Errorf("error closing: %s", err)
	}
	assert.Equal(t, true, ts.Closed)
}
