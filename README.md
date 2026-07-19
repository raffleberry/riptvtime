# riptvtime

## Essential Features list
| Features | Done |
| :--- | :---: |
| Watch Next Feed | 80% |
| Upcoming(To Be Aired) Episodes Feed | ❌ |
| Search & Add TV show, Stop TV show, Remove TV show | 60% |
| Episode Mark as watched/Pop up to mark all previous as watched | 0% |
| Import TVTime gdpr data | 0% |
| Total Time Watched | 0% |
| Episodes Watched/Rewatched that count towards total time | 0% |

## Optional Features list
| Features | Done |
| :--- | :---: |
| Recommendation Feed Based on Watched TV Shows | ❌ |
| Stats Dashboard(Time-range watched hours, genre of TV Shows Pie chart, etc) | ❌ |
| Comments/Bookmarking episodes | ❌ |
| loads posters/images and looks good | ❌ |


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
```

AND dont forget to always compare with
```go
errors.Is(err, foo.ErrNotFound)
```

2. Js File names

- Components - `MyComponent.js`
    - Camel Case with first letter capitalized
- Other Files such as logic, stores, etc - `myComponentStore.js` Camel Case but First letter is lowercase 
