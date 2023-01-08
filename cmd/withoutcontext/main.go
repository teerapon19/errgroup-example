package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/teerapon19/errgroup-example/pkg/model"
	"github.com/teerapon19/errgroup-example/pkg/utils"
	"golang.org/x/sync/errgroup"
)

func main() {
	productsData := map[string]string{
		"item": "https://gist.githubusercontent.com/teerapon19/11bdafef545766c22e065fa02c80d484/raw/3fa72709faa7871328759d93d208f3065e540ee2/products.json",
		"type": "https://gist.githubusercontent.com/teerapon19/c9f6bca3f7b6a7d3f52f2cd369193a86/raw/6afb433e1b77e12b82455c6345adc1f5e5cc83e5/types.json",
	}

	eg := new(errgroup.Group)

	productsCh := make(chan []model.ProductResponse)
	typesCh := make(chan []model.Type)

	eg.Go(func() error {
		p := []model.ProductResponse{}
		err := utils.FetchJson(&p, productsData["item"])
		productsCh <- p
		return err
	})

	eg.Go(func() error {
		types := []model.Type{}
		err := utils.FetchJson(&types, productsData["type"])
		typesCh <- types
		return err
	})

	productResponses, types := <-productsCh, <-typesCh

	if err := eg.Wait(); err != nil {
		log.Fatal(err)
	}

	products := make([]model.Product, 0)
	for _, productResponse := range productResponses {
		productType := model.Type{}
		for _, typeItem := range types {
			if productResponse.Type == typeItem.ID {
				productType = typeItem
				break
			}
		}
		products = append(products, model.Product{
			ID:    productResponse.ID,
			Name:  productResponse.Name,
			Price: productResponse.Price,
			Type:  productType,
		})
	}

	productsJSON, err := json.MarshalIndent(products, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Products: %s\n", string(productsJSON))
}
