package apihandler

import (
	"fmt"
	"github.com/go-redis/redis"
	"github.com/golang/protobuf/proto"
	pb "github.com/transavro/DetailService/proto"
	pbSch "github.com/transavro/ScheduleService/proto"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"math/rand"
	"strings"
	"sync"
	"time"
)

type InterResult struct {
	Title       string   `json:"title"`
	Synopsis    string   `json:"synopsis"`
	Languages   []string `json:"languages"`
	Genre       []string `json:"genre"`
	Categories  []string `json:"categories"`
	Year        int      `json:"year"`
	Sources     []string `json:"sources"`
	Directors   []string `json:"directors"`
	Cast        []string `json:"cast"`
	Releasedate string   `json:"releasedate"`
	Rating      float64  `json:"rating"`
	Runtime     string   `json:"runtime"`
	Backdrop    []string `json:"backdrop"`
	Portrait    []string `json:"portrait"`
	Poster      []string `json:"poster"`
	Button      []struct {
		Monetize int    `json:"monetize"`
		Targetid string `json:"targetid"`
		Source   string `json:"source"`
		Package  string `json:"package"`
		Type     string `json:"type"`
		Target   string `json:"target"`
	} `json:"button"`
}

type Executor struct {
	*mongo.Collection
	*redis.Client
	*sync.WaitGroup
}

func (e *Executor) GetDetailInfo(ctx context.Context, req *pb.TileInfoRequest) (*pb.DetailTileInfo, error) {
	result := new(pb.DetailTileInfo)

	targetRedisKey := fmt.Sprintf("%s:detail", req.GetTileId())
	//if e.Exists(targetRedisKey).Val() == 1 {
	//	redis_result, err := e.SMembers(targetRedisKey).Result()
	//	if err != nil {
	//		return nil, status.Error(codes.Unavailable, fmt.Sprintf("Failed to get Result from Cache ", err))
	//	}
	//	//Important lesson, challenge was to convert interface{} to byte. used  ([]byte(k.(string)))
	//	err = proto.Unmarshal([]byte(redis_result[0]), result)
	//	if err != nil {
	//		return nil, status.Error(codes.Internal, fmt.Sprintf("Failed to unMarshal result ", err))
	//	}
	//	return result, nil
	//} else

	{


		resultCur, err := e.Collection.Aggregate(ctx, makePL(req.GetTileId()))
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
		defer resultCur.Close(ctx)
		interResult := new(InterResult)

		for resultCur.Next(ctx) {
			err = resultCur.Decode(&interResult)
			if err != nil {
				return nil, err
			}
			e.Add(2)
			go makeSuggestion(result, interResult, e.Collection, req.GetTileId(), ctx, e.WaitGroup)
			go makeDetailInfo(result, interResult, e.WaitGroup)
			e.Wait()
			break
		}

		resultByteArray, err := proto.Marshal(result)
		if err != nil {
			return result, err
		}
		if _, err = e.SAdd(targetRedisKey, resultByteArray).Result(); err != nil {
			return result, err
		}
	}
	return result, nil
}

func makeDetailInfo(result *pb.DetailTileInfo, interResult *InterResult, wg *sync.WaitGroup) {
	//making general info
	var button pb.Button
	result.Title = interResult.Title
	result.Portrait = interResult.Portrait
	result.BackDrop = interResult.Backdrop
	result.Poster = interResult.Poster
	result.Synopsis = interResult.Synopsis

	//making metadata map
	metaSet := make(map[string]string)
	if len(interResult.Genre) > 0 {
		metaSet["genre"] = strings.Join(interResult.Genre, ",")
	}
	if len(interResult.Languages) > 0 {
		metaSet["language"] = strings.Join(interResult.Languages, ",")
	}
	if len(interResult.Runtime) != 0 {
		metaSet["runtime"] = fmt.Sprint(interResult.Runtime)
	}
	if interResult.Year > 0 {
		metaSet["year"] = fmt.Sprint(interResult.Year)
	}
	if interResult.Rating != 0 {
		metaSet["rating"] = fmt.Sprint(interResult.Rating)
	}
	if len(interResult.Sources) > 0 {
		metaSet["source"] = strings.Join(interResult.Sources, ",")
	}
	if len(interResult.Cast) > 0 {
		metaSet["cast"] = strings.Join(interResult.Cast, ",")
	}
	if len(interResult.Directors) > 0 {
		metaSet["director"] = strings.Join(interResult.Directors, ",")
	}
	result.Metadata = metaSet

	//making button

	for _, btn := range interResult.Button {
		button.Title = "Play on " + btn.Source
		button.Package = btn.Package
		button.Target = btn.Target
		result.Button = append(result.Button, &button)
	}
	wg.Done()
}

