package handler

import (
	"MircoServer/goods_srv/global"
	"MircoServer/goods_srv/model"
	"MircoServer/goods_srv/proto"
	"context"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

//// 商品接口
//GoodsList(ctx context.Context, in *GoodsFilterRequest, opts ...grpc.CallOption) (*GoodsListResponse, error)
//// 现在用户提交订单有多个商品，你得批量查询商品的信息吧
//BatchGetGoods(ctx context.Context, in *BatchGoodsIdInfo, opts ...grpc.CallOption) (*GoodsListResponse, error)
//CreateGoods(ctx context.Context, in *CreateGoodsInfo, opts ...grpc.CallOption) (*GoodsInfoResponse, error)
//DeleteGoods(ctx context.Context, in *DeleteGoodsInfo, opts ...grpc.CallOption) (*emptypb.Empty, error)
//UpdateGoods(ctx context.Context, in *CreateGoodsInfo, opts ...grpc.CallOption) (*emptypb.Empty, error)
//GetGoodsDetail(ctx context.Context, in *GoodInfoRequest, opts ...grpc.CallOption) (*GoodsInfoResponse, error)
//// 商品分类
//GetAllCategorysList(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*CategoryListResponse, error)
//// 获取子分类
//GetSubCategory(ctx context.Context, in *CategoryListRequest, opts ...grpc.CallOption) (*SubCategoryListResponse, error)
//CreateCategory(ctx context.Context, in *CategoryInfoRequest, opts ...grpc.CallOption) (*CategoryInfoResponse, error)
//DeleteCategory(ctx context.Context, in *DeleteCategoryRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
//UpdateCategory(ctx context.Context, in *CategoryInfoRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)

//// 轮播图
//BannerList(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*BannerListResponse, error)
//CreateBanner(ctx context.Context, in *BannerRequest, opts ...grpc.CallOption) (*BannerResponse, error)
//DeleteBanner(ctx context.Context, in *BannerRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
//UpdateBanner(ctx context.Context, in *BannerRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
//// 品牌分类
//CategoryBrandList(ctx context.Context, in *CategoryBrandFilterRequest, opts ...grpc.CallOption) (*CategoryBrandListResponse, error)
//// 通过category获取brands
//GetCategoryBrandList(ctx context.Context, in *CategoryInfoRequest, opts ...grpc.CallOption) (*BrandListResponse, error)
//CreateCategoryBrand(ctx context.Context, in *CategoryBrandRequest, opts ...grpc.CallOption) (*CategoryBrandResponse, error)
//DeleteCategoryBrand(ctx context.Context, in *CategoryBrandRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
//UpdateCategoryBrand(ctx context.Context, in *CategoryBrandRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)

type GoodsServer struct {
	*proto.UnimplementedGoodsServer
}

func ModelToResponse(goods model.Goods) proto.GoodsInfoResponse {
	return proto.GoodsInfoResponse{
		Id:              goods.ID,
		CategoryId:      goods.CategoryID,
		Name:            goods.Name,
		GoodsSn:         goods.GoodsSn,
		ClickNum:        goods.ClickNum,
		SoldNum:         goods.SoldNum,
		FavNum:          goods.FavNum,
		MarketPrice:     goods.MarketPrice,
		ShopPrice:       goods.ShopPrice,
		GoodsBrief:      goods.GoodsBrief,
		ShipFree:        goods.ShipFree,
		GoodsFrontImage: goods.GoodsFrontImage,
		IsNew:           goods.IsNew,
		IsHot:           goods.IsHot,
		OnSale:          goods.OnSale,
		DescImages:      goods.DescImages,
		Images:          goods.Images,
		Category: &proto.CategoryBriefInfoResponse{
			Id:   goods.Category.ID,
			Name: goods.Category.Name,
		},
		Brand: &proto.BrandInfoResponse{
			Id:   goods.Brands.ID,
			Name: goods.Brands.Name,
			Logo: goods.Brands.Logo,
		},
	}
}

func (s *GoodsServer) GoodsList(ctx context.Context, req *proto.GoodsFilterRequest) (*proto.GoodsListResponse, error) {
	//关键词搜索、查询新品、查询热门商品、通过价格区间筛选， 通过商品分类筛选
	goodsListResponse := &proto.GoodsListResponse{}
	var goods []model.Goods
	localMysqlDB := global.MysqlDB.Model(model.Goods{})
	if req.KeyWords != "" {
		localMysqlDB = localMysqlDB.Where("name like ?", "%"+req.KeyWords+"%")
	}
	if req.IsHot {
		localMysqlDB.Where(model.Goods{IsHot: true})
	}
	if req.IsNew {
		localMysqlDB.Where(model.Goods{IsNew: true})
	}

	if req.PriceMin > 0 {
		localMysqlDB.Where("shop_price >= ?", req.PriceMin)
	}
	if req.PriceMax > 0 {
		localMysqlDB.Where("shop_price <= ?", req.PriceMax)
	}

	if req.Brand > 0 {
		localMysqlDB.Where("brands_id = ?", req.Brand)
	}

	//通过category去查询商品
	var subQuery string
	//categoryIds := make([]interface{}, 0)
	if req.TopCategory > 0 {

		var category model.Category
		if result := global.MysqlDB.First(&category, req.TopCategory); result.RowsAffected == 0 {
			return nil, status.Errorf(codes.NotFound, "商品分类不存在")
		}

		if category.Level == 1 {
			subQuery = fmt.Sprintf("select id from categories where parent_category_id in ( select id from  categories  WHERE parent_category_id = %d )", req.TopCategory)
		} else if category.Level == 2 {
			subQuery = fmt.Sprintf("select id from categories WHERE parent_category_id = %d ", req.TopCategory)
		} else if category.Level == 3 {
			subQuery = fmt.Sprintf("select id from categories WHERE id = %d ", req.TopCategory)
		}

		localMysqlDB = localMysqlDB.Where(fmt.Sprintf("category_id in (%s)", subQuery))
	}

	var count int64
	localMysqlDB.Count(&count)
	goodsListResponse.Total = int32(count)

	result := localMysqlDB.Scopes(Paginate(int(req.Pages), int(req.PagePerNums))).Find(&goods)
	if result.Error != nil {
		return nil, result.Error
	}

	for _, good := range goods {
		goodsInfoResponse := ModelToResponse(good)
		goodsListResponse.Data = append(goodsListResponse.Data, &goodsInfoResponse)
	}

	return goodsListResponse, nil
}

func (s *GoodsServer) BatchGetGoods(ctx context.Context, req *proto.BatchGoodsIdInfo) (*proto.GoodsListResponse, error) {
	goodsListResponse := &proto.GoodsListResponse{}
	var goods []model.Goods

	//调用where并不会真正执行sql 只是用来生成sql的 当调用find， first才会去执行sql，
	result := global.MysqlDB.Where(req.Id).Find(&goods)
	for _, good := range goods {
		goodsInfoResponse := ModelToResponse(good)
		goodsListResponse.Data = append(goodsListResponse.Data, &goodsInfoResponse)
	}
	goodsListResponse.Total = int32(result.RowsAffected)
	return goodsListResponse, nil
}

func (s *GoodsServer) GetGoodsDetail(ctx context.Context, req *proto.GoodInfoRequest) (*proto.GoodsInfoResponse, error) {
	var goods model.Goods

	if result := global.MysqlDB.Preload("Category").Preload("Brands").First(&goods, req.Id); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "商品不存在")
	}
	goodsInfoResponse := ModelToResponse(goods)
	return &goodsInfoResponse, nil
}

