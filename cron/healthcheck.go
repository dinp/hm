package cron

import (
	"encoding/json"
	"fmt"
	"github.com/dinp/common/model"
	"github.com/dinp/hm/g"
	"github.com/fsouza/go-dockerclient"
	"github.com/go-av/curl"
	"log"
	"strconv"
	"strings"
	"time"
)

func HealthCheck() {
	duration := time.Duration(g.Config().CheckInterval) * time.Second
	time.Sleep(duration)
	for {
		time.Sleep(duration)
		healthCheck()
	}
}

func healthCheck() {
	err, httpBody := curl.String(g.Config().ServerHttpApi, "timeout=", time.Second*5)
	if err != nil {
		log.Printf("[ERROR] curl server http api fail: %s", err)
		return
	}

	json.Unmarshal([]byte(httpBody), &g.RealState)

	appNames := g.RealState.Keys()
	checkList, _ := getCheckList()
	for _, name := range appNames {
		if len(checkList) > 0 {
			for app, health := range checkList {
				if name == app {
					sa, _ := g.RealState.GetSafeApp(name)
					for _, c := range sa.M {
						rsAddress := c.Ip + ":" + strconv.Itoa(c.Ports[0].PublicPort)
						rsHealthAddress := "http://" + rsAddress + health

						err, responseBody := curl.String(rsHealthAddress, "timeout=", time.Second*time.Duration(g.Config().ResponseTimeout))
						if err != nil {
							log.Printf("[ERROR] curl app:%s rs:%s fail: %s", name, rsHealthAddress, err)
							log.Printf("[WARN] kill rs %s", rsAddress)
							dropContainer(c)

							continue
						}

						if strings.Contains(responseBody, g.Config().HealthSign) == false {
							log.Printf("[ERROR] app:%s rs:%s health check fail", name, rsHealthAddress)
							log.Printf("[WARN] kill rs %s", rsAddress)
							dropContainer(c)

							continue
						}
					}
				}
			}
		}
	}
}

func getCheckList() (map[string]string, error) {
	sql := "select name, health from app where health <> ''"
	rows, err := g.DB.Query(sql)
	if err != nil {
		log.Printf("[ERROR] exec %s fail: %s", sql, err)
		return nil, err
	}

	var checkList = make(map[string]string)
	for rows.Next() {
		var app, health string
		err = rows.Scan(&app, &health)
		if err != nil {
			log.Printf("[ERROR] %s scan fail: %s", sql, err)
			return nil, err
		}

		checkList[app] = health
	}

	return checkList, nil
}

func dropContainer(c *model.Container) {

	if g.Config().Debug {
		log.Println("drop container:", c)
	}

	addr := fmt.Sprintf("http://%s:%d", c.Ip, g.Config().DockerPort)
	client, err := docker.NewClient(addr)
	if err != nil {
		log.Println("docker.NewClient fail:", err)
		return
	}

	err = client.RemoveContainer(docker.RemoveContainerOptions{ID: c.Id, Force: true})
	if err != nil {
		log.Println("docker.RemoveContainer fail:", err)
		return
	}
}
