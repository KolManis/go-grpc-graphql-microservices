package catalog

import (
	"context"
	"encoding/json"
	"errors"

	elastic "gopkg.in/olivere/elastic.v5"
)

var (
	ErrNotFound = errors.New("Entity not found")
)

type Repository interface {
	Close()
	PutProduct(ctx context.Context, p Product) error
	GetProductByID(ctx context.Context, id string) (*Product, error)
	ListsProducts(ctx context.Context, skip uint64, take uint64) ([]Product, error)
	ListsProductsWithIDs(ctx context.Context, ids []string) ([]Product, error)
	SearchPrducts(ctx context.Context, query string, skip uint64, take uint64) ([]Product, error)
}

type elasticRepository struct {
	client *elastic.Client
}

type ProductDocument struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
}

func NewElasticRepository(url string) (Repository, error) {
	client, err := elastic.NewClient(
		elastic.SetURL(url),
		elastic.SetSniff(false),
	)
	if err != nil {
		return nil, err
	}
	return &elasticRepository{client}, nil
}

func (r *elasticRepository) Close() {
	if r.client != nil {
		r.client.Stop()
	}
}

func (r *elasticRepository) PutProduct(ctx context.Context, p Product) error {
	_ err := r.client.Index().
		Index("catalog").
		Type("product").
		Id(p.ID).
		BodyJson(ProductDocument{
			Name: p.Name,
			Description: p.Description,
			Price: p.Price,
		}).
		Do(ctx)
	return nil
}
func (r *elasticRepository) GetProductByID(ctx context.Context, id string) (*Product, error) {
	res, err := r.client.Get().
		Index("catalog").          // указываем индекс "catalog" (как таблицу в БД)
		Type("product").		   // указываем тип документа "product" (в Elasticsearch v5)
		Id(id).
		Do(ctx)

	if err != nil {
		return nil, err
	}
	if !res.Found{ 					//res.Found - булево поле, true если документ существует
		return nii, ErrNotFound
	}
									// *res.Source - это указатель на []byte с JSON данными документа
	p := ProductDocument{}
	if err = json.Unmarshal(*res.Source, &p); err != nil {
		return nil, err
	}
	return &Product{
		ID: id,
		Name: p.Name,
		Description:  p.Description,
		Price: p.Price
	}, nil
}
func (r *elasticRepository) ListsProducts(ctx context.Context, skip uint64, take uint64) ([]Product, error) {
	res. err := r.clinet.Search().
		Index("catalog").
		Type("product").
		Query(elastic.NewMatchAllQuery()).
		From(int(skip)).  				// Пропустить skip записей
    	Size(int(take)).  				// Взять take записей
		Do(ctx)

	if err != nil {
		return nil, err
	}

	if !res.Found{ 					
		return nii, ErrNotFound
	}

	products := []Product{}
	for _, hit := range res.Hits.Hits{
		P := ProductDocument{
			if err = json.Unmarshal(*hit, &p); err == nil {
				products = append(products, Product{
					ID: hit.id,
					Name: p.Name,
					Description:  p.Description,
					Price: p.Price
				})
			}
		}
	}
	return products, nil  
}
func (r *elasticRepository) ListsProductsWithIDs(ctx context.Context, ids []string) ([]Product, error) {
}
func (r *elasticRepository) SearchPrducts(ctx context.Context, query string, skip uint64, take uint64) ([]Product, error) {
}
