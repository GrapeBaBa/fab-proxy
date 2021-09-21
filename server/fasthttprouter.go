package server

import (
	"fab-proxy/pkg"
	"fab-proxy/pkg/chain"
	"fmt"
	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
	"strconv"
)

func NewFastHttpServer() *fasthttp.Server {
	r := router.New()
	rg := r.Group("/api")
	rg.GET("/test", func(c *fasthttp.RequestCtx) {
		//time.Sleep(5 * time.Second)
		_, _ = c.WriteString("Welcome Gin Server")
	})
	rg.GET("/block_height", func(ctx *fasthttp.RequestCtx) {
		height := chain.GetChainInfo().GetBlockHeight()
		resStr, err := height.MarshalJSON()
		if err != nil {
			panic(err.Error())
		}
		ctx.SetContentType("application/json; charset=utf-8")
		_, _ = ctx.WriteString(string(resStr))
	})
	rg.GET("/accepted_total", func(ctx *fasthttp.RequestCtx) {
		count := chain.GetChainInfo().GetAcceptedTxCount()
		resStr, err := count.MarshalJSON()
		if err != nil {
			panic(err.Error())
		}
		ctx.SetContentType("application/json; charset=utf-8")
		_, _ = ctx.WriteString(string(resStr))
	})
	rg.GET("/confirmed_total", func(ctx *fasthttp.RequestCtx) {
		count := chain.GetChainInfo().GetConfirmedTxCount()
		resStr, err := count.MarshalJSON()
		if err != nil {
			panic(err.Error())
		}
		ctx.SetContentType("application/json; charset=utf-8")
		_, _ = ctx.WriteString(string(resStr))
	})
	rg.GET("/block", func(ctx *fasthttp.RequestCtx) {
		block := ctx.QueryArgs().Peek("height")
		blockNum, err := strconv.Atoi(string(block))
		if err != nil {
			panic(err.Error())
		}
		bi := chain.GetChainInfo().GetBlockInfo(uint64(blockNum))
		resStr, err := bi.MarshalJSON()
		if err != nil {
			panic(err.Error())
		}
		ctx.SetContentType("application/json; charset=utf-8")
		_, _ = ctx.WriteString(string(resStr))
	})
	rg.GET("/node_count", func(ctx *fasthttp.RequestCtx) {
		caseNumStr := ctx.QueryArgs().Peek("case")
		caseNum, err := strconv.Atoi(string(caseNumStr))
		if err != nil {
			panic(err.Error())
		}
		var count uint
		switch caseNum {
		case 1, 3, 5, 6, 7, 9:
			count = 4
		default:
			count = 16
		}
		resp := &pkg.NodeCount{Count: count}
		resStr, err := resp.MarshalJSON()
		if err != nil {
			panic(err.Error())
		}
		ctx.SetContentType("application/json; charset=utf-8")
		_, _ = ctx.WriteString(string(resStr))
	})
	rg.GET("/tx", func(ctx *fasthttp.RequestCtx) {
		typeStr := ctx.QueryArgs().Peek("type")
		content := chain.GetChainInfo().CreateTx(string(typeStr))
		ctx.SetContentType("application/json; charset=utf-8")
		_, _ = ctx.WriteString(content)
	})
	rg.POST("/tx/invoke", func(ctx *fasthttp.RequestCtx) {
		content := ctx.PostBody()
		tokenTx := &pkg.InvokeTokenTx{}
		err := tokenTx.UnmarshalJSON(content)
		if err != nil {
			panic(err.Error())
		}

		data := fmt.Sprintf("%s,%s,%s,%s,%s,%d", tokenTx.ContractId, tokenTx.Method, tokenTx.Function, tokenTx.SourceAccount, tokenTx.DestAccount, tokenTx.Amount)
		resp := chain.GetChainInfo().SendTx(data)
		resStr, err := resp.MarshalJSON()
		if err != nil {
			panic(err.Error())
		}
		ctx.SetContentType("application/json; charset=utf-8")
		_, _ = ctx.WriteString(string(resStr))
	})
	rg.POST("/tx/query", func(ctx *fasthttp.RequestCtx) {
		content := ctx.PostBody()
		tokenTx := &pkg.QueryTokenTx{}
		err := tokenTx.UnmarshalJSON(content)
		if err != nil {
			panic(err.Error())
		}

		data := fmt.Sprintf("%s,%s,%s,%s", tokenTx.ContractId, tokenTx.Method, tokenTx.Function, tokenTx.Account)
		resp := chain.GetChainInfo().SendTxQuery(data)
		ctx.SetContentType("application/json; charset=utf-8")
		_, _ = ctx.WriteString(resp)
	})
	rg.POST("/tx/open", func(ctx *fasthttp.RequestCtx) {
		resp := chain.GetChainInfo().InitAccounts()
		resStr, err := resp.MarshalJSON()
		if err != nil {
			panic(err.Error())
		}
		ctx.SetContentType("application/json; charset=utf-8")
		_, _ = ctx.WriteString(string(resStr))
	})
	srv := &fasthttp.Server{Handler: r.Handler}

	return srv
}
