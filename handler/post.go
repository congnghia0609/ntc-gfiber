/**
 *
 * @author nghiatc
 * @since Jan 5, 2021
 */

package handler

import (
	"encoding/json"
	"fmt"
	"github.com/congnghia0609/ntc-gfiber/mdb"
	"github.com/congnghia0609/ntc-gfiber/post"
	"log"
	"strconv"
	"time"

	"github.com/congnghia0609/ntc-gconf/nconf"
	"github.com/gofiber/fiber/v2"
)

// AddPost api add post
func AddPost(c *fiber.Ctx) error {
	params := make(map[string]interface{})
	err := json.Unmarshal([]byte(c.Body()), &params)
	if err != nil {
		log.Println(err)
	}
	fmt.Println("param:", params)
	title := ""
	if params["title"] != nil {
		title = params["title"].(string)
	}
	body := ""
	if params["body"] != nil {
		body = params["body"].(string)
	}
	fmt.Println("title:", title)
	fmt.Println("body:", body)
	// Validate params
	if len(title) == 0 || len(body) == 0 {
		dataResp := DataResp{Err: -1, Msg: "Parameters invalid"}
		return c.JSON(dataResp)
	}

	id, _ := mdb.Next(post.TablePost)
	p := post.Post{
		ID:        id,
		Title:     title,
		Body:      body,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	err1 := post.InsertPost(p)
	if err1 != nil {
		fmt.Println("err1:", err1)
		dataResp := DataResp{Err: -1, Msg: "Add post fail"}
		return c.JSON(dataResp)
	}

	dataResp := DataResp{Err: 0, Msg: "Add post successfully", Data: p}
	return c.JSON(dataResp)
}

// UpdatePost api update post
func UpdatePost(c *fiber.Ctx) error {
	params := make(map[string]interface{})
	err := json.Unmarshal([]byte(c.Body()), &params)
	if err != nil {
		log.Println(err)
	}
	fmt.Println("param:", params)
	var id int64 = 0
	if params["id"] != nil {
		id = int64(params["id"].(float64))
	}
	title := ""
	if params["title"] != nil {
		title = params["title"].(string)
	}
	body := ""
	if params["body"] != nil {
		body = params["body"].(string)
	}
	fmt.Println("id:", id)
	fmt.Println("title:", title)
	fmt.Println("body:", body)
	// Validate params
	if id == 0 || len(title) == 0 || len(body) == 0 {
		dataResp := DataResp{Err: -1, Msg: "Parameters invalid"}
		return c.JSON(dataResp)
	}

	// id, _ := strconv.ParseInt(sid, 10, 64)
	p := post.GetPost(id)
	if p.ID <= 0 {
		dataResp := DataResp{Err: -1, Msg: "Post is not exist"}
		return c.JSON(dataResp)
	}
	// id, _ := mdb.Next(post.TablePost)
	np := post.Post{
		ID:        id,
		Title:     title,
		Body:      body,
		CreatedAt: p.CreatedAt,
		UpdatedAt: time.Now(),
	}
	count, err1 := post.UpdatePost(np)
	if err1 != nil || count < 1 {
		fmt.Println("err1:", err1)
		dataResp := DataResp{Err: -1, Msg: "Update post fail"}
		return c.JSON(dataResp)
	}

	dataResp := DataResp{Err: 0, Msg: "Update post successfully", Data: np}
	return c.JSON(dataResp)
}

// GetPost api get post
func GetPost(c *fiber.Ctx) error {
	sid := c.Params("id")
	id, _ := strconv.ParseInt(sid, 10, 64)
	p := post.GetPost(id)
	if p.ID <= 0 {
		dataResp := DataResp{Err: -1, Msg: "Post is not exist"}
		return c.JSON(dataResp)
	}

	dataResp := DataResp{Err: 0, Msg: "Get post successfully", Data: p}
	return c.JSON(dataResp)
}

// GetAllPosts api get all post
func GetAllPosts(c *fiber.Ctx) error {
	posts := post.GetAllPost()
	dataResp := DataResp{Err: 0, Msg: "Get all posts successfully", Data: posts}
	return c.JSON(dataResp)
}

// GetPosts api get page post
func GetPosts(c *fiber.Ctx) error {
	cfg := nconf.GetConfig()
	paging := cfg.GetInt64("system.paging")

	mapData := make(map[string]interface{})
	isMore := false
	posts := []post.Post{}
	var page int64 = 1
	pg := c.Query("page")
	if len(pg) > 0 {
		page, _ = strconv.ParseInt(pg, 10, 64)
	}
	log.Println("page:", page)
	mapData["page"] = page

	total, _ := post.GetTotalPost()
	// finish soon.
	if total == 0 {
		mapData["isMore"] = false
		mapData["posts"] = posts
		dataResp := DataResp{Err: 0, Msg: "Get all posts successfully", Data: mapData}
		return c.JSON(dataResp)
	}

	maxPage := (total-1)/paging + 1
	// finish soon.
	if page > maxPage {
		mapData["isMore"] = false
		mapData["posts"] = posts
		dataResp := DataResp{Err: 0, Msg: "Get all posts successfully", Data: mapData}
		return c.JSON(dataResp)
	}
	// Get data paging.
	if page < maxPage {
		isMore = true
	}
	skip := (page - 1) * paging

	posts = post.GetSlidePost(skip, paging)
	mapData["isMore"] = isMore
	mapData["posts"] = posts
	dataResp := DataResp{Err: 0, Msg: "Get all posts successfully", Data: mapData}
	return c.JSON(dataResp)
}

// DeletePost api get post
func DeletePost(c *fiber.Ctx) error {
	sid := c.Params("id")
	id, _ := strconv.ParseInt(sid, 10, 64)
	count, err := post.DeletePost(id)
	if err != nil || count < 1 {
		dataResp := DataResp{Err: -1, Msg: "Delete post fail"}
		return c.JSON(dataResp)
	}

	dataResp := DataResp{Err: 0, Msg: "Delete post successfully"}
	return c.JSON(dataResp)
}
