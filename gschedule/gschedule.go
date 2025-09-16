package gschedule

import (
	"context"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

// Job 任务接口
type Job interface {
	Run()
}

// JobFunc 函数式任务
type JobFunc func()

func (f JobFunc) Run() { f() }

// Task 封装任务
type Task struct {
	ID       string
	Schedule Schedule
	Job      Job
	ctx      context.Context
	cancel   context.CancelFunc
}

// Schedule 接口：返回下一次执行时间点
type Schedule interface {
	Next(t time.Time) time.Time
}

// Every 固定间隔
type Every time.Duration

func (e Every) Next(t time.Time) time.Time { return t.Add(time.Duration(e)) }

// Delay 单次延迟
type Delay time.Duration

func (d Delay) Next(t time.Time) time.Time {
	if t.IsZero() {
		return time.Now().Add(time.Duration(d))
	}
	return time.Time{} // 只执行一次
}

// Cron 5 位标准：min hour day month week
type Cron struct {
	expr string
	min  []int
	hour []int
	dom  []int
	mon  []int
	dow  []int
}

// NewCron 解析 5 位 cron
// 例："* 2-4 * * 1-5"
func NewCron(expr string) (*Cron, error) {
	fields := strings.Fields(expr)
	if len(fields) != 5 {
		return nil, fmt.Errorf("cron need 5 fields")
	}
	c := &Cron{expr: expr}
	var err error
	c.min, err = parseField(fields[0], 0, 59)
	if err != nil {
		return nil, err
	}
	c.hour, err = parseField(fields[1], 0, 23)
	if err != nil {
		return nil, err
	}
	c.dom, err = parseField(fields[2], 1, 31)
	if err != nil {
		return nil, err
	}
	c.mon, err = parseField(fields[3], 1, 12)
	if err != nil {
		return nil, err
	}
	c.dow, err = parseField(fields[4], 0, 6)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (c *Cron) Next(t time.Time) time.Time {
	// 简单实现：逐分钟加，直到匹配
	t = t.Truncate(time.Minute).Add(time.Minute)
	for !t.After(time.Now().Add(366 * 24 * time.Hour)) {
		if c.match(t) {
			return t
		}
		t = t.Add(time.Minute)
	}
	return time.Time{}
}

func (c *Cron) match(t time.Time) bool {
	return contains(c.min, t.Minute()) &&
		contains(c.hour, t.Hour()) &&
		contains(c.dom, t.Day()) &&
		contains(c.mon, int(t.Month())) &&
		contains(c.dow, int(t.Weekday()))
}

// Scheduler 调度器
type Scheduler struct {
	taskMu sync.RWMutex
	tasks  map[string]*Task
	pool   chan struct{} // 协程池令牌
	wg     sync.WaitGroup
	ctx    context.Context
	cancel context.CancelFunc
}

// New 创建调度器，poolSize 控制最大并发
func New(poolSize int) *Scheduler {
	ctx, cancel := context.WithCancel(context.Background())
	return &Scheduler{
		tasks:  make(map[string]*Task),
		pool:   make(chan struct{}, poolSize),
		ctx:    ctx,
		cancel: cancel,
	}
}

// Add 添加任务
func (s *Scheduler) Add(id string, schedule Schedule, job Job) {
	s.taskMu.Lock()
	defer s.taskMu.Unlock()
	if old, ok := s.tasks[id]; ok {
		old.cancel()
	}
	ctx, cancel := context.WithCancel(s.ctx)
	t := &Task{
		ID:       id,
		Schedule: schedule,
		Job:      job,
		ctx:      ctx,
		cancel:   cancel,
	}
	s.tasks[id] = t
	s.wg.Add(1)
	go s.runTask(t)
}

// Remove 移除任务
func (s *Scheduler) Remove(id string) {
	s.taskMu.Lock()
	t, ok := s.tasks[id]
	if ok {
		t.cancel()
		delete(s.tasks, id)
	}
	s.taskMu.Unlock()
	if ok {
		s.wg.Done()
	}
}

// Stop 优雅停止
func (s *Scheduler) Stop() {
	s.cancel()
	s.wg.Wait()
}

func (s *Scheduler) runTask(t *Task) {
	defer s.wg.Done()
	timer := time.NewTimer(0)
	<-timer.C
	for {
		next := t.Schedule.Next(time.Now())
		if next.IsZero() {
			return // 单次 Delay 结束
		}
		wait := time.Until(next)
		if wait > 0 {
			select {
			case <-time.After(wait):
			case <-t.ctx.Done():
				return
			}
		}
		select {
		case s.pool <- struct{}{}:
			// 拿到令牌
			go func() {
				t.Job.Run()
				<-s.pool
			}()
		case <-t.ctx.Done():
			return
		}
	}
}

// --------------- 辅助函数 ----------------
func parseField(s string, min, max int) ([]int, error) {
	var res []int
	if s == "*" {
		for i := min; i <= max; i++ {
			res = append(res, i)
		}
		return res, nil
	}
	parts := strings.SplitSeq(s, ",")
	for part := range parts {
		if strings.Contains(part, "/") {
			// 步进，简化实现：只支持 */step
			sp := strings.Split(part, "/")
			if len(sp) != 2 || sp[0] != "*" {
				return nil, fmt.Errorf("invalid step %s", part)
			}
			step, _ := strconv.Atoi(sp[1])
			for i := min; i <= max; i += step {
				res = append(res, i)
			}
			continue
		}
		if strings.Contains(part, "-") {
			rng := strings.Split(part, "-")
			if len(rng) != 2 {
				return nil, fmt.Errorf("invalid range %s", part)
			}
			start, _ := strconv.Atoi(rng[0])
			end, _ := strconv.Atoi(rng[1])
			for i := start; i <= end; i++ {
				res = append(res, i)
			}
			continue
		}
		n, err := strconv.Atoi(part)
		if err != nil {
			return nil, err
		}
		res = append(res, n)
	}
	sort.Ints(res)
	return res, nil
}

func contains(arr []int, v int) bool {
	i := sort.SearchInts(arr, v)
	return i < len(arr) && arr[i] == v
}
