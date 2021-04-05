package uniq

import (
	"context"
	"hash/fnv"
	"math/rand"
	"strconv"

	"github.com/borscht/backend/utils/errors"
	"github.com/borscht/backend/utils/logger"
)

func GetUniqFilename(ctx context.Context, filename string) (string, error) {
	// создаем рандомную последовательность чтобы точно названия не повторялись
	hashingSalt := strconv.Itoa(rand.Int() % 1000)

	// создаем хеш от названия файла
	hash := fnv.New32a()
	_, err := hash.Write([]byte(filename + hashingSalt))
	if err != nil {
		custErr := errors.FailServerError(err.Error())
		logger.UsecaseLevel().ErrorLog(ctx, custErr)
		return "", custErr
	}

	return strconv.Itoa(int(hash.Sum32())), nil
}
