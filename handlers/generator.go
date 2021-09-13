package handlers

import (
	"log"
	"net/http"

	"chameleon/generator"

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
// @Router /generators/{name} [get]
func GetGenerator(c *gin.Context) {
	name := c.Param("name")
	gm := c.MustGet("gm").(*generator.GeneratorManager)
	res := gm.Manager[name].Config
	c.JSON(http.StatusOK, res)
}

// DeleteGenerator godoc
// @Summary delete generator by name.
// @Description delete generator by name.
// @Tags generator
// @Accept json
// @Produce json
// @Param name path string true "configuration name"
// @Success 204
// @Router /generators/{name} [delete]
func DeleteGenerator(c *gin.Context) {
	name := c.Param("name")
	gm := c.MustGet("gm").(*generator.GeneratorManager)
	gm.DeleteGenerator(name)
	c.Status(http.StatusNoContent)
}

// StartGenerator godoc
// @Summary start to run a generator.
// @Description start to run a generator.
// @Tags generator
// @Accept json
// @Produce json
// @Param name path string true "configuration name"
// @Success 204
// @Router /generators/{name}/start [post]
func StartGenerator(c *gin.Context) {
	name := c.Param("name")
	gm := c.MustGet("gm").(*generator.GeneratorManager)
	g := gm.Manager[name]
	go func() {
		g.Run(1000 * 1000)
	}()
	c.Status(http.StatusNoContent)
}

// StopGenerator godoc
// @Summary stop a running generator.
// @Description stop a running generator.
// @Tags generator
// @Accept json
// @Produce json
// @Param name path string true "configuration name"
// @Success 204
// @Router /generators/{name}/stop [post]
func StopGenerator(c *gin.Context) {
	name := c.Param("name")
	gm := c.MustGet("gm").(*generator.GeneratorManager)
	g := gm.Manager[name]
	go func() {
		g.Stop()
	}()
	c.Status(http.StatusNoContent)
}
