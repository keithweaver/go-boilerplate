# Logging

The severity levels are fairly self explanatory:
- Info: Nice to know
- Warning: Error occurred but not mission critical
- Error: Mission critical error

## Context
We declare an initial context on each request, and add values to that context. This context will be passed through the stack trace. The most common use case will be the logger. This provides better reporting for where and how issues happened. Each instance of services, handler, repository, etc. should have a logger declared.

The handler will look like:

```
ctx := context.Background()
ctx = context.WithValue(ctx, logging.CtxDomain, "Cars")
ctx = context.WithValue(ctx, logging.CtxHandlerMethod, "GetAll")
ctx = context.WithValue(ctx, logging.CtxRequestID, uuid.New().String())

u.logger.Info(ctx, "Called")
```

This ctx will be the first argument as we call further down the stack (ie. The Service layer). The info log above has very little details, but it is backfilled by the context and would output:

```
INFO: 2021/03/11 07:52:54 logging.go:64: {"message": "Called", "requestId": "43cc7b70-f7e6-4646-a8d0-ec7e4af9a251", "domain": "Cars", "handlerMethod": "GetAll"}
```

The service layer would be:
```
func (c *CarsService) GetAll(ctx context.Context, session models.Session, query models.ListCarQuery) ([]models.Car, error) {
	ctx = context.WithValue(ctx, logging.CtxServiceMethod, "GetAll")
	
	...
	return cars, nil
}
```

Once again, we simply build onto the context. Understanding this context is only scope to this service. If we log on the next level up, it will not include the service layer information.

