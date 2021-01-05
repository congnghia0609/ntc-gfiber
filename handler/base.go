/**
 *
 * @author nghiatc
 * @since Jan 5, 2021
 */

package handler

// DataResp is struct data response
type DataResp struct {
	Err  int         `json:"err"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}
