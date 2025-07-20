# Virtual Memory Simulator

Simulador dos algoritmos de substituição de páginas FIFO e Ótimo, com modo didático para acompanhar passo a passo. Este programa simula dois algoritmos clássicos de substituição de páginas da memória virtual: FIFO (First-In, First-Out) que remove a página mais antiga na memória, e Ótimo, que remove a página que será usada mais tarde no futuro (ideal para comparação). O simulador lê um arquivo texto com a sequência de referências de páginas e simula a execução considerando a memória física especificada. O modo didático (`--didatico`) exibe passo a passo os acessos, faltas, páginas removidas e estado da memória.

## Requisitos

- Go 1.18+ instalado na sua máquina  
- Arquivo de referências de páginas (um arquivo texto com uma página por linha)  
- Memória física especificada em KB, MB ou GB (ex: 8MB)  

## Como rodar

1. Clone o repositório:  
   ```bash
   git clone https://github.com/seuusuario/seurepositorio.git
   cd seurepositorio
2. Compile o programa:
   ```bash
   go build -o simulator.exe main.go:  
3. Prepare seu arquivo de entrada (ex: pages.txt), com uma página por linha, exemplo:
    ```nginx
    I0
    D0
    I0
    D1
    I1
4. Execute o simulador:
   ```bash
    ./simulator <arquivo_de_entrada> <tamanho_memoria>

    Exemplo:
    ./simulator pages.txt 8MB
5. Para rodar no modo didático (passo a passo), use o flag --didatico:
   ```bash
    ./simulator --didatico pages.txt 8MB
6. Após a execução, o programa exibirá o resumo das faltas de página e perguntará se deseja listar o número de carregamentos por página.

## Sobre o simulador
Este programa simula dois algoritmos clássicos de substituição de páginas da memória virtual:

FIFO (First-In, First-Out): remove a página que está na memória há mais tempo.

Ótimo: remove a página que será usada mais tarde no futuro (ideal e impossível de implementar em sistemas reais, mas ótimo para comparação).

O simulador lê um arquivo texto com a sequência de referências de páginas e simula a execução dos dois algoritmos considerando a memória física especificada.

Modo didático (--didatico) exibe passo a passo os acessos, faltas de página, páginas removidas e estado atual da memória.

