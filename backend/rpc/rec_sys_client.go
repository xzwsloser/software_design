package rpc

import (
	"context"
	"fmt"
	"sync"
	"github.com/xzwsloser/software_design/backend/rpc/pb"
	"github.com/xzwsloser/software_design/backend/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type RecSysClient struct {
	client 	pb.RecSysServiceClient
	ctx  	context.Context
	conn 	*grpc.ClientConn
	mu 		*sync.Mutex
}

var (
	recSysClient *RecSysClient = nil
)

func NewGrpcClient(addr string, port int) {
	address := fmt.Sprintf("%s:%d", addr,
								 	port)
	
	conn, err := grpc.NewClient(address, 
							grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		utils.GetLogger().Error(err.Error())
		return 
	}

	client := pb.NewRecSysServiceClient(conn)

	recSysClient = &RecSysClient{
		client: client,
		ctx: context.Background(),
		conn: conn,
		mu: &sync.Mutex{},
	}
}

func GetRecSysClient() *RecSysClient {
	return recSysClient
}

func (r *RecSysClient) CloseGrpcClient() error {
	return r.conn.Close()
}

// GetRecResult 获取对应用户推荐结果, 线程安全
func (r *RecSysClient) GetRecResult(userId int,
									addressId int,
									touristType int,
									priceSensitive int,
									likeType []int,
									targetType []int,
									attentionType []int,
									update bool,
									limit int) ([]int, error) {
	
	r.mu.Lock()									
	defer r.mu.Unlock()

	result, err := r.client.GetRecResult(r.ctx, &pb.GetRecResultReq{
		UserId: int32(userId),
		AddressId: int32(addressId),
		TouristType: int32(touristType),
		PriceSensitive: int32(priceSensitive),
		LikeType: from_int32_to_int(likeType),
		TargetType: from_int32_to_int(targetType),
		AttentionType: from_int32_to_int(attentionType),
		Update: update,
		Limit: int32(limit),
	})

	if err != nil {
		utils.GetLogger().Error(err.Error())
		return nil, err
	}

	return from_int_to_int32(result.SiteIdxList), nil
}

func from_int32_to_int(raw []int) []int32 {
	transformed := make([]int32, 0, len(raw))
	for _, element := range raw {
		transformed = append(transformed, int32(element))
	}

	return transformed
}

func from_int_to_int32(raw []int32) []int {
	transformed := make([]int, 0, len(raw))
	for _, element := range raw {
		transformed = append(transformed, int(element))
	}

	return transformed
}



