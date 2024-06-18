# Redis CRUD Operations and Advanced Features

This guide provides an overview of CRUD operations, key commands, data structures, and other useful Redis features.

## CRUD Operations

### Create

- **SET key value**: Sets the value of a key.
- **MSET key1 value1 key2 value2 ...**: Sets multiple keys to multiple values.
- **HSET key field value**: Sets the value of a field in a hash.
- **LPUSH key value**: Prepends a value to a list.
- **SADD key member**: Adds a member to a set.
- **ZADD key score member**: Adds a member to a sorted set with a score.

### Read

- **GET key**: Gets the value of a key.
- **MGET key1 key2 ...**: Gets the values of multiple keys.
- **HGET key field**: Gets the value of a field in a hash.
- **LRANGE key start stop**: Gets a range of elements from a list.
- **SMEMBERS key**: Gets all the members in a set.
- **ZRANGE key start stop**: Gets a range of members in a sorted set.

### Update

- **SET key value**: Sets the value of a key.
- **HSET key field value**: Sets the value of a field in a hash.
- **LSET key index value**: Sets the value of an element in a list by its index.
- **SADD key member**: Adds a member to a set.
- **ZADD key score member**: Adds a member to a sorted set with a score.

### Delete

- **DEL key**: Deletes a key.
- **HDEL key field**: Deletes a field in a hash.
- **LPOP key**: Removes and returns the first element of a list.
- **SREM key member**: Removes a member from a set.
- **ZREM key member**: Removes a member from a sorted set.

## Key Commands

- **EXISTS key**: Checks if a key exists.
- **EXPIRE key seconds**: Sets a timeout on a key.
- **TTL key**: Gets the remaining time to live of a key.
- **TYPE key**: Gets the type of value stored at key.
- **RENAME key newkey**: Renames a key.
- **SCAN cursor [MATCH pattern] [COUNT count]**: Incrementally iterates the keys space.

## Data Structures

### Strings

- **SET key value**: Sets the value of a key.
- **GET key**: Gets the value of a key.
- **INCR key**: Increments the integer value of a key by one.
- **DECR key**: Decrements the integer value of a key by one.

### Hashes

- **HSET key field value**: Sets the value of a field in a hash.
- **HGET key field**: Gets the value of a field in a hash.
- **HGETALL key**: Gets all the fields and values in a hash.
- **HDEL key field**: Deletes a field in a hash.

### Lists

- **LPUSH key value**: Prepends a value to a list.
- **RPUSH key value**: Appends a value to a list.
- **LPOP key**: Removes and returns the first element of a list.
- **RPOP key**: Removes and returns the last element of a list.
- **LRANGE key start stop**: Gets a range of elements from a list.

### Sets

- **SADD key member**: Adds a member to a set.
- **SREM key member**: Removes a member from a set.
- **SMEMBERS key**: Gets all the members in a set.
- **SISMEMBER key member**: Checks if a value is a member of a set.

### Sorted Sets

- **ZADD key score member**: Adds a member to a sorted set with a score.
- **ZREM key member**: Removes a member from a sorted set.
- **ZRANGE key start stop**: Gets a range of members in a sorted set.
- **ZRANK key member**: Determines the index of a member in a sorted set.

## Transactions

- **MULTI**: Marks the start of a transaction block.
- **EXEC**: Executes all commands issued after MULTI.
- **DISCARD**: Discards all commands issued after MULTI.
- **WATCH key [key ...]**: Watches the given keys to determine execution of the MULTI/EXEC block.

## Scripting

- **EVAL script numkeys key [key ...] arg [arg ...]**: Evaluates a Lua script.

## Pub/Sub

- **PUBLISH channel message**: Posts a message to a channel.
- **SUBSCRIBE channel [channel ...]**: Subscribes to the given channels.
- **UNSUBSCRIBE [channel [channel ...]]**: Unsubscribes from the given channels.

## Important Considerations

- **Persistence**: Configure RDB and AOF for data persistence.
- **Replication**: Use replication for high availability.
- **Sharding**: Implement sharding for horizontal scalability.
- **Security**: Secure your Redis instance with ACLs and encryption.
- **Performance**: Monitor and optimize performance with Redis commands and tools.
