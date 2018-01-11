package scraper

import (
	"fmt"
	"github.com/akreal/knn"
)

func knnMain() {

	knn := knn.NewKNN()

	knn.Train("Hello world!", "class1")
	knn.Train("Hello, hello.", "class2")

	k := 1

	predictedClass := knn.Predict("Say hello!", k)

	fmt.Println(predictedClass)

}
