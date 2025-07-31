package delaytask

import (
	"context"
	"fmt"

	times "github.com/daozhonglee/go-util/times"
	"github.com/redis/go-redis/v9"
)

const (
	INTERVAL_MILLISECONDS = 1
	INTERVAL_SECONDS      = 1000
	INTERVAL_MUNITES      = 60 * 1000
	INTERVAL_HOUR         = 60 * 60 * 1000
)

func PushTaskInternal(r *redis.Client, ctx context.Context, taskname string, tick_time int64, content string, interval int64) error {
	next_tick := (tick_time + interval - 1) / interval
	keys := []string{fmt.Sprintf("dtaskq:{%s}:%d", taskname, next_tick)}
	args := []string{content, fmt.Sprintf("%d", 60*60*60*24)} // WARN: expire time

	_, err := r.Eval(ctx, "redis.call('rpush', KEYS[1], ARGV[1]); redis.call('expire', KEYS[1], ARGV[2]); return 0", keys, args).Result()
	return err
}

/* Returns next tick, tasklist */
func PullTaskInternal(r *redis.Client, ctx context.Context, taskname string, interval int64) ([]string, error) {
	keys := []string{fmt.Sprintf("dtaskt:{%s}", taskname)}

	tm := (times.GetCurrentMilliUnix() - 1) / interval // delayed than the push task
	args := []string{fmt.Sprintf("%d", tm), fmt.Sprintf("%d", 7200)}

	if ts, err := r.Eval(ctx, "local v = redis.call('get', KEYS[1]); if (not v) then redis.call('setex', KEYS[1], ARGV[2], ARGV[1]); return ARGV[1] end return v", keys, args).Text(); err != nil {
		return []string{}, err
	} else {
		keys = append(keys, fmt.Sprintf("dtaskq:{%s}:%s", taskname, ts))
		result := r.Eval(ctx, "local rt={} "+
			"for i=1,100,1 do "+
			"    local v1=redis.call(\"lpop\", KEYS[2]); "+
			"    if (not v1) then break; end "+
			"    rt[i] = v1; "+
			"end; "+
			"if #rt == 0 then "+
			"    local v = redis.call('get', KEYS[1]); "+
			"    if (not v) then "+
			"        redis.call('setex', KEYS[1], ARGV[2], ARGV[1]); return rt; "+
			"    end "+
			"    if tonumber(v) < tonumber(ARGV[1]) then redis.call('incr', KEYS[1]); redis.call('expire', KEYS[1], ARGV[2]) end "+
			"end return rt;", keys, args)

		if ret, err := result.StringSlice(); err != nil {
			return []string{}, err
		} else {
			return ret, nil
		}
	}
}

func PushTask(r *redis.Client, ctx context.Context, taskname string, tick_time int64, content string) error {
	return PushTaskInternal(r, ctx, taskname, tick_time, content, INTERVAL_SECONDS)
}

/* Returns next tick, tasklist */
func PullTask(r *redis.Client, ctx context.Context, taskname string) ([]string, error) {
	return PullTaskInternal(r, ctx, taskname, INTERVAL_SECONDS)
}
