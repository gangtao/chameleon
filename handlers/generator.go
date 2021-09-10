package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type GeneratorConfig struct {
	Name string `json:"name"`
}

// CreateGenerator godoc
// @Summary Create a generator.
// @Description create a generator.
// @Tags generator
// @Accept json
// @Produce json
// @Param config body GeneratorConfig true "generator configuration"
// @Success 201 {object} GeneratorConfig
// @Router /generators [post]
func CreateGenerator(c *gin.Context) {
	json := GeneratorConfig{}
	c.BindJSON(&json)

	log.Printf("%v", &json)

	c.JSON(http.StatusCreated, json)
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
	res := [3]string{"a", "b", "c"}

	c.JSON(http.StatusOK, res)
}

// GetGenerator godoc
// @Summary get generator by name.
// @Description get generator by name.
// @Tags generator
// @Accept json
// @Produce json
// @Param name path string true "configuration name"
// @Success 200 {object} GeneratorConfig
// @Router /generators/{name} [get]
func GetGenerator(c *gin.Context) {
	name := c.Param("name")
	log.Printf("%s", name)

	res := GeneratorConfig{Name: name}
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
	log.Printf("%s", name)
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
	log.Printf("%s", name)
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
	log.Printf("%s", name)
	c.Status(http.StatusNoContent)
}
