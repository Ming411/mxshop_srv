package handler

import (
	"mx_shop/goods_srv/proto"
	// ."mx_shop/goods_srv/proto"
	// 这种前面加.的方式，可以直接使用GoodsFilterRequest，而不用写proto.GoodsFilterRequest
)

type GoodsServer struct {
	proto.UnimplementedGoodsServer
}

// func (s *GoodsServer) GoodsList(context.Context, *proto.GoodsFilterRequest) (*proto.GoodsListResponse, error)

// // 现在用户提交订单有多个商品，你得批量查询商品的信息吧
// func (s *GoodsServer) BatchGetGoods(context.Context, *proto.BatchGoodsIdInfo) (*proto.GoodsListResponse, error)
// func (s *GoodsServer) CreateGoods(context.Context, *proto.CreateGoodsInfo) (*proto.GoodsInfoResponse, error)
// func (s *GoodsServer) DeleteGoods(context.Context, *proto.DeleteGoodsInfo) (*emptypb.Empty, error)
// func (s *GoodsServer) UpdateGoods(context.Context, *proto.CreateGoodsInfo) (*emptypb.Empty, error)
// func (s *GoodsServer) GetGoodsDetail(context.Context, *proto.GoodInfoRequest) (*proto.GoodsInfoResponse, error)