func (s *GoodsServer) CreateGoods(ctx context.Context, req *proto.CreateGoodsInfo) (*proto.GoodsInfoResponse, error) {
	var category model.Category
	if result := global.MysqlDB.First(&category, req.CategoryId); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "商品分类不存在")
	}

	var brand model.Brands
	if result := global.MysqlDB.First(&brand, req.BrandId); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "品牌不存在")
	}
	//先检查redis中是否有这个token
	//防止同一个token的数据重复插入到数据库中，如果redis中没有这个token则放入redis
	//这里没有看到图片文件是如何上传， 在微服务中 普通的文件上传已经不再使用
	goods := model.Goods{
		Brands:          brand,
		BrandsID:        brand.ID,
		Category:        category,
		CategoryID:      category.ID,
		Name:            req.Name,
		GoodsSn:         req.GoodsSn,
		MarketPrice:     req.MarketPrice,
		ShopPrice:       req.ShopPrice,
		GoodsBrief:      req.GoodsBrief,
		ShipFree:        req.ShipFree,
		Images:          req.Images,
		DescImages:      req.DescImages,
		GoodsFrontImage: req.GoodsFrontImage,
		IsNew:           req.IsNew,
		IsHot:           req.IsHot,
		OnSale:          req.OnSale,
	}

	//srv之间互相调用了
	tx := global.MysqlDB.Begin()
	result := tx.Save(&goods)
	if result.Error != nil {
		tx.Rollback()
		return nil, result.Error
	}
	tx.Commit()
	return &proto.GoodsInfoResponse{
		Id: goods.ID,
	}, nil
}

func (s *GoodsServer) DeleteGoods(ctx context.Context, req *proto.DeleteGoodsInfo) (*emptypb.Empty, error) {
	if result := global.MysqlDB.Delete(&model.Goods{BaseModel: model.BaseModel{ID: req.Id}}, req.Id); result.Error != nil {
		return nil, status.Errorf(codes.NotFound, "商品不存在")
	}
	return &emptypb.Empty{}, nil
}

func (s *GoodsServer) UpdateGoods(ctx context.Context, req *proto.CreateGoodsInfo) (*emptypb.Empty, error) {
	var goods model.Goods

	if result := global.MysqlDB.First(&goods, req.Id); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "商品不存在")
	}

	var category model.Category
	if result := global.MysqlDB.First(&category, req.CategoryId); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "商品分类不存在")
	}

	var brand model.Brands
	if result := global.MysqlDB.First(&brand, req.BrandId); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "品牌不存在")
	}

	goods.Brands = brand
	goods.BrandsID = brand.ID
	goods.Category = category
	goods.CategoryID = category.ID
	goods.Name = req.Name
	goods.GoodsSn = req.GoodsSn
	goods.MarketPrice = req.MarketPrice
	goods.ShopPrice = req.ShopPrice
	goods.GoodsBrief = req.GoodsBrief
	goods.ShipFree = req.ShipFree
	goods.Images = req.Images
	goods.DescImages = req.DescImages
	goods.GoodsFrontImage = req.GoodsFrontImage
	goods.IsNew = req.IsNew
	goods.IsHot = req.IsHot
	goods.OnSale = req.OnSale

	tx := global.MysqlDB.Begin()
	result := tx.Save(&goods)
	if result.Error != nil {
		tx.Rollback()
		return nil, result.Error
	}
	tx.Commit()
	return &emptypb.Empty{}, nil
}
