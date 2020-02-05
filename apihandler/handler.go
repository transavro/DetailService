package apihandler

import (
	"context"
	"fmt"
	pb "github.com/transavro/DetialService/proto"
	pbSch "github.com/transavro/ScheduleService/proto"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"math/rand"
	"strings"
	"time"
)


type  IntermidiateTile struct {
	Title    string   `json:"title"`
	Synopsis string   `json:"synopsis"`
	Poster   []string `json:"poster"`
	Portriat []string `json:"portriat"`
	Backdrop []string `json:"backdrop"`
	Genre    []string `json:"genre"`
	Languages    []string `json:"languages"`
	Categories    []string `json:"categories"`
	Runtime  string   `json:"runtime"`
	Year     int32  `json:"year"`
	Rating 	float64`json:"rating"`
	Cast      []string `json:"cast"`
	Sources   []string `json:"sources"`
	Directors []string `json:"directors"`
	API       []struct {
		Title   string `json:"title"`
		Icon    string `json:"icon"`
		Action  string `json:"action"`
		Package string `json:"package"`
	} `json:"api"`
}


type Server struct {
	TileCollection *mongo.Collection
}

func(s *Server) GetDetailInfo(ctx context.Context, tileInfo *pb.TileInfoRequest) (*pb.DetailTileInfo, error){

	log.Println("hit detail ")

	cur, err := s.TileCollection.Aggregate(ctx, makingDeliveryPipeline(tileInfo.GetTileId()))
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("Failed to fetch data from DB:  %s ",err.Error()))
	}

	var temp IntermidiateTile
	for cur.Next(ctx){
		err = cur.Decode(&temp)
		if err != nil {
			return nil, status.Error(codes.Internal, fmt.Sprintf("Error while decoding data from DB:  %s ",err.Error()))
		}
		break
	}

	var detailTileInfo pb.DetailTileInfo
	detailTileInfo.Title = temp.Title
	detailTileInfo.Synopsis = temp.Synopsis
	detailTileInfo.Portrait = temp.Portriat
	detailTileInfo.Poster  = temp.Poster
	detailTileInfo.BackDrop = temp.Backdrop

	var buttons []*pb.Button

	if len(temp.API) > 0 {
		for _, v := range temp.API {
			var btn pb.Button
			btn.Title = v.Title
			btn.Action = v.Action
			btn.Icon = v.Icon
			btn.Package = v.Package
			btn.ButtonType = pb.ButtonType_LocalApi
			buttons = append(buttons, &btn)
		}
		detailTileInfo.Button = buttons
	}

	cur, err = s.TileCollection.Aggregate(context.Background(), makeSuggestionPipeLine(temp, tileInfo.GetTileId()))
	if err != nil {
		return nil, status.Errorf(codes.NotFound, fmt.Sprintf("Error while getting Detial Recommended data from DB: %s ", err))
	}

	var relatedTiles []*pbSch.Content
	for cur.Next(context.Background()) {
		var content pbSch.Content
		err = cur.Decode(&content)
		if err != nil {
			return nil, status.Errorf(codes.Internal, fmt.Sprintf("Error while decoding data: %s ", err))
		}
		relatedTiles = append(relatedTiles, &content)
	}

	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(relatedTiles), func(i, j int) { relatedTiles[i], relatedTiles[j] = relatedTiles[j], relatedTiles[i] })
	detailTileInfo.ContentTile = relatedTiles

	metaSet := make(map[string]string)
	if len(temp.Genre) > 0 {
		metaSet["genre"] = strings.Join(temp.Genre, ",")
	}

	if len(temp.Languages) > 0 {
		metaSet["language"] = strings.Join(temp.Languages, ",")
	}

	if len(temp.Runtime) != 0  {
		metaSet["runtime"] = fmt.Sprint(temp.Runtime)
	}
	if temp.Year > 0 {
		metaSet["year"] = fmt.Sprint(temp.Year)
	}
	if temp.Rating != 0 {
		metaSet["rating"] = fmt.Sprint(temp.Rating)
	}
	if len(temp.Sources) > 0 {
		metaSet["source"] = strings.Join(temp.Sources, ",")
	}

	if len(temp.Cast) > 0 {
		metaSet["cast"] = strings.Join(temp.Cast, ",")
	}

	if len(temp.Directors) > 0 {
		metaSet["director"] = strings.Join(temp.Directors, ",")
	}
	detailTileInfo.Metadata = metaSet
	return &detailTileInfo, nil
}



