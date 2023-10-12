package handler

import (
	"MircoServer/inventory_srv/config"
	"MircoServer/inventory_srv/global"
	"MircoServer/inventory_srv/model"
	"MircoServer/inventory_srv/proto"
	"context"
	"fmt"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v9"
	goredislib "github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"sync"
)

type InventoryServer struct {
	*proto.UnimplementedInventoryServer
}

func (*InventoryServer) SetInv(ctx context.Context, req *proto.GoodsInvInfo) (*emptypb.Empty, error) {
	//设置库存， 如果我要更新库存
	var inv model.Inventory
	global.MysqlDB.Where(&model.Inventory{Goods: req.GoodsId}).First(&inv)
	inv.Goods = req.GoodsId
	inv.Stocks = req.Num

	global.MysqlDB.Save(&inv)
	return &emptypb.Empty{}, nil
}

func (*InventoryServer) InvDetail(ctx context.Context, req *proto.GoodsInvInfo) (*proto.GoodsInvInfo, error) {
	var inv model.Inventory
	if result := global.MysqlDB.Where(&model.Inventory{Goods: req.GoodsId}).First(&inv); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "没有库存信息")
	}
	return &proto.GoodsInvInfo{
		GoodsId: inv.Goods,
		Num:     inv.Stocks,
	}, nil
}

var m sync.Mutex

func (*InventoryServer) Sell(ctx context.Context, req *proto.SellInfo) (*emptypb.Empty, error) {
	//扣减库存， 本地事务 [1:10,  2:5, 3: 20]
	redisAddr := fmt.Sprintf("%s:%d", config.TheServerConfig.RedisConfig.Host, config.TheServerConfig.RedisConfig.Port)
	zap.S().Info(redisAddr)
	client := goredislib.NewClient(&goredislib.Options{
		Addr: redisAddr,
	})
	pool := goredis.NewPool(client)

	rs := redsync.New(pool)

	tx := global.MysqlDB.Begin()
	//m.Lock()'

	for _, goodInfo := range req.GoodsInfo {
		var inv model.Inventory
		////for update悲观锁
		//if result := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where(&model.Inventory{Goods: goodInfo.GoodsId}).First(&inv); result.RowsAffected == 0 {
		//	tx.Rollback() //回滚之前的操作
		//	zap.S().Info(result.Error)
		//	return nil, status.Errorf(codes.InvalidArgument, "没有库存信息")
		//}
		//if inv.Stocks < goodInfo.Num {
		//	tx.Rollback() //回滚之前的操作
		//	zap.S().Info("库存不足")
		//	return nil, status.Errorf(codes.InvalidArgument, "库存不足")
		//}
		//inv.Stocks -= goodInfo.Num
		//tx.Save(&inv)

		//for {
		//	if result := global.MysqlDB.Where(&model.Inventory{Goods: goodInfo.GoodsId}).First(&inv); result.RowsAffected == 0 {
		//		tx.Rollback()
		//		zap.S().Info("没有库存信息")
		//		return nil, status.Errorf(codes.InvalidArgument, "没有库存信息")
		//	}
		//
		//	if inv.Stocks < goodInfo.Num {
		//		tx.Rollback() //回滚之前的操作
		//		zap.S().Info("库存不足")
		//		return nil, status.Errorf(codes.InvalidArgument, "库存不足")
		//	}
		//	//扣减， 会出现数据不一致的问题 - 锁，分布式锁
		//	inv.Stocks -= goodInfo.Num
		//	zap.S().Info("Stocks Now:", inv.Stocks)
		//	//乐观锁
		//	//update inventory set stocks= stocks - 1,version = version +1 where goods = goodInfo.GoodsId and version = inv..Version
		//	if result := tx.Model(&model.Inventory{}).Select("Stocks", "Version").Where("goods = ? and version = ?", goodInfo.GoodsId, inv.Version).Updates(&model.Inventory{
		//		Stocks:  inv.Stocks,
		//		Version: inv.Version + 1,
		//	}); result.RowsAffected == 0 {
		//		zap.S().Info(">>>>>>>>>>>>>>", result.Error)
		//		zap.S().Info("库存扣减失败")
		//	} else {
		//		break
		//	}
		mutexname := "sell_mutex"
		mutex := rs.NewMutex(mutexname)
		if err := mutex.Lock(); err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "redis lock failed")
		}
		if result := global.MysqlDB.Where(&model.Inventory{Goods: goodInfo.GoodsId}).First(&inv); result.RowsAffected == 0 {
			tx.Rollback() //回滚之前的操作
			zap.S().Info(result.Error)
			return nil, status.Errorf(codes.InvalidArgument, "没有库存信息")
		}
		if inv.Stocks < goodInfo.Num {
			tx.Rollback() //回滚之前的操作
			zap.S().Info("库存不足")
			return nil, status.Errorf(codes.InvalidArgument, "库存不足")
		}
		inv.Stocks -= goodInfo.Num
		tx.Save(&inv)
		if ok, err := mutex.Unlock(); !ok || err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "redis unlock failed")
		}
	}
	tx.Commit() // 需要自己手动提交操作

	//m.Unlock()  //释放锁
	return &emptypb.Empty{}, nil
}

func (*InventoryServer) Reback(ctx context.Context, req *proto.SellInfo) (*emptypb.Empty, error) {
	//库存归还： 1：订单超时归还 2. 订单创建失败，归还之前扣减的库存 3. 手动归还
	tx := global.MysqlDB.Begin()
	for _, goodInfo := range req.GoodsInfo {
		var inv model.Inventory
		if result := global.MysqlDB.Where(&model.Inventory{Goods: goodInfo.GoodsId}).First(&inv); result.RowsAffected == 0 {
			tx.Rollback() //回滚之前的操作
			return nil, status.Errorf(codes.InvalidArgument, "没有库存信息")
		}

		//扣减， 会出现数据不一致的问题 - 锁，分布式锁
		inv.Stocks += goodInfo.Num
		tx.Save(&inv)
	}
	tx.Commit() // 需要自己手动提交操作
	return &emptypb.Empty{}, nil
}
