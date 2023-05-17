# mongonull

Usage example:

```go
package main

import (
	"github.com/Logistrack/mongonull"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var opts = options.Client().SetRegistry(mongonull.BuildDefaultRegistry())

```