func makingDeliveryPipeline(targetTileId string) mongo.Pipeline {
	pipeline := mongo.Pipeline{}
	pipeline = append(pipeline, bson.D{{"$match", bson.D{{"ref_id", targetTileId}}}})
	pipeline = append(pipeline,  bson.D{{"$lookup", bson.M{"from": "buttons", "localField": "ref_id", "foreignField": "ref_id", "as": "api"}}})
	pipeline = append(pipeline,  bson.D{{"$unwind", "$api"}})
	pipeline = append(pipeline, bson.D{{"$project", bson.D{
		{"_id", 0},
		{"title", "$metadata.title"},
		{"synopsis", "$metadata.synopsis"},
		{"poster", "$media.landscape"},
		{"portriat", "$media.portrait"},
		{"backdrop", "$media.backdrop"},
		{"genre", "$metadata.genre"},
		{"languages" , "$metdata.languages"},
		{"categories" , "$metdata.categories"},
		{"runtime", "$metadata.runtime"},
		{"year" ,"$metadata.year"},
		{"rating", "$metadata.rating"},
		{"cast" ,"$metadata.cast"},
		{"sources", "$content.sources"},
		{"directors", "$metadata.directors"},
		{"play", "$play.contentAvailable"},
		{"api", "$api.buttons"},
	}}})
	return pipeline
}

func makeSuggestionPipeLine(temp IntermidiateTile, targetId string) mongo.Pipeline {
	// creating pipes for mongo aggregation for recommedation
	myStages := mongo.Pipeline{}
	myStages = append(myStages, bson.D{{"$lookup", bson.M{"from": "monetizes", "localField": "ref_id", "foreignField": "ref_id", "as": "play"}}})
	myStages = append(myStages,  bson.D{{"$unwind", "$play"}})
	myStages = append(myStages,bson.D{{"$match", bson.D{{"ref_id", bson.D{{"$ne", targetId}}}}}},)
	myStages = append(myStages , bson.D{{"$match", bson.D{{"content.publishState", true}}}})

	if len(temp.Categories) > 0 {
		myStages = append(myStages, bson.D{{"$match", bson.D{{"metadata.categories", bson.D{{"$in", temp.Categories}}}}}})
	}

	if len(temp.Genre) > 0 {
		myStages = append(myStages, bson.D{{"$match", bson.D{{"metadata.genre", bson.D{{"$in", temp.Genre}}}}}})
	}

	if len(temp.Languages) > 0 {
		myStages = append(myStages, bson.D{{"$match", bson.D{{"metadata.languages", bson.D{{"$in", temp.Languages}}}}}})
	}

	myStages = append(myStages, bson.D{{"$group", bson.D{{"_id", bson.D{
		{"releaseDate", "$metadata.releaseDate"},
		{"year", "$metadata.year"},

	}}, {"contentTile", bson.D{{"$push", bson.D{
		{"title", "$metadata.title"},
		{"portrait", "$media.portrait",},
		{"poster", "$media.landscape"},
		{"video", "$media.video"},
		{"contentId", "$ref_id"},
		{"isDetailPage", "$content.detailPage"},
		{"type", "$tileType"},
		{"play", "$play.contentAvailable"},
	}}}}}}})
	myStages = append(myStages, bson.D{{"$unwind", "$contentTile"}})
	myStages = append(myStages,bson.D{{"$limit", 15}})
	myStages = append(myStages, bson.D{{"$project", bson.D{
		{"_id", 0},
		{"title", "$contentTile.title"},
		{"poster", "$contentTile.poster"},
		{"portriat", "$contentTile.portrait"},
		{"type", "$contentTile.type"},
		{"isDetailPage", "$contentTile.isDetailPage"},
		{"contentId", "$contentTile.contentId"},
		{"play", "$contentTile.play"},
		{"video", "$contentTile.video"},
	}}})

	return myStages
}