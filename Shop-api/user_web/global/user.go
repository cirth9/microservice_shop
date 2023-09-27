package global

import (
	"Shop-api/user_web/proto"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
	"time"
)

// TokenExpireDuration token过期时间
const TokenExpireDuration = time.Hour * 2

// Secret 密钥
var Secret = []byte("secret")

// Rdb redis client
var (
	Rdb         *redis.Client
	ExpiredTime = time.Minute * 2
)

var (
	//GrpcAddress GRPC地址
	GrpcAddress string

	UserConn   *grpc.ClientConn
	UserClient proto.UserClient
)
