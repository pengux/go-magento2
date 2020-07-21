# go-magento2
Magento 2 REST API client in Go

## Usage
```
c := magento2.NewClient("https://magento2-shop.com")

criteria := magento2.NewSearchCriteria()
c.Customers().Search(critera)
```
