package main

import (
	"fmt"
	"github.com/Knetic/govaluate"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
)

type calculation struct {
	ID         string `json:"id"`
	Expression string `json:"expression"`
	Result     string `json:"result"`
}
type CalculationRequest struct {
	Expression string `json:"expression"`
}

var calculations = []calculation{}

func calculateExpression(expression string) (string, error) {
	expr, err := govaluate.NewEvaluableExpression(expression)
	if err != nil {
		return "", err
	}
	result, err := expr.Evaluate(nil)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%v", result), err
}
func getCalculations(c echo.Context) error {
	return c.JSON(http.StatusOK, calculations)
}
func createCalculations(c echo.Context) error {
	var req CalculationRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid input"})
	}
	result, err := calculateExpression(req.Expression)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid expression"})
	}
	newCalculation := calculation{
		ID:         fmt.Sprintf("%d", len(calculations)+1),
		Expression: req.Expression,
		Result:     result,
	}
	calculations = append(calculations, newCalculation)
	return c.JSON(http.StatusCreated, newCalculation)
}

func main() {
	e := echo.New()
	e.Use(middleware.CORS())
	e.Use(middleware.Logger())
	e.GET("/calculations", getCalculations)
	e.POST("/calculations", createCalculations)
	e.Start("localhost:8080")

}
