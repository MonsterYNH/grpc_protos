mkdir -p gencode/user gencode/helloworld gencode/gateway

protoc --proto_path=. \
    --go_out=gencode/user/ \
    --go-grpc_out=gencode/user/ \
    --grpc-gateway_out=logtostderr=true:gencode \
    --graphql_out=gencode/user \
    --gateway_out=gencode/gateway \
    user/*.proto
protoc --proto_path=. \
    --go_out=gencode/helloworld/ \
    --go-grpc_out=gencode/helloworld/ \
    --grpc-gateway_out=logtostderr=true:gencode \
    --graphql_out=gencode/helloworld \
    --gateway_out=gencode/gateway \
    helloworld/*.proto
