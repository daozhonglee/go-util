package delaytask

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	times "github.com/daozhonglee/go-util/times"
	"github.com/redis/go-redis/v9"
)

// setupRedisClient 创建Redis客户端连接
func setupRedisClient() *redis.Client {
	// 支持通过环境变量配置Redis地址
	addr := os.Getenv("REDIS_ADDR")
	if addr == "" {
		addr = "127.0.0.1:6379"
	}

	return redis.NewClient(&redis.Options{
		Addr:     addr, // Redis服务器地址
		Password: "",   // 无密码
		DB:       0,    // 使用默认数据库
	})
}

// cleanupRedisKeys 清理测试用的Redis键
func cleanupRedisKeys(r *redis.Client, ctx context.Context, patterns ...string) {
	for _, pattern := range patterns {
		keys, err := r.Keys(ctx, pattern).Result()
		if err == nil && len(keys) > 0 {
			r.Del(ctx, keys...)
		}
	}
}

func TestPushTaskInternal(t *testing.T) {
	r := setupRedisClient()
	ctx := context.Background()

	// 测试Redis连接
	if err := r.Ping(ctx).Err(); err != nil {
		t.Skipf("Redis not available: %v", err)
	}

	taskname := "test_push_internal"
	content := "test content for internal push"
	tick_time := times.GetCurrentMilliUnix()
	interval := int64(INTERVAL_SECONDS)

	// 清理测试数据
	defer cleanupRedisKeys(r, ctx, fmt.Sprintf("dtaskq:{%s}:*", taskname))

	// 测试推送任务
	err := PushTaskInternal(r, ctx, taskname, tick_time, content, interval)
	if err != nil {
		t.Errorf("PushTaskInternal failed: %v", err)
	}

	// 验证任务是否被正确推送
	next_tick := (tick_time + interval - 1) / interval
	key := fmt.Sprintf("dtaskq:{%s}:%d", taskname, next_tick)

	// 检查键是否存在
	exists, err := r.Exists(ctx, key).Result()
	if err != nil {
		t.Errorf("Failed to check key existence: %v", err)
	}
	if exists != 1 {
		t.Errorf("Expected key %s to exist", key)
	}

	// 检查内容是否正确
	items, err := r.LRange(ctx, key, 0, -1).Result()
	if err != nil {
		t.Errorf("Failed to get list items: %v", err)
	}
	if len(items) != 1 || items[0] != content {
		t.Errorf("Expected content '%s', got %v", content, items)
	}
}

func TestPullTaskInternal(t *testing.T) {
	r := setupRedisClient()
	ctx := context.Background()

	// 测试Redis连接
	if err := r.Ping(ctx).Err(); err != nil {
		t.Skipf("Redis not available: %v", err)
	}

	taskname := "test_pull_internal"
	content1 := "test content 1"
	content2 := "test content 2"
	interval := int64(INTERVAL_SECONDS)

	// 清理测试数据
	defer cleanupRedisKeys(r, ctx,
		fmt.Sprintf("dtaskq:{%s}:*", taskname),
		fmt.Sprintf("dtaskt:{%s}", taskname))

	// 计算当前拉取窗口的时间
	currentPullTime := (times.GetCurrentMilliUnix() - 1) / interval * interval
	// 使用这个时间作为推送时间，确保任务在当前拉取窗口内
	pushTime := currentPullTime

	// 先推送几个任务
	err := PushTaskInternal(r, ctx, taskname, pushTime, content1, interval)
	if err != nil {
		t.Errorf("Failed to push task 1: %v", err)
	}

	err = PushTaskInternal(r, ctx, taskname, pushTime, content2, interval)
	if err != nil {
		t.Errorf("Failed to push task 2: %v", err)
	}

	// 等待一小段时间确保任务可以被拉取
	time.Sleep(50 * time.Millisecond)

	// 拉取任务
	tasks, err := PullTaskInternal(r, ctx, taskname, interval)
	if err != nil {
		t.Errorf("PullTaskInternal failed: %v", err)
	}

	// 验证拉取的任务
	if len(tasks) != 2 {
		t.Errorf("Expected 2 tasks, got %d: %v", len(tasks), tasks)
		return // 避免panic
	}

	// 验证任务内容（注意：Redis列表是FIFO，所以顺序应该保持）
	expectedTasks := []string{content1, content2}
	for i, expected := range expectedTasks {
		if i >= len(tasks) || tasks[i] != expected {
			t.Errorf("Expected task %d to be '%s', got '%s'", i, expected, tasks[i])
		}
	}
}

