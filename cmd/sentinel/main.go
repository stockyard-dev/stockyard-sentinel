package main
import ("fmt";"log";"net/http";"os";"github.com/stockyard-dev/stockyard-sentinel/internal/server";"github.com/stockyard-dev/stockyard-sentinel/internal/store")
func main(){port:=os.Getenv("PORT");if port==""{port="8680"};dataDir:=os.Getenv("DATA_DIR");if dataDir==""{dataDir="./sentinel-data"}
db,err:=store.Open(dataDir);if err!=nil{log.Fatalf("sentinel: %v",err)};defer db.Close();srv:=server.New(db,server.DefaultLimits())
fmt.Printf("\n  Sentinel — Self-hosted alert manager\n  ─────────────────────────────────\n  Dashboard:  http://localhost:%s/ui\n  API:        http://localhost:%s/api\n  Fire:       POST http://localhost:%s/api/fire\n  Data:       %s\n  ─────────────────────────────────\n\n",port,port,port,dataDir)
log.Printf("sentinel: listening on :%s",port);log.Fatal(http.ListenAndServe(":"+port,srv))}
