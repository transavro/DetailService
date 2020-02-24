package apihandler

import (
	"context"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/golang/protobuf/proto"
	pb "github.com/transavro/DetialService/proto"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"math/rand"
	"strings"
	"time"
)



type movieTile struct {
	Posters struct {
		Landscape        []string      `json:"landscape"`
		Portrait         []string      `json:"portrait"`
		Backdrop         []string      `json:"backdrop"`
	} `json:"posters"`
	Content struct {
		Source       string   `json:"source"`
		DetailPage 	 bool	  `json:"detailPage"`
		Package      string   `json:"package"`
		Target       []string `json:"target"`
	} `json:"content"`
	Metadata struct {
		Title           string      `json:"title"`
		Synopsis        string      `json:"synopsis"`
		Rating          float64         `json:"rating"`
		Runtime         string      `json:"runtime"`
		Year            string      `json:"year"`
		Cast            []string    `json:"cast"`
		Directors       []string    `json:"directors"`
		Genre           []string    `json:"genre"`
		Categories      []string    `json:"categories"`
		Languages       []string    `json:"languages"`
	} `json:"metadata"`
	Buttons []struct {
		Title       string `json:"title"`
		Icon        string `json:"icon"`
		PackageName string `json:"packageName"`
		Action      string `json:"action"`
	} `json:"buttons"`
}


type Server struct {
	TileCollection *mongo.Collection
	RedisConnection *redis.Client
}

