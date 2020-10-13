package handlers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/kurkop/gojob/cmd/gojob-api/config"
	"github.com/kurkop/gojob/internal/jobs/storage/inkube"
	"github.com/labstack/echo/v4"
)

type (
	user struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}
	job struct {
		Name      string `json:"name"`
		Namespace string `json:"namespace"`
	}
)

var (
	users = map[int]*user{}
	seq   = 1
)

//----------
// Handlers
//----------

func CreateJob(c echo.Context) error {
	u := &user{
		ID: seq,
	}
	if err := c.Bind(u); err != nil {
		return err
	}
	users[u.ID] = u
	seq++
	return c.JSON(http.StatusCreated, u)
}

func GetJob(c echo.Context) error {
	namespace := c.Param("namespace")
	name := c.Param("name")
	log.Printf("Getting job from namespace: %v name: %v", namespace, name)

	goJobRepo := inkube.NewGoJobsRepository(config.KubeClient)
	goJobGot, err := goJobRepo.Get(name, namespace)
	if err != nil {
		log.Printf("error getting job: %v", err)
	}
	log.Printf("%v", goJobGot)
	return c.JSON(http.StatusOK, goJobGot)
}

func UpdateJob(c echo.Context) error {
	u := new(user)
	if err := c.Bind(u); err != nil {
		return err
	}
	name, _ := strconv.Atoi(c.Param("name"))
	users[name].Name = u.Name
	return c.JSON(http.StatusOK, users[name])
}

func DeleteJob(c echo.Context) error {
	name, _ := strconv.Atoi(c.Param("name"))
	delete(users, name)
	return c.NoContent(http.StatusNoContent)
}
