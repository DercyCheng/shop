package handler

import (
	"context"
	"fmt"
	"goods_srv/global"
	"goods_srv/model"
	proto "goods_srv/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GoodsServer struct {
	proto.UnimplementedGoodsServer
}

var _ proto.GoodsServer = &GoodsServer{}

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
		//使用外键注意加载进来
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

// 商品接口
func (s *GoodsServer) GoodsList(ctx context.Context, req *proto.GoodsFilterRequest) (*proto.GoodsListResponse, error) {
	//关键词搜索，查询新品，查询热门商品，通过价格区间筛选，通过商品分类筛选
	//使用es的目的是搜索出商品的id来，通过id拿到具体的字段信息是通过mysql来完成的
	//我们使用es是用来做搜索的，是否应该将所有的mysql字段全部在es中保存一份？
	//es用来做搜索，这个时候我们一般只会把搜索和过滤的字段信息保存到es中
	//es可以用来当作mysql使用，但是实际上mysql和es之间是互补的关系 ，  一般mysql用来做存储使用，es用来做搜索使用
	//es想要提高性能，就要将es的内存设置的够大（有最大限制）  占内存1k  2k  没必要的字段不要存
	goodsListResponse := &proto.GoodsListResponse{}
	//match bool 复合查询
	//q = q.Must(NewTermQuery("tag", "wow"))
	//q = q.Filter(NewTermQuery("account", "1"))
	//q := elastic.NewBoolQuery()
	localDB := global.DB.Model(model.Goods{})
	if req.KeyWords != "" {
		//搜索
		localDB = localDB.Where("name LIKE ?", "%"+req.KeyWords+"%")
		//q = q.Must(elastic.NewMultiMatchQuery(req.KeyWords, "name", "goods_brief"))
	}
	if req.IsHot {
		localDB = localDB.Where(model.Goods{IsHot: true})
		//Filter不会算分  Must会参数得分
		//q = q.Filter(elastic.NewTermQuery("is_hot", req.IsHot))
	}
	if req.IsNew {
		localDB = localDB.Where(model.Goods{IsNew: true})
		//q = q.Filter(elastic.NewTermQuery("is_new", req.IsHot))
	}
	if req.PriceMin > 0 {
		localDB = localDB.Where("shop_price>=?", req.PriceMin)
		//q = q.Filter(elastic.NewRangeQuery("shop_price").Gte(req.PriceMin))
	}
	if req.PriceMax > 0 {
		localDB = localDB.Where("shop_price<=?", req.PriceMax)
		//q = q.Filter(elastic.NewRangeQuery("shop_price").Lte(req.PriceMax))
	}
	if req.Brand > 0 {
		localDB = localDB.Where("brands_id=?", req.Brand)
		//q = q.Filter(elastic.NewTermQuery("brands_id", req.Brand))
	}
	//通过category去查询商品
	//用mysql查询取id放到categoryIds用es查询
	var subQuery string
	categoryIds := make([]interface{}, 0)
	var goods []model.Goods
	if req.TopCategory > 0 {
		var category model.Category
		if result := global.DB.First(&category, req.TopCategory); result.RowsAffected == 0 {
			return nil, status.Errorf(codes.NotFound, "商品分类不存在")
		}
		if category.Level == 1 {
			subQuery = fmt.Sprintf("SELECT id FROM category WHERE parent_category_id IN (SELECT id FROM category WHERE parent_category_id=%d)", req.TopCategory)
		} else if category.Level == 2 {
			subQuery = fmt.Sprintf("SELECT id FROM category WHERE parent_category_id=%d", req.TopCategory)
		} else if category.Level == 3 {
			subQuery = fmt.Sprintf("SELECT id FROM category WHERE id=%d", req.TopCategory)
		}
		type Result struct {
			ID int32
		}
		var results []Result
		global.DB.Model(&model.Category{}).Raw(subQuery).Scan(&results)
		for _, re := range results {
			categoryIds = append(categoryIds, re.ID)
		}
		//生成terms查询
		//基础知识：函数参数是...interface{}的时候传递多个值没有问题   但就是不能传[]int{1,2,3}...   必须得是[]interface{}{1,2,3}...
		//q = q.Filter(elastic.NewTermsQuery("category_id", categoryIds...))
		//localDB = localDB.Where(fmt.Sprintf("category_id in (%s)", subQuery)).Find(&goods)
		localDB = localDB.Where(fmt.Sprintf("category_id in (%s)", subQuery))
	}
	//动词 名词 执行条件 确定执行
	//.Query(q).From().Size() - 分页
	if req.Pages == 0 {
		req.Pages = 1
	}
	switch {
	case req.PagePerNums > 100:
		req.PagePerNums = 100
	case req.PagePerNums <= 0:
		req.PagePerNums = 10
	}
	//result, err := global.EsClient.Search().Index(model.EsGoods{}.GetIndexName()).Query(q).From(int(req.Pages)).Size(int(req.PagePerNums)).Do(context.Background())
	//if err != nil {
	//	return nil, err
	//}
	//goodsIds := make([]int32, 0)
	//goodsListResponse.Total = int32(result.Hits.TotalHits.Value)
	//for _, value := range result.Hits.Hits {
	//	goods := model.EsGoods{}
	//	_ = json.Unmarshal(value.Source, &goods)
	//	goodsIds = append(goodsIds, goods.ID)
	//}
	//if len(goodsIds) == 0 {
	//	return &proto.GoodsListResponse{}, nil
	//}
	//var count int64
	//localDB.Count(&count)
	//goodsListResponse.Total = int32(count)

	//if result2 := localDB.Preload("Category").Preload("Brands").Scopes(Paginate(int(req.Pages), int(req.PagePerNums))); result2.Error != nil {
	//	return nil, result2.Error
	//}
	//查询id在某个数组中的值

	re := localDB.Preload("Category").Preload("Brands").Find(&goods)
	if re.Error != nil {
		return nil, re.Error
	}
	for _, good := range goods {
		goodsInfoResponse := ModelToResponse(good)
		goodsListResponse.Data = append(goodsListResponse.Data, &goodsInfoResponse)
	}
	return goodsListResponse, nil
}

