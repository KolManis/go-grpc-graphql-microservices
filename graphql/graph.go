package main

// type Server struct{
// 	accountClient *account.Client
// 	catalogClient *catalog.Client
// 	orderClient   *order.Client
// }

func NewGraphQlServer(accountUrl, catalogUrl, orderUrl string) (*Server, error) {
	accountClient, err := account.NewClient(accountUrl)
	if err != nil {
		return nil, err
	}

	catalogClient, err := catalog.NewClient(catalogUrl)
	if err != nil {
		return nil, err
	}

	orderClient, err := order.NewClient(orderUrl)
	if err != nil {
		return nil, err
	}
}
