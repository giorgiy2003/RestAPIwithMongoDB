package Logic

import (
	"context"
	"database/sql"
	Model "myapp/internal/model"
	Repository "myapp/internal/repository"
	"os"

	"github.com/sirupsen/logrus"
)

func Create(ctx context.Context, p Model.Person) (*Model.Person, error) {
	pid, err := Repository.Create(ctx, p)
	if err != nil {
		return nil, err
	}
	p.Id = *pid
	return &p, nil
}

func Read(ctx context.Context) ([]*Model.Person, error) {
	p, err := Repository.Read(ctx)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func ReadOne(ctx context.Context, id string) (*Model.Person, error) {
	
	person, err := Repository.ReadOne(ctx, id)
	if err != nil {
		return nil, err
	}
	p := Model.Person{}
	if p == *person {
		return nil, sql.ErrNoRows
	}
	return person, nil
}

func Update(ctx context.Context, id string, p Model.Person) error {
	_, err := ReadOne(ctx, id)
	if err != nil {
		return err
	}
	err = Repository.Update(ctx, id, p)
	if err != nil {
		return err
	}
	return nil
}

func Delete(ctx context.Context, id string) error {
	_, err := ReadOne(ctx, id)
	if err != nil {
		return err
	}
	err = Repository.Delete(ctx, id)
	if err != nil {
		return err
	}
	return nil
}

func InitLogger() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(os.Stdout)
}
