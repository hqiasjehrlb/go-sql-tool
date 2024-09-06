# sqltool

A simple tool to execute sql and scan returned rows into generic struct.

一個透過泛型處理sql query回傳資料的簡易工具

## Usage

```go

type MyData struct {
  ID   string `db:"id"`
  Col1 string `db:"col_1"`
}

func main() {
  var db *sql.DB // db can also be *sql.Conn
  var query string
  var args []any

  // ... neet to initialize DB connection

  // query id, col_1 and rows count from my_table
  query = `
    SELECT
      id, col_1
      count(*) over() AS count
    FROM my_table;
  `
  args = []any{"%"}
  datas, ext, err := sqltool.QueryRows[MyData](context.Background(), db, query, args...)
  if err != nil {
    // error handle
  }
  // datas will be type of []MyData
  // ext will be type of []map[string]any and `count` will be contained in this map

  var total int64 = 0
  for i, row := range datas {
    fmt.Println("data:", row)
    count := ext[i].GetInt("count")
    if count != nil {
      total = *count
    }
  }
}
```