// 现在用户提交订单有多个商品，你得批量查询商品的信息吧
func (s *GoodsServer) BatchGetGoods(ctx context.Context, req *proto.BatchGoodsIdInfo) (*proto.GoodsListResponse, error) {
	goodsListResponse := &proto.GoodsListResponse{}
	var goods []model.Goods
	result := global.DB.Where(req.Id).Find(&goods)
	if result.Error != nil {
		return nil, status.Errorf(codes.InvalidArgument, "参数无效")
	}
	for _, good := range goods {
		goodsInfoResponse := ModelToResponse(good)
		goodsListResponse.Data = append(goodsListResponse.Data, &goodsInfoResponse)
	}
	goodsListResponse.Total = int32(result.RowsAffected)
	return goodsListResponse, nil
}

// 商品详情
func (s *GoodsServer) GetGoodsDetail(ctx context.Context, req *proto.GoodInfoRequest) (*proto.GoodsInfoResponse, error) {
	var goods model.Goods

	if result := global.DB.Preload("Brands").Preload("Category").Where(req.Id).First(&goods); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "参数无效")
	}
	goodsInfoResponse := ModelToResponse(goods)
	return &goodsInfoResponse, nil
}

// 新建商品
func (s *GoodsServer) CreateGoods(ctx context.Context, req *proto.CreateGoodsInfo) (*proto.GoodsInfoResponse, error) {
	var category model.Category

	if result := global.DB.First(&category, req.CategoryId); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "商品分类不存在")
	}
	var brand model.Brands
	if result := global.DB.First(&brand, req.BrandId); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "品牌不存在")
	}
	var goods model.Goods
	if result := global.DB.First(&goods, req.Id); result.RowsAffected != 0 {
		return nil, status.Errorf(codes.InvalidArgument, "商品已存在")
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
	tx := global.DB.Begin()
	if result := tx.Save(&goods); result.Error != nil {
		tx.Rollback()
		return nil, result.Error
	}
	tx.Commit()
	goodsDetailResponse := ModelToResponse(goods)
	return &goodsDetailResponse, nil
}

// 删除商品
func (s *GoodsServer) DeleteGoods(ctx context.Context, req *proto.DeleteGoodsInfo) (*proto.Empty, error) {
	if result := global.DB.Delete(&model.Goods{BaseModel: model.BaseModel{ID: req.Id}}); result.Error != nil {
		return nil, result.Error
	}
	return &proto.Empty{}, nil
}

// 更新商品以及部分更新
func (s *GoodsServer) UpdateGoods(ctx context.Context, req *proto.CreateGoodsInfo) (*proto.Empty, error) {
	var goods model.Goods
	if req.CategoryId == 0 && req.BrandId == 0 {
		if result := global.DB.First(&goods, req.Id); result.RowsAffected == 0 {
			return nil, status.Errorf(codes.InvalidArgument, "商品不存在")
		}
		goods.IsNew = req.IsNew
		goods.IsHot = req.IsHot
		goods.OnSale = req.OnSale
		tx := global.DB.Begin()
		if result := tx.Save(&goods); result.Error != nil {
			tx.Rollback()
			return nil, result.Error
		}
		tx.Commit()
		return &proto.Empty{}, nil
	}
	var category model.Category
	if result := global.DB.First(&category, req.CategoryId); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "商品分类不存在")
	}
	var brand model.Brands
	if result := global.DB.First(&brand, req.BrandId); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "品牌不存在")
	}

	if result := global.DB.First(&goods, req.Id); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "商品不存在")
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
	tx := global.DB.Begin()
	if result := tx.Save(&goods); result.Error != nil {
		tx.Rollback()
		return nil, result.Error
	}
	tx.Commit()
	return &proto.Empty{}, nil
}
