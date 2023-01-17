package service

import (
	"context"
	"io"

	"github.com/diegofonseca/gRPC/internal/database"
	"github.com/diegofonseca/gRPC/internal/pb"
)

type CategoryService struct {
	pb.UnimplementedCategoryServiceServer
	CategoryDB database.Category
}

type Search struct {}

func NewCategoryService(categoryDB database.Category) *CategoryService {
	return &CategoryService{
		CategoryDB: categoryDB,
	}
}

func (c *CategoryService) CreateCategory(ctx context.Context, in *pb.CreateCategoryRequest) (*pb.CategoryResponse, error) {
	category, err := c.CategoryDB.Create(in.Name, in.Description)
	if err != nil {
		return nil, err
	}
	categoryResponse := &pb.Category{
		Id: category.ID,
		Name: category.Name,
		Description: category.Description,
	}

	return &pb.CategoryResponse{Category: categoryResponse}, nil
}

func (c *CategoryService) GetCategory(ctx context.Context, in *pb.GetCategoryRequest) (*pb.CategoryResponse, error) {
	category, err := c.CategoryDB.FindByID(in.Id)

	if err != nil {
		panic(err)
	}

	categoryResponse := &pb.Category{
		Id: category.ID,
		Name: category.Name,
		Description: category.Description,
	}

	return &pb.CategoryResponse{Category: categoryResponse}, nil
}

func (c *CategoryService) ListCategory(ctx context.Context, in *pb.Search) (*pb.Categories, error) {
	// category := &c.CategoryDB
	categories, err := c.CategoryDB.FindAll()

	if err != nil {
		panic(err)
	}

	response := []*pb.Category{}
	for i := 0; i < len(categories); i++ {
		c := categories[i]
		response = append(response, &pb.Category{
			Id: c.ID,
			Name: c.Name,
			Description: c.Description,
		})
	}

	return &pb.Categories{Category: response}, nil
}

func (c *CategoryService) CreateCategoryStream(stream pb.CategoryService_CreateCategoryStreamServer) error {
	categories := &pb.Categories{}
	for {
		category, err := stream.Recv()
        if err == io.EOF {
            return stream.SendAndClose(categories)
        }
        if err!= nil {
            return err
        }

		categoryResult, err := c.CategoryDB.Create(category.Name, category.Description)

		if err!= nil {
            return err
        }

		categories.Category = append(categories.Category, &pb.Category{
			Id: categoryResult.ID,
            Name: categoryResult.Name,
            Description: categoryResult.Description,
		})
	}
}

func (c *CategoryService) CreateCategoryStreamBidirectional(stream pb.CategoryService_CreateCategoryStreamBidirectionalServer) error {
	for {
		category, err := stream.Recv()
        if err == io.EOF {
            return  nil
        }
        if err!= nil {
            return err
        }

		categoryResult, err := c.CategoryDB.Create(category.Name, category.Description)

		if err!= nil {
            return err
        }

		err = stream.Send(&pb.Category{
			Id: categoryResult.ID,
            Name: categoryResult.Name,
            Description: categoryResult.Description,
		})

		if err!= nil {
            return err
        }
	}
}