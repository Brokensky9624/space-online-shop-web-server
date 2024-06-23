# Space online plan

## Function requirement

- Web UI & Restful api
- Member register/login
- Auth jwt
- Customize logger
- Admin/Member can create/query product
- Member can add product to order
- Member can manage order
- Testify


## System requirement

- Member edit/delete
- Product edit/delete
- Member can trace product and got notify when discount
- Deploy in k8s
- CI/CD

## Extra implement

- Implement pay system
- Auth Oath 2.0
- Kafuka/Redis

## DB model

Member

```go
type Member struct
    ID uint
    Account string
    Password string
    Name string
    Email string
    Phone string
    Address string
```

Product

```go
type Product struct {
    ID uint
    Name string
    Title string
    Desc string
    Brand string
    Store []StoreProduct `related`
}
type StoreProduct struct {
    ID uint
    Size string
    Color string
    Price uint
    ProductID uint
    Quantity uint
}
```

order

```go
type Order struct {
    ID uint
    StoreProductID uint
    Quantity uint
    MemberID uint 
}
```