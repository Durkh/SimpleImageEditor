package cli

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func Run() {

	if len(os.Args) < 2 {
		noParameters()
	}

	parameters()

}

func noParameters() {

	var (
		input string
		stdin = bufio.NewReader(os.Stdin)
	)

	for {
		fmt.Println("\nDigite o path para a imagem a ser editada")
		_, err := fmt.Scanf("%100s", &input)
		if err != nil {
			log.Fatal(err)
		}
		if input == "quit" {
			break
		}

		if n, err := stdin.Discard(stdin.Buffered()); err != nil && n > 0 {
			fmt.Println("você excedeu o limite de caracteres para o path, por favor mantenha menor que 100 caracteres")
			continue
		}

		fmt.Println(input)
	}

}

func parameters() {

}
