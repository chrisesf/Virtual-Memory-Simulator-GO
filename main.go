package main

import (
	"bufio"
	"container/list"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"
)

// SimulationResult agrupa os resultados de uma simulação de substituição de páginas.
type SimulationResult struct {
	pageFaults int            // Total de faltas de página
	loadCounts map[string]int // Quantas vezes cada página foi carregada na memória
}

// parseMemorySize converte uma string de tamanho de memória (ex: "8MB", "16KB") para o valor em bytes.
//
// :param sizeStr: String contendo o tamanho da memória (ex: "8MB")
// :return: Valor convertido em bytes e um erro, se houver
func parseMemorySize(sizeStr string) (int64, error) {
	sizeStr = strings.ToUpper(strings.TrimSpace(sizeStr))
	var multiplier int64 = 1

	if strings.HasSuffix(sizeStr, "KB") {
		multiplier = 1024
		sizeStr = strings.TrimSuffix(sizeStr, "KB")
	} else if strings.HasSuffix(sizeStr, "MB") {
		multiplier = 1024 * 1024
		sizeStr = strings.TrimSuffix(sizeStr, "MB")
	} else if strings.HasSuffix(sizeStr, "GB") {
		multiplier = 1024 * 1024 * 1024
		sizeStr = strings.TrimSuffix(sizeStr, "GB")
	}

	value, err := strconv.ParseInt(sizeStr, 10, 64)
	if err != nil {
		return 0, err
	}
	return value * multiplier, nil
}

// fifoAlgorithm implementa o algoritmo de substituição de páginas FIFO.
//
// Utiliza uma lista duplamente encadeada (`list.List`) para armazenar a ordem de chegada
// das páginas e um mapa (`map[string]struct{}`) para checagem rápida de presença.
//
// :param pageReferences: Sequência de referências de página
// :param numFrames: Número de quadros de memória disponíveis
// :param didacticMode: Se true, imprime passo a passo
// :return: Resultado da simulação (faltas de página e carregamentos)
func fifoAlgorithm(pageReferences []string, numFrames int, didacticMode bool) SimulationResult {
	pageInMmemorySet := make(map[string]struct{}) // Set de páginas atualmente na memória
	memoryFrames := list.New()                    // Lista FIFO das páginas
	result := SimulationResult{pageFaults: 0, loadCounts: make(map[string]int)}

	for i, page := range pageReferences {
		if didacticMode {
			fmt.Printf("\n[FIFO - Passo %d] Acessando página: %s\n", i+1, page)
		}

		if _, found := pageInMmemorySet[page]; !found {
			result.pageFaults++
			result.loadCounts[page]++

			var evictedPage string
			if memoryFrames.Len() == numFrames {
				oldestPageElement := memoryFrames.Front()
				evictedPage = oldestPageElement.Value.(string)
				delete(pageInMmemorySet, evictedPage)
				memoryFrames.Remove(oldestPageElement)
			}

			memoryFrames.PushBack(page)
			pageInMmemorySet[page] = struct{}{}

			if didacticMode {
				fmt.Printf("  -> FALTA DE PáGINA (FAULT)!\n")
				if evictedPage != "" {
					fmt.Printf("     Página removida: %s\n", evictedPage)
				}
				fmt.Printf("     Página inserida: %s\n", page)
			}
		} else if didacticMode {
			fmt.Printf("  -> Página encontrada (HIT)!\n")
		}

		if didacticMode {
			var framesState []string
			for e := memoryFrames.Front(); e != nil; e = e.Next() {
				framesState = append(framesState, e.Value.(string))
			}
			fmt.Printf("  Estado da memoria: %v\n", framesState)
		}
	}
	return result
}

// optimalAlgorithmOptimized implementa o algoritmo ótimo de substituição de páginas.
//
// Utiliza pré-processamento das posições futuras de uso para decidir qual página remover.
// Remove aquela que será usada mais tarde ou nunca mais usada.
//
// :param pageReferences: Sequência de referências de página
// :param numFrames: Número de quadros de memória disponíveis
// :param pagePositions: Mapa com posições futuras de uso de cada página
// :param didacticMode: Se true, imprime passo a passo
// :return: Resultado da simulação (faltas de página e carregamentos)
func optimalAlgorithmOptimized(pageReferences []string, numFrames int, pagePositions map[string][]int, didacticMode bool) SimulationResult {
	pageInMmemorySet := make(map[string]struct{}) // Set de páginas em memória
	memoryFrames := make([]string, 0, numFrames)  // Lista de páginas na memória
	result := SimulationResult{pageFaults: 0, loadCounts: make(map[string]int)}
	nextUseCursor := make(map[string]int) // Cursor de leitura para cada página

	for i, page := range pageReferences {
		if didacticMode {
			fmt.Printf("\n[Ótimo - Passo %d] Acessando página: %s\n", i+1, page)
		}

		if _, found := pageInMmemorySet[page]; !found {
			result.pageFaults++
			result.loadCounts[page]++

			var evictedPage string
			if len(memoryFrames) < numFrames {
				memoryFrames = append(memoryFrames, page)
				pageInMmemorySet[page] = struct{}{}
			} else {
				farthest := -1
				victimIndex := -1

				for frameIdx, framePage := range memoryFrames {
					positions := pagePositions[framePage]
					cursor := nextUseCursor[framePage]

					nextPos := -1
					for cursor < len(positions) && positions[cursor] <= i {
						cursor++
					}
					nextUseCursor[framePage] = cursor
					if cursor < len(positions) {
						nextPos = positions[cursor]
					}

					if nextPos == -1 {
						victimIndex = frameIdx
						break
					}
					if nextPos > farthest {
						farthest = nextPos
						victimIndex = frameIdx
					}
				}

				evictedPage = memoryFrames[victimIndex]
				delete(pageInMmemorySet, evictedPage)
				memoryFrames[victimIndex] = page
				pageInMmemorySet[page] = struct{}{}
			}

			if didacticMode {
				fmt.Printf("  -> FALTA DE PáGINA (FAULT)!\n")
				if evictedPage != "" {
					fmt.Printf("     Página removida: %s\n", evictedPage)
				}
				fmt.Printf("     Página inserida: %s\n", page)
			}
		} else if didacticMode {
			fmt.Printf("  -> Página encontrada (HIT)!\n")
		}

		if didacticMode {
			framesState := make([]string, len(memoryFrames))
			copy(framesState, memoryFrames)
			sort.Strings(framesState)
			fmt.Printf("  Estado da memoria: %v\n", framesState)
		}
	}
	return result
}

