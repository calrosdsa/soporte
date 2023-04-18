package util

import (
	// "fmt"
	"strings"
	"strconv"
)

func AppendQueries(condition,q string,columns []string) (string, error){
	qParts := make([]string, 0, 2)
	// args := make([]interface{}, 0, 2)
	for i := 0; i < len(columns	);i++{
			qParts = append(qParts,columns[i] + ` = $` + strconv.Itoa(i+1))
	}
	q += strings.Join(qParts, ",") + ` WHERE ` + condition +`= $` + strconv.Itoa(len(columns) +1)
     return q,nil
}

// package util

// import (
// 	// "fmt"
// 	"strings"
// 	"strconv"
// )

// func AppendQueries(condition,q string,columns []string) (string, error){
// 	qParts := make([]string, len(columns))
// 	// args := make([]interface{}, 0, 2)
// 	for i := 0; i < len(columns	);i++{
// 		    qParts[i] = columns[i] + ` = $` + strconv.Itoa(i+1)
// 			// qParts = append(qParts,columns[i] + ` = $` + strconv.Itoa(i+1))
// 	}
// 	q += strings.Join(qParts, ",") + ` WHERE ` + condition +`= $` + strconv.Itoa(len(columns) +1 )
//      return q,nil
// }

