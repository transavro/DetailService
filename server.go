package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/joho/godotenv"
	pbAuth "github.com/transavro/AuthService/proto"
	"github.com/transavro/DetailService/apihandler"
	pb "github.com/transavro/DetailService/proto"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/bson/bsonrw"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"log"
	"net"
	"net/http"
	"os"
	"reflect"
	"sync"
	"time"
)

var mongoDbHost, redisPort, grpcPort, restPort string

type nullawareStrDecoder struct{}

func (nullawareStrDecoder) DecodeValue(dctx bsoncodec.DecodeContext, vr bsonrw.ValueReader, val reflect.Value) error {
	if !val.CanSet() || val.Kind() != reflect.String {
		return errors.New("bad type or not settable")
	}
	var str string
	var err error
	switch vr.Type() {
	case bsontype.String:
		if str, err = vr.ReadString(); err != nil {
			return err
		}
	case bsontype.Null: // THIS IS THE MISSING PIECE TO HANDLE NULL!
		if err = vr.ReadNull(); err != nil {
			return err
		}
	default:
		return fmt.Errorf("cannot decode %v into a string type", vr.Type())
	}
	val.SetString(str)
	return nil
}

func unaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	log.Println("unaryInterceptor")
	err := checkingJWTToken(ctx)
	if err != nil {
		return nil, err
	}
	return handler(ctx, req)
}

func checkingJWTToken(ctx context.Context) error {
	meta, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return status.Error(codes.NotFound, fmt.Sprintf("no auth meta-data found in request"))
	}

	token := meta["token"]

	if len(token) == 0 {
		return status.Error(codes.NotFound, fmt.Sprintf("Token not found"))
	}

	// calling auth service
	conn, err := grpc.Dial(":7757", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	// Auth here
	authClient := pbAuth.NewAuthServiceClient(conn)
	_, err = authClient.ValidateToken(context.Background(), &pbAuth.Token{
		Token: token[0],
	})
	if err != nil {
		return status.Error(codes.NotFound, fmt.Sprintf("Invalid token:  %s ", err))
	} else {
		return nil
	}
}

// streamAuthIntercept intercepts to validate authorization
func streamIntercept(server interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler, ) error {
	err := checkingJWTToken(stream.Context())
	if err != nil {
		return err
	}
	return handler(server, stream)
}

func startGRPCServer(address string, server apihandler.Executor) error {
	// create a listener on TCP port
	lis, err := net.Listen("tcp", address)
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	} // create a server instance
	if err != nil {
		return err
	}

	//serverOptions := []grpc.ServerOption{grpc.UnaryInterceptor(unaryInterceptor), grpc.StreamInterceptor(streamIntercept)}
	// attach the Ping service to the server
	//grpcServer := grpc.NewServer(serverOptions...)
	grpcServer := grpc.NewServer()
	// attach the Ping service to the server
	pb.RegisterDetailPageServiceServer(grpcServer, &server) // start the server
	//log.Printf("starting HTTP/2 gRPC server on %s", address)
	if err := grpcServer.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve: %s", err)
	}
	return nil
}

func startRESTServer(address, grpcAddress string) error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	mux := runtime.NewServeMux(runtime.WithIncomingHeaderMatcher(runtime.DefaultHeaderMatcher), runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{OrigName: false, EnumsAsInts: true, EmitDefaults: true}))

	opts := []grpc.DialOption{grpc.WithInsecure()}

	// Register ping
	err := pb.RegisterDetailPageServiceHandlerFromEndpoint(ctx, mux, grpcAddress, opts)
	if err != nil {
		return fmt.Errorf("could not register service Ping: %s", err)
	}

	log.Printf("starting HTTP/1.1 REST server on %s", address)
	http.ListenAndServe(address, mux)
	return nil
}

func getMongoCollection(dbName, collectionName, mongoHost string) *mongo.Collection {
	// Register custom codecs for protobuf Timestamp and wrapper types
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoHost), options.Client().SetRegistry(bson.NewRegistryBuilder().
		RegisterDecoder(reflect.TypeOf(""), nullawareStrDecoder{}).
		Build(), ))

	if err != nil {
		log.Println("Error while making collection obj ")
		log.Fatal(err)
	}
	return mongoClient.Database(dbName).Collection(collectionName)
}

func main() {
	serverhandler := initializeProcess()

	// fire the gRPC server in a goroutine
	go func() {
		err := startGRPCServer(grpcPort, serverhandler)
		if err != nil {
			log.Fatalf("failed to start gRPC server: %s", err)
		}
	}()

	// fire the REST server in a goroutine
	go func() {
		err := startRESTServer(restPort, grpcPort)
		if err != nil {
			log.Fatalf("failed to start gRPC server: %s", err)
		}
	}()
	// infinite loop
	select {}
}

func loadEnv() {
	mongoDbHost = os.Getenv("MONGO_HOST")
	redisPort = os.Getenv("REDIS_PORT")
	grpcPort = os.Getenv("GRPC_PORT")
	restPort = os.Getenv("REST_PORT")
}

func initializeProcess() apihandler.Executor {
	err := godotenv.Load()
	if err != nil {
		log.Println(err.Error())
	}
	loadEnv()

	tileCollection := getMongoCollection("transavro", "optimus_content", mongoDbHost)
	return apihandler.Executor{
		tileCollection,
		getRedisClient(redisPort),
			new(sync.WaitGroup),
	}
}

func getRedisClient(redisHost string) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     redisHost,
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	_, err := client.Ping().Result()
	if err != nil {
		log.Fatalf("Could not connect to redis %v", err)
	}
	return client
}