func(s *Server) GetDetailInfo(ctx context.Context, tileInfo *pb.TileInfoRequest) (*pb.DetailTileInfo, error){

	log.Println("Detail hit ")
	detailKey := fmt.Sprintf("%s:detail", tileInfo.GetTileId())

	var detailTileInfo pb.DetailTileInfo


	if s.RedisConnection.Exists(detailKey).Val() == 1 {
		log.Println("From cache")
		result, err := s.RedisConnection.SMembers(detailKey).Result()
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
		for _, data := range result {
			//Important lesson, challenge was to convert interface{} to byte. used  ([]byte(k.(string)))
			err = proto.Unmarshal([]byte(data), &detailTileInfo)
			if err != nil {
				return nil, status.Error(codes.Internal, fmt.Sprintf("Failed to unMarshal result from cache ", err))
			}
			break
		}
		return &detailTileInfo, nil
	}else {
		log.Println("From DB ")
		findFilter := bson.D{{"ref_id", tileInfo.TileId}}
		findResult := s.TileCollection.FindOne(ctx, findFilter)

		var tile movieTile
		err := findResult.Decode(&tile)
		if err != nil {
			return nil, err
		}

		//making Tile detail Info
		detailTileInfo.Title = tile.Metadata.Title
		if len(tile.Posters.Backdrop) > 0 {
			detailTileInfo.BackDrop = tile.Posters.Backdrop[0]
		}
		if len(tile.Posters.Portrait) > 0 {
			detailTileInfo.Portrait = tile.Posters.Portrait[0]
		}
		metaSet := make(map[string]string)

		// creating pipes for mongo aggregation
		myStages := mongo.Pipeline{}


		myStages = append(myStages,bson.D{{"$match", bson.D{{"ref_id", bson.D{{"$ne", tileInfo.GetTileId()}}}}}},)


		// fetching related content
		//myStages = append(myStages , bson.D{{"$match", bson.D{{"content.publishState", true}}}})

		if len(tile.Metadata.Categories) > 0 {
			myStages = append(myStages, bson.D{{"$match", bson.D{{"metadata.categories", bson.D{{"$in", tile.Metadata.Categories}}}}}})
		}

		if len(tile.Metadata.Genre) > 0 {
			metaSet["genre"] = strings.Join(tile.Metadata.Genre, ",")
			myStages = append(myStages, bson.D{{"$match", bson.D{{"metadata.genre", bson.D{{"$in", tile.Metadata.Genre}}}}}})
		}

		if len(tile.Metadata.Languages) > 0 {
			metaSet["language"] = strings.Join(tile.Metadata.Languages, ",")
			myStages = append(myStages, bson.D{{"$match", bson.D{{"metadata.languages", bson.D{{"$in", tile.Metadata.Languages}}}}}})
		}


		if len(tile.Metadata.Runtime) > 0  {
			metaSet["runtime"] = tile.Metadata.Runtime
		}
		if len(tile.Metadata.Year) > 0 {
			metaSet["year"] = tile.Metadata.Year
		}
		if tile.Metadata.Rating != 0 {
			metaSet["rating"] = fmt.Sprint(tile.Metadata.Rating)
		}
		if len(tile.Content.Source) > 0 {
			metaSet["source"] = tile.Content.Source
		}

		if(len(tile.Metadata.Cast) > 0){
			metaSet["cast"] = strings.Join(tile.Metadata.Cast, ",")
		}

		if(len(tile.Metadata.Directors) > 0){
			metaSet["director"] = strings.Join(tile.Metadata.Directors, ",")
		}

		detailTileInfo.Metadata = metaSet
		detailTileInfo.Synopsis = tile.Metadata.Synopsis
		detailTileInfo.Target = tile.Content.Target

		var buttons []*pb.Button
		for _, v := range tile.Buttons {
			var button pb.Button
			button.Title = v.Title
			button.PackageName = v.PackageName
			button.Action = v.Action
			button.Icon = v.Icon
			buttons = append(buttons, &button)
		}
		detailTileInfo.Button = buttons

		myStages = append(myStages,bson.D{{"$sort", bson.D{{"created_at", -1}, {"updated_at", -1}, {"metadata.year", -1}}}},)

		myStages = append(myStages,bson.D{{"$limit", 15}})

		myStages = append(myStages, bson.D{{"$project", bson.D{
			{"_id", 0},
			{"ref_id", 1},
			{"metadata.title", 1},
			{"posters.landscape", 1},
			{"posters.portrait", 1},
			{"content.package", 1},
			{"content.source", 1},
			{"content.target", 1},
			{"created_at", 1},
			{"content.detailPage", 1},
			{"metadata.releaseDate", 1}}}} )

		// creating aggregation query
		cur, err := s.TileCollection.Aggregate(ctx, myStages)

		if err != nil {
			return nil, err
		}

		var relatedTiles []*pb.ContentTile
		defer cur.Close(ctx)
		for cur.Next(ctx) {
			var mTile movieTile
			// converting curors to movieTiles
			err := cur.Decode(&mTile)
			if err != nil {
				return nil,err
			}

			var contentTile pb.ContentTile

			contentTile.Title = mTile.Metadata.Title
			contentTile.IsDetailPage = mTile.Content.DetailPage
			if len(mTile.Posters.Portrait) > 0 {
				contentTile.Portrait = mTile.Posters.Portrait[0]
			}

			if len(mTile.Posters.Landscape) > 0 {
				contentTile.Poster = mTile.Posters.Landscape[0]
			}

			contentTile.ContentId = cur.Current.Lookup("ref_id").StringValue()
			contentTile.Target = mTile.Content.Target
			contentTile.TileType = pb.TileType_ImageTile
			contentTile.PackageName = mTile.Content.Package
			relatedTiles = append(relatedTiles, &contentTile)
		}

		rand.Seed(time.Now().UnixNano())
		rand.Shuffle(len(relatedTiles), func(i, j int) { relatedTiles[i], relatedTiles[j] = relatedTiles[j], relatedTiles[i] })
		detailTileInfo.ContentTile = relatedTiles

		contentByte, err := proto.Marshal(&detailTileInfo)
		if err != nil {
			return nil, err
		}
		if err = s.RedisConnection.SAdd(detailKey, contentByte).Err(); err != nil {
			return nil, err
		}
		return &detailTileInfo, nil
	}
}