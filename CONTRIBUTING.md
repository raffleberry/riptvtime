
## Style guides (So i don't forget)
1. Error Handling
When returning errors, wrap the parent error in this order if you need to add some details/variables from the current function like for example :
```go
package foo

import (
    // ...
    "db"
    // ...
)

var (
    ErrNotFound = errors.New("Not Found -_-")
)
// ...
// ...
_, parentErr := db.someDbQuery()
fmt.Errorf("%w: %w: Some Details or Variables x=%v", parentErr, ErrNotFound, x) 

```

OR
just use

```go
errors.Join(parentErr, ErrNotFound)
errors.Join(err, errors.New(utils.Jn("ep", ep, "sn", sn)))
```

AND dont forget to always compare with
```go
errors.Is(err, foo.ErrNotFound)
```

2. Js File names

- Components - `MyComponent.js`
    - Camel Case with first letter capitalized
- Other Files such as logic, stores, etc - `myComponentStore.js` Camel Case with First letter in lowercase 
