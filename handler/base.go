/**
 *
 * @author nghiatc
 * @since Jan 5, 2021
 */

package handler

import "github.com/valyala/fasthttp"

// DataResp is struct data response
type DataResp struct {
	Err  int         `json:"err"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

func printJSON(ctx *fasthttp.RequestCtx, json string) {
	ctx.Response.Header.Set("content-type", "application/json;charset=UTF-8")
	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetBodyString(json)
}
