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
		Name         string `json:"name"`
		GenerateName string `json:"generate_name"`
		Namespace    string `json:"namespace"`
		Image        string `json:"image"`
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
	j := &job{}
	log.Println("Creating new job")

	if err := c.Bind(j); err != nil {
		return err
	}
	goJobRepo := inkube.NewGoJobsRepository(config.KubeClient)

	goJobCreated, err := goJobRepo.Create(j.Name, j.GenerateName, j.Namespace, j.Image)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	log.Printf("goJob created %v", goJobCreated)

	return c.JSON(http.StatusCreated, goJobCreated)
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
	namespace := c.Param("namespace")
	name := c.Param("name")
	log.Printf("Getting job from namespace: %v name: %v", namespace, name)

	goJobRepo := inkube.NewGoJobsRepository(config.KubeClient)
	err := goJobRepo.Delete(name, namespace)
	if err != nil {
		log.Printf("error deleting job: %v", err)
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.NoContent(http.StatusNoContent)
}
