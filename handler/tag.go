/**
 *
 * @author nghiatc
 * @since Jan 5, 2021
 */

package handler

import (
	"encoding/json"
	"fmt"
	"github.com/congnghia0609/ntc-gfiber/tag"
	"log"
	"strconv"
	"time"

	"github.com/congnghia0609/ntc-gconf/nconf"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// AddTag api add tag
func AddTag(ctx *fiber.Ctx) error {
	params := make(map[string]interface{})
	err := json.Unmarshal([]byte(ctx.Body()), &params)
	if err != nil {
		log.Println(err)
	}
	fmt.Println("param:", params)
	name := ""
	if params["name"] != nil {
		name = params["name"].(string)
	}
	fmt.Println("name:", name)
	// Validate params
	if len(name) == 0 {
		dataResp := DataResp{Err: -1, Msg: "Parameters invalid"}
		return ctx.JSON(dataResp)
	}

	id := primitive.NewObjectID()
	t := tag.Tag{
		ID:        id,
		Name:      name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	err1 := tag.InsertTag(t)
	if err1 != nil {
		fmt.Println("err1:", err1)
		dataResp := DataResp{Err: -1, Msg: "Add tag fail"}
		return ctx.JSON(dataResp)
	}

	dataResp := DataResp{Err: 0, Msg: "Add tag successfully", Data: t}
	return ctx.JSON(dataResp)
}

// UpdateTag api update tag
func UpdateTag(ctx *fiber.Ctx) error {
	params := make(map[string]interface{})
	err := json.Unmarshal([]byte(ctx.Body()), &params)
	if err != nil {
		log.Println(err)
	}
	fmt.Println("param:", params)
	var id string = ""
	if params["id"] != nil {
		id = params["id"].(string)
	}
	name := ""
	if params["name"] != nil {
		name = params["name"].(string)
	}
	fmt.Println("id:", id)
	fmt.Println("name:", name)
	// Validate params
	if len(id) == 0 || len(name) == 0 {
		dataResp := DataResp{Err: -1, Msg: "Parameters invalid"}
		return ctx.JSON(dataResp)
	}

	oid, _ := primitive.ObjectIDFromHex(id)
	t := tag.GetTag(oid)
	if t == nil {
		dataResp := DataResp{Err: -1, Msg: "Tag is not exist"}
		return ctx.JSON(dataResp)
	}
	nt := tag.Tag{
		ID:        oid,
		Name:      name,
		CreatedAt: t.CreatedAt,
		UpdatedAt: time.Now(),
	}
	count, err1 := tag.UpdateTag(nt)
	if err1 != nil || count < 1 {
		fmt.Println("err1:", err1)
		dataResp := DataResp{Err: -1, Msg: "Update tag fail"}
		return ctx.JSON(dataResp)
	}

	dataResp := DataResp{Err: 0, Msg: "Update tag successfully", Data: nt}
	return ctx.JSON(dataResp)
}

// GetTag api get tag
func GetTag(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		dataResp := DataResp{Err: -1, Msg: "TagId invalid"}
		return ctx.JSON(dataResp)
	}
	t := tag.GetTag(oid)
	if t == nil {
		dataResp := DataResp{Err: -1, Msg: "Tag is not exist"}
		return ctx.JSON(dataResp)
	}

	dataResp := DataResp{Err: 0, Msg: "Get tag successfully", Data: t}
	return ctx.JSON(dataResp)
}

// GetAllTags api get all tag
func GetAllTags(ctx *fiber.Ctx) error {
	tags := tag.GetAllTag()
	dataResp := DataResp{Err: 0, Msg: "Get all tags successfully", Data: tags}
	return ctx.JSON(dataResp)
}

// GetTags api get page tag
func GetTags(ctx *fiber.Ctx) error {
	c := nconf.GetConfig()
	paging := c.GetInt64("system.paging")

	mapData := make(map[string]interface{})
	isMore := false
	tags := []tag.Tag{}
	var page int64 = 1
	pg := ctx.Query("page")
	if len(pg) > 0 {
		page, _ = strconv.ParseInt(pg, 10, 64)
	}
	log.Println("page:", page)
	mapData["page"] = page

	total, _ := tag.GetTotalTag()
	// finish soon.
	if total == 0 {
		mapData["isMore"] = false
		mapData["tags"] = tags
		dataResp := DataResp{Err: 0, Msg: "Get all tags successfully", Data: mapData}
		return ctx.JSON(dataResp)
	}

	maxPage := (total-1)/paging + 1
	// finish soon.
	if page > maxPage {
		mapData["isMore"] = false
		mapData["tags"] = tags
		dataResp := DataResp{Err: 0, Msg: "Get all tags successfully", Data: mapData}
		return ctx.JSON(dataResp)
	}
	// Get data paging.
	if page < maxPage {
		isMore = true
	}
	skip := (page - 1) * paging

	tags = tag.GetSlideTag(skip, paging)
	mapData["isMore"] = isMore
	mapData["tags"] = tags
	dataResp := DataResp{Err: 0, Msg: "Get all tags successfully", Data: mapData}
	return ctx.JSON(dataResp)
}

// DeleteTag api get tag
func DeleteTag(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		dataResp := DataResp{Err: -1, Msg: "TagId invalid"}
		return ctx.JSON(dataResp)
	}
	count, err := tag.DeleteTag(oid)
	if err != nil || count < 1 {
		dataResp := DataResp{Err: -1, Msg: "Delete tag fail"}
		return ctx.JSON(dataResp)
	}

	dataResp := DataResp{Err: 0, Msg: "Delete tag successfully"}
	return ctx.JSON(dataResp)
}
