package constants

// RedisKeys
const (
	// ActiveUsersRK is the Redis key for the set of all active users
	ActiveUsersRK = "active_users"

	// UserSessionRK is the prefix for storing session data
	UserSessionRK = "session:%s"
)
