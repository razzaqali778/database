package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

// Connect to Redis
func connectRedis() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "yourpassword", // leave it empty if no password
		DB:       0,              // use default DB
	})

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Could not connect to Redis: %v", err)
	}

	log.Println("Connected to Redis")
	return rdb
}

// Utility function to log results
func logResult(operation string, result interface{}) {
	fmt.Printf("%s: %v\n", operation, result)
}

// CRUD Operations Examples
func crudOperations(rdb *redis.Client) {
	// Create
	rdb.Set(ctx, "name", "Alice", 0)
	rdb.MSet(ctx, "age", "30", "city", "New York")
	rdb.HSet(ctx, "user:1000", "username", "bob", "email", "bob@example.com")
	rdb.LPush(ctx, "tasks", "task1", "task2")
	rdb.SAdd(ctx, "skills", "JavaScript", "TypeScript")
	rdb.ZAdd(ctx, "scores", &redis.Z{Score: 100, Member: "player1"}, &redis.Z{Score: 200, Member: "player2"})

	// Read
	name, _ := rdb.Get(ctx, "name").Result()
	logResult("GET name", name)
	values, _ := rdb.MGet(ctx, "age", "city").Result()
	logResult("MGET age, city", values)
	email, _ := rdb.HGet(ctx, "user:1000", "email").Result()
	logResult("HGET user:1000 email", email)
	tasks, _ := rdb.LRange(ctx, "tasks", 0, -1).Result()
	logResult("LRANGE tasks 0 -1", tasks)
	skills, _ := rdb.SMembers(ctx, "skills").Result()
	logResult("SMEMBERS skills", skills)
	scores, _ := rdb.ZRangeWithScores(ctx, "scores", 0, -1).Result()
	logResult("ZRANGE scores 0 -1 WITHSCORES", scores)

	// Update
	rdb.Set(ctx, "name", "Alice Smith", 0)
	rdb.HSet(ctx, "user:1000", "email", "alice@example.com")
	rdb.LSet(ctx, "tasks", 0, "task1-updated")
	rdb.SAdd(ctx, "skills", "Node.js")
	rdb.ZAdd(ctx, "scores", &redis.Z{Score: 150, Member: "player1"})

	// Delete
	rdb.Del(ctx, "name")
	rdb.HDel(ctx, "user:1000", "email")
	rdb.LPop(ctx, "tasks")
	rdb.SRem(ctx, "skills", "JavaScript")
	rdb.ZRem(ctx, "scores", "player1")
}

// Key Commands Examples
func keyCommands(rdb *redis.Client) {
	rdb.Set(ctx, "temp", "value", 0)
	exists, _ := rdb.Exists(ctx, "temp").Result()
	logResult("EXISTS temp", exists)
	rdb.Expire(ctx, "temp", 10*time.Second)
	ttl, _ := rdb.TTL(ctx, "temp").Result()
	logResult("TTL temp", ttl)
	typeTemp, _ := rdb.Type(ctx, "temp").Result()
	logResult("TYPE temp", typeTemp)
	rdb.Rename(ctx, "temp", "temp_new")
	newValue, _ := rdb.Get(ctx, "temp_new").Result()
	logResult("GET temp_new", newValue)
	rdb.Del(ctx, "temp_new")
}

// Data Structures Examples
func dataStructures(rdb *redis.Client) {
	// Strings
	rdb.Set(ctx, "counter", "1", 0)
	rdb.Incr(ctx, "counter")
	rdb.Decr(ctx, "counter")
	counter, _ := rdb.Get(ctx, "counter").Result()
	logResult("GET counter", counter)

	// Hashes
	rdb.HSet(ctx, "profile:1001", "name", "Charlie", "age", "25")
	profile, _ := rdb.HGetAll(ctx, "profile:1001").Result()
	logResult("HGETALL profile:1001", profile)
	rdb.HDel(ctx, "profile:1001", "age")

	// Lists
	rdb.RPush(ctx, "queue", "item1", "item2")
	queue, _ := rdb.LRange(ctx, "queue", 0, -1).Result()
	logResult("LRANGE queue 0 -1", queue)
	rdb.LPop(ctx, "queue")

	// Sets
	rdb.SAdd(ctx, "tags", "redis", "database")
	tags, _ := rdb.SMembers(ctx, "tags").Result()
	logResult("SMEMBERS tags", tags)
	isMember, _ := rdb.SIsMember(ctx, "tags", "redis").Result()
	logResult("SISMEMBER tags redis", isMember)
	rdb.SRem(ctx, "tags", "database")

	// Sorted Sets
	rdb.ZAdd(ctx, "leaderboard", &redis.Z{Score: 100, Member: "player1"}, &redis.Z{Score: 200, Member: "player2"})
	leaderboard, _ := rdb.ZRangeWithScores(ctx, "leaderboard", 0, -1).Result()
	logResult("ZRANGE leaderboard 0 -1 WITHSCORES", leaderboard)
	rank, _ := rdb.ZRank(ctx, "leaderboard", "player1").Result()
	logResult("ZRANK leaderboard player1", rank)
	rdb.ZRem(ctx, "leaderboard", "player2")
}

// Transactions Examples
func transactions(rdb *redis.Client) {
	_, err := rdb.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		pipe.Set(ctx, "foo", "bar", 0)
		pipe.Incr(ctx, "counter")
		return nil
	})
	if err != nil {
		log.Fatalf("Transaction failed: %v", err)
	}
	logResult("MULTI/EXEC transaction", "success")
}

// Scripting Examples
func scripting(rdb *redis.Client) {
	script := `
		return redis.call('SET', KEYS[1], ARGV[1])
	`
	result, err := rdb.Eval(ctx, script, []string{"mykey"}, "myvalue").Result()
	if err != nil {
		log.Fatalf("EVAL script failed: %v", err)
	}
	logResult("EVAL script", result)
}

// Pub/Sub Examples
func pubSub() {
	subscriber := connectRedis()
	publisher := connectRedis()

	pubsub := subscriber.Subscribe(ctx, "news")
	_, err := pubsub.Receive(ctx)
	if err != nil {
		log.Fatalf("Subscribe failed: %v", err)
	}

	ch := pubsub.Channel()

	go func() {
		for msg := range ch {
			logResult(fmt.Sprintf("Received message from %s", msg.Channel), msg.Payload)
			err := subscriber.Unsubscribe(ctx, "news")
			if err != nil {
				log.Fatalf("Unsubscribe failed: %v", err)
			}
			subscriber.Close()
			publisher.Close()
		}
	}()

	err = publisher.Publish(ctx, "news", "Hello, world!").Err()
	if err != nil {
		log.Fatalf("Publish failed: %v", err)
	}
}

// Main function to run the examples
func main() {
	rdb := connectRedis()
	defer rdb.Close()

	crudOperations(rdb)
	keyCommands(rdb)
	dataStructures(rdb)
	transactions(rdb)
	scripting(rdb)
	pubSub()
}
