package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	// Открываем файл proxy.txt для чтения
	file, err := os.Open("proxy.txt")
	if err != nil {
		fmt.Println("Ошибка при открытии файла:", err)
		return
	}
	defer file.Close()

	// Создаем файл results.txt для записи
	outputFile, err := os.Create("results.txt")
	if err != nil {
		fmt.Println("Ошибка при создании файла:", err)
		return
	}
	defer outputFile.Close()

	// Создаем буфер для чтения строк из файла
	scanner := bufio.NewScanner(file)
	writer := bufio.NewWriter(outputFile)

	for scanner.Scan() {
		line := scanner.Text()

		// Разбиваем строку на части по символу @
		parts := strings.Split(line, "@")
		if len(parts) != 2 {
			fmt.Println("Некорректный формат строки:", line)
			continue
		}

		// Извлекаем часть с IP-адресом и портом
		ipPort := parts[1]

		// Разбиваем часть с IP-адресом и портом по символу :
		ipParts := strings.Split(ipPort, ":")
		if len(ipParts) != 2 {
			fmt.Println("Некорректный формат IP:Port:", ipPort)
			continue
		}

		// Записываем IP-адрес в файл results.txt
		ip := ipParts[0]
		_, err := writer.WriteString(ip + "\n")
		if err != nil {
			fmt.Println("Ошибка при записи в файл:", err)
			return
		}
	}

	// Проверяем ошибки при сканировании файла
	if err := scanner.Err(); err != nil {
		fmt.Println("Ошибка при чтении файла:", err)
	}

	// Сбрасываем буфер в файл
	writer.Flush()
}
