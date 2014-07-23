package main

import (
	"strconv"
	"time"
)

// an alternative way if you do not need to use websocket to do realtime counting

// only counting last 5 minutes, here assuming a user will stay online for five minutes after login/landing
var last_n_minutes = 2 // 2 in test example

// easily to change to yours
func redis_key_prefix() string {
	return App.Name + ":" + App.Env + ":"
}

func redis_key(suffix string) string {
	return redis_key_prefix() + suffix
}

// online key for online users every minute
// online_users_minute_5
// online_users_minute_6 // 5th and 6th minutes online users
func online_key(minute_suffix string) string {
	key := "online_users_minute_" + minute_suffix
	return redis_key(key)
}

// the key for THIS minute
func current_key() string {
	key := strconv.Itoa(time.Now().Minute())
    return online_key(key)
}

// return keys of last n minutes online users
func keys_in_last_n_minutes(n int) []string {
    now := time.Now()
    var res []string
    for i := 0; i < n ; i++ {
        ago := now.Add(-time.Duration(i) * time.Minute).Minute()
        res = append(res, online_key(strconv.Itoa(ago)))
    }
    return res
}

// add a online user to the set.
// call this operation from a ajax long pull is recommended, so
// do not need to write to redis every time user click/open a page.
func add_online_username(name string) {
    new_key := false
    key := current_key()
    if ok, _ := Redis.Exists(key); ok == false {
        new_key = true
    }
    Redis.Sadd(key, []byte(name))
    if new_key {
		// assuming a user will be offline after last_n_minutes
		expiredin := int64((last_n_minutes+1)*60)
        Redis.Expire(key, expiredin)
    }
}

// the online usernames
func online_usernames() []string {
    keys := keys_in_last_n_minutes(last_n_minutes)
    users, err := Redis.Sunion(keys...)
    if err != nil {
        return nil
    }
    var res []string
    for _, u := range users {
        res = append(res, string(u))
    }
    return res
}

// counting how many online users
// just do it from redis
func online_users_count() int {
    current_online_key := redis_key("online_users_current")
    keys := keys_in_last_n_minutes(last_n_minutes)
    Redis.Sunionstore(current_online_key, keys...)
    n, err := Redis.Scard(current_online_key)
    if err != nil {
        return -1
    }
    // go set_online_max(n)
    set_online_max(n)
    return n
}

// the max value of online users
func set_online_max(curr int) {
    max_online_key := redis_key("online_users_max")
    orig, _ := Redis.Get(max_online_key)
    n, _ := strconv.Atoi(string(orig))
    if curr > n {
        Redis.Set(max_online_key, []byte(strconv.Itoa(curr)))
    }
}

func get_online_max() int {
    max_online_key := redis_key("online_users_max")
    orig, _ := Redis.Get(max_online_key)
    n, err := strconv.Atoi(string(orig))
    if err != nil {
        return -1
    }
    return n
}
