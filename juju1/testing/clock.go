// Copyright 2015 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package testing

import (
	"sort"
	"sync"
	"time"
)

// Clock implements a mock clock.Clock for testing purposes.
type Clock struct {
	mu           sync.Mutex
	now          time.Time
	alarms       []alarm
	notifyAlarms chan struct{}
}

// NewClock returns a new clock set to the supplied time.
func NewClock(now time.Time) *Clock {
	return &Clock{
		now:          now,
		notifyAlarms: make(chan struct{}, 1024),
	}
}

// Now is part of the clock.Clock interface.
func (clock *Clock) Now() time.Time {
	clock.mu.Lock()
	defer clock.mu.Unlock()
	return clock.now
}

// After is part of the clock.Clock interface.
func (clock *Clock) After(d time.Duration) <-chan time.Time {
	defer clock.notifyAlarm()
	clock.mu.Lock()
	defer clock.mu.Unlock()
	notify := make(chan time.Time, 1)
	if d <= 0 {
		notify <- clock.now
	} else {
		clock.alarms = append(clock.alarms, alarm{clock.now.Add(d), notify})
		sort.Sort(byTime(clock.alarms))
	}
	return notify
}

// Advance advances the result of Now by the supplied duration, and sends
// the "current" time on all alarms which are no longer "in the future".
func (clock *Clock) Advance(d time.Duration) {
	clock.mu.Lock()
	defer clock.mu.Unlock()
	clock.now = clock.now.Add(d)
	rung := 0
	for _, alarm := range clock.alarms {
		if clock.now.Before(alarm.time) {
			break
		}
		alarm.notify <- clock.now
		rung++
	}
	clock.alarms = clock.alarms[rung:]
}

// Alarms returns a channel on which you can read one value for every call to
// After and AfterFunc; and for every successful Timer.Reset backed by this
// clock. It might not be elegant but it's necessary when testing time logic
// that runs on a goroutine other than that of the test.
func (clock *Clock) Alarms() <-chan struct{} {
	return clock.notifyAlarms
}

// notifyAlarm sends a value on the channel exposed by Alarms().
func (clock *Clock) notifyAlarm() {
	select {
	case clock.notifyAlarms <- struct{}{}:
	default:
		panic("alarm notification buffer full")
	}
}

// alarm records the time at which we're expected to send on notify.
type alarm struct {
	time   time.Time
	notify chan time.Time
}

// byTime is used to sort alarms by time.
type byTime []alarm

func (a byTime) Len() int           { return len(a) }
func (a byTime) Less(i, j int) bool { return a[i].time.Before(a[j].time) }
func (a byTime) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
