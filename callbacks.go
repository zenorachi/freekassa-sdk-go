package freekassa

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ICallbacks interface {
	Listen(whitelist map[string]struct{}) error
}

type callbacks struct {
	Host           string
	Port           string
	TrustedProxies []string
	Confirmation   func() (string, func(c *gin.Context))
	Success        func() (string, func(c *gin.Context))
	Failure        func() (string, func(c *gin.Context))
}

func NewCallbacks(
	host string,
	port string,
	trustedProxies []string,
	confirmation func() (string, func(c *gin.Context)),
	success func() (string, func(c *gin.Context)),
	failure func() (string, func(c *gin.Context)),
) ICallbacks {
	return &callbacks{
		Host:           host,
		Port:           port,
		TrustedProxies: trustedProxies,
		Confirmation:   confirmation,
		Success:        success,
		Failure:        failure,
	}
}

func middleware(whitelist map[string]struct{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		if _, found := whitelist[c.ClientIP()]; !found {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"message": "access denied"})

			return
		}

		c.Next()
	}
}

func (c *callbacks) Listen(whitelist map[string]struct{}) error {
	router := gin.Default()

	if len(c.TrustedProxies) > 0 {
		if err := router.SetTrustedProxies(c.TrustedProxies); err != nil {
			return err
		}
	}

	router.Use(middleware(whitelist))

	confirmPath, confirmFunc := c.Confirmation()
	router.POST(confirmPath, confirmFunc)

	successPath, successFunc := c.Success()
	router.GET(successPath, successFunc)

	failurePath, failureFunc := c.Failure()
	router.GET(failurePath, failureFunc)

	if err := router.Run(c.address()); err != nil {
		return err
	}

	return nil
}

func (c *callbacks) address() string {
	return fmt.Sprintf("%s:%s", c.Host, c.Port)
}
