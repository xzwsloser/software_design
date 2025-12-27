from rpc import rec_sys_pb2_grpc
from rpc import rec_sys_pb2
from rec_sys import RecSys
import grpc
from concurrent import futures

GRPC_SERVER_PORT = 7777

class RecSysRpcServer(rec_sys_pb2_grpc.RecSysServiceServicer):
    rec_sys: RecSys
    def __init__(self):
        self.rec_sys = RecSys()

    def GetRecResult(self, request, context):
        """
        获取推荐系统
        """

        print('='*10 + "Rec Sys Params" + "="*10)
        print(f"userId = {request.userId}")
        print(f"addressId = {request.addressId}")
        print(f"touristTypeId = {request.touristType}")
        print(f"priceSensitive = {request.priceSensitive}")
        print(f"likeType = {request.likeType}")
        print(f"targetsType = {request.targetType}")
        print(f"attentionType = {request.attentionType}")
        print(f"update = {request.update}")
        print(f"limit = {request.limit}")
        print("="*34)

        rec_site_idxs, _ = self.rec_sys.recommand_for_current_user(
            user_id=request.userId,
            address_id=request.addressId,
            tourist_type_id=request.touristType,
            price_sensitive=request.priceSensitive,
            like_type=request.likeType,
            targets_type=request.targetType,
            attention_type=request.attentionType,
            update=request.update,
            limit=request.limit
        )

        return rec_sys_pb2.GetRecResultResp(siteIdxList=rec_site_idxs)

def serve():
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    rec_sys_pb2_grpc.add_RecSysServiceServicer_to_server(RecSysRpcServer(), server)
    server.add_insecure_port(f"[::]:{GRPC_SERVER_PORT}")
    print(f"rec sys server listen on port {GRPC_SERVER_PORT}")
    server.start()
    server.wait_for_termination()

if __name__ == '__main__':
    serve()



