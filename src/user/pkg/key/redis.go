package key

import (
	"strconv"
)

func RedisUserSnippet(user_id int64) string {
	return "user_snippet:" + strconv.FormatInt(user_id, 36)
}
