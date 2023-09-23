package task

import (
	"errors"
	"github.com/ProtobufMan/bufman/internal/core/timewheel"
	"sync"
	"sync/atomic"
	"time"
)

var AsyncTaskManager *Manager

func InitAsyncTask() error {
	AsyncTaskManager = NewManager(
		WithName("default_async_task_mgr"),
	)
	err := AsyncTaskManager.Start()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = AsyncTaskManager.Stop()
		}
	}()

	err = AsyncTaskManager.AddGroup(pullPluginImageGroup)
	if err != nil {
		return err
	}

	err = AsyncTaskManager.AddGroup(checkPluginAvailableGroup)
	if err != nil {
		return err
	}

	err = AddCheckPluginAvailableJob()
	if err != nil {
		return err
	}

	return nil
}

var (
	ErrInitTimeWheel   = errors.New("async task manager init time wheel error")
	ErrGroupDuplicated = errors.New("async task manager group name duplicated")
	ErrNoGroup         = errors.New("async task manager no group name")
	ErrHasBeenStarted  = errors.New("async task manager has already been started")
	ErrHasBeenStopped  = errors.New("async task manager has already been stopped")
)

var (
	defaultInterval = time.Second
	defaultSlotNum  = 60 * 60
)

type Manager struct {
	name string

	isStarted atomic.Bool
	isClosed  atomic.Bool

	mu         sync.RWMutex
	groupMap   map[string]*timewheel.TimeWheel
	timeWheels []*timewheel.TimeWheel // 时间轮
	interval   time.Duration          // 时间轮扫描间隔
	slotNum    int
}

func NewManager(options ...Option) *Manager {
	m := &Manager{
		groupMap: make(map[string]*timewheel.TimeWheel),
		interval: defaultInterval,
		slotNum:  defaultSlotNum,
	}
	for _, option := range options {
		option(m)
	}

	return m
}

type Option func(m *Manager)

func WithDefaultTimeWheelInterval(interval time.Duration) Option {
	return func(m *Manager) {
		m.interval = interval
	}
}

func WithDefaultTimeWheelSlotNum(num int) Option {
	return func(m *Manager) {
		m.slotNum = num
	}
}

func WithName(name string) Option {
	return func(m *Manager) {
		m.name = name
	}
}

func (m *Manager) Start() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.isStarted.Swap(true) {
		return ErrHasBeenStarted
	}

	for _, tw := range m.timeWheels {
		tw.Start()
	}

	return nil
}

func (m *Manager) Stop() error {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.isClosed.Swap(true) {
		return ErrHasBeenStopped
	}

	for _, tw := range m.timeWheels {
		tw.Stop()
	}

	return nil
}

func (m *Manager) AddGroup(group string) error {
	return m.AddGroupWithSetting(group, m.interval, m.slotNum)
}

func (m *Manager) AddGroupWithSetting(group string, interval time.Duration, slotNum int) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.isClosed.Load() {
		return ErrHasBeenStopped
	}

	tw := timewheel.NewTimeWheel(interval, slotNum)
	if tw == nil {
		return ErrInitTimeWheel
	}

	if _, ok := m.groupMap[group]; ok {
		return ErrGroupDuplicated
	}
	m.groupMap[group] = tw

	if m.isStarted.Load() {
		tw.Start()
	}

	m.timeWheels = append(m.timeWheels, tw)
	return nil
}

func (m *Manager) AddJob(group string, delay time.Duration, key string, job func()) error {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.isClosed.Load() {
		return ErrHasBeenStopped
	}

	if tw, ok := m.groupMap[group]; !ok {
		return ErrNoGroup
	} else {
		tw.AddJob(delay, key, job)
	}

	return nil
}

func (m *Manager) RemoveJob(group string, key string) error {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.isClosed.Load() {
		return ErrHasBeenStopped
	}

	if tw, ok := m.groupMap[group]; !ok {
		return ErrNoGroup
	} else {
		tw.RemoveJob(key)
	}

	return nil
}
