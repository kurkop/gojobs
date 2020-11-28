package handlers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/kurkop/gojobs/cmd/gojobs-api/config"
	"github.com/kurkop/gojobs/internal/jobs/storage/inkube"
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

	goJobCreated, err := goJobRepo.Create(j.Name, j.GenerateName, j.Image)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	log.Printf("goJob created %v", goJobCreated)

	return c.JSON(http.StatusCreated, goJobCreated)
}

func GetJob(c echo.Context) error {
	name := c.Param("name")
	log.Printf("Getting job - name: %v", name)

	goJobRepo := inkube.NewGoJobsRepository(config.KubeClient)
	goJobGot, err := goJobRepo.Get(name)
	if err != nil {
		log.Printf("error getting job: %v", err)
	}
	log.Printf("%v", goJobGot)
	return c.JSON(http.StatusOK, goJobGot)
}

func GetAllJob(c echo.Context) error {
	log.Printf("Getting job from gojobs")

	goJobRepo := inkube.NewGoJobsRepository(config.KubeClient)
	goJobsGot, err := goJobRepo.GetAll()
	if err != nil {
		log.Printf("error getting job: %v", err)
	}
	log.Printf("Go Jobs return %v", goJobsGot.Items)
	return c.JSON(http.StatusOK, goJobsGot.Items)
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
	name := c.Param("name")
	log.Printf("Deleting job - name: %v", name)

	goJobRepo := inkube.NewGoJobsRepository(config.KubeClient)
	err := goJobRepo.Delete(name)
	if err != nil {
		log.Printf("Error deleting job: %v", err)
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.NoContent(http.StatusNoContent)
}