func TestPushTask(t *testing.T) {
	r := setupRedisClient()
	ctx := context.Background()

	// 测试Redis连接
	if err := r.Ping(ctx).Err(); err != nil {
		t.Skipf("Redis not available: %v", err)
	}

	taskname := "test_push"
	content := "test content for push"
	tick_time := times.GetCurrentMilliUnix()

	// 清理测试数据
	defer cleanupRedisKeys(r, ctx, fmt.Sprintf("dtaskq:{%s}:*", taskname))

	// 测试推送任务
	err := PushTask(r, ctx, taskname, tick_time, content)
	if err != nil {
		t.Errorf("PushTask failed: %v", err)
	}

	// 验证任务是否被正确推送（使用默认间隔）
	next_tick := (tick_time + int64(INTERVAL_SECONDS) - 1) / int64(INTERVAL_SECONDS)
	key := fmt.Sprintf("dtaskq:{%s}:%d", taskname, next_tick)

	exists, err := r.Exists(ctx, key).Result()
	if err != nil {
		t.Errorf("Failed to check key existence: %v", err)
	}
	if exists != 1 {
		t.Errorf("Expected key %s to exist", key)
	}
}

func TestPullTask(t *testing.T) {
	r := setupRedisClient()
	ctx := context.Background()

	// 测试Redis连接
	if err := r.Ping(ctx).Err(); err != nil {
		t.Skipf("Redis not available: %v", err)
	}

	taskname := "test_pull"
	content := "test content for pull"

	// 清理测试数据
	defer cleanupRedisKeys(r, ctx,
		fmt.Sprintf("dtaskq:{%s}:*", taskname),
		fmt.Sprintf("dtaskt:{%s}", taskname))

	// 计算当前拉取窗口的时间
	currentPullTime := (times.GetCurrentMilliUnix() - 1) / int64(INTERVAL_SECONDS) * int64(INTERVAL_SECONDS)

	// 先推送一个任务
	err := PushTask(r, ctx, taskname, currentPullTime, content)
	if err != nil {
		t.Errorf("Failed to push task: %v", err)
	}

	// 等待一小段时间确保任务可以被拉取
	time.Sleep(50 * time.Millisecond)

	// 拉取任务
	tasks, err := PullTask(r, ctx, taskname)
	if err != nil {
		t.Errorf("PullTask failed: %v", err)
	}

	// 验证拉取的任务
	if len(tasks) != 1 {
		t.Errorf("Expected 1 task, got %d: %v", len(tasks), tasks)
	}

	if len(tasks) > 0 && tasks[0] != content {
		t.Errorf("Expected task content '%s', got '%s'", content, tasks[0])
	}
}

func TestPullTaskEmpty(t *testing.T) {
	r := setupRedisClient()
	ctx := context.Background()

	// 测试Redis连接
	if err := r.Ping(ctx).Err(); err != nil {
		t.Skipf("Redis not available: %v", err)
	}

	taskname := "test_pull_empty"

	// 清理测试数据
	defer cleanupRedisKeys(r, ctx,
		fmt.Sprintf("dtaskq:{%s}:*", taskname),
		fmt.Sprintf("dtaskt:{%s}", taskname))

	// 拉取不存在的任务
	tasks, err := PullTask(r, ctx, taskname)
	if err != nil {
		t.Errorf("PullTask failed: %v", err)
	}

	// 应该返回空列表
	if len(tasks) != 0 {
		t.Errorf("Expected empty task list, got %d tasks: %v", len(tasks), tasks)
	}
}

