package ollama

import (
	"context"
	"fmt"

	"github.com/PirateDreamer/going/ginx"
	"github.com/gin-gonic/gin"
	"github.com/tmc/langchaingo/llms/ollama"
)

func InitController() {
	g := ginx.R.Group("/ollama")
	g.POST("/get_resp", ginx.Run(GetResp))
}

type GatRespReq struct {
	Question string
}

type GatRespRes struct {
	Reply string `json:"reply"`
}

var prompt = "%s"

func GetResp(ctx context.Context, c *gin.Context, param GatRespReq) (resp *GatRespRes, err error) {
	param.Question = fmt.Sprintf(prompt, param.Question)
	var llm *ollama.LLM
	llm, err = ollama.New(ollama.WithModel("qwen"))
	if err != nil {
		return
	}
	var ollResp string
	ollResp, err = llm.Call(ctx, param.Question)
	if err != nil {
		return
	}
	resp = &GatRespRes{Reply: ollResp}

	return
}
