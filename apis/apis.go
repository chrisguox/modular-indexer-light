package apis

import (
	"fmt"
	"net/http"
	"time"

	"github.com/RiemaLabs/modular-indexer-light/constant"
	"github.com/RiemaLabs/modular-indexer-light/runtime"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func StartService(df *runtime.RuntimeState, enableDebug bool, port int) {

	if !enableDebug {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	// TODO: Medium. Add the TRUSTED_PROXIES to our config
	// trustedProxies := os.Getenv("TRUSTED_PROXIES")
	// if trustedProxies != "" {
	//     r.SetTrustedProxies([]string{trustedProxies})
	// }

	r.Use(gin.Recovery(), gin.Logger(), cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"POST", "GET"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.GET(constant.LightState, func(c *gin.Context) {
		c.JSON(http.StatusOK, Brc20VerifiableLightStateResponse{
			State: constant.ApiState.String(),
		})
	})
	serv := r.Group("v1")
	{
		serv.Use(CheckState())
		serv.GET(constant.LightBlockHeight, func(c *gin.Context) {
			c.Data(http.StatusOK, "text/plain", []byte(fmt.Sprintf("%d", df.CurrentHeight())))
		})

		serv.GET(constant.LightCurrentBalanceOfWallet, func(c *gin.Context) {
			ck := df.CurrentFirstCheckpoint().Checkpoint

			GetCurrentBalanceOfWallet(c, ck)
		})

		serv.GET(constant.LightCurrentBalanceOfPkscript, func(c *gin.Context) {
			ck := df.CurrentFirstCheckpoint().Checkpoint

			GetCurrentBalanceOfPkscript(c, ck)
		})

		serv.GET(constant.LightCurrentCheckpoints, func(c *gin.Context) {
			cur := df.CurrentCheckpoints()
			c.JSON(http.StatusOK, cur)
		})

		serv.GET(constant.LightLastCheckpoint, func(c *gin.Context) {
			lt := df.LastCheckpoint()
			c.JSON(http.StatusOK, lt)
		})
	}

	r.Run(fmt.Sprintf(":%d", port))
}