func TestDelayTaskWorkflow(t *testing.T) {
	r := setupRedisClient()
	ctx := context.Background()

	// 测试Redis连接
	if err := r.Ping(ctx).Err(); err != nil {
		t.Skipf("Redis not available: %v", err)
	}

	taskname := "test_workflow"

	// 清理测试数据
	defer cleanupRedisKeys(r, ctx,
		fmt.Sprintf("dtaskq:{%s}:*", taskname),
		fmt.Sprintf("dtaskt:{%s}", taskname))

	// 测试完整的工作流程：推送 -> 拉取 -> 再次拉取（应该为空）

	// 1. 推送多个任务
	tasks := []string{"task1", "task2", "task3"}
	currentTime := times.GetCurrentMilliUnix()
	currentPullTime := (currentTime - 1) / int64(INTERVAL_SECONDS) * int64(INTERVAL_SECONDS)

	for _, task := range tasks {
		err := PushTask(r, ctx, taskname, currentPullTime, task)
		if err != nil {
			t.Errorf("Failed to push task '%s': %v", task, err)
		}
	}

	// 2. 第一次拉取任务
	time.Sleep(50 * time.Millisecond) // 等待确保任务可以被拉取

	pulledTasks, err := PullTask(r, ctx, taskname)
	if err != nil {
		t.Errorf("Failed to pull tasks: %v", err)
	}

	if len(pulledTasks) != len(tasks) {
		t.Errorf("Expected %d tasks, got %d: %v", len(tasks), len(pulledTasks), pulledTasks)
	}

	// 验证任务内容
	for i, expected := range tasks {
		if i >= len(pulledTasks) || pulledTasks[i] != expected {
			t.Errorf("Expected task %d to be '%s', got '%s'", i, expected, pulledTasks[i])
		}
	}

	// 3. 第二次拉取任务（应该为空，因为任务已经被消费）
	secondPull, err := PullTask(r, ctx, taskname)
	if err != nil {
		t.Errorf("Failed to pull tasks second time: %v", err)
	}

	if len(secondPull) != 0 {
		t.Errorf("Expected empty task list on second pull, got %d tasks: %v", len(secondPull), secondPull)
	}
}

func TestConstants(t *testing.T) {
	// 测试时间间隔常量
	if INTERVAL_MILLISECONDS != 1 {
		t.Errorf("Expected INTERVAL_MILLISECONDS to be 1, got %d", INTERVAL_MILLISECONDS)
	}

	if INTERVAL_SECONDS != 1000 {
		t.Errorf("Expected INTERVAL_SECONDS to be 1000, got %d", INTERVAL_SECONDS)
	}

	if INTERVAL_MUNITES != 60*1000 {
		t.Errorf("Expected INTERVAL_MUNITES to be %d, got %d", 60*1000, INTERVAL_MUNITES)
	}

	if INTERVAL_HOUR != 60*60*1000 {
		t.Errorf("Expected INTERVAL_HOUR to be %d, got %d", 60*60*1000, INTERVAL_HOUR)
	}
}

func TestDifferentIntervals(t *testing.T) {
	r := setupRedisClient()
	ctx := context.Background()

	// 测试Redis连接
	if err := r.Ping(ctx).Err(); err != nil {
		t.Skipf("Redis not available: %v", err)
	}

	taskname := "test_intervals"
	content := "test content"
	tick_time := times.GetCurrentMilliUnix()

	// 清理测试数据
	defer cleanupRedisKeys(r, ctx, fmt.Sprintf("dtaskq:{%s}:*", taskname))

	// 测试不同的时间间隔
	intervals := []int64{int64(INTERVAL_MILLISECONDS), int64(INTERVAL_SECONDS), int64(INTERVAL_MUNITES), int64(INTERVAL_HOUR)}

	for i, interval := range intervals {
		taskname_with_interval := fmt.Sprintf("%s_%d", taskname, i)

		err := PushTaskInternal(r, ctx, taskname_with_interval, tick_time, content, interval)
		if err != nil {
			t.Errorf("Failed to push task with interval %d: %v", interval, err)
		}

		// 验证任务被正确推送
		next_tick := (tick_time + interval - 1) / interval
		key := fmt.Sprintf("dtaskq:{%s}:%d", taskname_with_interval, next_tick)

		exists, err := r.Exists(ctx, key).Result()
		if err != nil {
			t.Errorf("Failed to check key existence for interval %d: %v", interval, err)
		}
		if exists != 1 {
			t.Errorf("Expected key %s to exist for interval %d", key, interval)
		}
	}
}
