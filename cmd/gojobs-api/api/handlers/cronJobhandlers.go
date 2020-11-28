package handlers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/kurkop/gojobs/cmd/gojobs-api/config"
	"github.com/kurkop/gojobs/internal/cronjobs/storage/inkube"
	"github.com/labstack/echo/v4"
)

type (
	cronjob struct {
		Name     string `json:"name"`
		Image    string `json:"image"`
		Schedule string `json:"schedule"`
	}
)

//----------
// Handlers
//----------

func CreateCronJob(c echo.Context) error {
	j := &cronjob{}
	log.Println("Creating new cronjob")

	if err := c.Bind(j); err != nil {
		return err
	}
	goCronJobRepo := inkube.NewGoCronJobsRepository(config.KubeClient)

	goCronJobCreated, err := goCronJobRepo.Create(j.Name, j.Image, j.Schedule)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	log.Printf("goCronJob created %v", goCronJobCreated)

	return c.JSON(http.StatusCreated, goCronJobCreated)
}

func GetCronJob(c echo.Context) error {
	name := c.Param("name")
	log.Printf("Getting cronjob - name: %v", name)

	goCronJobRepo := inkube.NewGoCronJobsRepository(config.KubeClient)
	goCronJobGot, err := goCronJobRepo.Get(name)
	if err != nil {
		log.Printf("error getting cronjob: %v", err)
	}
	log.Printf("%v", goCronJobGot)
	return c.JSON(http.StatusOK, goCronJobGot)
}

func GetAllCronJob(c echo.Context) error {
	log.Printf("Getting job from gojobs")

	goCronJobRepo := inkube.NewGoCronJobsRepository(config.KubeClient)
	goCronJobsGot, err := goCronJobRepo.GetAll()
	if err != nil {
		log.Printf("error getting cronjob: %v", err)
	}
	log.Printf("Go CronJobs return %v", goCronJobsGot.Items)
	return c.JSON(http.StatusOK, goCronJobsGot.Items)
}

func UpdateCronJob(c echo.Context) error {
	u := new(user)
	if err := c.Bind(u); err != nil {
		return err
	}
	name, _ := strconv.Atoi(c.Param("name"))
	users[name].Name = u.Name
	return c.JSON(http.StatusOK, users[name])
}

func DeleteCronJob(c echo.Context) error {
	name := c.Param("name")
	log.Printf("Deleting cronjob - name: %v", name)

	goCronJobRepo := inkube.NewGoCronJobsRepository(config.KubeClient)
	err := goCronJobRepo.Delete(name)
	if err != nil {
		log.Printf("Error deleting cronjob: %v", err)
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.NoContent(http.StatusNoContent)
}
