# gojpegquality

Read the jpeg quality

Inpiration from [nukr/jpegquality](https://github.com/nukr/jpegquality)

# Use

```command
go get github.com/madeinfree/gojpegquality
```

```go
import (
  quality "github.com/madeinfree/gojpegquality"
)
data, _ := ioutil.ReadFile("filePath")
q := gojpegquality.GetQ(File []buffer) float64
fmt.Println(q)
```

# License

MIT License (MIT)
