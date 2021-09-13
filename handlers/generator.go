package handlers

import (
	"chameleon/generator"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// CreateGenerator godoc
// @Summary Create a generator.
// @Description create a generator.
// @Tags generator
// @Accept json
// @Produce json
// @Param config body generator.GeneratorConfig true "generator configuration"
// @Success 201 {object} generator.GeneratorConfig
// @Failure 409
// @Router /generators [post]
func CreateGenerator(c *gin.Context) {
	config := generator.GeneratorConfig{}

	if c.ShouldBind(&config) == nil {
		log.Printf("%v", &config)
		gm := c.MustGet("gm").(*generator.GeneratorManager)
		err := gm.CreateGenerator(&config)
		if err != nil {
			c.Status(http.StatusConflict)
		} else {
			c.JSON(http.StatusCreated, config)
		}
	}
}

// ListGenerator godoc
// @Summary list all generators.
// @Description list all generators.
// @Tags generator
// @Accept json
// @Produce json
// @Success 200 {array} string
// @Router /generators [get]
func ListGenerator(c *gin.Context) {
	gm := c.MustGet("gm").(*generator.GeneratorManager)
	c.JSON(http.StatusOK, gm.ListGenerator())
}

// GetGenerator godoc
// @Summary get generator by name.
// @Description get generator by name.
// @Tags generator
// @Accept json
// @Produce json
// @Param name path string true "configuration name"
// @Success 200 {object} generator.GeneratorConfig
// @Failure 409
// @Router /generators/{name} [get]
func GetGenerator(c *gin.Context) {
	name := c.Param("name")
	gm := c.MustGet("gm").(*generator.GeneratorManager)

	_, exists := gm.Manager[name]
	if exists {
		res := gm.Manager[name].Config
		c.JSON(http.StatusOK, res)
	} else {
		c.Status(http.StatusNotFound)
	}
}

// DeleteGenerator godoc
// @Summary delete generator by name.
// @Description delete generator by name.
// @Tags generator
// @Accept json
// @Produce json
// @Param name path string true "configuration name"
// @Success 204
// @Failure 404
// @Router /generators/{name} [delete]
func DeleteGenerator(c *gin.Context) {
	name := c.Param("name")
	gm := c.MustGet("gm").(*generator.GeneratorManager)
	err := gm.DeleteGenerator(name)
	if err != nil {
		c.Status(http.StatusNoContent)
	} else {
		c.Status(http.StatusNotFound)
	}

}

// StartGenerator godoc
// @Summary start to run a generator.
// @Description start to run a generator.
// @Tags generator
// @Accept json
// @Produce json
// @Param name path string true "configuration name"
// @Success 204
// @Failure 404
// @Router /generators/{name}/start [post]
func StartGenerator(c *gin.Context) {
	name := c.Param("name")
	gm := c.MustGet("gm").(*generator.GeneratorManager)
	g, exists := gm.Manager[name]
	if exists {
		// TODO : check generator status
		// if generator stopped, should re-create a instance
		go func() {
			g.Run(1000 * 1000)
		}()
		c.Status(http.StatusNoContent)
	} else {
		c.Status(http.StatusNotFound)
	}
}

// StopGenerator godoc
// @Summary stop a running generator.
// @Description stop a running generator.
// @Tags generator
// @Accept json
// @Produce json
// @Param name path string true "configuration name"
// @Success 204
// @Failure 404
// @Router /generators/{name}/stop [post]
func StopGenerator(c *gin.Context) {
	name := c.Param("name")
	gm := c.MustGet("gm").(*generator.GeneratorManager)
	g, exists := gm.Manager[name]
	if exists {
		go func() {
			g.Stop()
		}()
		c.Status(http.StatusNoContent)
	} else {
		c.Status(http.StatusNotFound)
	}
}

// GeneratorStatus godoc
// @Summary get status of a generator.
// @Description get status of a generator.
// @Tags generator
// @Accept json
// @Produce json
// @Param name path string true "configuration name"
// @Success 200 {object} generator.GeneratorStatus
// @Failure 404
// @Router /generators/{name}/status [post]
func GeneratorStatus(c *gin.Context) {
	name := c.Param("name")
	gm := c.MustGet("gm").(*generator.GeneratorManager)

	status, err := gm.GetGeneratorStatus(name)
	if err == nil {
		c.JSON(http.StatusOK, *status)
	} else {
		c.Status(http.StatusNotFound)
	}
}
