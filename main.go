package main
import (
	"fmt"
	"bufio"
	"io"
	"os"
	"log"
)


func bits( b byte){
	for i:=7; i>=0; i-- {
		fmt.Printf("%d", (b >> i) & 1)
	}
}

func main(){
	file, err := os.Open("test.data")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	buf := make([]byte,1)

	E := newEstimator8()

	for {
		_, err := reader.Read(buf)
		if err != nil {
			if err != io.EOF {
				log.Fatal(err)
			}
			break
		}
		b := buf[0]

		for i:=7 ; i>=0; i-- {
			fmt.Printf("%d %d\n", E.ctxt.b, E.ctxt.c)
			E.update( ((b >> i)&1)!=0 )
			fmt.Println( E.predict())
		}

		//bits(buf[0])
		//bits(buf[1])
		//fmt.Printf( " " + string(buf[0:n]) + "\n")
	}
}