// main coordena a execução do simulador. Lê o arquivo de entrada, calcula os parâmetros de memória,
// executa os algoritmos FIFO e Ótimo, imprime os resultados e pergunta se o usuário quer ver os detalhes.
func main() {
	if len(os.Args) < 3 {
		fmt.Fprintf(os.Stderr, "Uso: %s [--didatico] <arquivo_de_entrada> <tamanho_memoria>\n", os.Args[0])
		os.Exit(1)
	}

	didacticMode := false
	args := os.Args[1:]
	if args[0] == "--didatico" {
		didacticMode = true
		args = args[1:]
	}

	if len(args) != 2 {
		fmt.Fprintf(os.Stderr, "Argumentos invalidos. Uso: %s [--didatico] <arquivo_de_entrada> <tamanho_memoria>\n", os.Args[0])
		os.Exit(1)
	}

	memorySizeStr := args[1]

	const pageSizeBytes = 4 * 1024 // Cada página tem 4KB
	physicalMemoryBytes, _ := parseMemorySize(memorySizeStr)

	if physicalMemoryBytes < pageSizeBytes {
		fmt.Fprintf(os.Stderr, "Tamanho de memória deve ser maior que 4KB. Uso: %s [--didatico] <arquivo_de_entrada> <tamanho_memoria>\n", os.Args[0])
		os.Exit(1)
	}

	filePath := args[0]
	
	// Leitura do arquivo de referências
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro: O arquivo '%s' nao foi encontrado.\n", filePath)
		os.Exit(1)
	}
	defer file.Close()

	pageReferences := make([]string, 0)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		pageReferences = append(pageReferences, scanner.Text())
	}

	// Pré-processamento das posições das páginas
	fmt.Println("Iniciando pre-processamento do arquivo de referencias...")
	pagePositions := make(map[string][]int)
	for i, page := range pageReferences {
		pagePositions[page] = append(pagePositions[page], i)
	}
	fmt.Println("Pre-processamento concluido.")

	numFrames := int(physicalMemoryBytes / pageSizeBytes)

	// Contar páginas distintas para calcular tamanho da tabela
	distinctPagesSet := make(map[string]struct{})
	for _, page := range pageReferences {
		distinctPagesSet[page] = struct{}{}
	}
	distinctPagesCount := len(distinctPagesSet)

	const sizeOfPTE = 4
	tableSize := distinctPagesCount * sizeOfPTE

	var optimalResult, fifoResult SimulationResult

	if didacticMode && len(pageReferences) > 1000 {
		fmt.Println("AVISO: O modo didatico com muitas referencias pode gerar saida muito longa!")
	}

	optimalResult = optimalAlgorithmOptimized(pageReferences, numFrames, pagePositions, didacticMode)
	fifoResult = fifoAlgorithm(pageReferences, numFrames, didacticMode)

	// Impressão de estatísticas
	fmt.Println("\n--- RESULTADO DA SIMULACAO ---")
	fmt.Printf("A memória física comporta %d páginas.\n", numFrames)
	fmt.Printf("Ha %d páginas distintas no arquivo.\n", distinctPagesCount)
	fmt.Printf("Tamanho estimado da Tabela de Páginas (1 nivel): %d bytes (%d entradas * %d bytes/entrada)\n", tableSize, distinctPagesCount, sizeOfPTE)

	fmt.Printf("Com o algoritmo Ótimo ocorrem %d faltas de página.\n", optimalResult.pageFaults)
	fmt.Printf("Com o algoritmo FIFO ocorrem %d faltas de página,\n", fifoResult.pageFaults)

	efficiency := 100.0
	if fifoResult.pageFaults > 0 {
		efficiency = (float64(optimalResult.pageFaults) / float64(fifoResult.pageFaults)) * 100.0
	}
	fmt.Printf("atingindo %.2f%% do desempenho do Ótimo.\n", efficiency)

	// Pergunta se deseja imprimir estatísticas por página
	fmt.Print("Deseja listar o numero de carregamentos (s/n)? ")
	var choice string
	_, err = fmt.Scanln(&choice)
	if err != nil && err != io.EOF {
		choice = "n"
	}

	if strings.ToLower(choice) == "s" {
		distinctPages := make([]string, 0, len(distinctPagesSet))
		for page := range distinctPagesSet {
			distinctPages = append(distinctPages, page)
		}
		sort.Strings(distinctPages)

		fmt.Println("\nPágina\tÓtimo\tFIFO")
		fmt.Println("------\t-----\t----")
		for _, page := range distinctPages {
			optCount := optimalResult.loadCounts[page]
			fifoCount := fifoResult.loadCounts[page]
			fmt.Printf("%s\t%d\t%d\n", page, optCount, fifoCount)
		}
	}
}
