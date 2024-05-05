package main

import (
	"crypto/aes"
	"crypto/rand"
	"fmt"
)

func generateRandom(size int) ([]byte, error) {
	// генерируем криптостойкие случайные байты в b
	b := make([]byte, size)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func main() {
	// src := []byte("Супервважное") // данные, которые хотим зашифровать
	src := []byte{0x68, 0x51, 0x38, 0x4a, 0x14, 0x2d, 0xc5, 0x86, 0xbb, 0x20, 0x2f, 0x3, 0x52, 0xd1, 0x59, 0x19, 0xc7, 0xcf, 0xce, 0x48, 0x17, 0x3b, 0xd9, 0x8e}
	fmt.Println(len(src))
	fmt.Printf("original: %s %x\n", src, src)
	// expected: []byte{0x68, 0x51, 0x38, 0x4a, 0x14, 0x2d, 0xc5, 0x86, 0xbb, 0x20, 0x2f, 0x3, 0x52, 0xd1, 0x59, 0x19, 0xc7, 0xcf, 0xce, 0x48, 0x17, 0x3b, 0xd9, 0x8e}
	// actual  : []byte{0x68, 0x51, 0x38, 0x4a, 0x14, 0x2d, 0xc5, 0x86, 0xbb, 0x20, 0x2f, 0x3, 0x52, 0xd1, 0x59, 0x19, 0xaa, 0x1c, 0xa5, 0xa5, 0x49, 0xc, 0x50, 0x26, 0xf1, 0x33, 0xca, 0x60, 0xe3, 0x8c, 0x17, 0xba}

	// константа aes.BlockSize определяет размер блока, она равна 16 байтам
	key, err := generateRandom(32) // ключ шифрования
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}
	fmt.Printf("key: %x\n", key)
	// 6851384a142dc586bb202f0352d15919c7cfce48173bd98e
	// 6851384a142dc586bb202f0352d15919c7cfce48173bd98e
	// получаем cipher.Block
	aesblock, err := aes.NewCipher(key)
	aesblock2, err := aes.NewCipher(key)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}

	srcLen := len(src)

	blockCnt := (srcLen + aes.BlockSize - 1) / aes.BlockSize
	src = append(src, make([]byte, aes.BlockSize*blockCnt-srcLen)...)

	fmt.Println(blockCnt)

	dst := make([]byte, blockCnt*aes.BlockSize) // зашифровываем
	for i := 0; i < blockCnt; i += 1 {
		currSrc := src[i*aes.BlockSize : (i+1)*aes.BlockSize]
		currDst := dst[i*aes.BlockSize : (i+1)*aes.BlockSize]
		aesblock.Encrypt(currDst, currSrc)
	}

	src2 := make([]byte, blockCnt*aes.BlockSize) // расшифровываем
	for i := 0; i < blockCnt; i += 1 {
		currDst := dst[i*aes.BlockSize : (i+1)*aes.BlockSize]
		currSrc2 := src2[i*aes.BlockSize : (i+1)*aes.BlockSize]
		aesblock2.Decrypt(currSrc2, currDst)
	}

	// расшифровали 16 байт
	// fmt.Printf("encrypted: %x\n", dst[:srcLen])
	fmt.Printf("decrypted: %x\n", src2[:srcLen])
}

// func main() {
// 	src := []byte("Ключ от сердца") // данные, которые хотим зашифровать
// 	fmt.Printf("original: %s\n", src)

// 	// будем использовать AES-256, создав ключ длиной 32 байта
// 	key, err := generateRandom(2 * aes.BlockSize) // ключ шифрования
// 	if err != nil {
// 		fmt.Printf("error: %v\n", err)
// 		return
// 	}

// 	// NewCipher создает и возвращает новый cipher.Block.
// 	// Ключевым аргументом должен быть ключ AES, 16, 24 или 32 байта
// 	// для выбора AES-128, AES-192 или AES-256.
// 	aesblock, err := aes.NewCipher(key)
// 	if err != nil {
// 		fmt.Printf("error: %v\n", err)
// 		return
// 	}

// 	// NewGCM возвращает заданный 128-битный блочный шифр
// 	aesgcm, err := cipher.NewGCM(aesblock)
// 	if err != nil {
// 		fmt.Printf("error: %v\n", err)
// 		return
// 	}

// 	// создаём вектор инициализации
// 	nonce, err := generateRandom(aesgcm.NonceSize())
// 	if err != nil {
// 		fmt.Printf("error: %v\n", err)
// 		return
// 	}

// 	dst := aesgcm.Seal(nil, nonce, src, nil) // зашифровываем
// 	fmt.Printf("encrypted: %x\n", dst)

// 	nonce, err = generateRandom(aesgcm.NonceSize())
// 	if err != nil {
// 		fmt.Printf("error: %v\n", err)
// 		return
// 	}

// 	src2, err := aesgcm.Open(nil, nonce, dst, nil) // расшифровываем
// 	if err != nil {
// 		fmt.Printf("error: %v\n", err)
// 		return
// 	}
// 	fmt.Printf("decrypted: %s\n", src2)
// }
