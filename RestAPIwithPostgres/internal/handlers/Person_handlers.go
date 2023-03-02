package Handler

import (
	"context"
	"fmt"
	Logic "myapp/internal/logic"
	Model "myapp/internal/model"
	Repository "myapp/internal/repository"
	"net/http"
	"time"

	"github.com/labstack/echo"
	log "github.com/sirupsen/logrus"
)

func PostPerson(c echo.Context) error {
	var newPerson Model.Person
	if err := c.Bind(&newPerson); err != nil {
		return err
	}
	p, err := Logic.Create(c.Request().Context(), newPerson)
	if err != nil {
		log.Error(err)
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	log.Info("Добавлена новая запись")
	return c.JSON(http.StatusCreated, p)
}

func GetPersons(c echo.Context) error {
	persons, err := Logic.Read(c.Request().Context())
	if err != nil {
		log.Error(err)
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	log.Info("Выведены все записи")
	return c.JSON(http.StatusOK, persons)
}

func GetById(c echo.Context) error {
	id := c.Param("id")
	persons, err := Logic.ReadOne(c.Request().Context(), id)
	if err != nil {
		log.Error(err)
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	log.Infof("Выведена Запись с id = %s", id)
	return c.JSON(http.StatusOK, persons)
}

func DeleteById(c echo.Context) error {
	id := c.Param("id")
	err := Logic.Delete(c.Request().Context(), id)
	if err != nil {
		log.Error(err)
		return c.JSON(http.StatusBadRequest, fmt.Sprint(err))
	}
	log.Infof("Запись с id = %s  успешно удалена", id)
	return c.JSON(http.StatusOK, fmt.Sprintf("Запись с id = %s  успешно удалена", id))
}

func UpdatePersonById(c echo.Context) error {
	var newPerson Model.Person
	id := c.Param("id")
	if err := c.Bind(&newPerson); err != nil {
		return err
	}
	err := Logic.Update(c.Request().Context(), id, newPerson)
	if err != nil {
		log.Error(err)
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	log.Infof("Запись с id = %s  успешно обновлена", id)
	return c.JSON(http.StatusOK, fmt.Sprintf("Запись с id = %s  успешно обновлена", id))
}

//Middleware
func ConnectDB(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		errorCh := make(chan error)
		go check(ctx, errorCh)
		err := Repository.OpenTable()
		//time.Sleep(4 * time.Second) //Используем для имитации долгого подключения к БД
		errorCh <- err
		return next(c)
	}
}

func check(ctx context.Context, errorCh chan error) {
	select {
	case <-ctx.Done():
		log.Fatalf("Timed out: %v", ctx.Err())
		return
	case err := <-errorCh:
		if err != nil {
			log.Fatalf("Возникла ошибка... %v", err)
			return
		}
	}
}
