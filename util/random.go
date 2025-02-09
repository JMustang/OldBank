package util

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

const (
	alphabetLower = "abcdefghijklmnopqrstuvwxyz"
	alphabetUpper = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	minBalance    = 0
	maxBalance    = 1000000 // Increased range for better test coverage
)

var currencies = []string{"EUR", "USD", "CAD", "BRL"}

func init() {
	rand.Seed(time.Now().UnixNano())
}

// RandomInt generates a random integer between min and max
func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

// RandomString generates a random string of length n
func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabetLower)

	for i := 0; i < n; i++ {
		c := alphabetLower[rand.Intn(k)]
		sb.WriteByte(c)
	}
	return sb.String()
}

// RandomOwner generates a random owner name
// Now includes a prefix for better test identification
func RandomOwner() string {
	return "user_" + RandomString(6)
}

// RandomMoney generates a random amount of money
func RandomMoney() int64 {
	return RandomInt(minBalance, maxBalance)
}

// RandomCurrency selects a random currency from the predefined list
func RandomCurrency() string {
	currencies := []string{EUR, USD, CAD}
	n := len(currencies)
	return currencies[rand.Intn(n)]
}

// TestAccount represents a complete test account data structure
type TestAccount struct {
	Owner    string
	Balance  int64
	Currency string
}

// RandomAccount generates a complete random account for testing
func RandomAccount() TestAccount {
	return TestAccount{
		Owner:    RandomOwner(),
		Balance:  RandomMoney(),
		Currency: RandomCurrency(),
	}
}

// RandomEmail generates a random email for testing
func RandomEmail() string {
	return fmt.Sprintf("%s@email.com", RandomString(6))
}

// RandomPositiveMoney generates a random positive amount
func RandomPositiveMoney() int64 {
	return RandomInt(1, maxBalance)
}

// RandomNegativeMoney generates a random negative amount
func RandomNegativeMoney() int64 {
	return -RandomInt(1, maxBalance)
}

// package util

// import (
// 	"crypto/rand"
// 	"math/big"
// 	"strings"
// )

// // Constantes para uso em toda a package
// const (
// 	alphabetLower  = "abcdefghijklmnopqrstuvwxyz"
// 	alphabetUpper  = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
// 	alphabetDigits = "0123456789"
// )

// // Currencies disponíveis para geração aleatória
// var currencies = []string{"EUR", "USD", "CAD", "BRL", "GBP", "JPY"}

// // RandomInt gera um número aleatório entre min e max (inclusive)
// // Usa crypto/rand para maior segurança
// func RandomInt(min, max int64) (int64, error) {
// 	if min > max {
// 		min, max = max, min
// 	}

// 	diff := max - min + 1
// 	n, err := rand.Int(rand.Reader, big.NewInt(diff))
// 	if err != nil {
// 		return 0, err
// 	}

// 	return min + n.Int64(), nil
// }

// // RandomString gera uma string aleatória com o tamanho especificado
// // Permite customizar os caracteres usados e inclui validação de entrada
// func RandomString(n int, includeUpper, includeDigits bool) (string, error) {
// 	if n <= 0 {
// 		return "", nil
// 	}

// 	var chars string
// 	chars += alphabetLower
// 	if includeUpper {
// 		chars += alphabetUpper
// 	}
// 	if includeDigits {
// 		chars += alphabetDigits
// 	}

// 	var sb strings.Builder
// 	sb.Grow(n) // Pré-aloca o buffer para melhor performance

// 	for i := 0; i < n; i++ {
// 		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(chars))))
// 		if err != nil {
// 			return "", err
// 		}
// 		sb.WriteByte(chars[num.Int64()])
// 	}

// 	return sb.String(), nil
// }

// // RandomOwner gera um nome de proprietário aleatório
// // Agora com opção de incluir maiúsculas e números
// func RandomOwner() (string, error) {
// 	return RandomString(6, true, false)
// }

// // RandomMoney gera um valor monetário aleatório
// // Adiciona validação de valores negativos
// func RandomMoney() (float64, error) {
// 	cents, err := RandomInt(0, 100000) // Aumentado para permitir valores com centavos
// 	if err != nil {
// 		return 0, err
// 	}
// 	return float64(cents) / 100, nil // Retorna valor com duas casas decimais
// }

// // RandomCurrency seleciona uma moeda aleatória da lista
// func RandomCurrency() (string, error) {
// 	index, err := rand.Int(rand.Reader, big.NewInt(int64(len(currencies))))
// 	if err != nil {
// 		return "", err
// 	}
// 	return currencies[index.Int64()], nil
// }

// // RandomMoney WithCurrency gera um valor monetário com sua moeda
// type MoneyWithCurrency struct {
// 	Amount   float64
// 	Currency string
// }

// func RandomMoneyWithCurrency() (MoneyWithCurrency, error) {
// 	amount, err := RandomMoney()
// 	if err != nil {
// 		return MoneyWithCurrency{}, err
// 	}

// 	currency, err := RandomCurrency()
// 	if err != nil {
// 		return MoneyWithCurrency{}, err
// 	}

// 	return MoneyWithCurrency{
// 		Amount:   amount,
// 		Currency: currency,
// 	}, nil
// }

// // package util

// // import (
// // 	"math/rand"
// // 	"strings"
// // 	"time"
// // )

// // const alphabet = "abcdefghijklmnopqrstuvwxyz"

// // func init() {
// // 	rand.Seed(time.Now().UnixNano())
// // }

// // func RandomInt(min, max int64) int64 {
// // 	return min + rand.Int63n(max-min+1)
// // }

// // func RandomString(n int) string {
// // 	var sb strings.Builder
// // 	k := len(alphabet)

// // 	for i := 0; i < n; i++ {
// // 		c := alphabet[rand.Intn(k)]
// // 		sb.WriteByte(c)
// // 	}
// // 	return sb.String()
// // }

// // func RandomOwner() string {
// // 	return RandomString(6)
// // }

// // func RandomMoney() int64 {
// // 	return RandomInt(0, 1000)
// // }

// // func RandomCurrency() string {
// // 	currencies := []string{"EUR", "USD", "CAD"}
// // 	n := len(currencies)
// // 	return currencies[rand.Intn(n)]
// // }
