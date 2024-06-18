import Redis from 'ioredis'

// Connect to Redis
const redis = new Redis({
  host: 'localhost',
  port: 6379,
  password: 'yourpassword', // remove this line if you don't use password
})

redis.on('connect', () => {
  console.log('Connected to Redis')
})

redis.on('error', (err) => {
  console.error('Redis error:', err)
})

// Utility function to log results
const logResult = (operation: string, result: any) => {
  console.log(`${operation}:`, result)
}
// CRUD Operations Examples
const crudOperations = async () => {
  // Create
  await redis.set('name', 'Alice')
  await redis.mset('age', '30', 'city', 'New York')
  await redis.hset('user:1000', 'username', 'bob', 'email', 'bob@example.com')
  await redis.lpush('tasks', 'task1', 'task2')
  await redis.sadd('skills', 'JavaScript', 'TypeScript')
  await redis.zadd('scores', 100, 'player1', 200, 'player2')

  // Read
  const name = await redis.get('name')
  logResult('GET name', name)
  const values = await redis.mget('age', 'city')
  logResult('MGET age, city', values)
  const email = await redis.hget('user:1000', 'email')
  logResult('HGET user:1000 email', email)
  const tasks = await redis.lrange('tasks', 0, -1)
  logResult('LRANGE tasks 0 -1', tasks)
  const skills = await redis.smembers('skills')
  logResult('SMEMBERS skills', skills)
  const scores = await redis.zrange('scores', 0, -1, 'WITHSCORES')
  logResult('ZRANGE scores 0 -1 WITHSCORES', scores)

  // Update
  await redis.set('name', 'Alice Smith')
  await redis.hset('user:1000', 'email', 'alice@example.com')
  await redis.lset('tasks', 0, 'task1-updated')
  await redis.sadd('skills', 'Node.js')
  await redis.zadd('scores', 150, 'player1')

  // Delete
  await redis.del('name')
  await redis.hdel('user:1000', 'email')
  await redis.lpop('tasks')
  await redis.srem('skills', 'JavaScript')
  await redis.zrem('scores', 'player1')
}

// Key Commands Examples
const keyCommands = async () => {
  await redis.set('temp', 'value')
  const exists = await redis.exists('temp')
  logResult('EXISTS temp', exists)
  await redis.expire('temp', 10)
  const ttl = await redis.ttl('temp')
  logResult('TTL temp', ttl)
  const type = await redis.type('temp')
  logResult('TYPE temp', type)
  await redis.rename('temp', 'temp_new')
  const newValue = await redis.get('temp_new')
  logResult('GET temp_new', newValue)
  await redis.del('temp_new')
}

// Data Structures Examples
const dataStructures = async () => {
  // Strings
  await redis.set('counter', '1')
  await redis.incr('counter')
  await redis.decr('counter')
  const counter = await redis.get('counter')
  logResult('GET counter', counter)

  // Hashes
  await redis.hset('profile:1001', 'name', 'Charlie', 'age', '25')
  const profile = await redis.hgetall('profile:1001')
  logResult('HGETALL profile:1001', profile)
  await redis.hdel('profile:1001', 'age')

  // Lists
  await redis.rpush('queue', 'item1', 'item2')
  const queue = await redis.lrange('queue', 0, -1)
  logResult('LRANGE queue 0 -1', queue)
  await redis.lpop('queue')

  // Sets
  await redis.sadd('tags', 'redis', 'database')
  const tags = await redis.smembers('tags')
  logResult('SMEMBERS tags', tags)
  const isMember = await redis.sismember('tags', 'redis')
  logResult('SISMEMBER tags redis', isMember)
  await redis.srem('tags', 'database')

  // Sorted Sets
  await redis.zadd('leaderboard', 100, 'player1', 200, 'player2')
  const leaderboard = await redis.zrange('leaderboard', 0, -1, 'WITHSCORES')
  logResult('ZRANGE leaderboard 0 -1 WITHSCORES', leaderboard)
  const rank = await redis.zrank('leaderboard', 'player1')
  logResult('ZRANK leaderboard player1', rank)
  await redis.zrem('leaderboard', 'player2')
}

// Transactions Examples
const transactions = async () => {
  const results = await redis.multi().set('foo', 'bar').incr('counter').exec()
  logResult('MULTI/EXEC transaction', results)
}

// Scripting Examples
const scripting = async () => {
  const script = `
    return redis.call('SET', KEYS[1], ARGV[1])
  `
  const result = await redis.eval(script, 1, 'mykey', 'myvalue')
  logResult('EVAL script', result)
}

// Pub/Sub Examples
const pubSub = async () => {
  const subscriber = new Redis()
  const publisher = new Redis()

  subscriber.subscribe('news', (err, count) => {
    if (err) {
      console.error('Failed to subscribe: %s', err.message)
    } else {
      console.log(`Subscribed to ${count} channel(s).`)
    }
  })

  subscriber.on('message', (channel, message) => {
    console.log(`Received message from ${channel}: ${message}`)
    subscriber.unsubscribe()
    subscriber.quit()
    publisher.quit()
  })

  publisher.publish('news', 'Hello, world!')
}

// Main function to run the examples
const main = async () => {
  await crudOperations()
  await keyCommands()
  await dataStructures()
  await transactions()
  await scripting()
  await pubSub()
  redis.quit()
}

main().catch(console.error)
