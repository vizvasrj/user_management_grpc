protoc:
        @echo "Generating Go files"
        cd proto && protoc --go_out=. --go-grpc_out=. \
                --go-grpc_opt=paths=source_relative --go_opt=paths=source_relative *.proto