func makeSuggestion(result *pb.DetailTileInfo, interResult *InterResult, collection *mongo.Collection, id string, ctx context.Context, wg *sync.WaitGroup) {
	start := time.Now()
	cur, err := collection.Aggregate(context.Background(), makeSugPL(interResult, id))
	log.Println("pure mongo suggestion timing ==>",time.Since(start))
	if err != nil {
		log.Fatal(err.Error())
	}
	defer cur.Close(ctx)
	var relatedTiles []*pbSch.Content
	content := new(pbSch.Content)
	for cur.Next(context.Background()) {
		err = cur.Decode(content)
		if err != nil {
			log.Fatal(err.Error())
		}
		relatedTiles = append(relatedTiles, content)
	}
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(relatedTiles), func(i, j int) { relatedTiles[i], relatedTiles[j] = relatedTiles[j], relatedTiles[i] })
	result.ContentTile = relatedTiles
	wg.Done()
}

func makePL(id string) mongo.Pipeline {
	var stages mongo.Pipeline
	stages = append(stages, bson.D{{"$match", bson.M{"refid": id}}})                                                                                            //adding stage 1
	stages = append(stages, bson.D{{"$lookup", bson.M{"from": "optimus_monetize", "localField": "refid", "foreignField": "refid", "as": "play"}}})              //adding stage 2
	stages = append(stages, bson.D{{"$replaceRoot", bson.M{"newRoot": bson.M{"$mergeObjects": bson.A{bson.M{"$arrayElemAt": bson.A{"$play", 0}}, "$$ROOT"}}}}}) //adding stage 3  ==> https://docs.mongodb.com/manual/reference/operator/aggregation/mergeObjects/#exp._S_mergeObjects
	stages = append(stages, bson.D{{"$project", bson.M{"play": 0}}})                                                                                            // adding stage 4
	stages = append(stages, bson.D{{"$project", bson.M{"_id": 0,
		"title":       "$metadata.title",
		"synopsis":    "$metadata.synopsis",
		"languages":   "$metadata.languages",
		"genre":       "$metadata.genre",
		"categories":  "$metadata.categories",
		"year":        "$metadata.year",
		"sources":     "$content.sources",
		"directors":   "$metadata.directors",
		"cast":        "$metadata.cast",
		"releasedate": "$metadata.releasedate",
		"rating":      "$metadata.rating",
		"runtime":     "$metadata.runtime",
		"backdrop":    "$media.backdrop",
		"poster":      "$media.landscape",
		"portrait":    "$media.portrait",
		"button":      "$contentavailable"}}}) // adding stage 5

	return stages
}

func makeSugPL(temp *InterResult, targetId string) mongo.Pipeline {
	// creating pipes for mongo aggregation for recommedation
	stages := mongo.Pipeline{}
	if len(temp.Categories) > 0 {
		stages = append(stages, bson.D{{"$match", bson.M{"metadata.categories": bson.M{"$in": temp.Categories}}}})
	}
	if len(temp.Genre) > 0 {
		stages = append(stages, bson.D{{"$match", bson.M{"metadata.genre": bson.M{"$in": temp.Genre}}}})
	}
	if len(temp.Languages) > 0 {
		stages = append(stages, bson.D{{"$match", bson.M{"metadata.languages": bson.M{"$in": temp.Languages}}}})
	}
	stages = append(stages, bson.D{{"$match", bson.M{"$and": bson.A{bson.M{"refid": bson.M{"$ne": targetId}}, bson.M{"content.publishstate": bson.M{"$ne":false}}}}}})
	stages = append(stages, bson.D{{"$limit", 15}})
	stages = append(stages, bson.D{{"$lookup", bson.M{"from": "optimus_monetize", "localField": "refid", "foreignField": "refid", "as": "play"}}})
	stages = append(stages, bson.D{{"$replaceRoot", bson.M{"newRoot": bson.M{"$mergeObjects": bson.A{bson.M{"$arrayElemAt": bson.A{"$play", 0}}, "$$ROOT"}}}}}) //adding stage 3  ==> https://docs.mongodb.com/manual/reference/operator/aggregation/mergeObjects/#exp._S_mergeObjects
	stages = append(stages, bson.D{{"$project", bson.M{"play": 0}}})
	stages = append(stages, bson.D{{"$project", bson.M{"_id": 0,
		"title":        "$metadata.title",
		"poster":       "$media.landscape",
		"portriat":     "$media.portrait",
		"video":        "$media.video",
		"type":         "$tiletype",
		"isDetailPage": "$content.detailpage",
		"contentId":    "$refid",
		"play":         "$contentavailable",
	}}})
	return stages
}
