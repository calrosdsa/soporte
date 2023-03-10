package main
import (
    "fmt"
    "log"
    "net/http"

    // "golang.org/x/sync/errgroup"
)


func ServeHTTP(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "text/html; charset=utf-8")
    w.Header().Set("Access-Control-Allow-Origin", "*")
    body :=`[
        <>
        <div onClick={printInConsole} className="bg-gray-500 cursor-pointer">
        <input onChange={onchange}/>
        <button className=" bg-amber-300">Save</button>
        <button >Delete</button>
        </div>
        </>
        ]`
    fmt.Fprint(w, body)
}

func main() {
    http.HandleFunc("/", ServeHTTP)
    if err := http.ListenAndServe(":80", nil); err != nil {
        log.Fatal(err)
    }
}
// func Task(task int) error {
// 	fmt.Println(task)
//     if 2 < task {
//         return fmt.Errorf("Task %v failed", task)
//     }
//     fmt.Printf("Task %v completed", task)
//     return nil
// }
// func main() {
    // eg := &errgroup.Group{}
    // for i := 0; i < 11; i++ {
    //     task := i
    //     eg.Go(func() error {
    //         return Task(task)
    //     })
    // }
    // if err := eg.Wait(); err != nil {
    //     log.Fatal("Error", err)
    // }
    // fmt.Println("Completed successfully!")
// }